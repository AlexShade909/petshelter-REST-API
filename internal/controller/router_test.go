package controller_test

import (
	"PetShelter/internal/controller"
	"PetShelter/internal/models"
	"PetShelter/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

const validDog = `{"nickname":"Бобик","age":0,"weight_kg":2.5,"check_in_date":"2026-09-09","shelter_id":1,"clinic_id":2}`

func newHandler() http.Handler {
	store := service.NewStore()
	return controller.NewRouter(controller.NewDog(service.NewDog(store)),
		controller.NewClinic(service.NewClinic(store)), controller.NewShelter(service.NewShelter(store)))
}

func request(h http.Handler, method, path, body, contentType string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestDogLifecycle(t *testing.T) {
	h := newHandler()
	created := request(h, "POST", "/dogs", validDog, "application/json")
	if created.Code != 201 {
		t.Fatalf("create: %d %s", created.Code, created.Body.String())
	}
	var dog models.Dog
	if err := json.Unmarshal(created.Body.Bytes(), &dog); err != nil {
		t.Fatal(err)
	}
	path := fmt.Sprintf("/dogs/%d", dog.ID)
	if dog.ID <= 0 || dog.Age != 0 || created.Header().Get("Location") != path {
		t.Fatalf("invalid created dog or Location: %+v", dog)
	}
	duplicate := request(h, "POST", "/dogs", validDog, "application/json")
	var other models.Dog
	if err := json.Unmarshal(duplicate.Body.Bytes(), &other); err != nil {
		t.Fatal(err)
	}
	if duplicate.Code != 201 || other.ID == dog.ID {
		t.Fatal("same nickname must produce a separate dog")
	}
	updated := request(h, "PUT", path, strings.Replace(validDog, "Бобик", "Рекс", 1), "application/json")
	if updated.Code != 200 {
		t.Fatalf("update: %d %s", updated.Code, updated.Body.String())
	}
	found := request(h, "GET", path, "", "")
	if err := json.Unmarshal(found.Body.Bytes(), &dog); err != nil {
		t.Fatal(err)
	}
	if found.Code != 200 || dog.Nickname != "Рекс" {
		t.Fatalf("get updated: %d %+v", found.Code, dog)
	}
	// A rejected update must not change the stored object.
	bad := request(h, "PUT", path, `{}`, "application/json")
	if bad.Code != 400 || request(h, "GET", path, "", "").Body.String() != found.Body.String() {
		t.Fatal("invalid update mutated the dog")
	}
	deleted := request(h, "DELETE", path, "", "")
	if deleted.Code != 204 || deleted.Body.Len() != 0 {
		t.Fatalf("delete: %d %s", deleted.Code, deleted.Body.String())
	}
	for _, method := range []string{"GET", "DELETE", "PUT"} {
		if got := request(h, method, path, validDog, "application/json"); got.Code != 404 {
			t.Fatalf("%s deleted dog: %d", method, got.Code)
		}
	}
}

func TestRoutesAndErrors(t *testing.T) {
	h := newHandler()
	cases := []struct {
		name, method, path, body, contentType string
		status                                int
	}{
		{"dogs", "GET", "/dogs", "", "", 200},
		{"clinics", "GET", "/clinics/", "", "", 200},
		{"clinic", "GET", "/clinics/1", "", "", 200},
		{"shelters", "GET", "/shelters", "", "", 200},
		{"shelter", "GET", "/shelters/2", "", "", 200},
		{"missing clinic", "GET", "/clinics/999", "", "", 404},
		{"missing shelter", "GET", "/shelters/999", "", "", 404},
		{"unknown route", "GET", "/missing", "", "", 404},
		{"not a collection", "GET", "/clinics/1/extra", "", "", 404},
		{"invalid ID", "GET", "/dogs/abc", "", "", 400},
		{"negative ID", "GET", "/dogs/-1", "", "", 400},
		{"wrong method", "PATCH", "/dogs/1", "", "", 405},
		{"read-only clinic", "POST", "/clinics", "", "", 405},
		{"media type", "POST", "/dogs", validDog, "text/plain", 415},
		{"empty", "POST", "/dogs", "", "application/json", 400},
		{"null", "POST", "/dogs", "null", "application/json", 400},
		{"unknown field", "POST", "/dogs", `{"unknown":1}`, "application/json", 400},
		{"extra JSON", "POST", "/dogs", validDog + `{}`, "application/json", 400},
		{"missing age", "POST", "/dogs", strings.Replace(validDog, `"age":0,`, "", 1), "application/json", 400},
		{"negative age", "POST", "/dogs", strings.Replace(validDog, `"age":0`, `"age":-1`, 1), "application/json", 400},
		{"weight", "POST", "/dogs", strings.Replace(validDog, "2.5", "0", 1), "application/json", 400},
		{"date", "POST", "/dogs", strings.Replace(validDog, "2026-09-09", "2026-02-30", 1), "application/json", 400},
		{"reference", "POST", "/dogs", strings.Replace(validDog, `"shelter_id":1`, `"shelter_id":99`, 1), "application/json", 400},
		{"oversized", "POST", "/dogs", `{"nickname":"` + strings.Repeat("x", 1<<20) + `"}`, "application/json", 413},
		{"oversized whitespace", "POST", "/dogs", validDog + strings.Repeat(" ", 1<<20), "application/json", 413},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := request(h, tc.method, tc.path, tc.body, tc.contentType)
			if w.Code != tc.status {
				t.Fatalf("status: got %d want %d; %s", w.Code, tc.status, w.Body.String())
			}
			if !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") || !json.Valid(w.Body.Bytes()) {
				t.Fatalf("not a JSON response: %s", w.Body.String())
			}
			if tc.status == 405 && w.Header().Get("Allow") == "" {
				t.Fatal("missing Allow header")
			}
		})
	}
}

func TestConcurrentRequests(t *testing.T) {
	h := newHandler()
	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := request(h, "POST", "/dogs", validDog, "application/json")
			if w.Code != 201 {
				t.Errorf("create: %d", w.Code)
			}
			if w := request(h, "GET", "/dogs", "", ""); w.Code != 200 {
				t.Errorf("list: %d", w.Code)
			}
		}()
	}
	wg.Wait()
	var dogs []models.Dog
	w := request(h, "GET", "/dogs", "", "")
	if err := json.Unmarshal(w.Body.Bytes(), &dogs); err != nil {
		t.Fatal(err)
	}
	if len(dogs) != 33 {
		t.Fatalf("got %d dogs, want 33", len(dogs))
	}
	ids := make(map[int]bool)
	for _, dog := range dogs {
		if ids[dog.ID] {
			t.Fatalf("duplicate ID %d", dog.ID)
		}
		ids[dog.ID] = true
		if w := request(h, "DELETE", fmt.Sprintf("/dogs/%d", dog.ID), "", ""); w.Code != 204 {
			t.Fatalf("delete: %d", w.Code)
		}
	}
	if body := request(h, "GET", "/dogs", "", "").Body.String(); strings.TrimSpace(body) != "[]" {
		t.Fatalf("empty list should be [], got %s", body)
	}
}
