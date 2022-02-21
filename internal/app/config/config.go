package config

import (
	"flag"
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
)

var lock = &sync.Mutex{}

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	ShortLinkLength int    `env:"SHORT_LINK_LENGTH" envDefault:"15"`
}

var configInstance *config

func initConfig() *config {
	configInstance = &config{}

	err := env.Parse(configInstance)

	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&configInstance.ServerAddress, "a", configInstance.ServerAddress, "Server address")
	flag.StringVar(&configInstance.BaseURL, "b", configInstance.BaseURL, "Base URL")
	flag.StringVar(&configInstance.FileStoragePath, "f", configInstance.FileStoragePath, "File storage path")
	flag.Parse()

	return configInstance
}

func GetInstance() *config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		return initConfig()
	}

	return configInstance
}
