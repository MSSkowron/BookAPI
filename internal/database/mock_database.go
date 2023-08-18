package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/models"
)

// MockDatabase is a mock implementation of Database interface.
type MockDatabase struct {
	userMu sync.RWMutex
	bookMu sync.RWMutex
	users  []*models.User
	books  []*models.Book
}

// NewMockDatabase creates a new MockDatabase.
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

// Close closes the database connection.
func (db *MockDatabase) Close() {}

// InsertUser inserts a new user into the database.
func (db *MockDatabase) InsertUser(user *models.User) (int, error) {
	db.userMu.Lock()
	defer db.userMu.Unlock()

	user.ID = len(db.users) + 1

	for _, u := range db.users {
		if u.Email == user.Email {
			return -1, fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	db.users = append(db.users, user)

	return len(db.users), nil
}

// SelectUserByID selects a user with given ID from the database.
func (db *MockDatabase) SelectUserByID(id int) (*models.User, error) {
	db.userMu.RLock()
	defer db.userMu.RUnlock()

	for _, user := range db.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, nil
}

// SelectUserByEmail selects a user with given email from the database.
func (db *MockDatabase) SelectUserByEmail(email string) (*models.User, error) {
	db.userMu.RLock()
	defer db.userMu.RUnlock()

	for _, user := range db.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, nil
}

// InsertBook inserts a new book into the database.
func (db *MockDatabase) InsertBook(book *models.Book) (int, error) {
	db.bookMu.Lock()
	defer db.bookMu.Unlock()

	book.ID = len(db.books) + 1

	db.books = append(db.books, book)

	return len(db.books), nil
}

// SelectBookByID selects a book with given ID from the database.
func (db *MockDatabase) SelectBookByID(id int) (*models.Book, error) {
	db.bookMu.RLock()
	defer db.bookMu.RUnlock()

	for _, book := range db.books {
		if book.ID == id {
			return book, nil
		}
	}

	return nil, nil
}

// SelectAllBooks selects all books from the database.
func (db *MockDatabase) SelectAllBooks() ([]*models.Book, error) {
	db.bookMu.RLock()
	defer db.bookMu.RUnlock()

	return db.books, nil
}

// DeleteBook deletes a book with given ID from the database.
func (db *MockDatabase) DeleteBook(id int) error {
	db.bookMu.Lock()
	defer db.bookMu.Unlock()

	for i, book := range db.books {
		if book.ID == id {
			db.books = append(db.books[:i], db.books[i+1:]...)
			return nil
		}
	}

	return nil
}

// UpdateBook updates a book with given ID in the database.
func (db *MockDatabase) UpdateBook(id int, book *models.Book) error {
	db.bookMu.Lock()
	defer db.bookMu.Unlock()

	for i, b := range db.books {
		if b.ID == id {
			db.books[i].Author = book.Author
			db.books[i].Title = book.Title

			return nil
		}
	}

	return nil
}
