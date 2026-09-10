package service

import (
	"PetShelter/internal"
	"PetShelter/internal/models"
	"errors"
	"sort"
	"strconv"
	"sync"
	"time"
)

var ErrNotFound = errors.New("resource not found")

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

// Store is shared by all requests. Clinics and shelters are read-only.
// Dog access is protected by mu. Data is not persisted across restarts.
type Store struct {
	mu        sync.RWMutex
	dogs      map[int]models.Dog
	clinics   []models.Policlinic
	shelters  []models.Shelter
	nextDogID int
}

func NewStore() *Store {
	s := &Store{dogs: make(map[int]models.Dog), nextDogID: 1}
	clinics := internal.CreatePoliclinics()
	shelters := internal.CreateShelters()
	for i, c := range clinics {
		s.clinics = append(s.clinics, models.Policlinic{ID: i + 1, NumberClinic: c.NumberClinic,
			Address: c.Address, PhoneNumber: c.PhoneNumber, WorkingTime: c.WorkingTime})
	}
	for i, sh := range shelters {
		s.shelters = append(s.shelters, models.Shelter{ID: i + 1, NumberShelter: sh.NumberShelter,
			Address: sh.Address, Number: sh.Number, WorkingTime: sh.WorkingTime})
	}
	// Reuse the original examples, assigning deterministic IDs and ISO dates.
	dogs := internal.CreateDogs(shelters, clinics)
	names := make([]string, 0, len(dogs))
	for name := range dogs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		d := dogs[name]
		age, _ := strconv.Atoi(d.Age)
		weight, _ := strconv.ParseFloat(d.WeightKg, 64)
		date, _ := time.Parse("02.01.2006", d.CheckInDate)
		var shelterID, clinicID int
		for i := range shelters {
			if d.Shelter == &shelters[i] {
				shelterID = i + 1
			}
		}
		for i := range clinics {
			if d.Policlinic == &clinics[i] {
				clinicID = i + 1
			}
		}
		s.dogs[s.nextDogID] = models.Dog{ID: s.nextDogID, Nickname: d.Nickname,
			Age: age, WeightKg: weight, CheckInDate: date.Format(time.DateOnly),
			ShelterID: shelterID, ClinicID: clinicID}
		s.nextDogID++
	}
	return s
}
