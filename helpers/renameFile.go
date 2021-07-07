package helpers

import "os"

func RenameFile(originalPath string, newPath string) error {
	err := os.Rename(originalPath, newPath)
	if err != nil {
		return err
	}

	return nil
}
