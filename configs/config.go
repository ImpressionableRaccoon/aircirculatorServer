package configs

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress string
	DatabaseDSN   string
}

func NewConfig() *Config {
	cfg := &Config{
		ServerAddress: ":8080",
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
