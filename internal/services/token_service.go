package services

import (
	"errors"
	"time"

	"github.com/MSSkowron/BookRESTAPI/pkg/token"
)

var (
	// ErrInvalidToken is returned when an invalid token is provided.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when an expired token is provided.
	ErrExpiredToken = errors.New("token is expired")
)

// TokenService is an interface that defines the methods that the TokenService must implement.
type TokenService interface {
	GenerateToken(int, string) (string, error)
	ValidateToken(string) error
	GetUserIDFromToken(string) (int, error)
}

// TokenServiceImpl implements the TokenService interface.
type TokenServiceImpl struct {
	tokenSecret   string
	tokenDuration time.Duration
}

// NewTokenService creates a new TokenServiceImpl.
func NewTokenService(tokenSecret string, tokenDuration time.Duration) *TokenServiceImpl {
	return &TokenServiceImpl{
		tokenSecret:   tokenSecret,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken generates a token.
func (ts *TokenServiceImpl) GenerateToken(userID int, userEmail string) (string, error) {
	return token.Generate(userID, userEmail, ts.tokenSecret, ts.tokenDuration)
}

// ValidateToken validates a token.
func (ts *TokenServiceImpl) ValidateToken(tokenString string) error {
	if err := token.Validate(tokenString, ts.tokenSecret); err != nil {
		if errors.Is(err, token.ErrExpiredToken) {
			return ErrExpiredToken
		}

		return ErrInvalidToken
	}

	return nil
}

// GetUserIDFromToken retrieves the user ID from a token.
func (ts *TokenServiceImpl) GetUserIDFromToken(tokenString string) (int, error) {
	id, err := token.GetUserID(tokenString, ts.tokenSecret)
	if err != nil {
		return 0, err
	}

	return id, nil
}
