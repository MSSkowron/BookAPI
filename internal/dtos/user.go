package dtos

import "time"

// UserDTO represents a data transfer object (DTO) for a user
type UserDTO struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int64     `json:"age"`
}

// AccountCreateDTO represents a data transfer object (DTO) for creating a user account request
type AccountCreateDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int64  `json:"age"`
}

// UserLoginDTO represents a data transfer object (DTO) for user login request
type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenDTO represents a data transfer object (DTO) for a token
type TokenDTO struct {
	Token string `json:"token"`
}
