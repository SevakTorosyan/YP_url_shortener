package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v6"
)

var lock = &sync.Mutex{}

type config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

var configInstance *config

func InitConfig() *config {
	configInstance = &config{}

	err := env.Parse(configInstance)
	if err != nil {
		log.Fatal(err)
	}

	return configInstance
}

func GetInstance() *config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		return InitConfig()
	}

	return configInstance
}
