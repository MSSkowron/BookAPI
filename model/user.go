package model

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func NewUser(email, firstName, lastName string, age int) *User {
	return &User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}
