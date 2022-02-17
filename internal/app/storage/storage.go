package storage

type Storage interface {
	GetItem(shortLink string) (string, error)
	SaveItem(link string) (string, error)
}
