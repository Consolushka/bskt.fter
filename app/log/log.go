package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func Init() {
	log = logrus.New()

	log.SetLevel(logrus.DebugLevel)

	log.SetFormatter(&logrus.JSONFormatter{})

	// set app.log file as logs output
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Не удалось открыть файл логов, используется стандартный stderr")
	}
}

func GetLogger() *logrus.Logger {
	return log
}
