package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	data := []struct {
		name     string
		password string
	}{
		{"ValidPassword", "password123"},
		{"EmptyPassword", ""},
		{"ShortPassword", "123"},
		{"LongPassword", "averylongpasswordthatexceedsthelimitof72charactersandshouldcausethe"},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			assert.NoError(t, err, "HashPassword should not return an error")

			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
			assert.NoError(t, err, "CompareHashAndPassword should not return an error")
		})
	}
}

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"ValidPassword", "password123"},
		{"EmptyPassword", ""},
		{"ShortPassword", "123"},
		{"LongPassword", "averylongpasswordthatexceedsthelimitof72charactersandshouldcausethe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, _ := bcrypt.GenerateFromPassword([]byte(tt.password), 10)

			err := CheckPassword(tt.password, string(hash))
			assert.NoError(t, err, "CheckPassword should not return an error")

			err = CheckPassword("wrongpassword", string(hash))
			assert.Error(t, err, "CheckPassword should return an error for incorrect password")
		})
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"ValidPassword", "password123"},
		{"EmptyPassword", ""},
		{"ShortPassword", "123"},
		{"LongPassword", "averylongpasswordthatexceedsthelimitof72charactersandshouldcausethe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, errHash := HashPassword(tt.password)
			assert.NoError(t, errHash, "HashPassword should not return an error")

			errCheck := CheckPassword(tt.password, hash)
			assert.NoError(t, errCheck, "CheckPassword should not return an error")
		})
	}
}
