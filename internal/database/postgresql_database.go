package database

import (
	"context"
	"errors"

	"github.com/MSSkowron/BookRESTAPI/internal/models"
	"github.com/MSSkowron/BookRESTAPI/pkg/logger"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresqlDatabase is a Postgresql implementation of Database interface.
// It uses pgx as a Postgresql driver.
type PostgresqlDatabase struct {
	connPool *pgxpool.Pool
}

// NewPostgresqlDatabase creates a new PostgresqlDatabase.
func NewPostgresqlDatabase(connectionString string) (Database, error) {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	return &PostgresqlDatabase{
		connPool: pool,
	}, nil
}

// Close closes the database connection.
func (db *PostgresqlDatabase) Close() {
	db.connPool.Close()
}

// InsertUser inserts a new user into the database.
func (db *PostgresqlDatabase) InsertUser(user *models.User) (int, error) {
	var (
		query string = "INSERT INTO users (email, password, first_name, last_name, age) VALUES ($1, $2, $3, $4, $5) RETURNING id"
		id    int    = -1
	)

	if err := db.connPool.QueryRow(context.Background(), query, user.Email, user.Password, user.FirstName, user.LastName, user.Age).Scan(&id); err != nil {
		logger.Errorf("Error (%s) while inserting new user", err)

		return id, err
	}

	logger.Infof("Inserted new user with ID: %d", id)

	return id, nil
}

// SelectUserByID selects a user with given ID from the database.
func (db *PostgresqlDatabase) SelectUserByID(id int) (*models.User, error) {
	query := "SELECT * FROM users WHERE id=$1"

	user := &models.User{}
	if err := db.connPool.QueryRow(context.Background(), query, id).Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		logger.Errorf("Error (%s) while selecting user with ID: %d", err, id)

		return nil, err
	}

	logger.Infof("Selected user with ID: %d", id)

	return user, nil
}

// SelectUserByEmail selects a user with given email
func (db *PostgresqlDatabase) SelectUserByEmail(email string) (*models.User, error) {
	query := "SELECT * FROM users WHERE email=$1"

	user := &models.User{}
	if err := db.connPool.QueryRow(context.Background(), query, email).Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		logger.Errorf("Error (%s) while selecting user with email: %s", err, email)

		return nil, err
	}

	logger.Infof("Selected user with email: %s", email)

	return user, nil
}

// InsertBook inserts a new book into the database.
func (db *PostgresqlDatabase) InsertBook(book *models.Book) (int, error) {
	var (
		query string = "INSERT INTO books (author, title, created_by) VALUES ($1, $2, $3) RETURNING id"
		id    int    = -1
	)

	if err := db.connPool.QueryRow(context.Background(), query, book.Author, book.Title, book.CreatedBy).Scan(&id); err != nil {
		logger.Errorf("Error (%s) while inserting new book", err)

		return id, err
	}

	logger.Infof("Inserted new book with ID: %d", id)

	return id, nil
}

// SelectAllBooks selects all books from the database.
func (db *PostgresqlDatabase) SelectAllBooks() ([]*models.Book, error) {
	query := "SELECT * FROM books"

	rows, err := db.connPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*models.Book{}
	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author, &book.CreatedBy); err != nil {
			logger.Errorf("Error (%s) while selecting all books", err)

			return nil, err
		}

		books = append(books, book)
	}

	logger.Infoln("Selected all books")

	return books, nil
}

// SelectBookByID selects a book with given ID from the database.
func (db *PostgresqlDatabase) SelectBookByID(id int) (*models.Book, error) {
	query := "SELECT * FROM books WHERE id=$1"

	row := db.connPool.QueryRow(context.Background(), query, id)
	book := &models.Book{}
	if err := row.Scan(&book.ID, &book.CreatedAt, &book.Title, &book.Author, &book.CreatedBy); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		logger.Errorf("Error (%s) while selecting book with ID: %d", err, id)

		return nil, err
	}

	logger.Infof("Selected book with ID: %d", id)

	return book, nil
}

// DeleteBook deletes a book with given ID from the database.
func (db *PostgresqlDatabase) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id=$1"

	if _, err := db.connPool.Exec(context.Background(), query, id); err != nil {
		logger.Errorf("Error (%s) while deleting book with ID: %d", err, id)

		return err
	}

	logger.Infof("Deleted book with ID: %d", id)

	return nil
}

// UpdateBook updates a book with given ID in the database.
func (db *PostgresqlDatabase) UpdateBook(id int, book *models.Book) error {
	query := "UPDATE books SET author = $1, title = $2 WHERE id = $3"

	if _, err := db.connPool.Exec(context.Background(), query, book.Author, book.Title, id); err != nil {
		logger.Errorf("Error (%s) while updating book with ID: %d", err, id)

		return err
	}

	logger.Infof("Updated book with ID: %d", id)

	return nil
}
