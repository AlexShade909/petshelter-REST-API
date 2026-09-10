package service

import "PetShelter/internal/models"

type Shelter struct{ store *Store }

func NewShelter(store *Store) *Shelter { return &Shelter{store: store} }
func (s *Shelter) GetAllShelters() []models.Shelter {
	return append([]models.Shelter{}, s.store.shelters...)
}
func (s *Shelter) GetShelter(id int) (models.Shelter, error) {
	for _, shelter := range s.store.shelters {
		if shelter.ID == id {
			return shelter, nil
		}
	}
	return models.Shelter{}, ErrNotFound
}
