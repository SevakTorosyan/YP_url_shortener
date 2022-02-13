package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SevakTorosyan/YP_url_shortener/internal/app/storage/mapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, reqBody io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	client := http.Client{}
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)

	responseBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(responseBody)
}

func TestHandlers(t *testing.T) {
	r := NewHandler(mapper.NewStorageMap())
	r.Get("/{shortLink}", r.GetShortLink)
	r.Post("/", r.SaveShortLink)

	ts := httptest.NewServer(r)
	defer ts.Close()

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
				code:        http.StatusCreated,
				response:    ts.URL + "/0",
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method: http.MethodPost,
				target: "/",
				body:   strings.NewReader("https://jsonplaceholder.typicode.com/posts/1"),
			},
		},
		{
			name: "Short link getting",
			want: want{
				code:        http.StatusBadRequest,
				response:    "{\n  \"userId\": 1,\n  \"id\": 1,\n  \"title\": \"sunt aut facere repellat provident occaecati excepturi optio reprehenderit\",\n  \"body\": \"quia et suscipit\\nsuscipit recusandae consequuntur expedita et cum\\nreprehenderit molestiae ut ut quas totam\\nnostrum rerum est autem sunt rem eveniet architecto\"\n}",
				contentType: "text/plain; charset=utf-8",
			},
			request: request{
				method: http.MethodGet,
				target: "/0",
				body:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, tt.request.method, tt.request.target, tt.request.body)

			assert.Equal(t, tt.want.code, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			resp.Body.Close()
		})
	}
}
