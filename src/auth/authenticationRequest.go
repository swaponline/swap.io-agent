package auth

import (
	"errors"
	"net/http"
)

func AuthenticationRequest (request *http.Request) (string, error) {
	tokenInfo := request.URL.Query()["token"]
	if len(tokenInfo) == 0 {
		return "", errors.New("not exist token")
	}
	userId, err := DecodeAccessToken(tokenInfo[0])
	if err {
		return "", errors.New("not valid token")
	}
	return userId, nil
}