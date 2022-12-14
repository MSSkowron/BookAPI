package model

type Book struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
}

func NewBook(title string, author string) *Book {
	return &Book{
		Title:  title,
		Author: author,
	}
}
