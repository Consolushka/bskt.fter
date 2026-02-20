package logger

import (
	"os"
	"strconv"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func BuildSettingsFromEnv() []composite_logger.LoggerSetting {
	settings := make([]composite_logger.LoggerSetting, 0)

	if os.Getenv("LOGGER_CONSOLE_ENABLED") == "true" {
		level, _ := composite_logger.ParseLevel(os.Getenv("LOGGER_CONSOLE_LEVEL"))
		settings = append(settings, setting.ConsoleSetting{
			LowerLevel: level,
		})
	}

	if os.Getenv("LOGGER_FILE_ENABLED") == "true" {
		level, _ := composite_logger.ParseLevel(os.Getenv("LOGGER_FILE_LEVEL"))
		settings = append(settings, setting.FileSetting{
			Path:       os.Getenv("LOGGER_FILE_PATH"),
			LowerLevel: level,
		})
	}

	if os.Getenv("LOGGER_TELEGRAM_ENABLED") == "true" {
		chatId, _ := strconv.ParseInt(os.Getenv("LOGGER_TELEGRAM_CHAT_ID"), 10, 64)
		level, _ := composite_logger.ParseLevel(os.Getenv("LOGGER_TELEGRAM_LEVEL"))
		settings = append(settings, setting.TelegramSetting{
			Enabled:              true,
			BotKey:               os.Getenv("LOGGER_TELEGRAM_TOKEN"),
			ChatId:               chatId,
			LowerLevel:           level,
			UseLevelTitleWrapper: true,
		})
	}

	return settings
}
