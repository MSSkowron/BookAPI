package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MSSkowron/GoBankAPI/model"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*model.User) error
	GetUserByEmail(string) (*model.User, error)
	GetBooks() ([]*model.Book, error)
}

type PostgresSQLStorage struct {
	db *sql.DB
}

func NewPostgresSQLStorage() (*PostgresSQLStorage, error) {
	connStr := "user=gobookapiuser dbname=postgres password=gobookapipassword sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	postgresSQLStore := &PostgresSQLStorage{
		db: db,
	}

	return postgresSQLStore, nil
}

func (s *PostgresSQLStorage) CreateUser(user *model.User) error {
	query := `insert into users (email, password, first_name, last_name, age) values ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(query, user.Email, user.Password, user.FirstName, user.LastName, user.Age)
	if err != nil {
		log.Println("[PostgresSQLStorage] Error while inserting new user: " + err.Error())
		return err
	}

	log.Println("[PostgresSQLStorage] Inserted new user")

	return nil
}

func (s *PostgresSQLStorage) GetUserByEmail(email string) (*model.User, error) {
	query := `select * from users where email=$1`

	row := s.db.QueryRow(query, email)

	user := &model.User{}
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Age); err != nil {
		return nil, err
	}

	log.Println("[PostgresSQLStorage] User correctly pulled from database")

	return user, nil
}

func (s *PostgresSQLStorage) GetBooks() ([]*model.Book, error) {
	query := `select * from books`

	rows, err := s.db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(rows)

	books := []*model.Book{}
	for rows.Next() {
		book := &model.Book{}
		if err := rows.Scan(&book.ID, &book.Isbn, &book.Title, &book.Author); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	log.Println("[PostgresSQLStorage] Books correctly pulled from database")

	return books, nil
}
