package handlers

import (
	"fmt"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

const shortLinkLength = 15
const hostName = "localhost:8080"

type Handler struct {
	storage storage.Storage
	*chi.Mux
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage,
		chi.NewMux(),
	}
}

func (h *Handler) GetShortLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "shortLink")

	originalLink, err := h.storage.GetItem(id)
	if err != nil {
		http.Error(w, "Incorrect link", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) SaveShortLink(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortLink := utils.GenerateRandomString(shortLinkLength)
	fmt.Println(shortLink)
	_, err = h.storage.SaveItem(string(b), shortLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/%s", hostName, shortLink)
}
