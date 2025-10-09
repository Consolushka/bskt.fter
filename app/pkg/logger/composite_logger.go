package logger

import "IMP/app/internal/ports"

var instance CompositeLogger

type CompositeLogger struct {
	loggers []ports.Logger
}

func Init(loggers []ports.Logger) {
	instance = CompositeLogger{
		loggers: loggers,
	}
}

func Info(msg string, ctx map[string]interface{}) {
	for _, logger := range instance.loggers {
		logger.Info("[INFO] "+msg, ctx)
	}
}

func Warn(msg string, ctx map[string]interface{}) {
	for _, logger := range instance.loggers {
		logger.Warn("[WARNING] "+msg, ctx)
	}
}

func Error(msg string, ctx map[string]interface{}) {
	for _, logger := range instance.loggers {
		logger.Error("[ERROR] "+msg, ctx)
	}
}
