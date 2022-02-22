package main

import (
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/config"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/handlers"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/mapper"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetInstance()
	handler := handlers.NewHandler(mapper.NewStorageMap())

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, handler))
}
