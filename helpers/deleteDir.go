package helpers

import "os"

func DeleteDir(path string) {
	os.Remove(path)
}
