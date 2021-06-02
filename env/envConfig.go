package env

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

type single struct {
	PORT           string
	STORAGE_BUCKET string
	JWT_SECRET     string
}

var singleInstance *single

func Get() *single {
	if singleInstance == nil {
		once.Do(
			func() {
				envConfig, err := godotenv.Read(".env")
				if err != nil {
					log.Fatal("Error loading .env file")
				}

				singleInstance = new(single)
				singleInstance.PORT = envConfig["PORT"]
				singleInstance.STORAGE_BUCKET = envConfig["STORAGE_BUCKET"]
				singleInstance.JWT_SECRET = envConfig["JWT_SECRET"]
			})
	}
	return singleInstance
}
