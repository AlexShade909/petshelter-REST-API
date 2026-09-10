package controller

import (
	"PetShelter/internal/service"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("encode response: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(append(data, '\n')); err != nil {
		log.Printf("write response: %v", err)
	}
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
func serviceError(w http.ResponseWriter, err error) {
	var validation *service.ValidationError
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.As(err, &validation):
		writeError(w, http.StatusBadRequest, validation.Error())
	default:
		log.Printf("service error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}
func readJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(dst)
	if err == nil {
		var extra any
		err = decoder.Decode(&extra)
		if errors.Is(err, io.EOF) {
			return true
		}
	}
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeError(w, http.StatusRequestEntityTooLarge, "request body exceeds 1 MiB")
	} else {
		writeError(w, http.StatusBadRequest, "body must contain one JSON object with known fields and valid types")
	}
	return false
}
