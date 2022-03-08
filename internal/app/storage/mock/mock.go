package mock

import (
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage"
)

const Link = "https://jsonplaceholder.typicode.com/posts/1"
const ShortLink = "asdjnd3242"
const UserID = "sdfmsdkfmf34"

type StorageMock struct{}

func NewMockStorage() *StorageMock {
	return &StorageMock{}
}

func (s StorageMock) GetItem(shortLink string) (storage.ItemRepository, error) {
	user := auth.User{ID: UserID}

	return storage.ItemRepository{ShortURL: ShortLink, OriginalURL: Link, User: user}, nil
}

func (s *StorageMock) SaveItem(link string, user auth.User) (storage.ItemRepository, error) {
	return storage.ItemRepository{ShortURL: ShortLink, OriginalURL: Link, User: user}, nil
}

func (s *StorageMock) GetItemsByUserID(serverAddress string, user auth.User) ([]storage.ItemRepository, error) {
	return []storage.ItemRepository{
		{
			ShortURL:    serverAddress + "/sdfsdfsrfw",
			OriginalURL: "https://vk.com",
		},
	}, nil
}
