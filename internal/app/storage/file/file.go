package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SevakTorosyan/YP_url_shortener/internal/app/config"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
)

type StorageFile struct {
	file  *os.File
	items map[string]string
}

type Item struct {
	ShortLink string `json:"short_link"`
	Link      string `json:"link"`
}

func NewStorageFile(filename string) (*StorageFile, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}

	storageFile := &StorageFile{file: file}
	storageFile.items = make(map[string]string)
	storageFile.loadItems()

	return storageFile, nil
}

func (s *StorageFile) GetItem(shortLink string) (string, error) {
	link, ok := s.items[shortLink]

	if !ok {
		return "", fmt.Errorf("link not found")
	}

	return link, nil
}

func (s *StorageFile) SaveItem(link string) (string, error) {
	encoder := json.NewEncoder(s.file)
	shortLink := utils.GenerateRandomString(config.GetInstance().ShortLinkLength)
	item := Item{shortLink, link}
	s.items[shortLink] = link

	if err := encoder.Encode(item); err != nil {
		return "", err
	}

	return shortLink, nil
}

func (s *StorageFile) loadItems() {
	item := &Item{}
	decoder := json.NewDecoder(s.file)

	for {
		err := decoder.Decode(item)

		if err != nil {
			break
		}

		s.items[item.ShortLink] = item.Link
	}
}
