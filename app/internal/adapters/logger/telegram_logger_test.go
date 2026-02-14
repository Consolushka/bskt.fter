package logger

import (
	"encoding/json"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeLogContext_ConvertsErrorToString(t *testing.T) {
	context := map[string]interface{}{
		"taskType": "process_not_urgent_european_tournaments_task",
		"error":    errors.New("provider timeout"),
		"nested": map[string]interface{}{
			"err": errors.New("nested failure"),
		},
	}

	normalized := normalizeLogContext(context)
	raw, err := json.Marshal(normalized)

	assert.NoError(t, err)
	assert.Contains(t, string(raw), `"error":"provider timeout"`)
	assert.Contains(t, string(raw), `"err":"nested failure"`)
}

func TestStringToTelegramMarkdown_InfoLevel(t *testing.T) {
	text := stringToTelegramMarkdown("hello_[world]-!", map[string]interface{}{
		"taskType": "process_american_tournaments_task",
	}, logrus.InfoLevel)

	assert.Contains(t, text, "*INFO*")
	assert.Contains(t, text, "ℹ️ℹ️")
	assert.Contains(t, text, "hello\\_\\[world\\]\\-\\!")
	assert.Contains(t, text, "```json")
	assert.Contains(t, text, `"taskType": "process_american_tournaments_task"`)
}

func TestStringToTelegramMarkdown_ErrorLevelSerializesError(t *testing.T) {
	text := stringToTelegramMarkdown("processing failed", map[string]interface{}{
		"error": errors.New("provider timeout"),
	}, logrus.ErrorLevel)

	assert.Contains(t, text, "*ERROR*")
	assert.Contains(t, text, "‼️‼️")
	assert.Contains(t, text, `"error": "provider timeout"`)
}

func TestStringToTelegramMarkdown_DefaultLevelAndNilContext(t *testing.T) {
	text := stringToTelegramMarkdown("debug message", nil, logrus.DebugLevel)

	assert.Contains(t, text, "*DEBUG*")
	assert.Contains(t, text, "⁉️⁉️")
	assert.Contains(t, text, "```json\nnull\n```")
}

func TestTelegramLogger_IsEnabled(t *testing.T) {
	t.Setenv("LOGGER_TELEGRAM_ENABLED", "true")
	assert.True(t, TelegramLogger{}.IsEnabled())

	t.Setenv("LOGGER_TELEGRAM_ENABLED", "false")
	assert.False(t, TelegramLogger{}.IsEnabled())
}

func TestNewTelegramLogger_PanicsWhenChatIDIsMissing(t *testing.T) {
	t.Setenv("LOGGER_TELEGRAM_CHAT_ID", "")

	assert.PanicsWithValue(t, "LOGGER_TELEGRAM_CHAT_ID is not set", func() {
		_ = NewTelegramLogger()
	})
}

func TestNewTelegramLogger_PanicsWhenChatIDIsNotNumber(t *testing.T) {
	t.Setenv("LOGGER_TELEGRAM_CHAT_ID", "not-a-number")

	assert.PanicsWithValue(t, "LOGGER_TELEGRAM_CHAT_ID is not a number", func() {
		_ = NewTelegramLogger()
	})
}

func TestNewTelegramLogger_UsesInfoLevelWhenInvalidLevelProvided(t *testing.T) {
	originalFactory := newTelegramBotAPI
	defer func() { newTelegramBotAPI = originalFactory }()

	newTelegramBotAPI = func(token string) (*tgbotapi.BotAPI, error) {
		return &tgbotapi.BotAPI{}, nil
	}

	t.Setenv("LOGGER_TELEGRAM_CHAT_ID", "12345")
	t.Setenv("LOGGER_TELEGRAM_LEVEL", "invalid-level")
	t.Setenv("LOGGER_TELEGRAM_TOKEN", "fake-token")

	logger := NewTelegramLogger()

	assert.Equal(t, int64(12345), logger.logChatId)
	assert.Equal(t, logrus.InfoLevel, logger.level)
	require.NotNil(t, logger.botApi)
}

func TestNewTelegramLogger_ParsesConfiguredLevel(t *testing.T) {
	originalFactory := newTelegramBotAPI
	defer func() { newTelegramBotAPI = originalFactory }()

	newTelegramBotAPI = func(token string) (*tgbotapi.BotAPI, error) {
		return &tgbotapi.BotAPI{}, nil
	}

	t.Setenv("LOGGER_TELEGRAM_CHAT_ID", "12345")
	t.Setenv("LOGGER_TELEGRAM_LEVEL", "warn")
	t.Setenv("LOGGER_TELEGRAM_TOKEN", "fake-token")

	logger := NewTelegramLogger()

	assert.Equal(t, logrus.WarnLevel, logger.level)
}

func TestNewTelegramLogger_PanicsWhenTelegramFactoryReturnsError(t *testing.T) {
	originalFactory := newTelegramBotAPI
	defer func() { newTelegramBotAPI = originalFactory }()

	newTelegramBotAPI = func(token string) (*tgbotapi.BotAPI, error) {
		return nil, errors.New("network unavailable")
	}

	t.Setenv("LOGGER_TELEGRAM_CHAT_ID", "12345")
	t.Setenv("LOGGER_TELEGRAM_LEVEL", "info")
	t.Setenv("LOGGER_TELEGRAM_TOKEN", "fake-token")

	assert.PanicsWithValue(t, "Error creating telegram bot api. Error: network unavailable", func() {
		_ = NewTelegramLogger()
	})
}
