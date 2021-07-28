package handler

import (
	"errors"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var (
	ErrUnknown = errors.New("unknown error")
)

func getTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger
}
