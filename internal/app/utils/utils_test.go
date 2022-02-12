package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUtils_GetIdentifier(t *testing.T) {
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
			name:  "Get identifier",
			value: "/5",
			want: want{
				id:  5,
				err: nil,
			},
		},
		{
			name:  "Incorrect path",
			value: "/something",
			want: want{
				id:  0,
				err: fmt.Errorf("не удалось получить идентификатор"),
			},
		},
		{
			name:  "Too long path",
			value: "/5/hello",
			want: want{
				id:  0,
				err: fmt.Errorf("некорректная ссылка"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identifier, err := GetIdentifier(tt.value)
			assert.Equal(t, tt.want.id, identifier)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
