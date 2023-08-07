package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/model"
	"github.com/MSSkowron/BookRESTAPI/internal/storage"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestHandleRegister(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewServer("", "secret1234567890", 1*time.Minute, mockStorage)
	mux := mux.NewRouter()

	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	data := []struct {
		name               string
		input              any
		expectedStatusCode int
		expectedResponse   any
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
			expectedResponse: model.User{
				ID:        4,
				Email:     "test@test.com",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
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
			expectedResponse: model.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name: "missing required fields",
			input: model.CreateAccountRequest{
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.ErrorResponse{
				Error: "invalid request body",
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			inputJSON, err := json.Marshal(d.input)
			require.NoError(t, err)

			resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(inputJSON))
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				responseBody := model.User{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)

				require.NotEmpty(t, responseBody.Password, "Password should not be empty")
				require.Equal(t, d.expectedResponse.(model.User).ID, responseBody.ID, "ID should be equal")
				require.Equal(t, d.expectedResponse.(model.User).Email, responseBody.Email, "Email should be equal")
				require.Equal(t, d.expectedResponse.(model.User).FirstName, responseBody.FirstName, "First name should be equal")
				require.Equal(t, d.expectedResponse.(model.User).LastName, responseBody.LastName, "Last name should be equal")
				require.Equal(t, d.expectedResponse.(model.User).Age, responseBody.Age, "Age should be equal")
			case http.StatusBadRequest:
				responseError := model.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponse.(model.ErrorResponse).Error, responseError.Error)
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
	require.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	responseError := model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.Equal(t, "user already exists", responseError.Error)
}

func TestHandleLogin(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewServer("", "secret1234567890", 1*time.Minute, mockStorage)
	mux := mux.NewRouter()

	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")

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
	require.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	data := []struct {
		name               string
		input              any
		expectedStatusCode int
		expectedResponse   any
	}{
		{
			name: "valid",
			input: model.LoginRequest{
				Email:    "test@test.com",
				Password: "test",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.LoginResponse{
				Token: "<token-value>",
			},
		},
		{
			name: "invalid password",
			input: model.LoginRequest{
				Email:    "test@test.com",
				Password: "invalidPassword",
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: model.ErrorResponse{
				Error: "invalid credentials",
			},
		},
		{
			name: "invalid email",
			input: model.LoginRequest{
				Email:    "invalidEmail@test.com",
				Password: "test",
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: model.ErrorResponse{
				Error: "invalid credentials",
			},
		},
		{
			name: "no password",
			input: struct {
				Email string `json:"email"`
			}{
				Email: "test@test.com",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name: "no email",
			input: struct {
				Password string `json:"password"`
			}{
				Password: "test",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.ErrorResponse{
				Error: "invalid request body",
			},
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
			expectedResponse: model.ErrorResponse{
				Error: "invalid request body",
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			loginRequestJSON, err := json.Marshal(d.input)
			require.NoError(t, err)

			resp, err := http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				loginResponse := model.LoginResponse{}
				err = json.NewDecoder(resp.Body).Decode(&loginResponse)
				require.NoError(t, err)
				require.NotEmpty(t, loginResponse.Token)
			case http.StatusUnauthorized, http.StatusBadRequest:
				responseError := model.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponse.(model.ErrorResponse).Error, responseError.Error)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}
}

func TestHandlePostBook(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewServer("", "secret1234567890", 1*time.Minute, mockStorage)
	mux := mux.NewRouter()

	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books", validateJWT(makeHTTPHandlerFunc(server.handlePostBook), server.tokenSecret)).Methods("POST")

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
	require.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginRequest := model.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := model.LoginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.Token)

	data := []struct {
		name               string
		input              interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "valid",
			input: model.Book{
				Author: "test",
				Title:  "test",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: model.Book{
				ID:     4,
				Author: "test",
				Title:  "test",
			},
		},
		{
			name: "no title",
			input: model.Book{
				Author: "test",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name: "bad request",
			input: struct {
				Author string `json:"author"`
				Title  int    `json:"title"`
			}{
				Author: "test",
				Title:  123,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: model.ErrorResponse{
				Error: "invalid request body",
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			bookJSON, err := json.Marshal(d.input)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, testServer.URL+"/books", bytes.NewReader(bookJSON))
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				responseBodyBook := model.Book{}
				err = json.NewDecoder(resp.Body).Decode(&responseBodyBook)
				require.NoError(t, err)

				require.Equal(t, d.expectedResponse.(model.Book).ID, responseBodyBook.ID)
				require.Equal(t, d.expectedResponse.(model.Book).Author, responseBodyBook.Author)
				require.Equal(t, d.expectedResponse.(model.Book).Title, responseBodyBook.Title)
			case http.StatusBadRequest:
				responseError := model.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponse.(model.ErrorResponse).Error, responseError.Error)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}

	// test invalid token
	req, err := http.NewRequest(http.MethodPost, testServer.URL+"/books", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer invalid_token")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	responseError := model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)

	// test no token
	req, err = http.NewRequest(http.MethodPost, testServer.URL+"/books", nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)
}

func TestHandleGetBookByID(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewServer("", "secret1234567890", 1*time.Minute, mockStorage)
	mux := mux.NewRouter()

	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books/{id}", validateJWT(makeHTTPHandlerFunc(server.handleGetBookByID), server.tokenSecret)).Methods("GET")

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
	require.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginRequest := model.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := model.LoginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.Token)

	data := []struct {
		name                 string
		inputID              int
		expectedStatusCode   int
		expectedResponseBody any
	}{
		{
			name:               "valid",
			inputID:            1,
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: model.Book{
				ID:     1,
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings",
			},
		},
		{
			name:               "invalid id",
			inputID:            100,
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: model.ErrorResponse{
				Error: "not found",
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/books/%d", testServer.URL, d.inputID), nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				responseBody := model.Book{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)

				require.Equal(t, d.expectedResponseBody, responseBody)
			case http.StatusNotFound:
				responseError := model.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponseBody, responseError)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}

	// test id is not a number
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/books/abc", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	responseError := model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "invalid book id", responseError.Error)

	// test invalid token
	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/books/"+strconv.Itoa(data[0].inputID), nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer invalid_token")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	responseError = model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)

	// test no token
	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/books/"+strconv.Itoa(data[0].inputID), nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)
}

func TestHandleGetBooks(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewServer("", "secret1234567890", 1*time.Minute, mockStorage)
	mux := mux.NewRouter()

	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books", validateJWT(makeHTTPHandlerFunc(server.handleGetBooks), server.tokenSecret)).Methods("GET")

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
	require.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginRequest := model.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := model.LoginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.Token)

	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/books", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	responseBodyBooks := []model.Book{}
	err = json.NewDecoder(resp.Body).Decode(&responseBodyBooks)
	require.NoError(t, err)

	require.Len(t, responseBodyBooks, 3)

	// test invalid token
	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/books", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer invalid_token")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	responseError := model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)

	// test no token
	req, err = http.NewRequest(http.MethodGet, testServer.URL+"/books", nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)
}

func TestHandleDeleteBookByID(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewServer("", "secret1234567890", 1*time.Minute, mockStorage)
	mux := mux.NewRouter()

	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books/{id}", validateJWT(makeHTTPHandlerFunc(server.handleDeleteBookByID), server.tokenSecret)).Methods("DELETE")

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
	require.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginRequest := model.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := model.LoginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.Token)

	data := []struct {
		name                 string
		inputID              int
		expectedStatusCode   int
		expectedResponseBody any
	}{
		{
			name:                 "valid",
			inputID:              1,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "null",
		},
		{
			name:               "invalid id",
			inputID:            100,
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: model.ErrorResponse{
				Error: "not found",
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/books/"+strconv.Itoa(d.inputID), nil)
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				responseBody, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				require.Equal(t, d.expectedResponseBody, string(responseBody))
			case http.StatusBadRequest, http.StatusNotFound:
				responseError := model.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponseBody, responseError)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}

	// test id is not a number
	req, err := http.NewRequest(http.MethodDelete, testServer.URL+"/books/abc", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	responseError := model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "invalid book id", responseError.Error)

	// test invalid token
	req, err = http.NewRequest(http.MethodDelete, testServer.URL+"/books/2", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer invalid_token")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	responseError = model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)

	// test no token
	req, err = http.NewRequest(http.MethodDelete, testServer.URL+"/books/3", nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)
}

func TestHandlePutBookByID(t *testing.T) {
	mockStorage := storage.NewMockStorage()
	defer mockStorage.Reset()

	server := NewServer("", "secret1234567890", 1*time.Minute, mockStorage)
	mux := mux.NewRouter()

	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books/{id}", validateJWT(makeHTTPHandlerFunc(server.handlePutBookByID), server.tokenSecret)).Methods("PUT")

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
	require.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/register", "application/json", bytes.NewReader(createAccountRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginRequest := model.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := model.LoginResponse{}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.Token)

	data := []struct {
		name                 string
		input                any
		expectedStatusCode   int
		expectedResponseBody any
	}{
		{
			name: "valid",
			input: model.Book{
				ID:     1,
				Author: "J. K. Rowling",
				Title:  "Harry Potter and the Philosopher's Stone",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: model.Book{
				ID:     1,
				Author: "J. K. Rowling",
				Title:  "Harry Potter and the Philosopher's Stone",
			},
		},
		{
			name: "invalid id",
			input: model.Book{
				ID:     100,
				Author: "J. K. Rowling",
				Title:  "Harry Potter and the Philosopher's Stone",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: model.ErrorResponse{
				Error: "not found",
			},
		},
		{
			name: "invalid body",
			input: struct {
				ID     int    `json:"id"`
				Author string `json:"author"`
				Title  int    `json:"title"`
			}{
				ID:     1,
				Author: "J. K. Rowling",
				Title:  1,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: model.ErrorResponse{
				Error: "invalid request body",
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var (
				req *http.Request
				err error
			)

			requestBody, err := json.Marshal(d.input)
			require.NoError(t, err)

			if d.name == "invalid body" {
				req, err = http.NewRequest(http.MethodPut, testServer.URL+"/books/1", bytes.NewReader(requestBody))
			} else {
				req, err = http.NewRequest(http.MethodPut, testServer.URL+"/books/"+strconv.Itoa(d.input.(model.Book).ID), bytes.NewReader(requestBody))
			}
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				responseBody := model.Book{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)

				require.Equal(t, d.expectedResponseBody, responseBody)
			case http.StatusBadRequest, http.StatusNotFound:
				responseError := model.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponseBody, responseError)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}

	// test id is not a number
	requestBody, err := json.Marshal(model.Book{
		ID:     3,
		Author: "J. K. Rowling",
		Title:  "Harry Potter and the Philosopher's Stone",
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, testServer.URL+"/books/abc", bytes.NewReader(requestBody))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	responseError := model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "invalid book id", responseError.Error)

	// test invalid token
	req, err = http.NewRequest(http.MethodPut, testServer.URL+"/books/3", bytes.NewReader(requestBody))
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer invalid_token")

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	responseError = model.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)

	// test no token
	requestBody, err = json.Marshal(model.Book{
		ID:     3,
		Author: "J. K. Rowling",
		Title:  "Harry Potter and the Philosopher's Stone",
	})
	require.NoError(t, err)

	req, err = http.NewRequest(http.MethodPut, testServer.URL+"/books/3", bytes.NewReader(requestBody))
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)
}
