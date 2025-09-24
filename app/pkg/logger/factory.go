package logger

import (
	"IMP/app/internal/adapters/logger"
	"IMP/app/internal/ports"
)

func BuildLoggers() []ports.Logger {
	loggers := make([]ports.Logger, 1)

	fileLoggerEnabled := logger.FileLogger{}.IsEnabled()
	if fileLoggerEnabled {
		loggers = append(loggers, logger.NewFileLogger())
	}

	telegramLoggerEnabled := logger.TelegramLogger{}.IsEnabled()
	if telegramLoggerEnabled {
		loggers = append(loggers, logger.NewTelegramLogger())
	}

	consoleLoggerEnabled := logger.ConsoleLogger{}.IsEnabled()
	if consoleLoggerEnabled || len(loggers) == 0 {
		loggers = append(loggers, logger.NewConsoleLogger())
	}

	return loggers
}
