package services

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
	"github.com/MSSkowron/BookRESTAPI/pkg/crypto"
	"github.com/MSSkowron/BookRESTAPI/pkg/token"
)

var (
	// ErrInvalidEmail is returned when an invalid email address is provided.
	ErrInvalidEmail = errors.New("email must not be empty and must be a valid email address")
	// ErrInvalidPassword is returned when an invalid password is provided.
	// Password must have at least 6 characters, including 1 uppercase letter,
	// 1 lowercase letter, and 1 digit.
	ErrInvalidPassword = errors.New("password must not be empty and must be have at least 6 characters, including 1 uppercase letter, 1 lowercase letter, and 1 digit")
	// ErrEmptyPassword is returned when an empty password is provided.
	ErrEmptyPassword = errors.New("password must not be empty")
	// ErrInvalidFirstName is returned when an invalid first name is provided.
	// First name must consist of alphabetic characters and spaces, with at least 2 characters.
	ErrInvalidFirstName = errors.New("first name must must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters")
	// ErrInvalidLastName is returned when an invalid last name is provided.
	// Last name must consist of alphabetic characters and spaces, with at least 2 characters.
	ErrInvalidLastName = errors.New("last name must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters")
	// ErrInvalidAge is returned when an invalid age is provided.
	// Age must be between 18 and 120.
	ErrInvalidAge = errors.New("age must must not be empty and must be between 18 and 120")
	// ErrInvalidCredentials is returned when invalid user credentials are provided.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserAlreadyExists is returned when a user with the same details already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidToken is returned when an invalid token is provided.
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when an expired token is provided.
	ErrExpiredToken = errors.New("token is expired")
)

// UserService is an interface that defines the methods that the UserService must implement
type UserService interface {
	RegisterUser(email, password, firstName, lastName string, age int) (user *models.User, err error)
	LoginUser(email, password string) (token string, err error)
	ValidateToken(token string) (err error)
}

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	db            database.Database
	tokenSecret   string
	tokenDuration time.Duration
}

// NewUserService creates a new UserServiceImpl
func NewUserService(db database.Database, tokenSecret string, tokenDuration time.Duration) *UserServiceImpl {
	return &UserServiceImpl{
		db:            db,
		tokenSecret:   tokenSecret,
		tokenDuration: tokenDuration,
	}
}

// RegisterUser registers a user
func (us *UserServiceImpl) RegisterUser(email, password, firstName, lastName string, age int) (*models.User, error) {
	if !us.validateEmail(email) {
		return nil, ErrInvalidEmail
	}
	if !us.validatePassword(password) {
		return nil, ErrInvalidPassword
	}
	if !us.validateFirstName(firstName) {
		return nil, ErrInvalidFirstName
	}
	if !us.validateLastName(lastName) {
		return nil, ErrInvalidLastName
	}
	if !us.validateAge(age) {
		return nil, ErrInvalidAge
	}

	if user, _ := us.db.SelectUserByEmail(email); user != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
	id, err := us.db.InsertUser(user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

// LoginUser logs a user in and returns a token
func (us *UserServiceImpl) LoginUser(email, password string) (string, error) {
	if !us.validateEmail(email) {
		return "", ErrInvalidEmail
	}
	if password == "" {
		return "", ErrEmptyPassword
	}

	user, err := us.db.SelectUserByEmail(email)
	if err != nil || user == nil {
		return "", ErrInvalidCredentials
	}

	if err := crypto.CheckPassword(password, user.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := us.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GenerateToken generates a token
func (us *UserServiceImpl) GenerateToken(userID int, userEmail string) (string, error) {
	return token.Generate(userID, userEmail, us.tokenSecret, us.tokenDuration)
}

// ValidateToken validates a token
func (us *UserServiceImpl) ValidateToken(tokenString string) error {
	if err := token.Validate(tokenString, us.tokenSecret); err != nil {
		if errors.Is(err, token.ErrExpiredToken) {
			return ErrExpiredToken
		}

		return ErrInvalidToken
	}

	return nil
}

// validateEmail validates an email address
func (us *UserServiceImpl) validateEmail(email string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`).MatchString(email)
}

// validatePassword validates a password for at least 6 characters, at least 1 uppercase letter, 1 lowercase letter, and 1 digit
func (us *UserServiceImpl) validatePassword(password string) bool {
	return len(password) >= 6 &&
		strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") &&
		strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") &&
		strings.ContainsAny(password, "0123456789")
}

// validateFirstName validates a name field for alphabetic characters and spaces with a minimum length
func (us *UserServiceImpl) validateFirstName(firstName string) bool {
	return len(firstName) >= 2 && regexp.MustCompile(`^[a-zA-Z ]+$`).MatchString(firstName)
}

// validateLastName validates a name field for alphabetic characters and spaces with a minimum length
func (us *UserServiceImpl) validateLastName(lastName string) bool {
	return len(lastName) >= 2 && regexp.MustCompile(`^[a-zA-Z ]+$`).MatchString(lastName)
}

// validateAge validates an age to be between 18 and 120
func (us *UserServiceImpl) validateAge(age int) bool {
	return age >= 18 && age <= 120
}
