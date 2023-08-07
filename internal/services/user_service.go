package services

import "github.com/MSSkowron/BookRESTAPI/internal/database"

// UserService is an interface that defines the methods that the UserService struct must implement
type UserService interface {
	RegisterUser(email, password, firstName, lastName string, age int) (err error)
	LoginUser(email, password string) (token string, err error)
	ValidateToken(token string) (err error)
}

// UserServiceImpl is a struct that implements the UserService interface
type UserServiceImpl struct {
	db database.Database
}

func NewUserService(db database.Database) UserService {
	return &UserServiceImpl{db: db}
}

// RegisterUser registers a user in the database
func (us *UserServiceImpl) RegisterUser(email, password, firstName, lastName string, age int) (err error) {
	return nil
}

// LoginUser logs a user in and returns a token
func (us *UserServiceImpl) LoginUser(email, password string) (token string, err error) {
	return "", nil
}

// ValidateToken validates a token
func (us *UserServiceImpl) ValidateToken(token string) (err error) {
	return nil
}
