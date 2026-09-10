package models

// Dog is the API representation. IDs remain stable when a nickname changes.
type Dog struct {
	ID          int     `json:"id"`
	Nickname    string  `json:"nickname"`
	Age         int     `json:"age"`
	WeightKg    float64 `json:"weight_kg"`
	CheckInDate string  `json:"check_in_date"`
	ShelterID   int     `json:"shelter_id"`
	ClinicID    int     `json:"clinic_id"`
}

// Pointers distinguish zero from an omitted field in POST and PUT requests.
type DogInput struct {
	Nickname    string   `json:"nickname"`
	Age         *int     `json:"age"`
	WeightKg    *float64 `json:"weight_kg"`
	CheckInDate string   `json:"check_in_date"`
	ShelterID   int      `json:"shelter_id"`
	ClinicID    int      `json:"clinic_id"`
}
