package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func getExpirationTime(envVar string, defaultDuration time.Duration) time.Duration {
	expirationStr := os.Getenv(envVar)
	if expirationStr == "" {
		return defaultDuration
	}

	duration, err := time.ParseDuration(expirationStr)
	if err != nil {
		return defaultDuration
	}

	return duration
}

func GenerateToken(userID int, name string) (string, string, error) {
	accessSecret := []byte(os.Getenv("JWT_SECRET"))
	refreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))

	accessExpiration := time.Now().Add(getExpirationTime("JWT_EXPIRATION_TIME", 15*time.Minute))
	refreshExpiration := time.Now().Add(getExpirationTime("JWT_REFRESH_EXPIRATION_TIME", 7*24*time.Hour))

	accessClaims := &Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiration),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiration),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateToken(tokenString string, isRefresh bool) (*Claims, error) {
	var secretKey []byte
	if isRefresh {
		secretKey = []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET"))
	} else {
		secretKey = []byte(os.Getenv("JWT_SECRET"))
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}