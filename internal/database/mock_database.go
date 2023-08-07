package database

import (
	"fmt"

	"github.com/MSSkowron/BookRESTAPI/internal/models"
)

// MockDatabase is a mock implementation of Database interface
type MockDatabase struct {
	users []*models.User
	books []*models.Book
}

// NewMockDatabase creates a new MockDatabase
func NewMockDatabase() Database {
	return &MockDatabase{
		users: []*models.User{
			{
				ID:        1,
				Email:     "johndoe@net.eu",
				Password:  "johnpassword",
				FirstName: "John",
				LastName:  "Doe",
				Age:       30,
			},
			{
				ID:        2,
				Email:     "janedoe@net.eu",
				Password:  "janepassword",
				FirstName: "Jane",
				LastName:  "Doe",
				Age:       25,
			},
			{
				ID:        3,
				Email:     "jankowalski@net.pl",
				Password:  "janpassword",
				FirstName: "Jan",
				LastName:  "Kowalski",
				Age:       30,
			},
		},
		books: []*models.Book{
			{
				ID:     1,
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings",
			},
			{
				ID:     2,
				Author: "J.K. Rowling",
				Title:  "Harry Potter",
			},
			{
				ID:     3,
				Author: "Stephen King",
				Title:  "The Shining",
			},
		},
	}
}

// InsertUser inserts a new user
func (db *MockDatabase) InsertUser(user *models.User) (int, error) {
	for _, u := range db.users {
		if u.Email == user.Email {
			return -1, fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	db.users = append(db.users, user)

	return len(db.users), nil
}

// SelectUserByEmail selects a user with given email
func (db *MockDatabase) SelectUserByEmail(email string) (*models.User, error) {
	for _, user := range db.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, nil
}

// InsertBook inserts a new book
func (db *MockDatabase) InsertBook(book *models.Book) (int, error) {
	return 4, nil
}

// SelectBookByID selects a book with given ID
func (db *MockDatabase) SelectBookByID(id int) (*models.Book, error) {
	for _, book := range db.books {
		if book.ID == id {
			return book, nil
		}
	}

	return nil, nil
}

// SelectAllBooks selects all books
func (db *MockDatabase) SelectAllBooks() ([]*models.Book, error) {
	return db.books, nil
}

// DeleteBook deletes a book with given ID
func (db *MockDatabase) DeleteBook(id int) error {
	for i, book := range db.books {
		if book.ID == id {
			db.books = append(db.books[:i], db.books[i+1:]...)
			return nil
		}
	}

	return nil
}

// UpdateBook updates a book with given ID
func (db *MockDatabase) UpdateBook(book *models.Book) error {
	return nil
}
