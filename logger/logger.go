package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init() {
	if log != nil {

		return
	}
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
}

func L() *logrus.Entry {
	return logrus.NewEntry(log)
}
