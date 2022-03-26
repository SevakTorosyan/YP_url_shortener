package handlers

import (
	"compress/flate"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/config"
	myMiddleware "github.com/SevakTorosyan/YP_url_shortener/internal/app/middleware"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"io"
	"net/http"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type Handler struct {
	storage storage.Storage
	config  *config.Config
	*chi.Mux
}

func NewHandler(storage storage.Storage, config *config.Config) *Handler {
	handler := &Handler{
		storage,
		config,
		chi.NewMux(),
	}

	handler.Use(myMiddleware.Login(handler.config.SecretKey))
	handler.Use(myMiddleware.Decompress())
	compressor := middleware.NewCompressor(flate.BestCompression)
	handler.Use(compressor.Handler)
	handler.registerRoutes()

	return handler
}

func (h *Handler) GetShortLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "shortLink")

	item, err := h.storage.GetItem(id)
	if err != nil {
		http.Error(w, "Incorrect link", http.StatusBadRequest)
		return
	}

	responseItem := item.ToItemView()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", responseItem.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) SaveShortLink(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(myMiddleware.UserCtxValue).(auth.User)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := h.storage.SaveItem(string(b), user)
	if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	itemView := item.ToItemView()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if errors.Is(err, database.ErrItemAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	fmt.Fprintf(w, "%s/%s", h.config.BaseURL, itemView.ShortURL)
}

func (h *Handler) SaveShortLinkJSON(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(myMiddleware.UserCtxValue).(auth.User)
	requestBody := Request{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.storage.SaveItem(requestBody.URL, user)
	if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	itemView := item.ToItemView()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if errors.Is(err, database.ErrItemAlreadyExists) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode(GetResponse(itemView, h.config.BaseURL))
}

func (h *Handler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(myMiddleware.UserCtxValue).(auth.User)
	items, err := h.storage.GetItemsByUserID(h.config.BaseURL+"/", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if len(items) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		http.Error(w, "", http.StatusNoContent)
		return
	}

	viewItems := make([]storage.ItemView, 0)
	for _, item := range items {
		viewItems = append(viewItems, item.ToItemView())
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(viewItems)
}

func (h *Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	conn, err := pgx.Connect(r.Context(), h.config.DatabaseDSN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	defer conn.Close(r.Context())

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SaveBatch(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(myMiddleware.UserCtxValue).(auth.User)
	batchRequests := make([]storage.BatchRequest, 0)

	if err := json.NewDecoder(r.Body).Decode(&batchRequests); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items, err := h.storage.SaveBatch(batchRequests, user, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewItems := make([]storage.BatchItemView, 0)
	for _, item := range items {
		viewItems = append(viewItems, item.ToBatchItemView(h.config.BaseURL))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(viewItems)
}

func (h *Handler) registerRoutes() {
	h.Get("/{shortLink}", h.GetShortLink)
	h.Post("/", h.SaveShortLink)
	h.Post("/api/shorten/batch", h.SaveBatch)
	h.Post("/api/shorten", h.SaveShortLinkJSON)
	h.Get("/api/user/urls", h.GetAllItems)
	h.Get("/ping", h.PingDB)
}

func GetResponse(itemView storage.ItemView, baseURL string) Response {
	return Response{Result: fmt.Sprintf("%s/%s", baseURL, itemView.ShortURL)}
}
