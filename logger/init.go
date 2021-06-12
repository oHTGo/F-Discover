package logger

import (
	"f-discover/env"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var once sync.Once

type single struct {
	LOGGER *logrus.Logger
}

var singleInstance *single

func Init() *logrus.Logger {
	if singleInstance == nil {
		once.Do(
			func() {
				log := logrus.New()
				log.Out = os.Stdout

				if env.Get().LOG_MODE == "file" {
					file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
					if err == nil {
						log.Out = file
					} else {
						log.Warn("Failed to log to file, using default stderr")
					}
				}

				singleInstance = new(single)
				singleInstance.LOGGER = log
			})
	}
	return singleInstance.LOGGER
}
