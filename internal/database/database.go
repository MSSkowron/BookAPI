package database

import (
	"github.com/MSSkowron/BookRESTAPI/internal/models"
)

// Database is an interface for database operations.
type Database interface {
	InsertUser(*models.User) (int, error)
	SelectUserByID(int) (*models.User, error)
	SelectUserByEmail(string) (*models.User, error)
	InsertBook(*models.Book) (int, error)
	SelectBookByID(int) (*models.Book, error)
	SelectAllBooks() ([]*models.Book, error)
	DeleteBook(int) error
	UpdateBook(int, *models.Book) error
}
