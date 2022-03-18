package memory

import (
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/auth"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageSlice_SaveItem(t *testing.T) {
	type want struct {
		shortLink string
		err       error
	}

	tests := []struct {
		name  string
		value string
		want  want
	}{
		{
			name:  "Correct test",
			value: "https://jsonplaceholder.typicode.com/posts/1",
			want: want{
				shortLink: "asdjnd3242",
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMap := mock.StorageMock{}
			item, err := storageMap.SaveItem(tt.value, auth.User{})

			assert.Equal(t, tt.want.shortLink, item.ShortURL)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestStorageSlice_GetItem(t *testing.T) {
	type want struct {
		link string
		err  error
	}

	tests := []struct {
		name  string
		value string
		want  want
	}{
		{
			name:  "Correct test",
			value: "asdjnd3242",
			want: want{
				link: "https://jsonplaceholder.typicode.com/posts/1",
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		storageMap := mock.NewMockStorage()

		t.Run(tt.name, func(t *testing.T) {
			item, err := storageMap.GetItem(tt.value)
			assert.Equal(t, tt.want.link, item.OriginalURL)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
