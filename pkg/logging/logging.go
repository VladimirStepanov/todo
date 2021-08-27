package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func GetLogger(logLevel string, logFile string) (*logrus.Logger, error) {

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetFormatter(
		&logrus.JSONFormatter{},
	)
	logger.SetLevel(level)

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, file)

	logger.SetOutput(mw)

	return logger, nil
}
