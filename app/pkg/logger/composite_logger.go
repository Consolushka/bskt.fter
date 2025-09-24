package logger

import "IMP/app/internal/ports"

var Logger CompositeLogger

type CompositeLogger struct {
	loggers []ports.Logger
}

func Init(loggers []ports.Logger) {
	Logger = CompositeLogger{
		loggers: loggers,
	}
}

func (c CompositeLogger) Info(msg string, ctx map[string]interface{}) {
	for _, logger := range c.loggers {
		logger.Info(msg, ctx)
	}
}

func (c CompositeLogger) Warn(msg string, ctx map[string]interface{}) {
	for _, logger := range c.loggers {
		logger.Warn(msg, ctx)
	}
}

func (c CompositeLogger) Error(msg string, ctx map[string]interface{}) {
	for _, logger := range c.loggers {
		logger.Error(msg, ctx)
	}
}
