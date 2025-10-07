package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/revandpratama/edc-service/config"
)

type CustomClaims struct {
	TerminalID string
	jwt.RegisteredClaims
}

func GenerateToken(terminalID string) (string, error) {

	expirationSecond, err := strconv.Atoi(config.ENV.JWT_EXPIRATION_SECOND)
	if err != nil || expirationSecond == 0 {
		expirationSecond = 30
	}

	expirationTime := time.Now().Add(time.Second * time.Duration(expirationSecond))

	claims := &CustomClaims{
		TerminalID: terminalID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "edc-service",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.ENV.JWT_SECRET))
	if err != nil {
		return "", err
	}

	tokenString = fmt.Sprintf("Bearer %s", tokenString)

	return tokenString, nil
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.ENV.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
