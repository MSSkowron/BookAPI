package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MSSkowron/BookRESTAPI/crypto"
	"github.com/MSSkowron/BookRESTAPI/storage"
	"github.com/MSSkowron/BookRESTAPI/token"
	"github.com/MSSkowron/BookRESTAPI/types"
	"github.com/gorilla/mux"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
		}
	}
}

// BookRESTAPIServer is a server for handling REST API requests
type BookRESTAPIServer struct {
	listenAddr    string
	storage       storage.Storage
	tokenSecret   string
	tokenDuration time.Duration
}

// NewBookRESTAPIServer creates a new BookRESTAPIServer
func NewBookRESTAPIServer(listenAddr, tokenSecret string, tokenDuration time.Duration, storage storage.Storage) *BookRESTAPIServer {
	return &BookRESTAPIServer{
		listenAddr:    listenAddr,
		storage:       storage,
		tokenSecret:   tokenSecret,
		tokenDuration: tokenDuration,
	}
}

// Run runs the server
func (s *BookRESTAPIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/register", makeHTTPHandler(s.handleRegister)).Methods("POST")
	r.HandleFunc("/login", makeHTTPHandler(s.handleLogin)).Methods("POST")
	r.HandleFunc("/books", s.validateJWT(s.handleGetBooks)).Methods("GET")
	r.HandleFunc("/books", s.validateJWT(s.handlePostBook)).Methods("POST")
	r.HandleFunc("/books/{id}", s.validateJWT(s.handleGetBookByID)).Methods("GET")
	r.HandleFunc("/books/{id}", s.validateJWT(s.handlePutBookByID)).Methods("PUT")
	r.HandleFunc("/books/{id}", s.validateJWT(s.handleDeleteBookByID)).Methods("DELETE")

	log.Println("[BookRESTAPIServer] Server is running on: " + s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		log.Fatal("[BookRESTAPIServer] Error while running server: " + err.Error())
	}
}

func (s *BookRESTAPIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := &types.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return writeJSONResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	user, _ := s.storage.SelectUserByEmail(createAccountRequest.Email)
	if user != nil {
		return writeJSONResponse(w, http.StatusBadRequest, "user with this email already exists")
	}

	hashedPass, err := crypto.HashPassword(createAccountRequest.Password)
	if err != nil {
		return writeJSONResponse(w, http.StatusInternalServerError, fmt.Errorf("error while hashing password: %w", err))
	}

	newUser := types.NewUser(createAccountRequest.Email, hashedPass, createAccountRequest.FirstName, createAccountRequest.LastName, int(createAccountRequest.Age))
	id, err := s.storage.InsertUser(newUser)
	if err != nil {
		return writeJSONResponse(w, http.StatusInternalServerError, fmt.Errorf("error while creating new user: %w", err))
	}

	newUser.ID = id

	return writeJSONResponse(w, http.StatusOK, newUser)
}

func (s *BookRESTAPIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginRequest := &types.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		return writeJSONResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	user, err := s.storage.SelectUserByEmail(loginRequest.Email)
	if err != nil || user == nil {
		return writeJSONResponse(w, http.StatusUnauthorized, errors.New("invalid credentials"))
	}

	if err := crypto.CheckPassword(loginRequest.Password, user.Password); err != nil {
		return writeJSONResponse(w, http.StatusUnauthorized, errors.New("invalid credentials"))
	}

	token, err := token.Generate(user.ID, user.Email, s.tokenSecret, s.tokenDuration)
	if err != nil {
		return writeJSONResponse(w, http.StatusInternalServerError, fmt.Errorf("error while generating token: %w", err))
	}

	return writeJSONResponse(w, http.StatusOK, token)
}

func (s *BookRESTAPIServer) handleGetBooks(w http.ResponseWriter, r *http.Request) error {
	books, err := s.storage.SelectAllBooks()
	if err != nil {
		return writeJSONResponse(w, http.StatusInternalServerError, fmt.Errorf("error while getting books: %w", err))
	}

	return writeJSONResponse(w, http.StatusOK, books)
}

func (s *BookRESTAPIServer) handlePostBook(w http.ResponseWriter, r *http.Request) error {
	createBookRequest := &types.CreateBookRequest{}
	if err := json.NewDecoder(r.Body).Decode(createBookRequest); err != nil {
		return writeJSONResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	newBook := types.NewBook(createBookRequest.Title, createBookRequest.Author)
	id, err := s.storage.InsertBook(newBook)
	if err != nil {
		return writeJSONResponse(w, http.StatusInternalServerError, fmt.Errorf("error while creating new book: %w", err))
	}

	newBook.ID = id

	return writeJSONResponse(w, http.StatusOK, newBook)
}

func (s *BookRESTAPIServer) handleGetBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		return writeJSONResponse(w, http.StatusBadRequest, errors.New("invalid id"))
	}

	book, err := s.storage.SelectBookByID(id)
	if err != nil {
		return writeJSONResponse(w, http.StatusNotFound, errors.New("not found"))
	}

	return writeJSONResponse(w, http.StatusOK, book)
}

func (s *BookRESTAPIServer) handlePutBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		return writeJSONResponse(w, http.StatusBadRequest, errors.New("invalid id"))
	}

	_, err = s.storage.SelectBookByID(id)
	if err != nil {
		return writeJSONResponse(w, http.StatusNotFound, errors.New("not found"))
	}

	book := &types.Book{}
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		return writeJSONResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	if err := s.storage.UpdateBook(book); err != nil {
		return writeJSONResponse(w, http.StatusInternalServerError, errors.New("error while updating the book"))
	}

	return writeJSONResponse(w, http.StatusOK, nil)
}

func (s *BookRESTAPIServer) handleDeleteBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		return writeJSONResponse(w, http.StatusBadRequest, errors.New("invalid id"))
	}

	_, err = s.storage.SelectBookByID(id)
	if err != nil {
		return writeJSONResponse(w, http.StatusNotFound, errors.New("not found"))
	}

	if err := s.storage.DeleteBook(id); err != nil {
		return writeJSONResponse(w, http.StatusInternalServerError, errors.New("error while deleting the book"))
	}

	return writeJSONResponse(w, http.StatusOK, nil)
}

func (s *BookRESTAPIServer) validateJWT(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			if err := token.Validate(r.Header.Get("Token"), s.tokenSecret); err != nil {
				if err := writeJSONResponse(w, http.StatusUnauthorized, "unauthorized: "+err.Error()); err != nil {
					log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
					return
				}

				log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
				return
			}

			if err := f(w, r); err != nil {
				if err := writeJSONResponse(w, http.StatusInternalServerError, err.Error()); err != nil {
					log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
					return
				}

				log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
			}
		} else {
			if err := writeJSONResponse(w, http.StatusUnauthorized, "not authorized"); err != nil {
				log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
				return
			}
		}
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}
