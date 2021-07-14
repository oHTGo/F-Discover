package env

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

type single struct {
	LOG_MODE          string
	URL               string
	PORT              string
	STORAGE_BUCKET    string
	MAX_FILE_SIZE     int64
	VIDEO_COMPRESSION bool
	LOCAL_UPLOAD      bool
	JWT_SECRET        string
	APP_ID_ZALO       string
	APP_SECRET_ZALO   string
}

var singleInstance *single

func Get() *single {
	if singleInstance == nil {
		once.Do(
			func() {
				err := godotenv.Load(".env")
				if err != nil {
					log.Fatal("Error loading .env file")
				}

				singleInstance = new(single)
				singleInstance.LOG_MODE = os.Getenv("LOG_MODE")
				singleInstance.URL = os.Getenv("URL")
				singleInstance.PORT = os.Getenv("PORT")
				singleInstance.STORAGE_BUCKET = os.Getenv("STORAGE_BUCKET")
				singleInstance.MAX_FILE_SIZE, _ = strconv.ParseInt(os.Getenv("MAX_FILE_SIZE"), 10, 64)
				singleInstance.VIDEO_COMPRESSION, _ = strconv.ParseBool(os.Getenv("VIDEO_COMPRESSION"))
				singleInstance.LOCAL_UPLOAD, _ = strconv.ParseBool(os.Getenv("LOCAL_UPLOAD"))
				singleInstance.JWT_SECRET = os.Getenv("JWT_SECRET")
				singleInstance.APP_ID_ZALO = os.Getenv("APP_ID_ZALO")
				singleInstance.APP_SECRET_ZALO = os.Getenv("APP_SECRET_ZALO")
			})
	}
	return singleInstance
}
