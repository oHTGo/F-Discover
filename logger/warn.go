package logger

import (
	"html"

	"github.com/sirupsen/logrus"
)

func Warn(event string, message string, err error) {
	Init().WithFields(logrus.Fields{
		"event":  event,
		"errors": html.EscapeString(err.Error()),
	}).Warnln(message)
}
