package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetIdentifier(r *http.Request) (int, error) {
	p := strings.Split(r.URL.Path, "/")

	if len(p) != 2 {
		return 0, fmt.Errorf("некорректная ссылка")
	}

	id, err := strconv.Atoi(p[1])

	if err != nil {
		return 0, err
	}

	return id, nil
}
