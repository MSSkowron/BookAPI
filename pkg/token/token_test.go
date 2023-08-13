package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Test cases
	tests := []struct {
		name                  string
		userID                int
		userEmail             string
		secret                string
		expirationTime        time.Duration
		expectedGenerateError error
		expectedValidateError error
	}{
		{
			name:                  "valid token",
			userID:                1,
			userEmail:             "test@example.com",
			secret:                "secret",
			expirationTime:        time.Hour,
			expectedGenerateError: nil,
			expectedValidateError: nil,
		},
		{
			name:                  "expired token",
			userID:                1,
			userEmail:             "test@example.com",
			secret:                "secret",
			expirationTime:        -time.Hour, // Negative expiration time
			expectedGenerateError: nil,
			expectedValidateError: ErrExpiredToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenString, err := Generate(test.userID, test.userEmail, test.secret, test.expirationTime)
			assert.Equal(t, test.expectedGenerateError, err)

			err = Validate(tokenString, test.secret)
			assert.Equal(t, test.expectedValidateError, err)
		})
	}
}

func TestGetUserID(t *testing.T) {
	// Test cases
	tests := []struct {
		name                  string
		userID                int
		userEmail             string
		secret                string
		expirationTime        time.Duration
		expectedUserID        int
		expectedGenerateError error
	}{
		{
			name:                  "valid token",
			userID:                1,
			userEmail:             "test@net.com",
			secret:                "secret1234567890",
			expirationTime:        time.Minute,
			expectedUserID:        1,
			expectedGenerateError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenString, err := Generate(test.userID, test.userEmail, test.secret, test.expirationTime)
			assert.ErrorIs(t, err, test.expectedGenerateError)

			userID, err := GetUserID(tokenString, test.secret)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedUserID, userID)
		})
	}

	// Test invalid token
	userID, err := GetUserID("invalid_token123", "secret1234567890")
	assert.ErrorIs(t, err, ErrInvalidToken)
	assert.Equal(t, 0, userID)
}
