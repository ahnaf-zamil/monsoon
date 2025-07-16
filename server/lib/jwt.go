package lib

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJWTTokenHelper interface {
	CreateNewToken(data string, expireAfter int64) (string, error)
}

type JWTTokenHelper struct {
	SecretKey []byte
}

func (j *JWTTokenHelper) CreateNewToken(data string, expireAfter int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Unix() + expireAfter,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(j.SecretKey)
	return tokenString, err
}

func GetJWTTokenHelper() IJWTTokenHelper {
	return &JWTTokenHelper{SecretKey: []byte(config.JWTSecret)}
}
