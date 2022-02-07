package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
)

var links []string

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodPost {
		saveShortLink(w, r)

		return
	}

	if r.Method == http.MethodGet {
		getShortLink(w, r)

		return
	}

	http.NotFound(w, r)
}

func getShortLink(w http.ResponseWriter, r *http.Request) {
	id, err := utils.GetIdentifier(r)

	if err != nil {
		http.Error(w, "Произошла ошибка", http.StatusBadRequest)

		return
	}

	if id >= len(links) {
		http.Error(w, "Некорректный идентификатор", http.StatusBadRequest)

		return
	}

	w.Header().Set("Location", links[id])
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func saveShortLink(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	links = append(links, string(b))

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/%d", r.Host, len(links)-1)
}
