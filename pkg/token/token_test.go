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
	data := []struct {
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

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tokenString, err := Generate(d.tokenUserID, d.tokenUserEmail, d.generateSecret, d.tokenExpirationTime)
			require.Equal(t, d.expectedGenerateError, err)

			err = Validate(tokenString, d.validateSecret)
			require.Equal(t, d.expectedValidateError, err)

			if err == nil {
				userID, err := GetUserID(tokenString, d.getUserIDSecret)
				require.Equal(t, d.expectedGetUserIDError, err)
				require.Equal(t, d.expectedUserID, userID)
			}
		})
	}

	require.ErrorIs(t, Validate("123XDTOKEN", testSecret), ErrInvalidToken)
}
