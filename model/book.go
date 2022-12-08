package model

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

func NewBook(title string, author *Author) *Book {
	return &Book{
		Title:  title,
		Author: author,
	}
}
