package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func newPostgresStore() (*PostgresStore, error) {
	connectionStr := "user=postgres dbname=bankdb port=5000 password=mostafa sslmode=disable"

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	sqlStatment := `INSERT INTO account (firstName,lastName,number ,balance,createdat) VALUES($1,$2,$3,$4,$5)`
	res, err := s.db.Query(sqlStatment, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", res)
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query(`DELETE FROM account WHERE id=$1`, id)
	return err
}
func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE id= $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return ScanIntoAccount(rows)
	}
	return nil, fmt.Errorf("no account for this %d", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	sqlStatement := `SELECT * FROM account`
	rows, err := s.db.Query(sqlStatement)

	if err != nil {
		return nil, err
	}
	accounts := []*Account{}

	for rows.Next() {
		account, err := ScanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	return accounts, nil
}

func ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
