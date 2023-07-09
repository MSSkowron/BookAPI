package model

// User represents a user
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

// NewUser creates a new user
func NewUser(email string, password string, firstName string, lastName string, age int) *User {
	return &User{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}
