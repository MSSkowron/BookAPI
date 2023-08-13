package database

import (
	"fmt"
	"time"

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
				CreatedAt: time.Now(),
				Email:     "johndoe@net.eu",
				Password:  "johnpassword",
				FirstName: "John",
				LastName:  "Doe",
				Age:       30,
			},
			{
				ID:        2,
				CreatedAt: time.Now(),
				Email:     "janedoe@net.eu",
				Password:  "janepassword",
				FirstName: "Jane",
				LastName:  "Doe",
				Age:       25,
			},
			{
				ID:        3,
				CreatedAt: time.Now(),
				Email:     "jankowalski@net.pl",
				Password:  "janpassword",
				FirstName: "Jan",
				LastName:  "Kowalski",
				Age:       30,
			},
		},
		books: []*models.Book{
			{
				ID:        1,
				CreatedBy: 1,
				CreatedAt: time.Now(),
				Author:    "J.R.R. Tolkien",
				Title:     "The Lord of the Rings",
			},
			{
				ID:        2,
				CreatedBy: 2,
				CreatedAt: time.Now(),
				Author:    "J.K. Rowling",
				Title:     "Harry Potter",
			},
			{
				ID:        3,
				CreatedBy: 3,
				CreatedAt: time.Now(),
				Author:    "Stephen King",
				Title:     "The Shining",
			},
		},
	}
}

// InsertUser inserts a new user
func (db *MockDatabase) InsertUser(user *models.User) (int, error) {
	user.ID = len(db.users) + 1

	for _, u := range db.users {
		if u.Email == user.Email {
			return -1, fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	db.users = append(db.users, user)

	return len(db.users), nil
}

// SelectUserByID selects a user with given ID
func (db *MockDatabase) SelectUserByID(id int) (*models.User, error) {
	for _, user := range db.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, nil
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
	book.ID = len(db.books) + 1

	db.books = append(db.books, book)

	return len(db.books), nil
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
func (db *MockDatabase) UpdateBook(id int, book *models.Book) error {
	for i, b := range db.books {
		if b.ID == id {
			db.books[i].Author = book.Author
			db.books[i].Title = book.Title

			return nil
		}
	}

	return nil
}
