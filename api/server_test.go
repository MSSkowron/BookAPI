package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MSSkowron/BookRESTAPI/model"
	"github.com/MSSkowron/BookRESTAPI/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleRegister(t *testing.T) {
	server := NewBookRESTAPIServer("0.0.0.0:8080", "secret1234567890", 1*time.Minute, storage.NewMockStorage())
	testServer := httptest.NewServer(makeHTTPHandler(server.handleRegister))
	defer testServer.Close()

	createAccountRequest := model.CreateAccountRequest{
		Email:     "test@test.com",
		Password:  "test",
		FirstName: "test",
		LastName:  "test",
		Age:       30,
	}

	createAccountRequestJSON, err := json.Marshal(createAccountRequest)
	assert.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseBody := &model.User{}
	err = json.NewDecoder(resp.Body).Decode(responseBody)
	assert.NoError(t, err)

	assert.NotEmpty(t, responseBody.ID, "ID should not be empty")
	assert.Equal(t, responseBody.ID, 4, "ID should be equal to 4")
	assert.Equal(t, createAccountRequest.Email, responseBody.Email, "Email should be equal")
	assert.Equal(t, createAccountRequest.FirstName, responseBody.FirstName, "First name should be equal")
	assert.Equal(t, createAccountRequest.LastName, responseBody.LastName, "Last name should be equal")
	assert.Equal(t, createAccountRequest.Age, int64(responseBody.Age), "Age should be equal")
}

func TestHandleLogin(t *testing.T) {
	//TODO
}

func TestHandleGetBooks(t *testing.T) {
	//TODO
}

func TestHandlePostBook(t *testing.T) {
	//TODO
}

func TestHandleGetBookByID(t *testing.T) {
	//TODO
}

func TestHandlePutBookByID(t *testing.T) {
	//TODO
}

func TestHandleDeleteBookByID(t *testing.T) {
	//TODO
}
