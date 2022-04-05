package memory

import (
	"context"
	"fmt"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
)

type StorageMap struct {
	links map[string]storage.ItemRepository
}

func NewStorageMap() *StorageMap {
	return &StorageMap{links: make(map[string]storage.ItemRepository)}
}

func (s StorageMap) GetItem(shortLink string) (storage.ItemRepository, error) {
	item, ok := s.links[shortLink]

	if !ok {
		return storage.ItemRepository{}, fmt.Errorf("link not found")
	}

	return item, nil
}

func (s *StorageMap) SaveItem(originalURL string, user auth.User) (storage.ItemRepository, error) {
	shortURL := utils.GenerateRandomString(15)
	item := storage.ItemRepository{ShortURL: shortURL, OriginalURL: originalURL, User: user}
	s.links[shortURL] = item

	return item, nil
}

func (s *StorageMap) GetItemsByUserID(serverAddress string, user auth.User) ([]storage.ItemRepository, error) {
	items := make([]storage.ItemRepository, 0)

	for shortLink, item := range s.links {
		if item.User.ID == user.ID {
			items = append(items, storage.ItemRepository{ShortURL: serverAddress + shortLink, OriginalURL: item.OriginalURL})
		}
	}

	return items, nil
}

func (s *StorageMap) SaveBatch(ctx context.Context, batch []storage.BatchRequest, user auth.User) ([]storage.ItemRepository, error) {
	return []storage.ItemRepository{}, fmt.Errorf("method is not supported")
}

func (s StorageMap) Ping() error {
	return nil
}

func (s StorageMap) Close() error {
	return nil
}

func (s StorageMap) DeleteByIds(batch []string, user auth.User) error {
	return fmt.Errorf("method is not supported")
}
