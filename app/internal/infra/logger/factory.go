package logger

import (
	"IMP/app/internal/infra/config"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/Consolushka/golang.composite_logger/pkg/adapters/setting"
)

func BuildSettings(cfg config.LoggerConfig) []composite_logger.LoggerSetting {
	settings := make([]composite_logger.LoggerSetting, 0)

	if cfg.ConsoleEnabled {
		level, _ := composite_logger.ParseLevel(cfg.ConsoleLevel)
		settings = append(settings, setting.ConsoleSetting{
			LowerLevel: level,
		})
	}

	if cfg.FileEnabled {
		level, _ := composite_logger.ParseLevel(cfg.FileLevel)
		settings = append(settings, setting.FileSetting{
			Path:       cfg.FilePath,
			LowerLevel: level,
		})
	}

	if cfg.TelegramEnabled {
		level, _ := composite_logger.ParseLevel(cfg.TelegramLevel)
		settings = append(settings, setting.TelegramSetting{
			Enabled:              true,
			BotKey:               cfg.TelegramToken,
			ChatId:               cfg.TelegramChatId,
			LowerLevel:           level,
			UseLevelTitleWrapper: true,
		})
	}

	return settings
}
