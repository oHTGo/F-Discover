package helpers

import (
	"net/http"
	"os"
)

var TYPE_IMAGES = [...]string{"image/jpeg", "image/png", "image/bmp"}
var TYPE_VIDEOS = [...]string{"video/mp4"}

func IsImage(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	contentType, err := getFileContentType(file)
	if err != nil {
		return false
	}

	for _, typeFile := range TYPE_IMAGES {
		if typeFile == contentType {
			return true
		}
	}
	return false
}

func IsVideo(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	contentType, err := getFileContentType(file)
	if err != nil {
		return false
	}

	for _, typeFile := range TYPE_VIDEOS {
		if typeFile == contentType {
			return true
		}
	}
	return false
}

func getFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
