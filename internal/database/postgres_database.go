package database

import (
	"context"
	"log"

	"github.com/MSSkowron/BookRESTAPI/internal/models"
	pgx "github.com/jackc/pgx/v5"
)

// PostgresqlDatabase is a Postgresql implementation of Database interface
type PostgresqlDatabase struct {
	conn *pgx.Conn
}

// NewPostgresqlDatabase creates a new PostgresqlDatabase
func NewPostgresqlDatabase(connectionString string) (Database, error) {
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	PostgresqlDatabase := &PostgresqlDatabase{
		conn: conn,
	}

	return PostgresqlDatabase, nil
}

// InsertUser inserts a new user
func (db *PostgresqlDatabase) InsertUser(user *models.User) (int, error) {
	var (
		query string = "INSERT INTO users (email, password, first_name, last_name, age) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		id    int    = -1
	)

	err := db.conn.QueryRow(context.Background(), query, user.Email, user.Password, user.FirstName, user.LastName, user.Age).Scan(&id)
	if err != nil {
		log.Printf("[PostgresqlDatabase] Error while inserting new user: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresqlDatabase] Inserted new user with ID: %d", id)

	return id, nil
}

// SelectUserByEmail selects a user with given email
func (db *PostgresqlDatabase) SelectUserByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email=$1"

	row := db.conn.QueryRow(context.Background(), query, email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age)
	if err != nil {
		log.Printf("[PostgresqlDatabase] Error while selecting user with email %s: %s", email, err.Error())
		return nil, err
	}

	log.Printf("[PostgresqlDatabase] Selected user with email: %s", email)

	return user, nil
}

// InsertBook inserts a new book
func (db *PostgresqlDatabase) InsertBook(book *models.Book) (int, error) {
	var (
		query string = "INSERT INTO books (author, title) VALUES ($1, $2) RETURNING id"
		id    int    = -1
	)

	err := db.conn.QueryRow(context.Background(), query, book.Author, book.Title).Scan(&id)
	if err != nil {
		log.Printf("[PostgresqlDatabase] Error while inserting new book: %s", err.Error())
		return id, err
	}

	log.Printf("[PostgresqlDatabase] Inserted new book with ID %d", id)

	return id, nil
}

// SelectAllBooks selects all books
func (db *PostgresqlDatabase) SelectAllBooks() ([]*models.Book, error) {
	query := "SELECT * FROM books"

	rows, err := db.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*models.Book{}
	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
			log.Printf("[PostgresqlDatabase] Error while selecting all books: %s", err.Error())
			return nil, err
		}

		books = append(books, book)
	}

	log.Println("[PostgresqlDatabase] Selected all books")

	return books, nil
}

// SelectBookByID selects a book with given ID
func (db *PostgresqlDatabase) SelectBookByID(id int) (*models.Book, error) {
	query := "SELECT * FROM books WHERE id=$1"

	row := db.conn.QueryRow(context.Background(), query, id)
	book := &models.Book{}
	if err := row.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
		log.Printf("[PostgresqlDatabase] Error while selecting book with ID %d: %s", id, err.Error())
		return nil, err
	}

	log.Printf("[PostgresqlDatabase] Selected book with ID: %d", id)

	return book, nil
}

// DeleteBook deletes a book with given ID
func (db *PostgresqlDatabase) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id=$1"

	if _, err := db.conn.Exec(context.Background(), query, id); err != nil {
		log.Printf("[PostgresqlDatabase] Error while deleting book with ID %d: %s", id, err.Error())
		return err
	}

	log.Printf("[PostgresqlDatabase] Deleted book with ID %d", id)

	return nil
}

// UpdateBook updates a book with given ID
func (db *PostgresqlDatabase) UpdateBook(book *models.Book) error {
	query := "UPDATE books SET author = $1, title = $2 WHERE id = $3"

	_, err := db.conn.Exec(context.Background(), query, book.Author, book.Title, book.ID)
	if err != nil {
		log.Printf("[PostgresqlDatabase] Error while updating book with ID %d: %s", book.ID, err.Error())
		return err
	}

	log.Printf("[PostgresqlDatabase] Updated book with ID %d", book.ID)

	return nil
}
