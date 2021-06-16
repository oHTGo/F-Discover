package location

func CheckValidation(id string) bool {
	if _, ok := LOCATIONS[id]; !ok {
		return false
	}

	return true
}
