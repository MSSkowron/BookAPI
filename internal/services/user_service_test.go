package services

import (
	"testing"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/dtos"
	"github.com/MSSkowron/BookRESTAPI/pkg/crypto"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	mockDB := database.NewMockDatabase()

	us := NewUserService(mockDB, "", 0)

	hashedPassword, _ := crypto.HashPassword("Password1")

	data := []struct {
		name     string
		input    *dtos.AccountCreateDTO
		expected struct {
			user *dtos.UserDTO
			err  error
		}
	}{
		{
			name: "valid user",
			input: &dtos.AccountCreateDTO{
				Email:     "user@email.com",
				Password:  "Password1",
				FirstName: "John",
				LastName:  "Doe",
				Age:       20,
			},
			expected: struct {
				user *dtos.UserDTO
				err  error
			}{
				user: &dtos.UserDTO{
					ID:        4,
					Email:     "user@email.com",
					Password:  hashedPassword,
					FirstName: "John",
					LastName:  "Doe",
					Age:       20,
				},
				err: nil,
			},
		},
		{
			name: "invalid email",
			input: &dtos.AccountCreateDTO{
				Email:     "invalid email",
				Password:  "Password1",
				FirstName: "John",
				LastName:  "Doe",
				Age:       20,
			},
			expected: struct {
				user *dtos.UserDTO
				err  error
			}{
				user: nil,
				err:  ErrInvalidEmail,
			},
		},
		{
			name: "invalid password",
			input: &dtos.AccountCreateDTO{
				Email:     "user@email.com",
				Password:  "invalid password",
				FirstName: "John",
				LastName:  "Doe",
				Age:       20,
			},
			expected: struct {
				user *dtos.UserDTO
				err  error
			}{
				user: nil,
				err:  ErrInvalidPassword,
			},
		},
		{
			name: "invalid first name",
			input: &dtos.AccountCreateDTO{
				Email:     "user@email.com",
				Password:  "Password1",
				FirstName: "X07.@",
				LastName:  "Doe",
				Age:       20,
			},
			expected: struct {
				user *dtos.UserDTO
				err  error
			}{
				user: nil,
				err:  ErrInvalidFirstName,
			},
		},
		{
			name: "invalid last name",
			input: &dtos.AccountCreateDTO{
				Email:     "user@email.com",
				Password:  "Password1",
				FirstName: "John",
				LastName:  "X07.@",
				Age:       20,
			},
			expected: struct {
				user *dtos.UserDTO
				err  error
			}{
				user: nil,
				err:  ErrInvalidLastName,
			},
		},
		{
			name: "invalid age",
			input: &dtos.AccountCreateDTO{
				Email:     "user@email.com",
				Password:  "Password1",
				FirstName: "John",
				LastName:  "Doe",
				Age:       -1,
			},
			expected: struct {
				user *dtos.UserDTO
				err  error
			}{
				user: nil,
				err:  ErrInvalidAge,
			},
		},
		{
			name: "user already exists",
			input: &dtos.AccountCreateDTO{
				Email:     "johndoe@net.eu",
				Password:  "Password1",
				FirstName: "John",
				LastName:  "Doe",
				Age:       20,
			},
			expected: struct {
				user *dtos.UserDTO
				err  error
			}{
				user: nil,
				err:  ErrUserAlreadyExists,
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			user, err := us.RegisterUser(d.input)
			if d.expected.user != nil {
				require.NotNil(t, user)
				require.Equal(t, d.expected.user.ID, user.ID)
				require.LessOrEqual(t, user.CreatedAt, time.Now())
				require.Equal(t, d.expected.user.Email, user.Email)
				require.Nil(t, crypto.CheckPassword("Password1", user.Password))
				require.Equal(t, d.expected.user.FirstName, user.FirstName)
				require.Equal(t, d.expected.user.LastName, user.LastName)
				require.Equal(t, d.expected.user.Age, user.Age)
			} else {
				require.Nil(t, user)
			}
			require.Equal(t, d.expected.err, err)
		})
	}
}

func TestLoginUser(t *testing.T) {
	mockDB := database.NewMockDatabase()

	us := NewUserService(mockDB, "secret12345", 3*time.Second)

	user, err := us.RegisterUser(&dtos.AccountCreateDTO{
		Email:     "johntestdoe@net.eu",
		Password:  "Password1",
		FirstName: "John",
		LastName:  "Doe",
		Age:       20,
	})
	require.NotNil(t, user)
	require.Nil(t, err)

	data := []struct {
		name     string
		input    *dtos.UserLoginDTO
		expected struct {
			token bool
			err   error
		}
	}{
		{
			name: "valid user",
			input: &dtos.UserLoginDTO{
				Email:    "johntestdoe@net.eu",
				Password: "Password1",
			},
			expected: struct {
				token bool
				err   error
			}{
				token: true,
				err:   nil,
			},
		},
		{
			name: "invalid email",
			input: &dtos.UserLoginDTO{
				Email:    "invalid email",
				Password: "Password1",
			},
			expected: struct {
				token bool
				err   error
			}{
				token: false,
				err:   ErrInvalidEmail,
			},
		},
		{
			name: "empty password",
			input: &dtos.UserLoginDTO{
				Email:    "johntestdoe@net.eu",
				Password: "",
			},
			expected: struct {
				token bool
				err   error
			}{
				token: false,
				err:   ErrEmptyPassword,
			},
		},
		{
			name: "not existing user",
			input: &dtos.UserLoginDTO{
				Email:    "notexisting@net.eu",
				Password: "Password1",
			},
			expected: struct {
				token bool
				err   error
			}{
				token: false,
				err:   ErrInvalidCredentials,
			},
		},
		{
			name: "invalid password",
			input: &dtos.UserLoginDTO{
				Email:    "johntestdoe@net.eu",
				Password: "invalid password",
			},
			expected: struct {
				token bool
				err   error
			}{
				token: false,
				err:   ErrInvalidCredentials,
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			token, err := us.LoginUser(d.input)
			if d.expected.token {
				require.NotEmpty(t, token)
			} else {
				require.Empty(t, token)
			}
			require.Equal(t, d.expected.err, err)
		})
	}
}

func TestGenerateValidateToken(t *testing.T) {
	us := NewUserService(nil, "secret12345", 3*time.Second)

	token, err := us.GenerateToken(1, "email@net.com")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NoError(t, us.ValidateToken(token))

	time.Sleep(4 * time.Second)
	require.ErrorIs(t, us.ValidateToken(token), ErrExpiredToken)

	require.ErrorIs(t, us.ValidateToken("invalid token"), ErrInvalidToken)
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
			require.Equal(t, d.expected, us.validateEmail(d.email))
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
			require.Equal(t, d.expected, us.validatePassword(d.password))
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
			require.Equal(t, d.expected, us.validateFirstName(d.firstName))
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
			require.Equal(t, d.expected, us.validateLastName(d.lastName))
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
			require.Equal(t, d.expected, us.validateAge(d.age))
		})
	}
}
