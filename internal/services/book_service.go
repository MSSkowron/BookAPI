package services

import (
	"errors"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
)

var (
	ErrInvalidAuthor        = errors.New("invalid author")
	ErrInvalidTitle         = errors.New("invalid title")
	ErrInvalidAuthorOrTitle = errors.New("invalid author or title")
	ErrBookNotFound         = errors.New("book not found")
)

// BookService is an interface that defines the methods that the BookService struct must implement
type BookService interface {
	GetBooks() (books []*models.Book, err error)
	GetBook(id int) (book *models.Book, err error)
	AddBook(author, title string) (book *models.Book, err error)
	UpdateBook(id int, author, title string) (updatedBook *models.Book, err error)
	DeleteBook(id int) (err error)
}

// BookServiceImpl is a struct that implements the BookService interface
type BookServiceImpl struct {
	db database.Database
}

func NewBookService(db database.Database) BookService {
	return &BookServiceImpl{db: db}
}

// GetBooks returns all books from the database
func (bs *BookServiceImpl) GetBooks() (books []*models.Book, err error) {
	return nil, nil
}

// GetBook returns a book with the given id from the database
func (bs *BookServiceImpl) GetBook(id int) (book *models.Book, err error) {
	return nil, nil
}

// AddBook adds a book to the database
func (bs *BookServiceImpl) AddBook(author, title string) (book *models.Book, err error) {
	return nil, nil
}

// UpdateBook updates a book with the given id in the database
func (bs *BookServiceImpl) UpdateBook(id int, author, title string) (updatedBook *models.Book, err error) {
	return nil, nil
}

// DeleteBook deletes a book with the given id from the database
func (bs *BookServiceImpl) DeleteBook(id int) (err error) {
	return nil
}
