package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logger interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Fatalln(args ...interface{})
}

type logger struct {
	log *logrus.Logger
}

func NewLogger() Logger {
	loggerEntity := logrus.New()

	loggerEntity.SetLevel(logrus.DebugLevel)

	loggerEntity.SetFormatter(&logrus.JSONFormatter{})

	// set app.log file as logs output
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		loggerEntity.SetOutput(file)
	} else {
		loggerEntity.Info("Не удалось открыть файл логов, используется стандартный stderr")
	}

	return &logger{
		log: loggerEntity,
	}
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

// Fatalln is equivalent to [Println] followed by a call to [os.Exit](1).
func (l *logger) Fatalln(args ...interface{}) {
	l.log.Fatalln(args...)
}
