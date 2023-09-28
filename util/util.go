package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtkey = []byte("nishad")

type Claims struct {
	Username string
	jwt.StandardClaims
}

// Token generation
func GenerateJWT(username string) (string, error) {

	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtkey)
}

// Token verification
func VerifyJWT(tokenstring string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenstring, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
