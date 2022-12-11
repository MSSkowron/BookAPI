package model

type Author struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewAuthor(firstName, lastName string, age int) *Author {
	return &Author{
		FirstName: firstName,
		LastName:  lastName,
	}
}
