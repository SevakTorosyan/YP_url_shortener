package mapper

import "fmt"

type StorageMap struct {
	links map[string]string
}

func NewStorageMap() *StorageMap {
	return &StorageMap{links: make(map[string]string)}
}

func (s StorageMap) GetItem(shortLink string) (string, error) {
	link, ok := s.links[shortLink]

	if !ok {
		return "", fmt.Errorf("link not found")
	}

	return link, nil
}

func (s *StorageMap) SaveItem(link, shortLink string) (string, error) {
	s.links[shortLink] = link

	return shortLink, nil
}
