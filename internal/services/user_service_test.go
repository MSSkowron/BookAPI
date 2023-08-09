package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {

}

func TestLoginUser(t *testing.T) {
}

func TestGenerateValidateToken(t *testing.T) {
	us := NewUserService(nil, "secret12345", 3*time.Second)

	token, err := us.GenerateToken(1, "email@net.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	assert.NoError(t, us.ValidateToken(token))

	time.Sleep(4 * time.Second)
	assert.ErrorIs(t, us.ValidateToken(token), ErrExpiredToken)

	assert.ErrorIs(t, us.ValidateToken("invalid token"), ErrInvalidToken)
}

func TestValidateEmail(t *testing.T) {
	us := NewUserService(nil, "", 0)

	data := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "valid email",
			email:    "email@net.com",
			expected: true,
		},
		{
			name:     "invalid email - empty",
			email:    "",
			expected: false,
		},
		{
			name:     "invalid email - no @",
			email:    "email.net.com",
			expected: false,
		},
		{
			name:     "invalid email - no domain",
			email:    "email@net",
			expected: false,
		},
		{
			name:     "invalid email - no username",
			email:    "@net.com",
			expected: false,
		},
		{
			name:     "invalid email - no extension",
			email:    "email@net.",
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, us.validateEmail(d.email))
		})
	}
}

func TestValidatePassword(t *testing.T) {
	us := NewUserService(nil, "", 0)

	data := []struct {
		name     string
		password string
		expected bool
	}{
		{
			name:     "valid password",
			password: "Password1",
			expected: true,
		},
		{
			name:     "invalid password - empty",
			password: "",
			expected: false,
		},
		{
			name:     "invalid password - too short",
			password: "Pass1",
			expected: false,
		},
		{
			name:     "invalid password - no uppercase letter",
			password: "password1",
			expected: false,
		},
		{
			name:     "invalid password - no lowercase letter",
			password: "PASSWORD1",
			expected: false,
		},
		{
			name:     "invalid password - no digit",
			password: "Password",
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, us.validatePassword(d.password))
		})
	}
}

func TestValidateFirstName(t *testing.T) {
	us := NewUserService(nil, "", 0)

	data := []struct {
		name      string
		firstName string
		expected  bool
	}{
		{
			name:      "valid first name",
			firstName: "John",
			expected:  true,
		},
		{
			name:      "invalid first name - empty",
			firstName: "",
			expected:  false,
		},
		{
			name:      "invalid first name - too short",
			firstName: "J",
			expected:  false,
		},
		{
			name:      "invalid first name - contains numbers",
			firstName: "John1",
			expected:  false,
		},
		{
			name:      "invalid first name - contains special characters",
			firstName: "John@",
			expected:  false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, us.validateFirstName(d.firstName))
		})
	}
}

func TestValidateLastName(t *testing.T) {
	us := NewUserService(nil, "", 0)

	data := []struct {
		name     string
		lastName string
		expected bool
	}{
		{
			name:     "valid first name",
			lastName: "Doe",
			expected: true,
		},
		{
			name:     "invalid first name - empty",
			lastName: "",
			expected: false,
		},
		{
			name:     "invalid first name - too short",
			lastName: "D",
			expected: false,
		},
		{
			name:     "invalid first name - contains numbers",
			lastName: "Doe1",
			expected: false,
		},
		{
			name:     "invalid first name - contains special characters",
			lastName: "Doe@",
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, us.validateLastName(d.lastName))
		})
	}
}

func TestValidateAge(t *testing.T) {
	us := NewUserService(nil, "", 0)

	data := []struct {
		name     string
		age      int
		expected bool
	}{
		{
			name:     "valid age",
			age:      25,
			expected: true,
		},
		{
			name:     "invalid age - negative",
			age:      -1,
			expected: false,
		},
		{
			name:     "invalid age - zero",
			age:      0,
			expected: false,
		},
		{
			name:     "invalid age - too old",
			age:      121,
			expected: false,
		},
		{
			name:     "invalid age - too young",
			age:      17,
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, us.validateAge(d.age))
		})
	}
}
