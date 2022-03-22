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
}

type ItemRepository struct {
	ShortURL      string    `json:"short_url"`
	OriginalURL   string    `json:"original_url"`
	User          auth.User `json:"user"`
	CorrelationID string    `json:"correlation_id"`
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
	SaveBatch(batch []BatchRequest, user auth.User, ctx context.Context) ([]ItemRepository, error)
}

func (ir ItemRepository) ToItemView() ItemView {
	return ItemView(ir)
}

func (ir ItemRepository) ToBatchItemView() BatchItemView {
	return BatchItemView{
		CorrelationID: ir.CorrelationID,
		ShortURL:      ir.ShortURL,
	}
}
