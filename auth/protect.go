package auth

import (
	"fmt"

	jwt "github.com/golang-jwt/jwt/v5"
)

func Protect(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte("==signature=="), nil
	})

	return err
}
