package token

import (
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	signingMethod jwt.SigningMethod = jwt.SigningMethodHS256
	tokens        map[int]*Token
)

func init() {
	tokens = make(map[int]*Token)
}
