package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MSSkowron/GoBankAPI/model"
	"github.com/MSSkowron/GoBankAPI/storage"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type GoBookAPIServer struct {
	listenAddr string
	storage    storage.Storage
}

func NewGoBookAPIServer(listenAddr string, storage storage.Storage) *GoBookAPIServer {
	return &GoBookAPIServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *GoBookAPIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/users/register", s.handlePostUserRegister).Methods("POST")
	r.HandleFunc("/users/login", s.handlePostUserLogin).Methods("POST")
	r.HandleFunc("/books", validateJWT(s.handleGetBooks)).Methods("GET")
	r.HandleFunc("/books", validateJWT(s.handlePostBook)).Methods("POST")
	r.HandleFunc("/books/{id}", validateJWT(s.handleGetBookByID)).Methods("GET")
	r.HandleFunc("/books/{id}", validateJWT(s.handlePutBookByID)).Methods("PUT")
	r.HandleFunc("/books/{id}", validateJWT(s.handleDeleteBookByID)).Methods("DELETE")

	log.Println("[GoBookAPIServer] Server is running on: " + s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		log.Fatal("[GoBookAPIServer] Error while running server: " + err.Error())
	}
}

func (s *GoBookAPIServer) handlePostUserRegister(w http.ResponseWriter, r *http.Request) {
	createAccountRequest := &model.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createAccountRequest); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := s.storage.CreateUser(model.NewUser(createAccountRequest.Email, createAccountRequest.Password, createAccountRequest.FirstName, createAccountRequest.LastName, int(createAccountRequest.Age))); err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, "error while creating new user")
	}

	writeJSONResponse(w, http.StatusOK, "registered successfully")
}

func (s *GoBookAPIServer) handlePostUserLogin(w http.ResponseWriter, r *http.Request) {
	loginRequest := &model.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := s.storage.GetUserByEmail(loginRequest.Email)
	if err != nil || user == nil {
		writeJSONResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if user.Password != loginRequest.Password {
		writeJSONResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := generateToken(user.Email)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	writeJSONResponse(w, http.StatusOK, token)
}

func (s *GoBookAPIServer) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := s.storage.GetBooks()
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, nil)
		return
	}

	writeJSONResponse(w, http.StatusNotImplemented, books)
}

func (s *GoBookAPIServer) handlePostBook(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusNotImplemented, "TODO")
}

func (s *GoBookAPIServer) handleGetBookByID(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusNotImplemented, "TODO")
}

func (s *GoBookAPIServer) handlePutBookByID(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusNotImplemented, "TODO")
}

func (s *GoBookAPIServer) handleDeleteBookByID(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusNotImplemented, "TODO")
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

func validateJWT(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			err := validateToken(r.Header.Get("Token"))
			if err != nil {
				writeJSONResponse(w, http.StatusUnauthorized, "not authorized: "+err.Error())
				return
			}

			handlerFunc(w, r)

		} else {
			writeJSONResponse(w, http.StatusUnauthorized, "not authorized")
		}
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}
