package services

import (
	"testing"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/dtos"
	"github.com/stretchr/testify/require"
)

func TestGetBooks(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	books, err := bs.GetBooks()
	require.Nil(t, err)
	require.NotNil(t, books)
	require.Equal(t, 3, len(books))

	require.Equal(t, int64(1), books[0].ID)
	require.LessOrEqual(t, books[0].CreatedAt, time.Now())
	require.Equal(t, "J.R.R. Tolkien", books[0].Author)
	require.Equal(t, "The Lord of the Rings", books[0].Title)

	require.Equal(t, int64(2), books[1].ID)
	require.LessOrEqual(t, books[1].CreatedAt, time.Now())
	require.Equal(t, "J.K. Rowling", books[1].Author)
	require.Equal(t, "Harry Potter", books[1].Title)

	require.Equal(t, int64(3), books[2].ID)
	require.LessOrEqual(t, books[2].CreatedAt, time.Now())
	require.Equal(t, "Stephen King", books[2].Author)
	require.Equal(t, "The Shining", books[2].Title)
}

func TestGetBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	data := []struct {
		name         string
		id           int
		expectedErr  error
		expectedBook *dtos.BookDTO
	}{
		{
			name:        "valid",
			id:          1,
			expectedErr: nil,
			expectedBook: &dtos.BookDTO{
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
			require.Equal(t, d.expectedErr, err)

			if d.expectedErr != nil {
				require.Nil(t, book)
			} else {
				require.NotNil(t, book)
				require.Equal(t, d.expectedBook.ID, book.ID)
				require.LessOrEqual(t, book.CreatedAt, time.Now())
				require.Equal(t, d.expectedBook.Author, book.Author)
				require.Equal(t, d.expectedBook.Title, book.Title)
			}
		})
	}
}

func TestAddBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)

	data := []struct {
		name             string
		inputCreatedByID int
		inputBook        *dtos.BookCreateDTO
		expectedErr      error
		expectedBook     *dtos.BookDTO
	}{
		{
			name:             "valid",
			inputCreatedByID: 1,
			inputBook: &dtos.BookCreateDTO{
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings - The Fellowship of the Ring",
			},
			expectedErr: nil,
			expectedBook: &dtos.BookDTO{
				ID:     4,
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings - The Fellowship of the Ring",
			},
		},
		{
			name:             "invalid author - empty author",
			inputCreatedByID: 1,
			inputBook: &dtos.BookCreateDTO{
				Author: "",
				Title:  "The Lord of the Rings - The Fellowship of the Ring",
			},
			expectedErr:  ErrInvalidAuthor,
			expectedBook: nil,
		},
		{
			name:             "invalid title - empty title",
			inputCreatedByID: 1,
			inputBook: &dtos.BookCreateDTO{
				Author: "J.R.R. Tolkien",
				Title:  "",
			},
			expectedErr:  ErrInvalidTitle,
			expectedBook: nil,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			book, err := bs.AddBook(d.inputCreatedByID, d.inputBook)
			require.Equal(t, d.expectedErr, err)

			if d.expectedErr != nil {
				require.Nil(t, book)
			} else {
				require.NotNil(t, book)
				require.Equal(t, d.expectedBook.ID, book.ID)
				require.LessOrEqual(t, book.CreatedAt, time.Now())
				require.Equal(t, d.expectedBook.Author, book.Author)
				require.Equal(t, d.expectedBook.Title, book.Title)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	bs := NewBookService(mockDB)
	time.Sleep(1 * time.Millisecond)

	data := []struct {
		name         string
		id           int
		inputBook    *dtos.BookDTO
		expectedErr  error
		expectedBook *dtos.BookDTO
	}{
		{
			name: "valid",
			id:   1,
			inputBook: &dtos.BookDTO{
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings - The Fellowship of the Ring",
			},
			expectedErr: nil,
			expectedBook: &dtos.BookDTO{
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
			name: "invalid author - empty author",
			id:   1,
			inputBook: &dtos.BookDTO{
				Author: "",
			},
			expectedErr: ErrInvalidAuthor,
		},
		{
			name: "invalid title - empty title",
			id:   1,
			inputBook: &dtos.BookDTO{
				Author: "J.R.R. Tolkien",
				Title:  "",
			},
			expectedErr: ErrInvalidTitle,
		},
		{
			name: "invalid id - non-existent id",
			id:   100,
			inputBook: &dtos.BookDTO{
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings - The Fellowship of the Ring",
			},
			expectedErr: ErrBookNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			book, err := bs.UpdateBook(d.id, d.inputBook)
			require.Equal(t, d.expectedErr, err)

			if d.expectedErr != nil {
				require.Nil(t, book)
			} else {
				require.NotNil(t, book)
				require.Equal(t, d.expectedBook.ID, book.ID)
				require.LessOrEqual(t, book.CreatedAt, time.Now())
				require.Equal(t, d.expectedBook.Author, book.Author)
				require.Equal(t, d.expectedBook.Title, book.Title)
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
			require.Equal(t, d.expected, bs.DeleteBook(d.id))
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
			require.Equal(t, d.expected, bs.validateID(d.id))
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
			require.Equal(t, d.expected, bs.validateAuthor(d.author))
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
			require.Equal(t, d.expected, bs.validateTitle(d.title))
		})
	}
}
