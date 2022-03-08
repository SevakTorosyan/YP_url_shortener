package main

import (
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/config"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/handlers"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/file"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/mapper"
	"log"
	"net/http"
)

func main() {
	cfg := config.InitConfig()
	var handler *handlers.Handler

	if cfg.FileStoragePath == "" {
		handler = handlers.NewHandler(mapper.NewStorageMap(), cfg)
	} else {
		storage, err := file.NewStorageFile(cfg.FileStoragePath)

		if err != nil {
			log.Fatal("An error occurred during storage initialization")
		}
		handler = handlers.NewHandler(storage, cfg)
	}

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, handler))
}
