package log

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func New() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	log.Level = logrus.DebugLevel

	return log
}
