package mock

const Link = "https://jsonplaceholder.typicode.com/posts/1"
const ShortLink = "asdjnd3242"

type StorageMock struct{}

func NewMockStorage() *StorageMock {
	return &StorageMock{}
}

func (s StorageMock) GetItem(shortLink string) (string, error) {
	return Link, nil
}

func (s *StorageMock) SaveItem(link string) (string, error) {
	return ShortLink, nil
}
