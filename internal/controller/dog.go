package controller

import (
	"PetShelter/internal/models"
	"net/http"
	"strconv"
)

type dogService interface {
	GetAllDogs() []models.Dog
	GetDog(int) (models.Dog, error)
	CreateDog(models.DogInput) (models.Dog, error)
	UpdateDog(int, models.DogInput) (models.Dog, error)
	DeleteDog(int) error
}
type Dog struct{ dogService dogService }

func NewDog(s dogService) *Dog { return &Dog{dogService: s} }
func (c *Dog) GetAll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, c.dogService.GetAllDogs())
}
func (c *Dog) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	dog, err := c.dogService.GetDog(id)
	if err != nil {
		serviceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dog)
}
func (c *Dog) Create(w http.ResponseWriter, r *http.Request) {
	var input models.DogInput
	if !readJSON(w, r, &input) {
		return
	}
	dog, err := c.dogService.CreateDog(input)
	if err != nil {
		serviceError(w, err)
		return
	}
	w.Header().Set("Location", "/dogs/"+strconv.Itoa(dog.ID))
	writeJSON(w, http.StatusCreated, dog)
}
func (c *Dog) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	var input models.DogInput
	if !readJSON(w, r, &input) {
		return
	}
	dog, err := c.dogService.UpdateDog(id, input)
	if err != nil {
		serviceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dog)
}
func (c *Dog) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := c.dogService.DeleteDog(id); err != nil {
		serviceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
