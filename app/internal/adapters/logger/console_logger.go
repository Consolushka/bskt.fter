package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type ConsoleLogger struct {
	logrus *logrus.Logger
}

func (c ConsoleLogger) IsEnabled() bool {
	isEnabled := os.Getenv("LOGGER_CONSOLE_ENABLED")

	return isEnabled == "true"
}

func NewConsoleLogger() ConsoleLogger {
	loggerLevel, err := logrus.ParseLevel(os.Getenv("LOGGER_CONSOLE_LEVEL"))
	if err != nil {
		loggerLevel = logrus.InfoLevel
	}

	logrusInstance := logrus.New()
	logrusInstance.SetLevel(loggerLevel)

	return ConsoleLogger{
		logrus: logrusInstance,
	}
}

func (c ConsoleLogger) Info(message string, context map[string]interface{}) {
	c.logrus.WithFields(context).Info(message)
}

func (c ConsoleLogger) Warn(message string, context map[string]interface{}) {
	c.logrus.WithFields(context).Warn(message)
}

func (c ConsoleLogger) Error(message string, context map[string]interface{}) {
	c.logrus.WithFields(context).Error(message)
}
