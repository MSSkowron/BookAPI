package database

import (
	"context"
	"log"

	"github.com/MSSkowron/BookRESTAPI/internal/models"
	"github.com/MSSkowron/BookRESTAPI/pkg/logger"
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
		logger.Errorf("Error (%s) while inserting new user", err)
		return id, err
	}

	logger.Infof("Inserted new user with ID: %d", id)

	return id, nil
}

// SelectUserByEmail selects a user with given email
func (db *PostgresqlDatabase) SelectUserByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email=$1"

	row := db.conn.QueryRow(context.Background(), query, email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age)
	if err != nil {
		logger.Errorf("Error (%s) while selecting user with email: %s", err, email)
		return nil, err
	}

	logger.Infof("Selected user with email: %s", email)

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
		logger.Errorf("Error (%s) while inserting new book", err)
		return id, err
	}

	logger.Infof("Inserted new book with ID: %d", id)

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
			logger.Errorf("Error (%s) while selecting all books", err)
			return nil, err
		}

		books = append(books, book)
	}

	log.Println("Selected all books")

	return books, nil
}

// SelectBookByID selects a book with given ID
func (db *PostgresqlDatabase) SelectBookByID(id int) (*models.Book, error) {
	query := "SELECT * FROM books WHERE id=$1"

	row := db.conn.QueryRow(context.Background(), query, id)
	book := &models.Book{}
	if err := row.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author); err != nil {
		logger.Errorf("Error (%s) while selecting book with ID: %d", err, id)
		return nil, err
	}

	logger.Infof("Selected book with ID: %d", id)

	return book, nil
}

// DeleteBook deletes a book with given ID
func (db *PostgresqlDatabase) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id=$1"

	if _, err := db.conn.Exec(context.Background(), query, id); err != nil {
		logger.Errorf("Error (%s) while deleting book with ID: %d", err, id)
		return err
	}

	logger.Infof("Deleted book with ID: %d", id)

	return nil
}

// UpdateBook updates a book with given ID
func (db *PostgresqlDatabase) UpdateBook(book *models.Book) error {
	query := "UPDATE books SET author = $1, title = $2 WHERE id = $3"

	_, err := db.conn.Exec(context.Background(), query, book.Author, book.Title, book.ID)
	if err != nil {
		logger.Errorf("Error (%s) while updating book with ID: %d", err, book.ID)
		return err
	}

	logger.Infof("Updated book with ID: %d", book.ID)

	return nil
}
