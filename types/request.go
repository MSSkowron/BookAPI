package types

// CreateAccountRequest represents a create account request
type CreateAccountRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int64  `json:"age"`
}

// NewCreateAccountRequest creates a new create account request
func NewCreateAccountRequest(email string, password string, firstName string, lastName string, age int64) *CreateAccountRequest {
	return &CreateAccountRequest{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewLoginRequest creates a new login request
func NewLoginRequest(email string, password string) *LoginRequest {
	return &LoginRequest{
		Email:    email,
		Password: password,
	}
}

// CreateBookRequest represents a create book request
type CreateBookRequest struct {
	CreatedBy int64  `json:"created_by"`
	Author    string `json:"author"`
	Title     string `json:"title"`
}

// NewCreateBookRequest creates a new create book request
func NewCreateBookRequest(author string, title string) *CreateBookRequest {
	return &CreateBookRequest{
		Author: author,
		Title:  title,
	}
}
