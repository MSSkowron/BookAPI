package services

import (
	"testing"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetBooks(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	books, err := bs.GetBooks()
	assert.Nil(t, err)
	assert.NotNil(t, books)
	assert.Equal(t, 3, len(books))

	assert.Equal(t, 1, books[0].ID)
	assert.Equal(t, "J.R.R. Tolkien", books[0].Author)
	assert.Equal(t, "The Lord of the Rings", books[0].Title)

	assert.Equal(t, 2, books[1].ID)
	assert.Equal(t, "J.K. Rowling", books[1].Author)
	assert.Equal(t, "Harry Potter", books[1].Title)

	assert.Equal(t, 3, books[2].ID)
	assert.Equal(t, "Stephen King", books[2].Author)
	assert.Equal(t, "The Shining", books[2].Title)
}

func TestGetBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	data := []struct {
		name         string
		id           int
		expectedErr  error
		expectedBook *models.Book
	}{
		{
			name:        "valid",
			id:          1,
			expectedErr: nil,
			expectedBook: &models.Book{
				ID:     1,
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings",
			},
		},
		{
			name:         "invalid id - negative id",
			id:           -1,
			expectedErr:  ErrInvalidID,
			expectedBook: nil,
		},
		{
			name:         "invalid id - zero id",
			id:           0,
			expectedErr:  ErrInvalidID,
			expectedBook: nil,
		},
		{
			name:         "invalid id - non-existent id",
			id:           100,
			expectedErr:  ErrBookNotFound,
			expectedBook: nil,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			book, err := bs.GetBook(d.id)
			assert.Equal(t, d.expectedErr, err)

			if d.expectedErr != nil {
				assert.Nil(t, book)
			} else {
				assert.NotNil(t, book)
				assert.Equal(t, d.expectedBook.ID, book.ID)
				assert.Equal(t, d.expectedBook.Author, book.Author)
				assert.Equal(t, d.expectedBook.Title, book.Title)
			}
		})
	}
}

func TestAddBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	data := []struct {
		name         string
		author       string
		title        string
		expectedErr  error
		expectedBook *models.Book
	}{
		{
			name:        "valid",
			author:      "J.R.R. Tolkien",
			title:       "The Lord of the Rings - The Fellowship of the Ring",
			expectedErr: nil,
			expectedBook: &models.Book{
				ID:     4,
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings - The Fellowship of the Ring",
			},
		},
		{
			name:         "invalid author - empty author",
			author:       "",
			title:        "The Lord of the Rings - The Fellowship of the Ring",
			expectedErr:  ErrInvalidAuthor,
			expectedBook: nil,
		},
		{
			name:         "invalid title - empty title",
			author:       "J.R.R. Tolkien",
			title:        "",
			expectedErr:  ErrInvalidTitle,
			expectedBook: nil,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			book, err := bs.AddBook(d.author, d.title)
			assert.Equal(t, d.expectedErr, err)

			if d.expectedErr != nil {
				assert.Nil(t, book)
			} else {
				assert.NotNil(t, book)
				assert.Equal(t, d.expectedBook.ID, book.ID)
				assert.Equal(t, d.expectedBook.Author, book.Author)
				assert.Equal(t, d.expectedBook.Title, book.Title)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	data := []struct {
		name         string
		id           int
		author       string
		title        string
		expectedErr  error
		expectedBook *models.Book
	}{
		{
			name:        "valid",
			id:          1,
			author:      "J.R.R. Tolkien",
			title:       "The Lord of the Rings - The Fellowship of the Ring",
			expectedErr: nil,
			expectedBook: &models.Book{
				ID:     1,
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings - The Fellowship of the Ring",
			},
		},
		{
			name:        "invalid id - negative id",
			id:          -1,
			expectedErr: ErrInvalidID,
		},
		{
			name:        "invalid id - zero id",
			id:          0,
			expectedErr: ErrInvalidID,
		},
		{
			name:        "invalid author - empty author",
			id:          1,
			author:      "",
			expectedErr: ErrInvalidAuthor,
		},
		{
			name:        "invalid title - empty title",
			id:          1,
			author:      "J.R.R. Tolkien",
			title:       "",
			expectedErr: ErrInvalidTitle,
		},
		{
			name:        "invalid id - non-existent id",
			id:          100,
			author:      "J.R.R. Tolkien",
			title:       "The Lord of the Rings - The Fellowship of the Ring",
			expectedErr: ErrBookNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			book, err := bs.UpdateBook(d.id, d.author, d.title)
			assert.Equal(t, d.expectedErr, err)

			if d.expectedErr != nil {
				assert.Nil(t, book)
			} else {
				assert.NotNil(t, book)
				assert.Equal(t, d.expectedBook.ID, book.ID)
				assert.Equal(t, d.expectedBook.Author, book.Author)
				assert.Equal(t, d.expectedBook.Title, book.Title)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	data := []struct {
		name     string
		id       int
		expected error
	}{
		{
			name:     "valid id",
			id:       3,
			expected: nil,
		},
		{
			name:     "invalid id - negative id",
			id:       -1,
			expected: ErrInvalidID,
		},
		{
			name:     "invalid id - zero id",
			id:       0,
			expected: ErrInvalidID,
		},
		{
			name:     "invalid id - non-existent id",
			id:       100,
			expected: ErrBookNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, bs.DeleteBook(d.id))
		})
	}
}

func TestValidateID(t *testing.T) {
	bs := NewBookService(nil)

	data := []struct {
		name     string
		id       int
		expected bool
	}{
		{
			name:     "valid id",
			id:       1,
			expected: true,
		},
		{
			name:     "invalid id - negative id",
			id:       -1,
			expected: false,
		},
		{
			name:     "invalid id - zero id",
			id:       0,
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, bs.validateID(d.id))
		})
	}
}

func TestValidateAuthor(t *testing.T) {
	bs := NewBookService(nil)

	data := []struct {
		name     string
		author   string
		expected bool
	}{
		{
			name:     "valid author",
			author:   "J.R.R. Tolkien",
			expected: true,
		},
		{
			name:     "empty author",
			author:   "",
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, bs.validateAuthor(d.author))
		})
	}
}

func TestValidateTitle(t *testing.T) {
	bs := NewBookService(nil)

	data := []struct {
		name     string
		title    string
		expected bool
	}{
		{
			name:     "valid title",
			title:    "The Hobbit",
			expected: true,
		},
		{
			name:     "empty title",
			title:    "",
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			assert.Equal(t, d.expected, bs.validateTitle(d.title))
		})
	}
}
