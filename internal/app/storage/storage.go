package storage

type Storage interface {
	GetItem(linkID int) (string, error)
	SaveItem(link string) (int, error)
}
