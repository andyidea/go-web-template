package util

import (
	"github.com/dgrijalva/jwt-go"
)

//生成token
func GenerateToken(data map[string]interface{}, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	for key, value := range data {
		token.Claims.(jwt.MapClaims)[key] = value
	}
	return token.SignedString([]byte(secret))
}

//解析token
func ParseToken(token_string string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err == nil && token.Valid {
		return token, nil
	} else {
		return nil, err
	}
}
