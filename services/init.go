package services

import (
	"f-discover/env"
	"f-discover/instance"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"google.golang.org/api/option"
)

var once sync.Once

type single struct {
	FirebaseApp   *firebase.App
	AuthClient    *auth.Client
	StoreClient   *firestore.Client
	StorageClient *storage.Client
	StorageBucket *storage.BucketHandle
}

var singleInstance *single

func GetInstance() *single {
	if singleInstance == nil {
		once.Do(
			func() {
				singleInstance = new(single)
				singleInstance.FirebaseApp = initialize()
				singleInstance.AuthClient = initialAuth(singleInstance.FirebaseApp)
				singleInstance.StoreClient = initialFirestore(singleInstance.FirebaseApp)
				singleInstance.StorageClient, singleInstance.StorageBucket = initialStorage()
			})
	}
	return singleInstance
}

func initialize() *firebase.App {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(instance.CtxBackground, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}

	return app
}

func initialFirestore(app *firebase.App) *firestore.Client {
	client, err := app.Firestore(instance.CtxBackground)
	if err != nil {
		log.Fatalf("Error initializing firestore: %v\n", err)
	}

	return client
}

func initialAuth(app *firebase.App) *auth.Client {
	client, err := app.Auth(instance.CtxBackground)
	if err != nil {
		log.Fatalf("Error initializing auth: %v\n", err)
	}

	return client
}

func initialStorage() (*storage.Client, *storage.BucketHandle) {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	client, err := storage.NewClient(instance.CtxBackground, opt)
	if err != nil {
		log.Fatalf("Error initializing storage: %v\n", err)
	}

	return client, client.Bucket(env.Get().STORAGE_BUCKET)
}
