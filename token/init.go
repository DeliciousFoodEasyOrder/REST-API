package token

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	signingMethod jwt.SigningMethod = jwt.SigningMethodHS256
	tokens        map[int]*Token
	secret        = []byte("secret")
)

func init() {
	tokens = make(map[int]*Token)
}

// ParseClaims parse a token to claims
func ParseClaims(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v",
					token.Header["alg"])
			}

			return secret, nil
		})

	if err != nil {
		panic(err)
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	return claims
}
