package auth

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"swap.io-agent/src/env"
)

func GenerateAccessToken(id string) (string,error) {
	token := jwt.New()
	token.Set("id", id)

	key, err := jwt.Sign(
		token,
		jwa.HS256,
		[]byte(env.SECRET_TOKEN),
	)

	return string(key), err
}