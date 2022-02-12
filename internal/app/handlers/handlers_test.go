package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlers(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	type request struct {
		method string
		target string
		body   io.Reader
	}
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name    string
		want    want
		request request
	}{
		{
			name: "Short link creating",
			want: want{
				code:        201,
				response:    "http://example.com/0",
				contentType: "",
			},
			request: request{
				method: http.MethodPost,
				target: "/",
				body:   strings.NewReader("https://ya.ru"),
			},
		},
		{
			name: "Short link getting",
			want: want{
				code:        307,
				response:    "",
				contentType: "",
			},
			request: request{
				method: http.MethodGet,
				target: "/0",
				body:   strings.NewReader("https://ya.ru"),
			},
		},
	}

	mux := &MyMux{}
	mux.InitStorage()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.target, tt.request.body)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(mux.ServeHTTP)

			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			responseBody, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(responseBody))
		})
	}
}
