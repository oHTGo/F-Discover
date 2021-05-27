package firebase

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var once sync.Once

type single struct {
	FirebaseApp *firebase.App
	AuthClient  *auth.Client
	StoreClient *firestore.Client
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
			})
	}
	return singleInstance
}

func initialize() *firebase.App {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing app: %v\n", err)
	}

	return app
}

func initialFirestore(app *firebase.App) *firestore.Client {
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Error initializing firestore: %v\n", err)
	}

	return client
}

func initialAuth(app *firebase.App) *auth.Client {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error initializing auth: %v\n", err)
	}

	return client
}
