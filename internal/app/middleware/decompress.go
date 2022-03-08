package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

type gzipReader struct {
	*http.Request
	Reader io.Reader
}

func (r gzipReader) Read(b []byte) (int, error) {
	return r.Reader.Read(b)
}

func Decompress() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Content-Encoding") != "gzip" {
				next.ServeHTTP(w, r)
				return
			}

			compressed, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			gr, err := gzip.NewReader(bytes.NewBuffer(compressed))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer gr.Close()

			r.Body = gr
			next.ServeHTTP(w, r)
		})
	}
}
