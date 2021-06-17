package location

import (
	"f-discover/helpers"
	"strings"
)

func FindByName(name string) []Location {
	var locations []Location

	for key, value := range LOCATIONS {
		valueLower := strings.ToLower(value)
		nameLower := strings.ToLower(name)
		if strings.Contains(
			valueLower,
			nameLower,
		) || strings.Contains(
			helpers.ConvertUnicodeToASCII(valueLower),
			helpers.ConvertUnicodeToASCII(nameLower),
		) {
			locations = append(locations, Location{
				ID:   key,
				Name: value,
			})
		}
	}

	return locations
}
