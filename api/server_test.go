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
	mockStorage := storage.NewMockStorage()
	server := NewBookRESTAPIServer("0.0.0.0:8080", "secret1234567890", 1*time.Minute, mockStorage)
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

	responseBodyUser := model.User{}
	err = json.NewDecoder(resp.Body).Decode(&responseBodyUser)
	assert.NoError(t, err)

	assert.NotEmpty(t, responseBodyUser.ID, "ID should not be empty")
	assert.Equal(t, 4, responseBodyUser.ID, "ID should be equal to 4")
	assert.Equal(t, createAccountRequest.Email, responseBodyUser.Email, "Email should be equal")
	assert.Equal(t, createAccountRequest.FirstName, responseBodyUser.FirstName, "First name should be equal")
	assert.Equal(t, createAccountRequest.LastName, responseBodyUser.LastName, "Last name should be equal")
	assert.Equal(t, createAccountRequest.Age, int64(responseBodyUser.Age), "Age should be equal")

	resp, err = http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseBodyError struct {
		Error string `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBodyError)
	assert.NoError(t, err)

	assert.Equal(t, "user with this email already exists", responseBodyError.Error)
}

func TestHandleLogin(t *testing.T) {
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
