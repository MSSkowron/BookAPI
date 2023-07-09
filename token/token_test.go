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
		userId                int
		userEmail             string
		secret                string
		expirationTime        time.Duration
		expectedGenerateError error
		expectedValidateError error
	}{
		{
			name:                  "Valid token",
			userId:                1,
			userEmail:             "test@example.com",
			secret:                "secret",
			expirationTime:        time.Hour,
			expectedGenerateError: nil,
			expectedValidateError: nil,
		},
		{
			name:                  "Expired token",
			userId:                1,
			userEmail:             "test@example.com",
			secret:                "secret",
			expirationTime:        -time.Hour, // Negative expiration time
			expectedGenerateError: nil,
			expectedValidateError: ErrExpiredToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenString, err := Generate(test.userId, test.userEmail, test.secret, test.expirationTime)
			assert.Equal(t, test.expectedGenerateError, err)

			err = Validate(tokenString, test.secret)
			assert.Equal(t, test.expectedValidateError, err)
		})
	}
}
