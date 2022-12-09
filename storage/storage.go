package storage

import (
	"database/sql"

	"github.com/MSSkowron/GoBankAPI/model"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*model.User) error
	GetUserByUsername(string) (*model.User, error)
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

func (s *PostgresSQLStorage) CreateUser(*model.User) error {
	return nil
}

func (s *PostgresSQLStorage) GetUserByUsername(string) (*model.User, error) {
	return nil, nil
}
