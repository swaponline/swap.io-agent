package auth

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"swap.io-agent/src/models"
)

func VerifyAccessToken(tokenString string) (int,bool) {
	tk := &models.Token{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		tk,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		},
	)
	if err != nil || !token.Valid {
		return 0, true
	}

	id, err := strconv.Atoi(tk.Id)
	if err != nil {
		return 0, true
	}

	return id, false
}