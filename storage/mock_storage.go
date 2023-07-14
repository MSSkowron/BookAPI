package storage

import (
	"fmt"

	"github.com/MSSkowron/BookRESTAPI/model"
)

var (
	users []*model.User = []*model.User{
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
	}
	books []*model.Book = []*model.Book{
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
	}
)

// MockStorage is a mock implementation of Storage interface
type MockStorage struct{}

// NewMockStorage creates a new MockStorage
func NewMockStorage() *MockStorage {
	return &MockStorage{}
}

// InsertUser inserts a new user
func (s *MockStorage) InsertUser(user *model.User) (int, error) {
	for _, u := range users {
		if u.Email == user.Email {
			return -1, fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	users = append(users, user)

	return len(users), nil
}

// SelectUserByEmail selects a user with given email
func (s *MockStorage) SelectUserByEmail(email string) (*model.User, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, nil
}

// InsertBook inserts a new book
func (s *MockStorage) InsertBook(book *model.Book) (int, error) {
	return 4, nil
}

// SelectBookByID selects a book with given ID
func (s *MockStorage) SelectBookByID(id int) (*model.Book, error) {
	for _, book := range books {
		if book.ID == id {
			return book, nil
		}
	}

	return nil, nil
}

// SelectAllBooks selects all books
func (s *MockStorage) SelectAllBooks() ([]*model.Book, error) {
	return books, nil
}

// DeleteBook deletes a book with given ID
func (s *MockStorage) DeleteBook(id int) error {
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			return nil
		}
	}

	return nil
}

// UpdateBook updates a book with given ID
func (s *MockStorage) UpdateBook(book *model.Book) error {
	return nil
}

// Reset resets the storage to its initial state
func (s *MockStorage) Reset() {
	users = []*model.User{
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
	}
	books = []*model.Book{
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
	}
}
