package slice

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageSlice_SaveItem(t *testing.T) {
	type want struct {
		id  int
		err error
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
				id:  0,
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageSlice := &StorageSlice{}

			id, err := storageSlice.SaveItem(tt.value)
			assert.Equal(t, tt.want.id, id)
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
		value int
		want  want
	}{
		{
			name:  "Correct test",
			value: 0,
			want: want{
				link: "https://ya.ru",
				err:  nil,
			},
		},
		{
			name:  "Incorrect test",
			value: 1,
			want: want{
				link: "",
				err:  fmt.Errorf("некорректный идентификатор"),
			},
		},
	}

	for _, tt := range tests {
		storageSlice := &StorageSlice{}
		storageSlice.SaveItem("https://ya.ru")

		t.Run(tt.name, func(t *testing.T) {
			link, err := storageSlice.GetItem(tt.value)
			assert.Equal(t, tt.want.link, link)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
