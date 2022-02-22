package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/config"
	"io"
	"net/http"

	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type Handler struct {
	storage storage.Storage
	*chi.Mux
}

func NewHandler(storage storage.Storage) *Handler {
	handler := &Handler{
		storage,
		chi.NewMux(),
	}
	handler.registerRoutes()

	return handler
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

	shortLink, err := h.storage.SaveItem(string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", config.GetInstance().BaseURL, shortLink)
}

func (h *Handler) SaveShortLinkJSON(w http.ResponseWriter, r *http.Request) {
	requestBody := Request{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortLink, err := h.storage.SaveItem(requestBody.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(GetResponse(shortLink))
}

func (h *Handler) registerRoutes() {
	h.Get("/{shortLink}", h.GetShortLink)
	h.Post("/", h.SaveShortLink)
	h.Post("/api/shorten", h.SaveShortLinkJSON)
}

func GetResponse(shortLink string) Response {
	return Response{Result: fmt.Sprintf("%s/%s", config.GetInstance().BaseURL, shortLink)}
}
