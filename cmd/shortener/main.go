package main

import (
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/handlers"
	"net/http"
)

func main() {
	mux := &handlers.MyMux{}
	mux.InitStorage()

	http.ListenAndServe("localhost:8080", mux)
}
