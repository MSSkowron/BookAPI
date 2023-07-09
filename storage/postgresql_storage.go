package storage

import (
	"database/sql"
	"log"

	"github.com/MSSkowron/BookRESTAPI/model"
	_ "github.com/lib/pq"
)

const (
	driverName       = "postgres"
	connectionString = "host=database user=postgres password=postgres dbname=postgres sslmode=disable"
)

type PostgresSQLStorage struct {
	db *sql.DB
}

func NewPostgresSQLStorage() (*PostgresSQLStorage, error) {
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	postgresSQLStorage := &PostgresSQLStorage{
		db: db,
	}

	return postgresSQLStorage, nil
}

func (s *PostgresSQLStorage) InsertUser(user *model.User) (int, error) {
	var (
		query = "insert into users (email, password, first_name, last_name, age) values ($1, $2, $3, $4, $5) RETURNING id"
		id    = -1
	)

	if err := s.db.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, user.Age).Scan(&id); err != nil {
		log.Printf("[PostgresSQLStorage] Error while inserting new user: %s", err.Error())

		return id, err
	}

	log.Printf("[PostgresSQLStorage] Inserted new user with ID: %d", id)

	return id, nil
}

func (s *PostgresSQLStorage) SelectUserByEmail(email string) (*model.User, error) {
	query := "select * from users where email=$1"

	row := s.db.QueryRow(query, email)
	user := &model.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age); err != nil {
		log.Printf("[PostgresSQLStorage] Error while selecting user with email %s: %s", email, err.Error())

		return nil, err
	}

	log.Printf("[PostgresSQLStorage] Selected user with email: %s", email)

	return user, nil
}

func (s *PostgresSQLStorage) InsertBook(book *model.Book) (int, error) {
	var (
		query = "insert into books (author, title) values ($1, $2) RETURNING id"
		id    = 0
	)

	if err := s.db.QueryRow(query, book.Author, book.Title).Scan(&id); err != nil {
		log.Printf("[PostgresSQLStorage] Error while inserting new book: %s", err.Error())

		return id, err
	}

	log.Printf("[PostgresSQLStorage] Inserted new book with ID %d", id)

	return id, nil
}

func (s *PostgresSQLStorage) SelectAllBooks() ([]*model.Book, error) {
	query := "select * from books"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	books := []*model.Book{}
	for rows.Next() {
		book := &model.Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			log.Printf("[PostgresSQLStorage] Error while selecting all books: %s", err.Error())

			return nil, err
		}

		books = append(books, book)
	}

	log.Println("[PostgresSQLStorage] Selected all books")

	return books, nil
}

func (s *PostgresSQLStorage) SelectBookByID(id int) (*model.Book, error) {
	query := "select * from books where id=$1"

	row := s.db.QueryRow(query, id)
	book := &model.Book{}
	if err := row.Scan(&book.ID, &book.Title, &book.Author); err != nil {
		log.Printf("[PostgresSQLStorage] Error while selecting book with ID %d: %s", id, err.Error())

		return nil, err
	}

	log.Printf("[PostgresSQLStorage] Selected book with ID: %d", id)

	return book, nil
}

func (s *PostgresSQLStorage) DeleteBook(id int) error {
	query := "delete from books where id=$1"

	if _, err := s.db.Exec(query, id); err != nil {
		log.Printf("[PostgresSQLStorage] Error while deleting book with ID %d: %s", id, err.Error())

		return err
	}

	log.Printf("[PostgresSQLStorage] Deleted book with ID %d", id)

	return nil
}

func (s *PostgresSQLStorage) UpdateBook(book *model.Book) error {
	query := "UPDATE books SET author = $1, title= $2 WHERE id = $3"

	if _, err := s.db.Exec(query, book.Author, book.Title, book.ID); err != nil {
		log.Printf("[PostgresSQLStorage] Error while updating book with ID %d: %s", book.ID, err.Error())

		return err
	}

	log.Printf("[PostgresSQLStorage] Updated book with ID %d", book.ID)

	return nil
}
