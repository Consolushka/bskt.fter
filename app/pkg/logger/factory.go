package logger

import (
	loggerAdapters "IMP/app/internal/adapters/logger"
	"IMP/app/internal/ports"
)

func BuildLoggers() []ports.Logger {
	loggers := make([]ports.Logger, 1)

	fileLoggerEnabled := loggerAdapters.FileLogger{}.IsEnabled()
	if fileLoggerEnabled {
		loggers = append(loggers, loggerAdapters.NewFileLogger())
	}

	telegramLoggerEnabled := loggerAdapters.TelegramLogger{}.IsEnabled()
	if telegramLoggerEnabled {
		loggers = append(loggers, loggerAdapters.NewTelegramLogger())
	}

	consoleLoggerEnabled := loggerAdapters.ConsoleLogger{}.IsEnabled()
	if consoleLoggerEnabled || len(loggers) == 0 {
		loggers = append(loggers, loggerAdapters.NewConsoleLogger())
	}

	return loggers
}
