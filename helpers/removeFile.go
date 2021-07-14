package helpers

import (
	"os"
)

func RemoveFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}
