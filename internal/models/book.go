package models

import "time"

// Book represents a model for a book
type Book struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
}
