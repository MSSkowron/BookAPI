package services

import (
	"errors"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
)

var (
	ErrInvalidID            = errors.New("id must be a positive integer")
	ErrInvalidAuthor        = errors.New("author must not be empty")
	ErrInvalidTitle         = errors.New("title must not be empty")
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

// NewBookService returns an implementation of BookService interface
func NewBookService(db database.Database) BookService {
	return &BookServiceImpl{db: db}
}

// GetBooks returns all books from the database
func (bs *BookServiceImpl) GetBooks() ([]*models.Book, error) {
	books, err := bs.db.SelectAllBooks()
	if err != nil {
		return nil, err
	}

	return books, nil
}

// GetBook returns a book with the given id from the database
func (bs *BookServiceImpl) GetBook(id int) (*models.Book, error) {
	book, err := bs.db.SelectBookByID(id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, ErrBookNotFound
	}

	return book, nil
}

// AddBook adds a book to the database
func (bs *BookServiceImpl) AddBook(author, title string) (*models.Book, error) {
	if !bs.validateAuthor(author) {
		return nil, ErrInvalidAuthor
	}
	if !bs.validateTitle(title) {
		return nil, ErrInvalidTitle
	}

	book := &models.Book{
		Author: author,
		Title:  title,
	}
	id, err := bs.db.InsertBook(book)
	if err != nil {
		return nil, err
	}

	book.ID = id

	return book, nil
}

// UpdateBook updates a book with the given id in the database
func (bs *BookServiceImpl) UpdateBook(id int, author, title string) (*models.Book, error) {
	if !bs.validateID(id) {
		return nil, ErrInvalidID
	}
	if !bs.validateAuthor(author) {
		return nil, ErrInvalidAuthor
	}
	if !bs.validateTitle(title) {
		return nil, ErrInvalidTitle
	}

	book, err := bs.db.SelectBookByID(id)
	if err != nil {
		return nil, ErrBookNotFound
	}

	book.Author = author
	book.Title = title

	if err := bs.db.UpdateBook(book); err != nil {
		return nil, err
	}

	return book, nil
}

// DeleteBook deletes a book with the given id from the database
func (bs *BookServiceImpl) DeleteBook(id int) error {
	if !bs.validateID(id) {
		return ErrInvalidID
	}

	book, err := bs.db.SelectBookByID(id)
	if err != nil || book == nil {
		return ErrBookNotFound
	}

	return bs.db.DeleteBook(id)
}

func (bs *BookServiceImpl) validateID(id int) bool {
	return id > 0
}

func (bs *BookServiceImpl) validateAuthor(author string) bool {
	return author != ""
}

func (bs *BookServiceImpl) validateTitle(title string) bool {
	return title != ""
}
