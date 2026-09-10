package service

import "PetShelter/internal/models"

type Clinic struct{ store *Store }

func NewClinic(store *Store) *Clinic { return &Clinic{store: store} }
func (c *Clinic) GetAllClinics() []models.Policlinic {
	return append([]models.Policlinic{}, c.store.clinics...)
}
func (c *Clinic) GetClinic(id int) (models.Policlinic, error) {
	for _, clinic := range c.store.clinics {
		if clinic.ID == id {
			return clinic, nil
		}
	}
	return models.Policlinic{}, ErrNotFound
}
