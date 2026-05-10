package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJwt(issuer string, secretKey string) (string, int64, error) {
	expireAt := time.Now().Add(time.Hour * 24).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: expireAt, // 1 day
	})
	jwt, err := claims.SignedString([]byte(secretKey))
	return jwt, expireAt, err
}

func GetIssuerFromJwt(tokenString string, secretKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("Invalid JWT claims")
	}

	return claims.Issuer, nil
}

func ParseJwt(cookie string, secretKey string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims.Issuer, nil
}
