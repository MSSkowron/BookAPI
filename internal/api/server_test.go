package api

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

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
	"github.com/MSSkowron/BookRESTAPI/internal/services"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

const (
	testTokenSecret   = "test1234567890"
	testTokenDuration = 1 * time.Minute
)

func TestHandleRegister(t *testing.T) {
	mockDB := database.NewMockDatabase()

	userService := services.NewUserService(mockDB, testTokenSecret, testTokenDuration)
	bookService := services.NewBookService(mockDB)

	server := NewServer("", userService, bookService)

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
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: models.User{
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
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name: "invalid password - too short",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Sh0rt",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:password must not be empty and must be have at least 6 characters, including 1 uppercase letter, 1 lowercase letter, and 1 digit",
			},
		},
		{
			name: "invalid password - no capital lettetr",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "nocapitalletters123",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:password must not be empty and must be have at least 6 characters, including 1 uppercase letter, 1 lowercase letter, and 1 digit",
			},
		},
		{
			name: "invalid password - no uppercase lettetr",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "nouppercaseletter123",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:password must not be empty and must be have at least 6 characters, including 1 uppercase letter, 1 lowercase letter, and 1 digit",
			},
		},
		{
			name: "invalid password - no lowercase lettetr",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "NOLOWERCASELETTER123",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:password must not be empty and must be have at least 6 characters, including 1 uppercase letter, 1 lowercase letter, and 1 digit",
			},
		},
		{
			name: "invalid password - no digit",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "NODIGIT",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:password must not be empty and must be have at least 6 characters, including 1 uppercase letter, 1 lowercase letter, and 1 digit",
			},
		},
		{
			name: "invalid password - empty",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:password must not be empty and must be have at least 6 characters, including 1 uppercase letter, 1 lowercase letter, and 1 digit",
			},
		},
		{
			name: "invalid email - invalid format",
			input: models.CreateAccountRequest{
				Email:     "test-test.com",
				Password:  "Test123@",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:email must not be empty and must be a valid email address",
			},
		},
		{
			name: "invalid email - empty",
			input: models.CreateAccountRequest{
				Email:     "",
				Password:  "Test123@",
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:email must not be empty and must be a valid email address",
			},
		},
		{
			name: "invalid first name - too short",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "X",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:first name must must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters",
			},
		},
		{
			name: "invalid first name - no letters",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "123098",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:first name must must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters",
			},
		},
		{
			name: "invalid first name - empty",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:first name must must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters",
			},
		},
		{
			name: "invalid last name - too short",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "test",
				LastName:  "X",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:last name must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters",
			},
		},
		{
			name: "invalid last name - no letters",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "test",
				LastName:  "123098",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:last name must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters",
			},
		},
		{
			name: "invalid last name - empty",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "test",
				LastName:  "",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:last name must not be empty and must consists of alphabetic characters and spaces, with at least 2 characters",
			},
		},
		{
			name: "invalid age - too young",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "test",
				LastName:  "test",
				Age:       10,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:age must must not be empty and must be between 18 and 120",
			},
		},
		{
			name: "invalid age - too old",
			input: models.CreateAccountRequest{
				Email:     "test@test.com",
				Password:  "Test123@",
				FirstName: "test",
				LastName:  "test",
				Age:       250,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:age must must not be empty and must be between 18 and 120",
			},
		},
		{
			name: "missing required fields",
			input: models.CreateAccountRequest{
				FirstName: "test",
				LastName:  "test",
				Age:       30,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:email must not be empty and must be a valid email address",
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
				responseBody := models.User{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)

				require.NotEmpty(t, responseBody.Password, "Password should not be empty")
				require.Equal(t, d.expectedResponse.(models.User).ID, responseBody.ID, "ID should be equal")
				require.Equal(t, d.expectedResponse.(models.User).Email, responseBody.Email, "Email should be equal")
				require.Equal(t, d.expectedResponse.(models.User).FirstName, responseBody.FirstName, "First name should be equal")
				require.Equal(t, d.expectedResponse.(models.User).LastName, responseBody.LastName, "Last name should be equal")
				require.Equal(t, d.expectedResponse.(models.User).Age, responseBody.Age, "Age should be equal")
			case http.StatusBadRequest:
				responseError := models.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponse.(models.ErrorResponse).Error, responseError.Error)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}

	// Test if user with this email already exists
	createAccountRequest := models.CreateAccountRequest{
		Email:     "test@test.com",
		Password:  "Test123",
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

	responseError := models.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.Equal(t, "user already exists", responseError.Error)
}

func TestHandleLogin(t *testing.T) {
	mockDB := database.NewMockDatabase()

	userService := services.NewUserService(mockDB, testTokenSecret, testTokenDuration)
	bookService := services.NewBookService(mockDB)

	server := NewServer("", userService, bookService)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	createAccountRequest := models.CreateAccountRequest{
		Email:     "test@test.com",
		Password:  "Test123@#",
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
			input: models.LoginRequest{
				Email:    "test@test.com",
				Password: "Test123@#",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: models.LoginResponse{
				Token: "<token-value>",
			},
		},
		{
			name: "invalid email - empty",
			input: models.LoginRequest{
				Email:    "",
				Password: "Test123@#",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:email must not be empty and must be a valid email address",
			},
		},
		{
			name: "invalid email",
			input: models.LoginRequest{
				Email:    "invalidEmail@test.com",
				Password: "Test123@#",
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: models.ErrorResponse{
				Error: "invalid credentials",
			},
		},
		{
			name: "invalid password",
			input: models.LoginRequest{
				Email:    "test@test.com",
				Password: "invalidPassword0#@",
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: models.ErrorResponse{
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
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:password must not be empty",
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
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body:email must not be empty and must be a valid email address",
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
			expectedResponse: models.ErrorResponse{
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
				loginResponse := models.LoginResponse{}
				err = json.NewDecoder(resp.Body).Decode(&loginResponse)
				require.NoError(t, err)
				require.NotEmpty(t, loginResponse.Token)
			case http.StatusUnauthorized, http.StatusBadRequest:
				responseError := models.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponse.(models.ErrorResponse).Error, responseError.Error)
			default:
				t.Fatalf("unexpected status code: %d", d.expectedStatusCode)
			}
		})
	}
}

func TestHandlePostBook(t *testing.T) {
	mockDB := database.NewMockDatabase()

	userService := services.NewUserService(mockDB, testTokenSecret, testTokenDuration)
	bookService := services.NewBookService(mockDB)

	server := NewServer("", userService, bookService)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books", server.validateJWT(makeHTTPHandlerFunc(server.handlePostBook))).Methods("POST")

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	createAccountRequest := models.CreateAccountRequest{
		Email:     "test@test.com",
		Password:  "Test123@#",
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

	loginRequest := models.LoginRequest{
		Email:    "test@test.com",
		Password: "Test123@#",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := models.LoginResponse{}
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
			input: models.Book{
				Author: "test",
				Title:  "test",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: models.Book{
				ID:     4,
				Author: "test",
				Title:  "test",
			},
		},
		{
			name: "invalid author - empty",
			input: models.Book{
				Author: "",
				Title:  "test",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name: "invalid title - empty",
			input: models.Book{
				Author: "test",
				Title:  "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name: "no title",
			input: models.Book{
				Author: "test",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name: "no author",
			input: models.Book{
				Title: "test",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
				Error: "invalid request body",
			},
		},
		{
			name:               "no fields",
			input:              models.Book{},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: models.ErrorResponse{
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
			expectedResponse: models.ErrorResponse{
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
				responseBodyBook := models.Book{}
				err = json.NewDecoder(resp.Body).Decode(&responseBodyBook)
				require.NoError(t, err)

				require.Equal(t, d.expectedResponse.(models.Book).ID, responseBodyBook.ID)
				require.Equal(t, d.expectedResponse.(models.Book).Author, responseBodyBook.Author)
				require.Equal(t, d.expectedResponse.(models.Book).Title, responseBodyBook.Title)
			case http.StatusBadRequest:
				responseError := models.ErrorResponse{}
				err = json.NewDecoder(resp.Body).Decode(&responseError)
				require.NoError(t, err)

				require.NotEmpty(t, responseError.Error)
				require.Equal(t, d.expectedResponse.(models.ErrorResponse).Error, responseError.Error)
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

	responseError := models.ErrorResponse{}
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
	mockDB := database.NewMockDatabase()

	userService := services.NewUserService(mockDB, testTokenSecret, testTokenDuration)
	bookService := services.NewBookService(mockDB)

	server := NewServer("", userService, bookService)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books/{id}", server.validateJWT(makeHTTPHandlerFunc(server.handleGetBookByID))).Methods("GET")

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	createAccountRequest := models.CreateAccountRequest{
		Email:     "test@test.com",
		Password:  "Test123!",
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

	loginRequest := models.LoginRequest{
		Email:    "test@test.com",
		Password: "Test123!",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := models.LoginResponse{}
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
			expectedResponseBody: models.Book{
				ID:     1,
				Author: "J.R.R. Tolkien",
				Title:  "The Lord of the Rings",
			},
		},
		{
			name:               "not existing id",
			inputID:            100,
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: models.ErrorResponse{
				Error: "not found",
			},
		},
		{
			name:               "negative id",
			inputID:            -200,
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: models.ErrorResponse{
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
				responseBody := models.Book{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)

				require.Equal(t, d.expectedResponseBody, responseBody)
			case http.StatusNotFound:
				responseError := models.ErrorResponse{}
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

	responseError := models.ErrorResponse{}
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

	responseError = models.ErrorResponse{}
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
	mockDB := database.NewMockDatabase()

	userService := services.NewUserService(mockDB, testTokenSecret, testTokenDuration)
	bookService := services.NewBookService(mockDB)

	server := NewServer("", userService, bookService)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books", server.validateJWT(makeHTTPHandlerFunc(server.handleGetBooks))).Methods("GET")

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	createAccountRequest := models.CreateAccountRequest{
		Email:     "test@test.com",
		Password:  "Test123!",
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

	loginRequest := models.LoginRequest{
		Email:    "test@test.com",
		Password: "Test123!",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := models.LoginResponse{}
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

	responseBodyBooks := []models.Book{}
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

	responseError := models.ErrorResponse{}
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
	mockDB := database.NewMockDatabase()

	userService := services.NewUserService(mockDB, testTokenSecret, testTokenDuration)
	bookService := services.NewBookService(mockDB)

	server := NewServer("", userService, bookService)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books/{id}", server.validateJWT(makeHTTPHandlerFunc(server.handleDeleteBookByID))).Methods("DELETE")

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	createAccountRequest := models.CreateAccountRequest{
		Email:     "test@test.com",
		Password:  "Test123!",
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

	loginRequest := models.LoginRequest{
		Email:    "test@test.com",
		Password: "Test123!",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := models.LoginResponse{}
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
			name:               "invalid id - negative",
			inputID:            -100,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: models.ErrorResponse{
				Error: "invalid request body:id must be a positive integer",
			},
		},
		{
			name:               "invalid id - not existing",
			inputID:            100,
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: models.ErrorResponse{
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
				responseError := models.ErrorResponse{}
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

	responseError := models.ErrorResponse{}
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

	responseError = models.ErrorResponse{}
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
	mockDB := database.NewMockDatabase()

	userService := services.NewUserService(mockDB, testTokenSecret, testTokenDuration)
	bookService := services.NewBookService(mockDB)

	server := NewServer("", userService, bookService)

	mux := mux.NewRouter()
	mux.HandleFunc("/register", makeHTTPHandlerFunc(server.handleRegister)).Methods("POST")
	mux.HandleFunc("/login", makeHTTPHandlerFunc(server.handleLogin)).Methods("POST")
	mux.HandleFunc("/books/{id}", server.validateJWT(makeHTTPHandlerFunc(server.handlePutBookByID))).Methods("PUT")

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	createAccountRequest := models.CreateAccountRequest{
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

	loginRequest := models.LoginRequest{
		Email:    "test@test.com",
		Password: "test",
	}

	loginRequestJSON, err := json.Marshal(loginRequest)
	require.NoError(t, err)

	resp, err = http.Post(testServer.URL+"/login", "application/json", bytes.NewReader(loginRequestJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := models.LoginResponse{}
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
			input: models.Book{
				ID:     1,
				Author: "J. K. Rowling",
				Title:  "Harry Potter and the Philosopher's Stone",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: models.Book{
				ID:     1,
				Author: "J. K. Rowling",
				Title:  "Harry Potter and the Philosopher's Stone",
			},
		},
		{
			name: "invalid id",
			input: models.Book{
				ID:     100,
				Author: "J. K. Rowling",
				Title:  "Harry Potter and the Philosopher's Stone",
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponseBody: models.ErrorResponse{
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
			expectedResponseBody: models.ErrorResponse{
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
				req, err = http.NewRequest(http.MethodPut, testServer.URL+"/books/"+strconv.Itoa(d.input.(models.Book).ID), bytes.NewReader(requestBody))
			}
			require.NoError(t, err)

			req.Header.Set("Authorization", "Bearer "+loginResponse.Token)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, d.expectedStatusCode, resp.StatusCode)

			switch d.expectedStatusCode {
			case http.StatusOK:
				responseBody := models.Book{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				require.NoError(t, err)

				require.Equal(t, d.expectedResponseBody, responseBody)
			case http.StatusBadRequest, http.StatusNotFound:
				responseError := models.ErrorResponse{}
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
	requestBody, err := json.Marshal(models.Book{
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

	responseError := models.ErrorResponse{}
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

	responseError = models.ErrorResponse{}
	err = json.NewDecoder(resp.Body).Decode(&responseError)
	require.NoError(t, err)

	require.NotEmpty(t, responseError.Error)
	require.Equal(t, "unauthorized", responseError.Error)

	// test no token
	requestBody, err = json.Marshal(models.Book{
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
