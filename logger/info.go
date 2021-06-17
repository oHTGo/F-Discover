package logger

import (
	"github.com/sirupsen/logrus"
)

func Info(event string, message string) {
	Init().WithFields(logrus.Fields{
		"event": event,
	}).Infoln(message)
}
