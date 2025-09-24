package logger

import "IMP/app/internal/ports"

var logger CompositeLogger

type CompositeLogger struct {
	loggers []ports.Logger
}

func Init(loggers []ports.Logger) {
	logger = CompositeLogger{
		loggers: loggers,
	}
}

func Info(msg string, ctx map[string]interface{}) {
	for _, logger := range logger.loggers {
		logger.Info(msg, ctx)
	}
}

func Warn(msg string, ctx map[string]interface{}) {
	for _, logger := range logger.loggers {
		logger.Warn(msg, ctx)
	}
}

func Error(msg string, ctx map[string]interface{}) {
	for _, logger := range logger.loggers {
		logger.Error(msg, ctx)
	}
}
