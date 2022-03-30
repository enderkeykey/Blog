package main

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte("AllYourBase")

func GenerateToken(username string, uid int) (string, error) {
	expireAt := time.Now().Add(time.Hour).Unix()

	claims := JwtClalims{
		uid,
		username,
		jwt.StandardClaims{
			ExpiresAt: expireAt,
			Issuer:    "ender",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*JwtClalims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &JwtClalims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*JwtClalims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
