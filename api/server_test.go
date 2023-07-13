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
	defer mockStorage.Reset()

	server := NewBookRESTAPIServer("0.0.0.0:8080", "secret1234567890", 1*time.Minute, mockStorage)
	testServer := httptest.NewServer(makeHTTPHandler(server.handleRegister))
	defer testServer.Close()

	data := []struct {
		name               string
		input              any
		expectedStatusCode int
	}{
		{
			name: "valid request",
			input: model.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "test",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "invalid request body",
			input: struct {
				Email     string `json:"email"`
				Password  int    `json:"password"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Age       int    `json:"age"`
			}{
				Email:     "test@test.com",
				Password:  123,
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "no all required fields",
			input: model.CreateAccountRequest{
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			inputJSON, err := json.Marshal(d.input)
			assert.NoError(t, err)

			resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(inputJSON))
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				responseBodyUser := model.User{}
				err = json.NewDecoder(resp.Body).Decode(&responseBodyUser)
				assert.NoError(t, err)

				assert.NotEmpty(t, responseBodyUser.ID, "ID should not be empty")
				assert.Equal(t, 4, responseBodyUser.ID, "ID should be equal to 4")
				assert.Equal(t, "test@test.com", responseBodyUser.Email, "Email should be equal")
				assert.Equal(t, "test", responseBodyUser.FirstName, "First name should be equal")
				assert.Equal(t, "test", responseBodyUser.LastName, "Last name should be equal")
				assert.Equal(t, 30, responseBodyUser.Age, "Age should be equal")
			case http.StatusBadRequest:
				var responseBodyError struct {
					Error string `json:"error"`
				}

				err = json.NewDecoder(resp.Body).Decode(&responseBodyError)
				assert.NoError(t, err)

				assert.NotEmpty(t, responseBodyError.Error)
				assert.Equal(t, "invalid request body", responseBodyError.Error)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}

	// Test if user with this email already exists
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

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseBodyError struct {
		Error string `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseBodyError)
	assert.NoError(t, err)

	assert.Equal(t, "user with this email already exists", responseBodyError.Error)
}

func TestHandleLogin(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewBookRESTAPIServer("0.0.0.0:8080", "secret1234567890", 1*time.Minute, mockStorage)
	mux := http.NewServeMux()

	mux.HandleFunc("/register", makeHTTPHandler(server.handleRegister))
	mux.HandleFunc("/login", makeHTTPHandler(server.handleLogin))

	testServer := httptest.NewServer(mux)
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

	data := []struct {
		name               string
		input              any
		expectedStatusCode int
	}{
		{
			name: "valid",
			input: model.LoginRequest{
				Email:    "test@test.com",
				Password: "test",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "invalid password",
			input: model.LoginRequest{
				Email:    "test@test.com",
				Password: "invalidPassword",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "invalid email",
			input: model.LoginRequest{
				Email:    "invalidEmail@test.com",
				Password: "test",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "no password",
			input: struct {
				Email string `json:"email"`
			}{
				Email: "test@test.com",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "no email",
			input: struct {
				Password string `json:"password"`
			}{
				Password: "test",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "bad request",
			input: struct {
				Email    string `json:"email"`
				Password int    `json:"password"`
			}{
				Email:    "test@test.com",
				Password: 123,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			loginRequestJSON, err := json.Marshal(d.input)
			assert.NoError(t, err)

			resp, err := http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				var responseBodyLoginResponse struct {
					Token string `json:"token"`
				}

				err = json.NewDecoder(resp.Body).Decode(&responseBodyLoginResponse)
				assert.NoError(t, err)

				assert.NotEmpty(t, responseBodyLoginResponse.Token)
			case http.StatusUnauthorized:
				var responseBodyError struct {
					Error string `json:"error"`
				}

				err = json.NewDecoder(resp.Body).Decode(&responseBodyError)
				assert.NoError(t, err)

				assert.NotEmpty(t, responseBodyError.Error)
				assert.Equal(t, "invalid credentials", responseBodyError.Error)
			case http.StatusBadRequest:
				var responseBodyError struct {
					Error string `json:"error"`
				}

				err = json.NewDecoder(resp.Body).Decode(&responseBodyError)
				assert.NoError(t, err)

				assert.NotEmpty(t, responseBodyError.Error)
				assert.Equal(t, "invalid request body", responseBodyError.Error)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}
}

func TestHandlePostBook(t *testing.T) {
	//TODO
}

func TestHandleGetBookByID(t *testing.T) {
	//TODO
}

func TestHandleGetBooks(t *testing.T) {
	//TODO
}

func TestHandleDeleteBookByID(t *testing.T) {
	//TODO
}

func TestHandlePutBookByID(t *testing.T) {
	//TODO
}
