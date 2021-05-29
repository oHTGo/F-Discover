package helpers

import "os"

func CreateDir(dirPath string) {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		os.Mkdir(dirPath, 0755)
	}
}
