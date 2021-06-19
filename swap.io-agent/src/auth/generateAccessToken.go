package auth

import (
	"github.com/dgrijalva/jwt-go"
	"math"
	"os"
	"strconv"
	"swap.io-agent/src/models"
)

func GenerateAccessToken(id int) (string,error) {
	tk := &models.Token{}
	tk.Id = strconv.Itoa(id)
	tk.ExpiresAt = math.MaxInt64

	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		tk,
	)

	return token.SignedString(
		[]byte(os.Getenv("TOKEN_SECRET")),
	)
}