package storage

import (
	"github.com/MSSkowron/BookRESTAPI/model"
)

// Storage is an interface for storage
type Storage interface {
	InsertUser(user *model.User) (int, error)
	SelectUserByEmail(email string) (*model.User, error)
	InsertBook(book *model.Book) (int, error)
	SelectBookByID(id int) (*model.Book, error)
	SelectAllBooks() ([]*model.Book, error)
	DeleteBook(id int) error
	UpdateBook(book *model.Book) error
}
