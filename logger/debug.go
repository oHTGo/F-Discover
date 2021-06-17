package logger

import (
	"github.com/sirupsen/logrus"
)

func Debug(event string, data string) {
	Init().WithFields(logrus.Fields{
		"event": event,
		"data":  data,
	}).Debugln()
}
