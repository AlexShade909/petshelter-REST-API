package service

import (
	"PetShelter/internal/models"
	"math"
	"sort"
	"strings"
	"time"
)

type Dog struct{ store *Store }

func NewDog(store *Store) *Dog { return &Dog{store: store} }

func (s *Dog) GetAllDogs() []models.Dog {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	dogs := make([]models.Dog, 0, len(s.store.dogs))
	for _, dog := range s.store.dogs {
		dogs = append(dogs, dog)
	}
	sort.Slice(dogs, func(i, j int) bool { return dogs[i].ID < dogs[j].ID })
	return dogs
}

func (s *Dog) GetDog(id int) (models.Dog, error) {
	s.store.mu.RLock()
	defer s.store.mu.RUnlock()
	dog, ok := s.store.dogs[id]
	if !ok {
		return models.Dog{}, ErrNotFound
	}
	return dog, nil
}

func (s *Dog) validate(input models.DogInput) (models.Dog, error) {
	input.Nickname = strings.TrimSpace(input.Nickname)
	if input.Nickname == "" {
		return models.Dog{}, &ValidationError{"nickname is required"}
	}
	if input.Age == nil || *input.Age < 0 {
		return models.Dog{}, &ValidationError{"age must be a non-negative integer"}
	}
	if input.WeightKg == nil || *input.WeightKg <= 0 || math.IsNaN(*input.WeightKg) || math.IsInf(*input.WeightKg, 0) {
		return models.Dog{}, &ValidationError{"weight_kg must be a positive finite number"}
	}
	if _, err := time.Parse(time.DateOnly, input.CheckInDate); err != nil {
		return models.Dog{}, &ValidationError{"check_in_date must be a valid date in YYYY-MM-DD format"}
	}
	if input.ShelterID < 1 || input.ShelterID > len(s.store.shelters) {
		return models.Dog{}, &ValidationError{"shelter_id must refer to an existing shelter"}
	}
	if input.ClinicID < 1 || input.ClinicID > len(s.store.clinics) {
		return models.Dog{}, &ValidationError{"clinic_id must refer to an existing clinic"}
	}
	return models.Dog{Nickname: input.Nickname, Age: *input.Age, WeightKg: *input.WeightKg,
		CheckInDate: input.CheckInDate, ShelterID: input.ShelterID, ClinicID: input.ClinicID}, nil
}

func (s *Dog) CreateDog(input models.DogInput) (models.Dog, error) {
	dog, err := s.validate(input)
	if err != nil {
		return models.Dog{}, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	dog.ID = s.store.nextDogID
	s.store.nextDogID++
	s.store.dogs[dog.ID] = dog
	return dog, nil
}

// UpdateDog fully replaces the editable fields of an existing dog (PUT).
func (s *Dog) UpdateDog(id int, input models.DogInput) (models.Dog, error) {
	dog, err := s.validate(input)
	if err != nil {
		return models.Dog{}, err
	}
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, ok := s.store.dogs[id]; !ok {
		return models.Dog{}, ErrNotFound
	}
	dog.ID = id
	s.store.dogs[id] = dog
	return dog, nil
}

func (s *Dog) DeleteDog(id int) error {
	s.store.mu.Lock()
	defer s.store.mu.Unlock()
	if _, ok := s.store.dogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.store.dogs, id)
	return nil
}
