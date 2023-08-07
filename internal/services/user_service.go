package services

import (
	"errors"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
	"github.com/MSSkowron/BookRESTAPI/pkg/token"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidFirstName   = errors.New("invalid first name")
	ErrInvalidLastName    = errors.New("invalid last name")
	ErrInvalidAge         = errors.New("invalid age")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when a token is expired
	ErrExpiredToken = errors.New("token is expired")
)

// UserService is an interface that defines the methods that the UserService struct must implement
type UserService interface {
	RegisterUser(email, password, firstName, lastName string, age int) (user *models.User, err error)
	LoginUser(email, password string) (token string, err error)
	ValidateToken(token string) (err error)
}

// UserServiceImpl is a struct that implements the UserService interface
type UserServiceImpl struct {
	db            database.Database
	tokenSecret   string
	tokenDuration time.Duration
}

func NewUserService(db database.Database, tokenSecret string, tokenDuration time.Duration) UserService {
	return &UserServiceImpl{
		db:            db,
		tokenSecret:   tokenSecret,
		tokenDuration: tokenDuration,
	}
}

// RegisterUser registers a user in the database
func (us *UserServiceImpl) RegisterUser(email, password, firstName, lastName string, age int) (*models.User, error) {
	return nil, nil
}

// LoginUser logs a user in and returns a token
func (us *UserServiceImpl) LoginUser(email, password string) (string, error) {
	return "", nil
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
		if errors.Is(err, token.ErrInvalidToken) || errors.Is(err, token.ErrInvalidSignature) {
			return ErrInvalidToken
		}

		return err
	}

	return nil
}
