package model

import (
	"math/rand"
	"strconv"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func NewUser(email, password, firstName, lastName string, age int) *User {
	return &User{
		ID:        strconv.Itoa(rand.Intn(10000)),
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}
