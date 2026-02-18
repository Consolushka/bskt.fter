package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type FileLogger struct {
	logrus *logrus.Logger
}

func (f FileLogger) IsEnabled() bool {
	isEnabled := os.Getenv("LOGGER_FILE_ENABLED")

	return isEnabled == "true"
}

func NewFileLogger() FileLogger {
	filePath := os.Getenv("LOGGER_FILE_PATH")
	if filePath == "" {
		panic("LOGGER_FILE_PATH is not set")
	}

	loggerLevel, err := logrus.ParseLevel(os.Getenv("LOGGER_FILE_LEVEL"))
	if err != nil {
		loggerLevel = logrus.InfoLevel
	}

	logrusInstance := logrus.New()
	logrusInstance.SetLevel(loggerLevel)

	logDir := filepath.Dir(filePath)
	if err = os.MkdirAll(logDir, 0755); err != nil {
		logrusInstance.Fatalf("Failed to create logrusInstance directory: %v", err)
	}

	// Открываем файл для записи логов
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrusInstance.Fatalf("Failed to open logrusInstance file: %v", err)
	}

	// Настраиваем вывод в файл и консоль
	mw := io.MultiWriter(os.Stdout, logFile)
	logrusInstance.SetOutput(mw)

	return FileLogger{
		logrus: logrusInstance,
	}
}

func (f FileLogger) Info(message string, context map[string]interface{}) {
	f.logrus.WithFields(context).Info(message)
}

func (f FileLogger) Warn(message string, context map[string]interface{}) {
	f.logrus.WithFields(context).Warn(message)
}

func (f FileLogger) Error(message string, context map[string]interface{}) {
	f.logrus.WithFields(context).Error(message)
}

func (f FileLogger) Fatal(message string, context map[string]interface{}) {
	f.logrus.WithFields(context).Log(logrus.FatalLevel, message)
}
