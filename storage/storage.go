package storage

import (
	"database/sql"

	"github.com/MSSkowron/GoBankAPI/model"
)

type Storage interface {
	CreateUser(*model.User) error
	GetUserByUsername(string) (*model.User, error)
}

type PostgresSQLStorage struct {
	db *sql.DB
}

func NewPostgresSQLStorage() (*PostgresSQLStorage, error) {
	connStr := "TODO"

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

	if err := postgresSQLStore.init(); err != nil {
		return nil, err
	}

	return postgresSQLStore, nil
}

func (s *PostgresSQLStorage) init() error {
	if err := s.createUsersTable(); err != nil {
		return err
	}

	return nil
}

func (s *PostgresSQLStorage) createUsersTable() error {
	query := ` CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email varchar(50),
		password varchar(256),
		first_name varchar(50),
		last_name varchar(50),
		age smallint,
	)
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresSQLStorage) CreateUser(*model.User) error {
	return nil
}

func (s *PostgresSQLStorage) GetUserByUsername(string) (*model.User, error) {
	return nil, nil
}
