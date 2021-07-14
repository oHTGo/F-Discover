package helpers

import (
	"f-discover/env"
	"path/filepath"
	"strings"
)

const LOCAL_STORAGE_DIR = "files"

func UploadFileLocal(path string, object string) (string, error) {
	localPath := filepath.Join(LOCAL_STORAGE_DIR, object)

	var dirPath string = ""
	for _, value := range strings.Split(filepath.Dir(localPath), "/") {
		dirPath += value + "/"
		CreateDir(dirPath)
	}

	if err := CopyFile(path, localPath); err != nil {
		return "", err
	}

	return env.Get().URL + "/" + LOCAL_STORAGE_DIR + "/" + object, nil
}
