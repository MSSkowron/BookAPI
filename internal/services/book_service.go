package services

import (
	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
)

// BookService is an interface that defines the methods that the BookService struct must implement
type BookService interface {
	GetBooks() (books []*models.Book, err error)
	GetBook(id int) (book *models.Book, err error)
	AddBook(book *models.Book) (err error)
	UpdateBook(id int, book *models.Book) (err error)
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
func (bs *BookServiceImpl) AddBook(book *models.Book) (err error) {
	return nil
}

// UpdateBook updates a book with the given id in the database
func (bs *BookServiceImpl) UpdateBook(id int, book *models.Book) (err error) {
	return nil
}

// DeleteBook deletes a book with the given id from the database
func (bs *BookServiceImpl) DeleteBook(id int) (err error) {
	return nil
}
