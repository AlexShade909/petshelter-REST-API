package controller

import "net/http"

func NewRouter(dog *Dog, clinic *Clinic, shelter *Shelter) http.Handler {
	mux := http.NewServeMux()
	for _, path := range []string{"/dogs", "/dogs/{$}"} {
		mux.HandleFunc("GET "+path, dog.GetAll)
		mux.HandleFunc("POST "+path, dog.Create)
		mux.HandleFunc(path, methodNotAllowed("GET, HEAD, POST"))
	}
	mux.HandleFunc("GET /dogs/{id}", dog.Get)
	mux.HandleFunc("PUT /dogs/{id}", dog.Update)
	mux.HandleFunc("DELETE /dogs/{id}", dog.Delete)
	mux.HandleFunc("/dogs/{id}", methodNotAllowed("GET, HEAD, PUT, DELETE"))
	for _, path := range []string{"/clinics", "/clinics/{$}"} {
		mux.HandleFunc("GET "+path, clinic.GetAll)
		mux.HandleFunc(path, methodNotAllowed("GET, HEAD"))
	}
	mux.HandleFunc("GET /clinics/{id}", clinic.Get)
	mux.HandleFunc("/clinics/{id}", methodNotAllowed("GET, HEAD"))
	for _, path := range []string{"/shelters", "/shelters/{$}"} {
		mux.HandleFunc("GET "+path, shelter.GetAll)
		mux.HandleFunc(path, methodNotAllowed("GET, HEAD"))
	}
	mux.HandleFunc("GET /shelters/{id}", shelter.Get)
	mux.HandleFunc("/shelters/{id}", methodNotAllowed("GET, HEAD"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "route not found")
	})
	return mux
}
func methodNotAllowed(allow string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", allow)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
