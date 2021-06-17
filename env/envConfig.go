package env

import (
	"log"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

type single struct {
	LOG_MODE        string
	PORT            string
	STORAGE_BUCKET  string
	MAX_FILE_SIZE   int64
	JWT_SECRET      string
	APP_ID_ZALO     string
	APP_SECRET_ZALO string
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
				singleInstance.LOG_MODE = envConfig["LOG_MODE"]
				singleInstance.PORT = envConfig["PORT"]
				singleInstance.STORAGE_BUCKET = envConfig["STORAGE_BUCKET"]
				singleInstance.MAX_FILE_SIZE, _ = strconv.ParseInt(envConfig["MAX_FILE_SIZE"], 10, 64)
				singleInstance.JWT_SECRET = envConfig["JWT_SECRET"]
				singleInstance.APP_ID_ZALO = envConfig["APP_ID_ZALO"]
				singleInstance.APP_SECRET_ZALO = envConfig["APP_SECRET_ZALO"]
			})
	}
	return singleInstance
}
