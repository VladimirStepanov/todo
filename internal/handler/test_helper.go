package handler

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func getTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger
}
