package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when the token is expired
	ErrExpiredToken = errors.New("token is expired")
	// ErrInvalidSignature is returned when the token signature is invalid
	ErrInvalidSignature = errors.New("invalid signature")
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
			return nil, ErrInvalidSignature
		}

		return []byte(secret), nil
	})
	if err != nil {
		return ErrInvalidToken
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	expiresAt, ok := token.Claims.(jwt.MapClaims)["expiresAt"].(float64)
	if !ok {
		return ErrInvalidToken
	}

	if int64(expiresAt) < time.Now().Local().Unix() {
		return ErrExpiredToken
	}

	return nil
}

// GetUserID returns the user ID from the JWT token
func GetUserID(tokenString, secret string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidSignature
		}

		return []byte(secret), nil
	})
	if err != nil {
		return 0, ErrInvalidToken
	}

	userID, ok := token.Claims.(jwt.MapClaims)["id"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	return int(userID), nil
}
