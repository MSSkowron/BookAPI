package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testSecret         = "testsecret123"
	testUserID         = 1
	testUserEmail      = "test@example.com"
	testExpirationTime = time.Hour
)

func TestToken(t *testing.T) {
	tests := []struct {
		name                   string
		tokenUserID            int
		tokenUserEmail         string
		tokenExpirationTime    time.Duration
		generateSecret         string
		validateSecret         string
		getUserIDSecret        string
		expectedGenerateError  error
		expectedValidateError  error
		expectedGetUserIDError error
		expectedUserID         int
	}{
		{
			name:                   "valid token",
			tokenUserID:            testUserID,
			tokenUserEmail:         testUserEmail,
			tokenExpirationTime:    testExpirationTime,
			generateSecret:         testSecret,
			validateSecret:         testSecret,
			getUserIDSecret:        testSecret,
			expectedGenerateError:  nil,
			expectedValidateError:  nil,
			expectedGetUserIDError: nil,
			expectedUserID:         testUserID,
		},
		{
			name:                   "expired token",
			tokenUserID:            testUserID,
			tokenUserEmail:         testUserEmail,
			tokenExpirationTime:    -testExpirationTime,
			generateSecret:         testSecret,
			validateSecret:         testSecret,
			getUserIDSecret:        testSecret,
			expectedGenerateError:  nil,
			expectedValidateError:  ErrExpiredToken,
			expectedGetUserIDError: nil,
			expectedUserID:         0,
		},
		{
			name:                   "validate token with incorrect secret",
			tokenUserID:            testUserID,
			tokenUserEmail:         testUserEmail,
			tokenExpirationTime:    testExpirationTime,
			generateSecret:         testSecret,
			validateSecret:         "invalidsecret321",
			getUserIDSecret:        testSecret,
			expectedGenerateError:  nil,
			expectedValidateError:  ErrInvalidToken,
			expectedGetUserIDError: nil,
			expectedUserID:         0,
		},
		{
			name:                   "get user id from token with incorrect secret",
			tokenUserID:            testUserID,
			tokenUserEmail:         testUserEmail,
			tokenExpirationTime:    testExpirationTime,
			generateSecret:         testSecret,
			validateSecret:         testSecret,
			getUserIDSecret:        "invalidsecret321",
			expectedGenerateError:  nil,
			expectedValidateError:  nil,
			expectedGetUserIDError: ErrInvalidToken,
			expectedUserID:         0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenString, err := Generate(test.tokenUserID, test.tokenUserEmail, test.generateSecret, test.tokenExpirationTime)
			require.Equal(t, test.expectedGenerateError, err)

			err = Validate(tokenString, test.validateSecret)
			require.Equal(t, test.expectedValidateError, err)

			if err == nil {
				userID, err := GetUserID(tokenString, test.getUserIDSecret)
				require.Equal(t, test.expectedGetUserIDError, err)
				require.Equal(t, test.expectedUserID, userID)
			}
		})
	}

	// Invalid token format
	require.ErrorIs(t, Validate("123XDTOKEN", testSecret), ErrInvalidToken)
}
