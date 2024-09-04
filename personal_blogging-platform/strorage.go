package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	GetUserByEmail(email string) (*User, error)
}
type PostgresStore struct {
	db *sql.DB
}

func newPostgresConnection() (*PostgresStore, error) {
	connectionString := "user=postgres dbname=blogdb password=123456 port=4040 sslmode=disable"

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Database Connected")
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateUser(user *User) error {
	SQLstatement := `INSERT INTO users(first_name,last_name,email,user_password) VALUES($1,$2,$3,$4)`
	res, err := s.db.Query(SQLstatement, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res)
	return nil
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	SQLStatment := `SELECT * FROM users WHERE email=$1`
	res := s.db.QueryRow(SQLStatment, email)
	
	user:= &User{}
	err := res.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}
