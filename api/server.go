package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MSSkowron/BookRESTAPI/crypto"
	"github.com/MSSkowron/BookRESTAPI/model"
	"github.com/MSSkowron/BookRESTAPI/storage"
	"github.com/MSSkowron/BookRESTAPI/token"
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
	r.HandleFunc("/books", validateJWT(s.handleGetBooks, s.tokenSecret)).Methods("GET")
	r.HandleFunc("/books", validateJWT(s.handlePostBook, s.tokenSecret)).Methods("POST")
	r.HandleFunc("/books/{id}", validateJWT(s.handleGetBookByID, s.tokenSecret)).Methods("GET")
	r.HandleFunc("/books/{id}", validateJWT(s.handlePutBookByID, s.tokenSecret)).Methods("PUT")
	r.HandleFunc("/books/{id}", validateJWT(s.handleDeleteBookByID, s.tokenSecret)).Methods("DELETE")

	log.Println("[BookRESTAPIServer] Server is running on: " + s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		log.Fatal("[BookRESTAPIServer] Error while running server: " + err.Error())
	}
}

func (s *BookRESTAPIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := &model.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return nil
	}

	if createAccountRequest.Email == "" || createAccountRequest.Password == "" || createAccountRequest.FirstName == "" || createAccountRequest.LastName == "" || createAccountRequest.Age == 0 {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return nil
	}

	user, _ := s.storage.SelectUserByEmail(createAccountRequest.Email)
	if user != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "user already exists"})
		return nil
	}

	hashedPass, err := crypto.HashPassword(createAccountRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error while hashing password"})
		return fmt.Errorf("error while hashing password: %w", err)
	}

	newUser := model.NewUser(createAccountRequest.Email, hashedPass, createAccountRequest.FirstName, createAccountRequest.LastName, int(createAccountRequest.Age))
	id, err := s.storage.InsertUser(newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error while inserting user"})
		return fmt.Errorf("error while inserting user: %w", err)
	}

	newUser.ID = id

	respondWithJSON(w, http.StatusOK, newUser)

	return nil
}

func (s *BookRESTAPIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginRequest := &model.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return nil
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return nil
	}

	user, err := s.storage.SelectUserByEmail(loginRequest.Email)
	if err != nil || user == nil {
		respondWithError(w, http.StatusUnauthorized, model.ErrorResponse{Error: "invalid credentials"})
		return nil
	}

	if err := crypto.CheckPassword(loginRequest.Password, user.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, model.ErrorResponse{Error: "invalid credentials"})
		return nil
	}

	token, err := token.Generate(user.ID, user.Email, s.tokenSecret, s.tokenDuration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error while generating token"})
		return fmt.Errorf("error while generating token: %w", err)
	}

	respondWithJSON(w, http.StatusOK, model.LoginResponse{Token: token})

	return nil
}

func (s *BookRESTAPIServer) handleGetBooks(w http.ResponseWriter, r *http.Request) error {
	books, err := s.storage.SelectAllBooks()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error while getting books"})
		return fmt.Errorf("error while getting books: %w", err)
	}

	respondWithJSON(w, http.StatusOK, books)

	return nil
}

func (s *BookRESTAPIServer) handlePostBook(w http.ResponseWriter, r *http.Request) error {
	createBookRequest := &model.CreateBookRequest{}
	if err := json.NewDecoder(r.Body).Decode(createBookRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return nil
	}

	if createBookRequest.Title == "" || createBookRequest.Author == "" {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return nil
	}

	newBook := model.NewBook(createBookRequest.Title, createBookRequest.Author)
	id, err := s.storage.InsertBook(newBook)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error while creating new book"})
		return fmt.Errorf("error while creating new book: %w", err)
	}

	newBook.ID = id

	respondWithJSON(w, http.StatusOK, newBook)

	return nil
}

func (s *BookRESTAPIServer) handleGetBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return nil
	}

	book, err := s.storage.SelectBookByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, model.ErrorResponse{Error: "not found"})
		return nil
	}

	respondWithJSON(w, http.StatusOK, book)

	return nil
}

func (s *BookRESTAPIServer) handlePutBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return nil
	}

	_, err = s.storage.SelectBookByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, model.ErrorResponse{Error: "not found"})
		return nil
	}

	book := &model.Book{}
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid request body"})
		return nil
	}

	if err := s.storage.UpdateBook(book); err != nil {
		respondWithError(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error while updating the book"})
		return fmt.Errorf("error while updating the book: %s", err)
	}

	respondWithJSON(w, http.StatusOK, book)

	return nil
}

func (s *BookRESTAPIServer) handleDeleteBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return nil
	}

	_, err = s.storage.SelectBookByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, model.ErrorResponse{Error: "not found"})
		return nil
	}

	if err := s.storage.DeleteBook(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, model.ErrorResponse{Error: "error while deleting the book"})
		return fmt.Errorf("error while deleting the book: %s", err.Error())
	}

	respondWithJSON(w, http.StatusOK, nil)

	return nil
}

func validateJWT(f apiFunc, tokenSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			if err := token.Validate(r.Header.Get("Token"), tokenSecret); err != nil {
				respondWithError(w, http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
				return
			}

			if err := f(w, r); err != nil {
				log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
			}
		} else {
			respondWithError(w, http.StatusUnauthorized, model.ErrorResponse{Error: "unauthorized"})
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, errorResponse model.ErrorResponse) {
	respondWithJSON(w, code, errorResponse)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// nolint:errcheck
	_, _ = w.Write(response)
}
