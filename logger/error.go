package logger

import (
	"html"

	"github.com/sirupsen/logrus"
)

func Error(event string, message string, err error) {
	if err != nil {
		Init().WithFields(logrus.Fields{
			"event":  event,
			"errors": html.EscapeString(err.Error()),
		}).Errorln(message)

		return
	}

	Init().WithFields(logrus.Fields{
		"event": event,
	}).Errorln(message)
}
