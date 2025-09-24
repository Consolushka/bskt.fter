package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type TelegramLogger struct {
	botApi    *tgbotapi.BotAPI
	logChatId int64
	level     logrus.Level
}

func (t TelegramLogger) IsEnabled() bool {
	isEnabled := os.Getenv("LOGGER_TELEGRAM_ENABLED")

	return isEnabled == "true"
}

func NewTelegramLogger() TelegramLogger {
	logChatId := os.Getenv("LOGGER_TELEGRAM_CHAT_ID")
	if logChatId == "" {
		panic("LOGGER_TELEGRAM_CHAT_ID is not set")
	}

	intLogChatId, err := strconv.ParseInt(logChatId, 10, 64)
	if err != nil {
		panic("LOGGER_TELEGRAM_CHAT_ID is not a number")
	}

	loggerLevel, err := logrus.ParseLevel(os.Getenv("LOGGER_TELEGRAM_LEVEL"))
	if err != nil {
		loggerLevel = logrus.InfoLevel
	}

	botApi, err := tgbotapi.NewBotAPI(os.Getenv("LOGGER_TELEGRAM_TOKEN"))
	if err != nil {
		panic("Error creating telegram bot api. Error: " + err.Error())
	}

	return TelegramLogger{
		botApi:    botApi,
		logChatId: intLogChatId,
		level:     loggerLevel,
	}
}

func (t TelegramLogger) Info(message string, context map[string]interface{}) {
	level := logrus.InfoLevel

	t.send(message, context, level)
}

func (t TelegramLogger) Warn(message string, context map[string]interface{}) {
	level := logrus.WarnLevel

	t.send(message, context, level)
}

func (t TelegramLogger) Error(message string, context map[string]interface{}) {
	level := logrus.ErrorLevel

	t.send(message, context, level)
}

func (t TelegramLogger) send(message string, context map[string]interface{}, level logrus.Level) {
	if t.level > level {
		return
	}

	text := stringToTelegramMarkdown(message, context, level)

	tgMessage := tgbotapi.NewMessage(t.logChatId, text)

	_, err := t.botApi.Send(tgMessage)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func stringToTelegramMarkdown(message string, context map[string]interface{}, level logrus.Level) string {
	escapeMarkdownV2 := func(text string) string {
		var markdownV2Regex = regexp.MustCompile(`([\[\]\-_*~` + "`" + `>#+=|{}.!])`)
		return markdownV2Regex.ReplaceAllString(text, "\\$1")
	}

	var emoji string
	now := time.Now().Format("[2006-01-02 15:04:05]")

	switch level {
	case logrus.ErrorLevel:
		emoji = "‼️‼️"
	case logrus.WarnLevel:
		emoji = "⚠️⚠️"
	case logrus.InfoLevel:
		emoji = "ℹ️ℹ️"
	default:
		emoji = "⁉️⁉️"
	}

	jsonContext, err := json.MarshalIndent(context, "", "    ")
	if err != nil {
		panic(err)
	}

	text := fmt.Sprintf("%s *%s* %s\n%s %s\n\n```json\n%s\n```",
		emoji,
		level,
		emoji,
		escapeMarkdownV2(now),
		escapeMarkdownV2(message),
		string(jsonContext))

	return text
}
