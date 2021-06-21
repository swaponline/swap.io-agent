package auth

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"os"
)

func GenerateAccessToken(id int) (string,error) {
	token := jwt.New()
	token.Set("id", id)

	key, err := jwt.Sign(
		token,
		jwa.HS256,
		[]byte(os.Getenv("TOKEN_SECRET")),
	)

	return string(key), err
}