package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database  DatabaseConfig
	Scheduler SchedulerConfig
	Providers ProvidersConfig
	Logger    LoggerConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
}

type SchedulerConfig struct {
	PollInterval    int `env:"SCHEDULER_POLL_INTERVAL" env-default:"30"`
	StaggerInterval int `env:"SCHEDULER_STAGGER_INTERVAL_MINUTES" env-default:"5"`
	RefreshInterval int `env:"SCHEDULER_REFRESH_INTERVAL_MINUTES" env-default:"5"`
}

type ProvidersConfig struct {
	ApiSportApiKey      string `env:"API_SPORT_API_KEY" env-required:"true"`
	ApiNbaRateLimit        int    `env:"API_NBA_RATE_LIMIT_PER_MINUTE" env-default:"10"`
	ApiBasketballRateLimit int    `env:"API_BASKETBALL_RATE_LIMIT_PER_MINUTE" env-default:"10"`
	InfobasketRateLimit    int    `env:"INFOBASKET_RATE_LIMIT_PER_MINUTE" env-default:"25"`
	SportotekaRateLimit int    `env:"SPORTOTEKA_RATE_LIMIT_PER_MINUTE" env-default:"25"`
}

type LoggerConfig struct {
	ConsoleEnabled  bool   `env:"LOGGER_CONSOLE_ENABLED" env-default:"true"`
	ConsoleLevel    string `env:"LOGGER_CONSOLE_LEVEL" env-default:"info"`
	FileEnabled     bool   `env:"LOGGER_FILE_ENABLED" env-default:"false"`
	FileLevel       string `env:"LOGGER_FILE_LEVEL" env-default:"info"`
	FilePath        string `env:"LOGGER_FILE_PATH" env-default:"tmp/logs/app.log"`
	TelegramEnabled bool   `env:"LOGGER_TELEGRAM_ENABLED" env-default:"false"`
	TelegramToken   string `env:"LOGGER_TELEGRAM_TOKEN"`
	TelegramChatId  int64  `env:"LOGGER_TELEGRAM_CHAT_ID"`
	TelegramLevel   string `env:"LOGGER_TELEGRAM_LEVEL" env-default:"warn"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}
	return &cfg, nil
}
