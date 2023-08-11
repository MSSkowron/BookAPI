package dtos

import "time"

// BookDTO represents a data transfer object (DTO) for a book
type BookDTO struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
}

// BookCreateDTO represents a data transfer object (DTO) for creating a book request
type BookCreateDTO struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}
