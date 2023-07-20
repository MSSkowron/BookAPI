package storage

import (
	"context"
	"log"

	"github.com/MSSkowron/BookRESTAPI/model"
	pgx "github.com/jackc/pgx/v5"
)

// PostgresStorage is a storage for PostgreSQL
type PostgresStorage struct {
	conn *pgx.Conn
}

// NewPostgresStorage creates a new PostgresStorage
func NewPostgresStorage(connectionString string) (*PostgresStorage, error) {
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	PostgresStorage := &PostgresStorage{
		conn: conn,
	}

	return PostgresStorage, nil
}

// InsertUser inserts a new user
func (s *PostgresStorage) InsertUser(user *model.User) (int, error) {
	var (
		query string = "INSERT INTO users (email, password, first_name, last_name, age) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		id    int    = -1
	)

	err := s.conn.QueryRow(context.Background(), query, user.Email, user.Password, user.FirstName, user.LastName, user.Age).Scan(&id)
	if err != nil {
		log.Printf("[PostgresStorage] Error while inserting new user: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresStorage] Inserted new user with ID: %d", id)

	return id, nil
}

// SelectUserByEmail selects a user with given email
func (s *PostgresStorage) SelectUserByEmail(email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email=$1"

	row := s.conn.QueryRow(context.Background(), query, email)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age)
	if err != nil {
		log.Printf("[PostgresStorage] Error while selecting user with email %s: %s", email, err.Error())
		return nil, err
	}

	log.Printf("[PostgresStorage] Selected user with email: %s", email)

	return user, nil
}

// InsertBook inserts a new book
func (s *PostgresStorage) InsertBook(book *model.Book) (int, error) {
	var (
		query string = "INSERT INTO books (author, title) VALUES ($1, $2) RETURNING id"
		id    int    = -1
	)

	err := s.conn.QueryRow(context.Background(), query, book.Author, book.Title).Scan(&id)
	if err != nil {
		log.Printf("[PostgresStorage] Error while inserting new book: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresStorage] Inserted new book with ID %d", id)

	return id, nil
}

// SelectAllBooks selects all books
func (s *PostgresStorage) SelectAllBooks() ([]*model.Book, error) {
	query := "SELECT * FROM books"

	rows, err := s.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*model.Book{}
	for rows.Next() {
		book := &model.Book{}
		if err := rows.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
			log.Printf("[PostgresStorage] Error while selecting all books: %s", err.Error())
			return nil, err
		}

		books = append(books, book)
	}

	log.Println("[PostgresStorage] Selected all books")

	return books, nil
}

// SelectBookByID selects a book with given ID
func (s *PostgresStorage) SelectBookByID(id int) (*model.Book, error) {
	query := "SELECT * FROM books WHERE id=$1"

	row := s.conn.QueryRow(context.Background(), query, id)
	book := &model.Book{}
	if err := row.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
		log.Printf("[PostgresStorage] Error while selecting book with ID %d: %s", id, err.Error())
		return nil, err
	}

	log.Printf("[PostgresStorage] Selected book with ID: %d", id)

	return book, nil
}

// DeleteBook deletes a book with given ID
func (s *PostgresStorage) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id=$1"

	if _, err := s.conn.Exec(context.Background(), query, id); err != nil {
		log.Printf("[PostgresStorage] Error while deleting book with ID %d: %s", id, err.Error())
		return err
	}

	log.Printf("[PostgresStorage] Deleted book with ID %d", id)

	return nil
}

// UpdateBook updates a book with given ID
func (s *PostgresStorage) UpdateBook(book *model.Book) error {
	query := "UPDATE books SET author = $1, title = $2 WHERE id = $3"

	_, err := s.conn.Exec(context.Background(), query, book.Author, book.Title, book.ID)
	if err != nil {
		log.Printf("[PostgresStorage] Error while updating book with ID %d: %s", book.ID, err.Error())
		return err
	}

	log.Printf("[PostgresStorage] Updated book with ID %d", book.ID)

	return nil
}
