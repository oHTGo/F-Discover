package location

import "errors"

func CheckValidation(id string) bool {
	if _, ok := LOCATIONS[id]; !ok {
		return false
	}

	return true
}

func CheckValidationForValidator(value interface{}) error {
	str, _ := value.(string)
	if !CheckValidation(str) {
		return errors.New("is invalid")
	}
	return nil
}
