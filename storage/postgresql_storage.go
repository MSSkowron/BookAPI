package storage

import (
	"database/sql"
	"log"

	"github.com/MSSkowron/BookRESTAPI/types"
	_ "github.com/lib/pq" // postgres driver
)

const (
	driverName = "postgres"
)

// PostgresSQLStorage is a storage for PostgreSQL
type PostgresSQLStorage struct {
	db *sql.DB
}

// NewPostgresSQLStorage creates a new PostgresSQLStorage
func NewPostgresSQLStorage(connectionString string) (*PostgresSQLStorage, error) {
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	postgresSQLStorage := &PostgresSQLStorage{
		db: db,
	}

	return postgresSQLStorage, nil
}

// InsertUser inserts a new user
func (s *PostgresSQLStorage) InsertUser(user *types.User) (int, error) {
	var (
		query string = "INSERT INTO users (email, password, first_name, last_name, age) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		id    int    = -1
	)

	err := s.db.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, user.Age).Scan(&id)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while inserting new user: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresSQLStorage] Inserted new user with ID: %d", id)

	return id, nil
}

// SelectUserByEmail selects a user by email
func (s *PostgresSQLStorage) SelectUserByEmail(email string) (*types.User, error) {
	query := "SELECT * FROM users WHERE email=$1"

	row := s.db.QueryRow(query, email)
	user := &types.User{}
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while selecting user with email %s: %s", email, err.Error())
		return nil, err
	}

	log.Printf("[PostgresSQLStorage] Selected user with email: %s", email)

	return user, nil
}

// InsertBook inserts a new book
func (s *PostgresSQLStorage) InsertBook(book *types.Book) (int, error) {
	var (
		query string = "INSERT INTO books (author, title) VALUES ($1, $2) RETURNING id"
		id    int    = -1
	)

	err := s.db.QueryRow(query, book.Author, book.Title).Scan(&id)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while inserting new book: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresSQLStorage] Inserted new book with ID %d", id)

	return id, nil
}

// SelectAllBooks selects all books
func (s *PostgresSQLStorage) SelectAllBooks() ([]*types.Book, error) {
	query := "SELECT * FROM books"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*types.Book{}
	for rows.Next() {
		book := &types.Book{}
		if err := rows.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
			log.Printf("[PostgresSQLStorage] Error while selecting all books: %s", err.Error())
			return nil, err
		}

		books = append(books, book)
	}

	log.Println("[PostgresSQLStorage] Selected all books")

	return books, nil
}

// SelectBookByID selects a book by ID
func (s *PostgresSQLStorage) SelectBookByID(id int) (*types.Book, error) {
	query := "SELECT * FROM books WHERE id=$1"

	row := s.db.QueryRow(query, id)
	book := &types.Book{}
	if err := row.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
		log.Printf("[PostgresSQLStorage] Error while selecting book with ID %d: %s", id, err.Error())
		return nil, err
	}

	log.Printf("[PostgresSQLStorage] Selected book with ID: %d", id)

	return book, nil
}

// DeleteBook deletes a book by ID
func (s *PostgresSQLStorage) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id=$1"

	if _, err := s.db.Exec(query, id); err != nil {
		log.Printf("[PostgresSQLStorage] Error while deleting book with ID %d: %s", id, err.Error())
		return err
	}

	log.Printf("[PostgresSQLStorage] Deleted book with ID %d", id)

	return nil
}

// UpdateBook updates a book
func (s *PostgresSQLStorage) UpdateBook(book *types.Book) error {
	query := "UPDATE books SET author = $1, title = $2 WHERE id = $3"

	_, err := s.db.Exec(query, book.Author, book.Title, book.ID)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while updating book with ID %d: %s", book.ID, err.Error())
		return err
	}

	log.Printf("[PostgresSQLStorage] Updated book with ID %d", book.ID)

	return nil
}
