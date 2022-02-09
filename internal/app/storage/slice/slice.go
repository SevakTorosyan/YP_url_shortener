package slice

import "fmt"

type StorageSlice struct {
	links []string
}

func (s StorageSlice) GetItem(linkID int) (string, error) {
	if linkID >= len(s.links) {
		return "", fmt.Errorf("некорректный идентификатор")
	}

	return s.links[linkID], nil
}

func (s *StorageSlice) SaveItem(link string) (int, error) {
	s.links = append(s.links, link)

	return len(s.links) - 1, nil
}
