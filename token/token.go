package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when the token is expired
	ErrExpiredToken = errors.New("token is expired")
)

// Generate generates a new JWT token
func Generate(userID int, userEmail, secret string, expirationTime time.Duration) (tokenString string, err error) {
	claims := &jwt.MapClaims{
		"id":        userID,
		"email":     userEmail,
		"expiresAt": time.Now().Add(expirationTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// Validate validates a JWT token
func Validate(tokenString, secret string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	if int64(token.Claims.(jwt.MapClaims)["expiresAt"].(float64)) < time.Now().Local().Unix() {
		return ErrExpiredToken
	}

	return nil
}
