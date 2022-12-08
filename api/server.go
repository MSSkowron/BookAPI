package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	model "github.com/MSSkowron/GoBankAPI/model"
	jwt "github.com/dgrijalva/jwt-go"
	mux "github.com/gorilla/mux"
)

type GoBookAPIServer struct {
	listenAddr string
}

func NewGoBookAPIServer(listenAddr string) *GoBookAPIServer {
	return &GoBookAPIServer{
		listenAddr: listenAddr,
	}
}

func (s *GoBookAPIServer) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.handlePostUser).Methods("POST")
	r.HandleFunc("/books", validateJWT(s.handleGetBooks)).Methods("GET")
	r.HandleFunc("/books", validateJWT(s.handlePostBook)).Methods("POST")
	r.HandleFunc("/books/{id}", validateJWT(s.handleGetBookByID)).Methods("GET")
	r.HandleFunc("/books/{id}", validateJWT(s.handlePutBookByID)).Methods("PUT")
	r.HandleFunc("/books/{id}", validateJWT(s.handleDeleteBookByID)).Methods("DELETE")

	log.Println("[GoBookAPI] Server is running on: " + s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		log.Fatal("[GoBookAPI] Error while running server: " + err.Error())
	}
}

func (s *GoBookAPIServer) handlePostUser(w http.ResponseWriter, r *http.Request) {
	user := model.NewUser("test@gmail.com", "test", "test", 18)

	token, err := createJWT(user)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, nil)
	}

	writeJSONResponse(w, http.StatusNotImplemented, token)
}

func (s *GoBookAPIServer) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusNotImplemented, "TODO")
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

func createJWT(user *model.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": time.Now().Add(time.Hour).Unix(),
		"userID":    user.ID,
		"userEmail": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(SECRET)
}

func validateJWT(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					writeJSONResponse(w, http.StatusUnauthorized, "not authorized")
				}

				return SECRET, nil
			})

			if err != nil {
				writeJSONResponse(w, http.StatusUnauthorized, "not authorized: "+err.Error())
			}

			if token.Valid {
				handlerFunc(w, r)
			}
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
