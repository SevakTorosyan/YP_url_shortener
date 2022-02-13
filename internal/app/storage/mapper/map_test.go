package mapper

import (
	"fmt"
	"github.com/SevakTorosyan/YP_url_shortener/internal/app/utils"
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
			value: "https://ya.ru",
			want: want{
				shortLink: utils.GenerateMockString(),
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMap := NewStorageMap()
			shortLink := utils.GenerateMockString()
			shortLink, err := storageMap.SaveItem(tt.value, shortLink)

			assert.Equal(t, tt.want.shortLink, shortLink)
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
			value: utils.GenerateMockString(),
			want: want{
				link: "https://ya.ru",
				err:  nil,
			},
		},
		{
			name:  "Incorrect test",
			value: "adfasdasd",
			want: want{
				link: "",
				err:  fmt.Errorf("link not found"),
			},
		},
	}

	for _, tt := range tests {
		storageMap := NewStorageMap()
		shortLink := utils.GenerateMockString()
		storageMap.SaveItem("https://ya.ru", shortLink)

		t.Run(tt.name, func(t *testing.T) {
			link, err := storageMap.GetItem(tt.value)
			assert.Equal(t, tt.want.link, link)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
