package storage

import "github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"

type ItemView struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	User        auth.User `json:"-"`
}

type ItemRepository struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	User        auth.User `json:"user"`
}

type Storage interface {
	GetItem(shortURL string) (ItemRepository, error)
	SaveItem(originalURL string, user auth.User) (ItemRepository, error)
	GetItemsByUserID(serverAddress string, user auth.User) ([]ItemRepository, error)
}

func (ir ItemRepository) ToItemView() ItemView {
	return ItemView(ir)
}
