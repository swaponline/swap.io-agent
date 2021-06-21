package auth

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"os"
)

func DecodeAccessToken(tokenString string) (int,bool) {
	info, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithVerify(
			jwa.HS256,
			[]byte(os.Getenv("TOKEN_SECRET")),
		),
	)
	if err != nil {
		return -1, true
	}

	if id, ok := info.Get("id"); ok {
		if idF, ok := id.(float64); ok {
			return int(idF), false
		}
	}

	return -1, true
}