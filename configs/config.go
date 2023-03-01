package configs

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	ServerAddress               string
	DatabaseDSN                 string
	PasswordSalt                string
	TokenTTL                    time.Duration
	TokenSigningKey             []byte
	FirmwarePath                string
	JournalTTL                  time.Duration
	JournalsWorkerInterval      time.Duration
	NotificationsWorkerInterval time.Duration
	TelegramToken               string
	TelegramChatID              string
	DeviceOfflineDuration       time.Duration
}

// todo: env
func NewConfig() *Config {
	cfg := &Config{
		ServerAddress:               ":8080",
		PasswordSalt:                "juaT9OLosPlhUhDj",
		TokenTTL:                    time.Hour * 12,
		TokenSigningKey:             []byte("qYqx2APnPhDHBl2AW3OjUYeWWFAtzF7d"),
		FirmwarePath:                "assets/firmware.bin",
		JournalTTL:                  time.Minute,
		JournalsWorkerInterval:      time.Second * 15,
		NotificationsWorkerInterval: time.Second * 5,
		DeviceOfflineDuration:       time.Minute * 15,
	}

	loadEnv(cfg)
	loadArgs(cfg)

	return cfg
}

func loadEnv(cfg *Config) {
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		cfg.ServerAddress = s
	}

	if s, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = s
	}

	if s, ok := os.LookupEnv("TELEGRAM_TOKEN"); ok {
		cfg.TelegramToken = s
	}

	if s, ok := os.LookupEnv("TELEGRAM_CHAT_ID"); ok {
		cfg.TelegramChatID = s
	}
}

func loadArgs(cfg *Config) {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database data source name")
	flag.StringVar(&cfg.TelegramToken, "t", cfg.TelegramToken, "telegram bot token")
	flag.StringVar(&cfg.TelegramChatID, "c", cfg.TelegramChatID, "telegram chat id")

	flag.Parse()
}
