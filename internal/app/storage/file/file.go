package file

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
	"os"

	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
)

type StorageFile struct {
	file  *os.File
	items map[string]storage.ItemRepository
}

func NewStorageFile(filename string) (*StorageFile, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}

	storageFile := &StorageFile{file: file}
	storageFile.items = make(map[string]storage.ItemRepository)
	storageFile.loadItems()

	return storageFile, nil
}

func (s *StorageFile) GetItem(shortURL string) (storage.ItemRepository, error) {
	link, ok := s.items[shortURL]

	if !ok {
		return storage.ItemRepository{}, fmt.Errorf("link not found")
	}

	return link, nil
}

func (s *StorageFile) SaveItem(originalURL string, user auth.User) (storage.ItemRepository, error) {
	encoder := json.NewEncoder(s.file)
	shortURL := utils.GenerateRandomString(15)
	item := storage.ItemRepository{ShortURL: shortURL, OriginalURL: originalURL, User: user}
	s.items[shortURL] = item

	if err := encoder.Encode(item); err != nil {
		return storage.ItemRepository{}, err
	}

	return item, nil
}

func (s *StorageFile) GetItemsByUserID(serverAddress string, user auth.User) ([]storage.ItemRepository, error) {
	items := make([]storage.ItemRepository, 0)

	for shortLink, item := range s.items {
		if item.User.ID == user.ID {
			items = append(items, storage.ItemRepository{ShortURL: serverAddress + shortLink, OriginalURL: item.OriginalURL})
		}
	}

	return items, nil
}

func (s *StorageFile) SaveBatch(batch []storage.BatchRequest, user auth.User, ctx context.Context) ([]storage.ItemRepository, error) {
	return []storage.ItemRepository{}, fmt.Errorf("method is not supported")
}

func (s *StorageFile) loadItems() {
	item := &storage.ItemRepository{}
	decoder := json.NewDecoder(s.file)

	for {
		err := decoder.Decode(item)

		if err != nil {
			break
		}

		s.items[item.ShortURL] = *item
	}
}
