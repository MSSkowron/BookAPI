package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MSSkowron/BookRESTAPI/internal/dtos"
	"github.com/MSSkowron/BookRESTAPI/internal/services"
	"github.com/MSSkowron/BookRESTAPI/pkg/logger"
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
	// ErrMsgUnauthorizedInvalidToken is a message for unauthorized with invalid token
	ErrMsgUnauthorizedInvalidToken = "unauthorized"
	// ErrMsgUnauthorizedInvalidCredentials is a message for unauthorized with invalid credentials
	ErrMsgUnauthorizedInvalidCredentials = "invalid credentials"
	// ErrMsgNotFound is a message for not found
	ErrMsgNotFound = "not found"
	// ErrMsgInternalError is a message for internal error
	ErrMsgInternalError = "internal server error"
	// ContextKeyUserID is a context key for user id
	ContextKeyUserID = contextKey("user_id")
	// DefaultAddress
	DefaultAddress = ":8080"
	// DefaultWriteTimeout
	DefaultWriteTimeout = 15 * time.Second
	// DefaultReadTimeout
	DefaultReadTimeout = 15 * time.Second
)

var (
	// ErrUserIDNotSetInContext is returned when user id is not set in context
	ErrUserIDNotSetInContext = errors.New("user id not set in context")
)

type contextKey string

type serverHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandlerFunc(f serverHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			logger.Errorf("Error (%s) while handling request from client with IP address: %s ", err)
		}
	}
}

// Server is a HTTP server for handling REST API requests
type Server struct {
	*http.Server
	userService  services.UserService
	tokenService services.TokenService
	bookService  services.BookService
}

// NewServer creates a new Server
func NewServer(userService services.UserService, bookService services.BookService, tokenService services.TokenService, opts ...ServerOpt) *Server {
	server := &Server{
		Server: &http.Server{
			Addr:         DefaultAddress,
			WriteTimeout: DefaultWriteTimeout,
			ReadTimeout:  DefaultReadTimeout,
		},
		userService:  userService,
		tokenService: tokenService,
		bookService:  bookService,
	}

	for _, opt := range opts {
		opt(server)
	}

	server.initRoutes()

	return server
}

type ServerOpt func(*Server)

func WithAddress(addr string) func(*Server) {
	return func(s *Server) {
		s.Addr = addr
	}
}

func WithReadTimeout(timeout time.Duration) func(*Server) {
	return func(s *Server) {
		s.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) func(*Server) {
	return func(s *Server) {
		s.WriteTimeout = timeout
	}
}

// Run runs the Server
func (s *Server) initRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/register", makeHTTPHandlerFunc(s.handleRegister)).Methods("POST")
	r.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin)).Methods("POST")

	bookRouter := r.PathPrefix("/books").Subrouter()
	bookRouter.Use(s.validateJWT)
	bookRouter.HandleFunc("", makeHTTPHandlerFunc(s.handleGetBooks)).Methods("GET")
	bookRouter.HandleFunc("", makeHTTPHandlerFunc(s.handlePostBook)).Methods("POST")
	bookRouter.HandleFunc("/{id}", makeHTTPHandlerFunc(s.handleGetBookByID)).Methods("GET")
	bookRouter.HandleFunc("/{id}", makeHTTPHandlerFunc(s.handlePutBookByID)).Methods("PUT")
	bookRouter.HandleFunc("/{id}", makeHTTPHandlerFunc(s.handleDeleteBookByID)).Methods("DELETE")

	s.Handler = r
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) error {
	logger.Infof("Received POST /register from %s", r.RemoteAddr)

	accountCreateDTO := &dtos.AccountCreateDTO{}
	if err := json.NewDecoder(r.Body).Decode(accountCreateDTO); err != nil {
		s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	userDTO, err := s.userService.RegisterUser(accountCreateDTO)
	if err != nil {
		if errors.Is(err, services.ErrInvalidEmail) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrInvalidPassword) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrInvalidFirstName) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrInvalidLastName) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrInvalidAge) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrUserAlreadyExists) {
			s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestUserAlreadyExists)
			return nil
		}

		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("register user: %w", err)
	}

	s.respondWithJSON(w, http.StatusOK, userDTO)

	return nil
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	logger.Infof("Received POST /login from %s", r.RemoteAddr)

	userLoginDTO := &dtos.UserLoginDTO{}
	if err := json.NewDecoder(r.Body).Decode(userLoginDTO); err != nil {
		s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	tokenDTO, err := s.userService.LoginUser(userLoginDTO)
	if err != nil {
		if errors.Is(err, services.ErrInvalidEmail) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrEmptyPassword) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrInvalidCredentials) {
			s.respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorizedInvalidCredentials)
			return nil
		}

		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("login user: %w", err)
	}

	s.respondWithJSON(w, http.StatusOK, tokenDTO)

	return nil
}

func (s *Server) handleGetBooks(w http.ResponseWriter, r *http.Request) error {
	logger.Infof("Received GET /books from %s", r.RemoteAddr)

	booksDTO, err := s.bookService.GetBooks()
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("get books: %w", err)
	}

	s.respondWithJSON(w, http.StatusOK, booksDTO)

	return nil
}

func (s *Server) handlePostBook(w http.ResponseWriter, r *http.Request) error {
	logger.Infof("Received POST /books from %s", r.RemoteAddr)

	bookCreateDTO := &dtos.BookCreateDTO{}
	if err := json.NewDecoder(r.Body).Decode(bookCreateDTO); err != nil {
		s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	userID := r.Context().Value(ContextKeyUserID).(int)
	if userID == 0 {
		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return ErrUserIDNotSetInContext
	}

	bookDTO, err := s.bookService.AddBook(userID, bookCreateDTO)
	if err != nil {
		if errors.Is(err, services.ErrInvalidAuthor) {
			s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
			return nil
		}
		if errors.Is(err, services.ErrInvalidTitle) {
			s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
			return nil
		}

		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("add book: %w", err)
	}

	s.respondWithJSON(w, http.StatusOK, bookDTO)

	return nil
}

func (s *Server) handleGetBookByID(w http.ResponseWriter, r *http.Request) error {
	logger.Infof("Received GET /books/{id} from %s", r.RemoteAddr)

	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidBookID)
		return nil
	}

	bookDTO, err := s.bookService.GetBook(id)
	if err != nil {
		if errors.Is(err, services.ErrInvalidID) {
			s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidBookID)
			return nil
		}
		if errors.Is(err, services.ErrBookNotFound) {
			s.respondWithError(w, http.StatusNotFound, ErrMsgNotFound)
			return nil
		}

		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("get book: %w", err)
	}

	s.respondWithJSON(w, http.StatusOK, bookDTO)

	return nil
}

func (s *Server) handlePutBookByID(w http.ResponseWriter, r *http.Request) error {
	logger.Infof("Received PUT /books/{id} from %s", r.RemoteAddr)

	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidBookID)
		return nil
	}

	bookDTO := &dtos.BookDTO{}
	if err := json.NewDecoder(r.Body).Decode(bookDTO); err != nil {
		s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidRequestBody)
		return nil
	}

	updatedBookDTO, err := s.bookService.UpdateBook(id, bookDTO)
	if err != nil {
		if errors.Is(err, services.ErrInvalidID) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrInvalidAuthor) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrInvalidTitle) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrBookNotFound) {
			s.respondWithError(w, http.StatusNotFound, ErrMsgNotFound)
			return nil
		}

		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("update book: %w", err)
	}

	s.respondWithJSON(w, http.StatusOK, updatedBookDTO)

	return nil
}

func (s *Server) handleDeleteBookByID(w http.ResponseWriter, r *http.Request) error {
	logger.Infof("Received DELETE /books{id} from %s", r.RemoteAddr)

	idString := mux.Vars(r)["id"]
	defer r.Body.Close()

	id, err := strconv.Atoi(idString)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, ErrMsgBadRequestInvalidBookID)
		return nil
	}

	if err := s.bookService.DeleteBook(id); err != nil {
		if errors.Is(err, services.ErrInvalidID) {
			s.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s:%s", ErrMsgBadRequestInvalidRequestBody, err))
			return nil
		}
		if errors.Is(err, services.ErrBookNotFound) {
			s.respondWithError(w, http.StatusNotFound, ErrMsgNotFound)
			return nil
		}

		s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
		return fmt.Errorf("delete book: %w", err)
	}

	s.respondWithJSON(w, http.StatusOK, nil)
	return nil
}

func (s *Server) validateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		logger.Infof("Validating JWT for client with IP address: %s", clientIP)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Infof("Authorization header missing for client with IP address: %s", clientIP)
			s.respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorized)
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			logger.Infof("Invalid authorization format for client with IP address: %s", clientIP)
			s.respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorized)
			return
		}

		tokenString := authHeaderParts[1]
		if err := s.tokenService.ValidateToken(tokenString); err != nil {
			if errors.Is(err, services.ErrExpiredToken) {
				logger.Infof("Expired JWT detected for client with IP address: %s", clientIP)
				s.respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorizedExpiredToken)
				return
			}
			if errors.Is(err, services.ErrInvalidToken) {
				logger.Infof("Invalid JWT detected for client with IP address: %s", clientIP)
				s.respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorizedInvalidToken)
				return
			}

			logger.Errorf("Error (%s) encountered during JWT validation for client with IP address: %s", err, clientIP)
			s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
			return
		}

		logger.Infof("JWT validation successful for client with IP address: %s", clientIP)

		userID, err := s.tokenService.GetUserIDFromToken(tokenString)
		if err != nil {
			logger.Errorf("Error (%s) encountered while retrieving user ID from JWT for client with IP address: %s", err, clientIP)

			if errors.Is(err, services.ErrInvalidToken) {
				s.respondWithError(w, http.StatusUnauthorized, ErrMsgUnauthorizedInvalidToken)
				return
			}

			s.respondWithError(w, http.StatusInternalServerError, ErrMsgInternalError)
			return
		}

		logger.Infof("User ID (%d) retrieved from JWT for client with IP address: %s", userID, clientIP)

		ctx := context.WithValue(r.Context(), ContextKeyUserID, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) respondWithError(w http.ResponseWriter, errCode int, errMessage string) {
	s.respondWithJSON(w, errCode, dtos.ErrorDTO{Error: errMessage})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("Error (%s) while marshalling JSON response", err)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(ErrMsgInternalError))

		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(response)
}
