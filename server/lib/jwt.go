package lib

import (
	"fmt"
	"monsoon/util"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJWTTokenHelper interface {
	CreateNewToken(data string, expireAfter int64) (string, error)
	VerifyToken(tokenString string) (string, error)
}

type JWTTokenHelper struct {
	SecretKey []byte
}

func (j *JWTTokenHelper) CreateNewToken(data string, expireAfter int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Unix() + expireAfter,
		"iat":  time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(j.SecretKey)
	return tokenString, err
}

func (j *JWTTokenHelper) VerifyToken(tokenString string) (string, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.SecretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("token parsing error: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if data, ok := claims["data"].(string); ok {
			return data, nil
		}
		return "", fmt.Errorf("data claim missing or not a string")
	}

	return "", fmt.Errorf("invalid token")
}

func GetJWTTokenHelper() IJWTTokenHelper {
	config := util.GetConfig()
	return &JWTTokenHelper{SecretKey: []byte(config.JWTSecret)}
}
