package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/database"
	"github.com/MSSkowron/BookRESTAPI/internal/models"
	"github.com/MSSkowron/BookRESTAPI/pkg/crypto"
	"github.com/MSSkowron/BookRESTAPI/pkg/token"
	"github.com/gorilla/mux"
)

const (
	// ErrMsgBadRequestInvalidRequestBody is a message for bad request with invalid request body
	ErrMsgBadRequestInvalidRequestBody = "invalid request body"
	// ErrMsgBadRequestUserAlreadyExists is a message for bad request with user already exists
	ErrMsgBadRequestUserAlreadyExists = "user already exists"
	// ErrMsgBadRequestInvalidBookID is a message for bad request with invalid book id
	ErrMsgBadRequestInvalidBookID = "invalid book id"
	// ErrMsgUnauthorized is a message for unauthorized
	ErrMsgUnauthorized = "unauthorized"
	// ErrMsgUnauthorizedExpiredToken is a message for unauthorized with expired token
	ErrMsgUnauthorizedExpiredToken = "expired token"
	// ErrMsgUnauthorizedInvalidCredentials is a message for unauthorized with invalid credentials
	ErrMsgUnauthorizedInvalidCredentials = "invalid credentials"
	// ErrMsgNotFound is a message for not found
	ErrMsgNotFound = "not found"
	// ErrMsgInternalError is a message for internal error
	ErrMsgInternalError = "internal error"
)

type ServerHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandlerFunc(f ServerHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("[Server] Error while handling request: %v", err)
		}
	}
}

// Server is a HTTP server for handling REST API requests
type Server struct {
	listenAddr    string
	database      database.Database
	tokenSecret   string
	tokenDuration time.Duration
}

// NewServer creates a new Server
func NewServer(listenAddr, tokenSecret string, tokenDuration time.Duration, database database.Database) *Server {
	return &Server{
		listenAddr:    listenAddr,
		database:      database,
		tokenSecret:   tokenSecret,
		tokenDuration: tokenDuration,
	}
}

// Run runs the Server
func (s *Server) Run() error {
	r := mux.NewRouter()
	r.HandleFunc("/register", makeHTTPHandlerFunc(s.handleRegister)).Methods("POST")
	r.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin)).Methods("POST")
	r.HandleFunc("/books", validateJWT(makeHTTPHandlerFunc(s.handleGetBooks), s.tokenSecret)).Methods("GET")
	r.HandleFunc("/books", validateJWT(makeHTTPHandlerFunc(s.handlePostBook), s.tokenSecret)).Methods("POST")
	r.HandleFunc("/books/{id}", validateJWT(makeHTTPHandlerFunc(s.handleGetBookByID), s.tokenSecret)).Methods("GET")
	r.HandleFunc("/books/{id}", validateJWT(makeHTTPHandlerFunc(s.handlePutBookByID), s.tokenSecret)).Methods("PUT")
	r.HandleFunc("/books/{id}", validateJWT(makeHTTPHandlerFunc(s.handleDeleteBookByID), s.tokenSecret)).Methods("DELETE")

	log.Println("[Server] Server is running on: " + s.listenAddr)

	return http.ListenAndServe(s.listenAddr, r)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) error {
	log.Println("[Server] Called POST /register")

	createAccountRequest := &models.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	if createAccountRequest.Email == "" || createAccountRequest.Password == "" || createAccountRequest.FirstName == "" || createAccountRequest.LastName == "" || createAccountRequest.Age == 0 {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	user, _ := s.database.SelectUserByEmail(createAccountRequest.Email)
	if user != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestUserAlreadyExists)
		return nil
	}

	hashedPassword, err := crypto.HashPassword(createAccountRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("error while hashing password: %w", err)
	}

	newUser := models.NewUser(createAccountRequest.Email, hashedPassword, createAccountRequest.FirstName, createAccountRequest.LastName, int(createAccountRequest.Age))
	id, err := s.database.InsertUser(newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("error while inserting user: %w", err)
	}

	newUser.ID = id
	respondWithJSON(w, http.StatusOK, newUser)

	return nil
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	log.Println("[Server] Called POST /login")

	loginRequest := &models.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	user, err := s.database.SelectUserByEmail(loginRequest.Email)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorizedInvalidCredentials)
		return nil
	}

	if err := crypto.CheckPassword(loginRequest.Password, user.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorizedInvalidCredentials)
		return nil
	}

	token, err := token.Generate(user.ID, user.Email, s.tokenSecret, s.tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("error while generating token: %w", err)
	}

	respondWithJSON(w, http.StatusOK, models.LoginResponse{Token: token})

	return nil
}

func (s *Server) handleGetBooks(w http.ResponseWriter, r *http.Request) error {
	log.Println("[Server] Called GET /books")

	books, err := s.database.SelectAllBooks()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("error while getting books: %w", err)
	}

	respondWithJSON(w, http.StatusOK, books)

	return nil
}

func (s *Server) handlePostBook(w http.ResponseWriter, r *http.Request) error {
	log.Println("[Server] Called POST /books")

	createBookRequest := &models.CreateBookRequest{}
	if err := json.NewDecoder(r.Body).Decode(createBookRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	if createBookRequest.Title == "" || createBookRequest.Author == "" {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	newBook := models.NewBook(createBookRequest.Title, createBookRequest.Author)
	id, err := s.database.InsertBook(newBook)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("error while creating new book: %w", err)
	}

	newBook.ID = id
	respondWithJSON(w, http.StatusOK, newBook)

	return nil
}

func (s *Server) handleGetBookByID(w http.ResponseWriter, r *http.Request) error {
	log.Println("[Server] Called GET /books/{id}")

	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidBookID)
		return nil
	}

	book, err := s.database.SelectBookByID(id)
	if err != nil || book == nil {
		respondWithError(w, http.StatusNotFound, ErrMsgNotFound)
		return nil
	}

	respondWithJSON(w, http.StatusOK, book)

	return nil
}

func (s *Server) handlePutBookByID(w http.ResponseWriter, r *http.Request) error {
	log.Println("[Server] Called PUT /books/{id}")

	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidBookID)
		return nil
	}

	book, err := s.database.SelectBookByID(id)
	if err != nil || book == nil {
		respondWithError(w, http.StatusNotFound, ErrMsgNotFound)
		return nil
	}

	book = &models.Book{
		ID:        id,
		CreatedAt: book.CreatedAt,
	}
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	if err := s.database.UpdateBook(book); err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("error while updating the book: %s", err)
	}

	respondWithJSON(w, http.StatusOK, book)

	return nil
}

func (s *Server) handleDeleteBookByID(w http.ResponseWriter, r *http.Request) error {
	log.Println("[Server] Called DELETE /books/{id}")

	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidBookID)
		return nil
	}

	book, err := s.database.SelectBookByID(id)
	if err != nil || book == nil {
		respondWithError(w, http.StatusNotFound, ErrMsgNotFound)
		return nil
	}

	if err := s.database.DeleteBook(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("error while deleting the book: %s", err.Error())
	}

	respondWithJSON(w, http.StatusOK, nil)
	return nil
}

func validateJWT(f http.HandlerFunc, tokenSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorized)
			return
		}

		tokenString := authHeaderParts[1]
		if err := token.Validate(tokenString, tokenSecret); err != nil {
			if errors.Is(err, token.ErrExpiredToken) {
				respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorizedExpiredToken)
			} else {
				respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorized)
			}

			return
		}

		f(w, r)
	}
}

func respondWithError(w http.ResponseWriter, errCode int, errMessage string) {
	respondWithJSON(w, errCode, models.ErrorResponse{Error: errMessage})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[Server] Error while marshaling JSON response: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrMsgInternalError))

		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(response)
}
