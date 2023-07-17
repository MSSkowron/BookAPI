package storage

import (
	"context"
	"log"

	"github.com/MSSkowron/BookRESTAPI/model"
	pgx "github.com/jackc/pgx/v5"
)

// PostgresSQLStorage is a storage for PostgreSQL
type PostgresSQLStorage struct {
	conn *pgx.Conn
}

// NewPostgresSQLStorage creates a new PostgresSQLStorage
func NewPostgresSQLStorage(connectionString string) (*PostgresSQLStorage, error) {
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	postgresSQLStorage := &PostgresSQLStorage{
		conn: conn,
	}

	return postgresSQLStorage, nil
}

// InsertUser inserts a new user
func (s *PostgresSQLStorage) InsertUser(user *model.User) (int, error) {
	var (
		query string = "INSERT INTO users (email, password, first_name, last_name, age) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		id    int    = -1
	)

	err := s.conn.QueryRow(context.Background(), query, user.Email, user.Password, user.FirstName, user.LastName, user.Age).Scan(&id)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while inserting new user: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresSQLStorage] Inserted new user with ID: %d", id)

	return id, nil
}

// SelectUserByEmail selects a user with given email
func (s *PostgresSQLStorage) SelectUserByEmail(email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email=$1"

	row := s.conn.QueryRow(context.Background(), query, email)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while selecting user with email %s: %s", email, err.Error())
		return nil, err
	}

	log.Printf("[PostgresSQLStorage] Selected user with email: %s", email)

	return user, nil
}

// InsertBook inserts a new book
func (s *PostgresSQLStorage) InsertBook(book *model.Book) (int, error) {
	var (
		query string = "INSERT INTO books (author, title) VALUES ($1, $2) RETURNING id"
		id    int    = -1
	)

	err := s.conn.QueryRow(context.Background(), query, book.Author, book.Title).Scan(&id)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while inserting new book: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresSQLStorage] Inserted new book with ID %d", id)

	return id, nil
}

// SelectAllBooks selects all books
func (s *PostgresSQLStorage) SelectAllBooks() ([]*model.Book, error) {
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
			log.Printf("[PostgresSQLStorage] Error while selecting all books: %s", err.Error())
			return nil, err
		}

		books = append(books, book)
	}

	log.Println("[PostgresSQLStorage] Selected all books")

	return books, nil
}

// SelectBookByID selects a book with given ID
func (s *PostgresSQLStorage) SelectBookByID(id int) (*model.Book, error) {
	query := "SELECT * FROM books WHERE id=$1"

	row := s.conn.QueryRow(context.Background(), query, id)
	book := &model.Book{}
	if err := row.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
		log.Printf("[PostgresSQLStorage] Error while selecting book with ID %d: %s", id, err.Error())
		return nil, err
	}

	log.Printf("[PostgresSQLStorage] Selected book with ID: %d", id)

	return book, nil
}

// DeleteBook deletes a book with given ID
func (s *PostgresSQLStorage) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id=$1"

	if _, err := s.conn.Exec(context.Background(), query, id); err != nil {
		log.Printf("[PostgresSQLStorage] Error while deleting book with ID %d: %s", id, err.Error())
		return err
	}

	log.Printf("[PostgresSQLStorage] Deleted book with ID %d", id)

	return nil
}

// UpdateBook updates a book with given ID
func (s *PostgresSQLStorage) UpdateBook(book *model.Book) error {
	query := "UPDATE books SET author = $1, title = $2 WHERE id = $3"

	_, err := s.conn.Exec(context.Background(), query, book.Author, book.Title, book.ID)
	if err != nil {
		log.Printf("[PostgresSQLStorage] Error while updating book with ID %d: %s", book.ID, err.Error())
		return err
	}

	log.Printf("[PostgresSQLStorage] Updated book with ID %d", book.ID)

	return nil
}
