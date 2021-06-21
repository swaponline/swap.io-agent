package socketServer

import (
	"errors"
	"net/http"
	"swap.io-agent/src/auth"
)

func AuthenticationSocketConnect (request *http.Request) (http.Header, error) {
	tokenInfo := request.URL.Query()["token"]
	if len(tokenInfo) == 0 {
		return nil, errors.New("not exist token")
	}
	_, err := auth.DecodeAccessToken(tokenInfo[0])
	if err {
		return nil, errors.New("not valid token")
	}
	return nil, nil
}