package internal

func AddDog(dogs map[string]Dog, nickname string, age string, weightKg string, checkInDate string, shelter *Shelter, policlinic *Policlinic) *Dog {
	d := &Dog{
		Nickname:    nickname,
		Age:         age,
		WeightKg:    weightKg,
		CheckInDate: checkInDate,
		Shelter:     shelter,
		Policlinic:  policlinic,
	}
	dogs[nickname] = *d
	return d
}

func RemoveDog(dogs map[string]Dog, nickname string) bool {
	_, ok := dogs[nickname]
	if !ok {
		return false
	}
	delete(dogs, nickname)
	return true
}

func FindDog(dogs map[string]Dog, nickname string) (Dog, bool) {
	d, ok := dogs[nickname]
	return d, ok
}
