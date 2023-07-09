package model

// Book is a model for a book
type Book struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
}

// NewBook creates a new book
func NewBook(title string, author string) *Book {
	return &Book{
		Title:  title,
		Author: author,
	}
}
