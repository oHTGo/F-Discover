package helpers

import (
	"f-discover/env"
	"f-discover/instance"
	"f-discover/services"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

func UploadFile(filePath string, object string) (string, error) {
	client := services.GetInstance().StorageClient
	bucket := services.GetInstance().StorageBucket

	// Open local file.
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bucket.Object(object).NewWriter(instance.CtxBackground)
	if _, err = io.Copy(writer, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	acl := client.Bucket(env.Get().STORAGE_BUCKET).Object(object).ACL()
	if err := acl.Set(instance.CtxBackground, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}
	url := "https://storage.googleapis.com/" + env.Get().STORAGE_BUCKET + "/" + writer.Name
	return url, nil
}
