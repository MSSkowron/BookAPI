package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBooks(t *testing.T) {
	// TODO: Implement
}

func TestGetBook(t *testing.T) {
	// TODO: Implement
}

func TestAddBook(t *testing.T) {
	// TODO: Implement
}

func TestUpdateBook(t *testing.T) {
	// TODO: Implement
}

func TestDeleteBook(t *testing.T) {
	// TODO: Implement
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
