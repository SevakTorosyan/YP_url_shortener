package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	ShortLinkLength int    `env:"SHORT_LINK_LENGTH" envDefault:"15"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"52fdfc072182654f163f5f0f9a621d72"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	DatabaseTimeout int64  `env:"DATABASE_TIMEOUT" envDefault:"5"`
}

func InitConfig() *Config {
	configInstance := &Config{}

	err := env.Parse(configInstance)

	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&configInstance.ServerAddress, "a", configInstance.ServerAddress, "Server address")
	flag.StringVar(&configInstance.BaseURL, "b", configInstance.BaseURL, "Base URL")
	flag.StringVar(&configInstance.FileStoragePath, "f", configInstance.FileStoragePath, "File storage path")
	flag.StringVar(&configInstance.DatabaseDSN, "d", configInstance.DatabaseDSN, "Database connection properties")
	flag.Int64Var(&configInstance.DatabaseTimeout, "dt", configInstance.DatabaseTimeout, "Database connection timeout")

	flag.Parse()

	return configInstance
}
