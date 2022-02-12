package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GetIdentifier(path string) (int, error) {
	p := strings.Split(path, "/")

	if len(p) != 2 {
		return 0, fmt.Errorf("некорректная ссылка")
	}

	id, err := strconv.Atoi(p[1])

	if err != nil {
		return 0, fmt.Errorf("не удалось получить идентификатор")
	}

	return id, nil
}
