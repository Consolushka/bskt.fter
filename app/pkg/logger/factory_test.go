package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildLoggers(t *testing.T) {
	t.Run("Default to console when nothing is enabled", func(t *testing.T) {
		t.Setenv("LOGGER_CONSOLE_ENABLED", "false")
		t.Setenv("LOGGER_FILE_ENABLED", "false")
		t.Setenv("LOGGER_TELEGRAM_ENABLED", "false")

		loggers := BuildLoggers()

		assert.Len(t, loggers, 1)
	})

	t.Run("Enable only console", func(t *testing.T) {
		t.Setenv("LOGGER_CONSOLE_ENABLED", "true")
		t.Setenv("LOGGER_FILE_ENABLED", "false")
		t.Setenv("LOGGER_TELEGRAM_ENABLED", "false")

		loggers := BuildLoggers()

		assert.Len(t, loggers, 1)
	})

	t.Run("Enable only file logger", func(t *testing.T) {
		tmpFile := t.TempDir() + "/test.log"
		t.Setenv("LOGGER_CONSOLE_ENABLED", "false")
		t.Setenv("LOGGER_FILE_ENABLED", "true")
		t.Setenv("LOGGER_FILE_PATH", tmpFile)
		t.Setenv("LOGGER_TELEGRAM_ENABLED", "false")

		loggers := BuildLoggers()

		// Based on logic:
		// 1. File added (len=1)
		// 2. Telegram skipped
		// 3. Console check: consoleEnabled(false) || len(loggers)==0(false) -> skipped
		assert.Len(t, loggers, 1)
	})

	t.Run("Enable file and console loggers", func(t *testing.T) {
		tmpFile := t.TempDir() + "/test.log"
		t.Setenv("LOGGER_CONSOLE_ENABLED", "true")
		t.Setenv("LOGGER_FILE_ENABLED", "true")
		t.Setenv("LOGGER_FILE_PATH", tmpFile)
		t.Setenv("LOGGER_TELEGRAM_ENABLED", "false")

		loggers := BuildLoggers()

		assert.Len(t, loggers, 2)
	})
}
