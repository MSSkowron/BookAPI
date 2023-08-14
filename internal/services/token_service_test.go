package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenService(t *testing.T) {
	ts := NewTokenService("secret12345", 3*time.Second)

	// Generate Token
	token, err := ts.GenerateToken(1, "email@net.com")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NoError(t, ts.ValidateToken(token))

	// Validate Token
	time.Sleep(4 * time.Second)
	require.Equal(t, ts.ValidateToken(token), ErrExpiredToken)

	require.Equal(t, ts.ValidateToken("invalid token"), ErrInvalidToken)
	require.Equal(t, ts.ValidateToken(""), ErrInvalidToken)

	// Get userID from Token
	id, err := ts.GetUserIDFromToken(token)
	require.NoError(t, err)
	require.Equal(t, 1, id)

	id, err = ts.GetUserIDFromToken("invalid token")
	require.Equal(t, err, ErrInvalidToken)
	require.Equal(t, 0, id)
}
