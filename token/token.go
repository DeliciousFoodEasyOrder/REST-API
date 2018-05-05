package token

import (
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Token is a struct for token authorization
type Token struct {
	AccessToken string     `json:"access_token"`
	TokenType   string     `json:"token_type"`
	ExpiresIn   int        `json:"expires_in"`
	IssuedAt    *time.Time `json:"issued_at"`
}

// NewJWTToken issues a new token with type "JWT"
func NewJWTToken(aud int, expiresIn int) *Token {
	_, exist := tokens[aud]
	if exist {
		if !tokens[aud].Expired() {
			return tokens[aud]
		}
	}

	now := time.Now()
	duration, _ := time.ParseDuration(strconv.Itoa(expiresIn) + "s")
	accessToken := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
		"iss": "DFEO",
		"sub": "access_token",
		"iat": now.Unix(),
		"exp": now.Add(duration).Unix(),
		"aud": aud,
	})
	accessSigned, _ := accessToken.SignedString(secret)

	token := &Token{
		AccessToken: accessSigned,
		TokenType:   "JWT",
		ExpiresIn:   expiresIn,
		IssuedAt:    &now,
	}
	tokens[aud] = token
	return token
}

// Expired determines if token is expired
func (t *Token) Expired() bool {
	duration, _ := time.ParseDuration(strconv.Itoa(t.ExpiresIn) + "s")
	return t.IssuedAt.Add(duration).Before(time.Now())
}
