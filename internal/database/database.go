package database

import (
	"github.com/MSSkowron/BookRESTAPI/internal/models"
)

// Database is an interface for database operations
type Database interface {
	InsertUser(user *models.User) (int, error)
	SelectUserByEmail(email string) (*models.User, error)
	InsertBook(book *models.Book) (int, error)
	SelectBookByID(id int) (*models.Book, error)
	SelectAllBooks() ([]*models.Book, error)
	DeleteBook(id int) error
	UpdateBook(book *models.Book) error
}
