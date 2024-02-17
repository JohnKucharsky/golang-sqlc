package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info")
	}

	values := strings.Split(val, " ")

	if len(values) != 2 {
		return "", errors.New("can't get header")
	}

	if values[0] != "ApiKey" {
		return "", errors.New("no ApiKey part in auth string")
	}

	return values[1], nil
}
