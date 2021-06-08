package helpers

import (
	"errors"
	"f-discover/env"
	"mime/multipart"
	"strconv"

	"github.com/kataras/iris/v12"
)

func UploadFiles(ctx iris.Context) ([]*multipart.FileHeader, error) {
	// Create dir uploads if it is inexistent
	CreateDir("uploads")

	files, _, err := ctx.UploadFormFiles("./uploads")
	if err != nil {
		return nil, errors.New("Upload files failed")
	}

	for _, file := range files {
		if file.Size > (env.Get().MAX_FILE_SIZE * iris.MB) {
			return nil, errors.New("File exceeds the allowed size (" + strconv.FormatInt(env.Get().MAX_FILE_SIZE, 10) + " MB)")
		}
	}

	return files, nil
}
