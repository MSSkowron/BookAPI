package storage

import (
	"github.com/MSSkowron/BookRESTAPI/types"
)

// Storage is an interface for storage
type Storage interface {
	InsertUser(user *types.User) (int, error)
	SelectUserByEmail(email string) (*types.User, error)
	InsertBook(book *types.Book) (int, error)
	SelectBookByID(id int) (*types.Book, error)
	SelectAllBooks() ([]*types.Book, error)
	DeleteBook(id int) error
	UpdateBook(book *types.Book) error
}
