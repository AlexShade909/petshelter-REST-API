package controller

import (
	"PetShelter/internal/models"
	"net/http"
)

type shelterService interface {
	GetAllShelters() []models.Shelter
	GetShelter(int) (models.Shelter, error)
}
type Shelter struct{ service shelterService }

func NewShelter(s shelterService) *Shelter { return &Shelter{service: s} }
func (c *Shelter) GetAll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, c.service.GetAllShelters())
}
func (c *Shelter) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	item, err := c.service.GetShelter(id)
	if err != nil {
		serviceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
