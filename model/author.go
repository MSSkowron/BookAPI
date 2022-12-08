package model

type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func NewAuthor(firstName, lastName string, age int) *Author {
	return &Author{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}
