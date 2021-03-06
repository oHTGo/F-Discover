package helpers

import (
	"f-discover/env"
	"f-discover/instance"
	"f-discover/services"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

func UploadFileStorage(filePath string, object string) (string, error) {
	client := services.GetInstance().StorageClient
	bucket := services.GetInstance().StorageBucket

	// Open local file.
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write to Google Storage
	writer := bucket.Object(object).NewWriter(instance.CtxBackground)
	if _, err = io.Copy(writer, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	objectStorage := client.Bucket(env.Get().STORAGE_BUCKET).Object(object)
	// Set public access
	acl := objectStorage.ACL()
	if err := acl.Set(instance.CtxBackground, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	// No cache
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		CacheControl: "no-store",
	}
	if _, err := objectStorage.Update(instance.CtxBackground, objectAttrsToUpdate); err != nil {
		return "", err
	}

	url := "https://storage.googleapis.com/" + env.Get().STORAGE_BUCKET + "/" + writer.Name
	return url, nil
}
