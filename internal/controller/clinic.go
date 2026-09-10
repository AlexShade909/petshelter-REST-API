package controller

import (
	"PetShelter/internal/models"
	"net/http"
)

type clinicService interface {
	GetAllClinics() []models.Policlinic
	GetClinic(int) (models.Policlinic, error)
}
type Clinic struct{ service clinicService }

func NewClinic(s clinicService) *Clinic { return &Clinic{service: s} }
func (c *Clinic) GetAll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, c.service.GetAllClinics())
}
func (c *Clinic) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	item, err := c.service.GetClinic(id)
	if err != nil {
		serviceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
