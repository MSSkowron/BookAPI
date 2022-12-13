package model

type CreateAccountRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int64  `json:"age"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateBookRequest struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}
