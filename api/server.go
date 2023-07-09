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
	"github.com/MSSkowron/BookRESTAPI/model"
	"github.com/MSSkowron/BookRESTAPI/storage"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if err := writeJSONResponse(w, http.StatusInternalServerError, err); err != nil {
				log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
				return
			}

			log.Printf("[BookRESTAPIServer] Error: %s", err.Error())
		}
	}
}

type BookRESTAPIServer struct {
	listenAddr string
	storage    storage.Storage
}

func NewBookRESTAPIServer(listenAddr string, storage storage.Storage) *BookRESTAPIServer {
	return &BookRESTAPIServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *BookRESTAPIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/register", makeHTTPHandler(s.handleRegister)).Methods("POST")
	r.HandleFunc("/login", makeHTTPHandler(s.handleLogin)).Methods("POST")
	r.HandleFunc("/books", validateJWT(s.handleGetBooks)).Methods("GET")
	r.HandleFunc("/books", validateJWT(s.handlePostBook)).Methods("POST")
	r.HandleFunc("/books/{id}", validateJWT(s.handleGetBookByID)).Methods("GET")
	r.HandleFunc("/books/{id}", validateJWT(s.handlePutBookByID)).Methods("PUT")
	r.HandleFunc("/books/{id}", validateJWT(s.handleDeleteBookByID)).Methods("DELETE")

	log.Println("[BookRESTAPIServer] Server is running on: " + s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		log.Fatal("[BookRESTAPIServer] Error while running server: " + err.Error())
	}
}

func (s *BookRESTAPIServer) handleRegister(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := &model.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return errors.New("invalid request body")
	}

	hashedPass, err := crypto.HashPassword(createAccountRequest.Password)
	if err != nil {
		return errors.New("error while creating new user")
	}

	newUser := model.NewUser(createAccountRequest.Email, hashedPass, createAccountRequest.FirstName, createAccountRequest.LastName, int(createAccountRequest.Age))
	id, err := s.storage.InsertUser(newUser)
	if err != nil {
		return errors.New("error while creating new user")
	}

	newUser.ID = id

	return writeJSONResponse(w, http.StatusOK, newUser)
}

func (s *BookRESTAPIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	loginRequest := &model.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		return errors.New("invalid request body")
	}

	user, err := s.storage.SelectUserByEmail(loginRequest.Email)
	if err != nil || user == nil {
		return errors.New("invalid credentials")
	}

	if err := crypto.CheckPassword(loginRequest.Password, user.Password); err != nil {
		return errors.New("invalid credentials")
	}

	token, err := generateToken(user.Email)
	if err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, token)
}

func (s *BookRESTAPIServer) handleGetBooks(w http.ResponseWriter, r *http.Request) error {
	books, err := s.storage.SelectAllBooks()
	if err != nil {
		return err
	}

	return writeJSONResponse(w, http.StatusOK, books)
}

func (s *BookRESTAPIServer) handlePostBook(w http.ResponseWriter, r *http.Request) error {
	createBookRequest := &model.CreateBookRequest{}
	if err := json.NewDecoder(r.Body).Decode(createBookRequest); err != nil {
		return errors.New("invalid request body")
	}

	newBook := model.NewBook(createBookRequest.Title, createBookRequest.Author)
	id, err := s.storage.InsertBook(newBook)
	if err != nil {
		return errors.New("error while creating new book")
	}

	newBook.ID = id

	return writeJSONResponse(w, http.StatusOK, newBook)
}

func (s *BookRESTAPIServer) handleGetBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		return errors.New("invalid id")
	}

	book, err := s.storage.SelectBookByID(id)
	if err != nil {
		return errors.New("not found")
	}

	return writeJSONResponse(w, http.StatusOK, book)
}

func (s *BookRESTAPIServer) handlePutBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		return errors.New("invalid id")
	}

	_, err = s.storage.SelectBookByID(id)
	if err != nil {
		return errors.New("not found")
	}

	book := &model.Book{}
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		return errors.New("invalid request body")
	}

	if err := s.storage.UpdateBook(book); err != nil {
		return errors.New("error while deleting the book")
	}

	return writeJSONResponse(w, http.StatusOK, nil)
}

func (s *BookRESTAPIServer) handleDeleteBookByID(w http.ResponseWriter, r *http.Request) error {
	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		return errors.New("invalid id")
	}

	_, err = s.storage.SelectBookByID(id)
	if err != nil {
		return errors.New("not found")
	}

	if err := s.storage.DeleteBook(id); err != nil {
		return errors.New("error while deleting the book")
	}

	return writeJSONResponse(w, http.StatusOK, nil)
}

var SECRET = []byte("super-secret-auth-key")

func generateToken(email string) (tokenString string, err error) {
	claims := &jwt.MapClaims{
		"email":     email,
		"expiresAt": time.Now().Add(10 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SECRET)
}

func validateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return SECRET, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token is not valid")
	}

	if int64(token.Claims.(jwt.MapClaims)["expiresAt"].(float64)) < time.Now().Local().Unix() {
		return errors.New("token expired")
	}

	return nil
}

func validateJWT(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			if err := validateToken(r.Header.Get("Token")); err != nil {
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
