package models

type Policlinic struct {
	ID           int `json:"id"`
	NumberClinic string
	Address      string
	PhoneNumber  string
	WorkingTime  string
}
