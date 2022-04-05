package storage

import (
	"context"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
)

type ItemView struct {
	ShortURL      string    `json:"short_url"`
	OriginalURL   string    `json:"original_url"`
	User          auth.User `json:"-"`
	CorrelationID string    `json:"-"`
	IsDeleted     bool      `json:"-"`
}

type ItemRepository struct {
	ShortURL      string    `json:"short_url"`
	OriginalURL   string    `json:"original_url"`
	User          auth.User `json:"user"`
	CorrelationID string    `json:"correlation_id"`
	IsDeleted     bool      `json:"is_deleted"`
}

type BatchItemView struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type Storage interface {
	GetItem(shortURL string) (ItemRepository, error)
	SaveItem(originalURL string, user auth.User) (ItemRepository, error)
	GetItemsByUserID(serverAddress string, user auth.User) ([]ItemRepository, error)
	SaveBatch(ctx context.Context, batch []BatchRequest, user auth.User) ([]ItemRepository, error)
	DeleteByIds(batch []string, user auth.User) error
	Ping() error
	Close() error
}

func (ir ItemRepository) ToItemView() ItemView {
	return ItemView(ir)
}

func (ir ItemRepository) ToBatchItemView(baseURL string) BatchItemView {
	return BatchItemView{
		CorrelationID: ir.CorrelationID,
		ShortURL:      baseURL + "/" + ir.ShortURL,
	}
}
