package helpers

import "os"

func CreateDir(dirPath string) {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}
}
