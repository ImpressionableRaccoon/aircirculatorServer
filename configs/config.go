package configs

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	ServerAddress   string
	DatabaseDSN     string
	PasswordSalt    string
	TokenTTL        time.Duration
	TokenSigningKey []byte
}

func NewConfig() *Config {
	cfg := &Config{
		ServerAddress:   ":8080",
		PasswordSalt:    "juaT9OLosPlhUhDj",
		TokenTTL:        time.Hour * 12,
		TokenSigningKey: []byte("qYqx2APnPhDHBl2AW3OjUYeWWFAtzF7d"),
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
}

func loadArgs(cfg *Config) {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database data source name")

	flag.Parse()
}
