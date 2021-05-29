package helpers

import (
	"context"
	"f-discover/services"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
)

func UploadFile(filePath string, object string) (string, error) {
	ctx := context.Background()
	envConfig, _ := godotenv.Read(".env")

	client := services.GetInstance().StorageClient
	bucket := services.GetInstance().StorageBucket

	// Open local file.
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := bucket.Object(object).NewWriter(ctx)
	if _, err = io.Copy(writer, file); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}
	acl := client.Bucket(envConfig["STORAGE_BUCKET"]).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}
	url := "https://storage.googleapis.com/" + envConfig["STORAGE_BUCKET"] + "/" + writer.Name
	return url, nil
}
