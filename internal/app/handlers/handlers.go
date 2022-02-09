package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/slice"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
)

type MyMux struct {
	s storage.Storage
}

// InitStorage временное решение пока нет DI
func (p *MyMux) InitStorage() {
	p.s = &slice.StorageSlice{}
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodPost {
		saveShortLink(w, r, p.s)

		return
	}

	if r.Method == http.MethodGet {
		getShortLink(w, r, p.s)

		return
	}

	http.NotFound(w, r)
}

func getShortLink(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	id, err := utils.GetIdentifier(r)
	if err != nil {
		http.Error(w, "Произошла ошибка", http.StatusBadRequest)

		return
	}

	originalLink, err := s.GetItem(id)
	if err != nil {
		http.Error(w, "Некорректный идентификатор", http.StatusBadRequest)

		return
	}

	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func saveShortLink(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	linkID, err := s.SaveItem(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/%d", r.Host, linkID)
}
