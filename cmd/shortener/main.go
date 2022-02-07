package main

import (
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/handlers"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/mapper"
	"log"
	"net/http"
)

func main() {
	handler := handlers.NewHandler(mapper.NewStorageMap())

	handler.Get("/{shortLink}", handler.GetShortLink)
	handler.Post("/", handler.SaveShortLink)

	log.Fatal(http.ListenAndServe("localhost:8080", handler))
}
