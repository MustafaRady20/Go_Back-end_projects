package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	GetUserByEmail(email string) (*User, error)
	getAllArticles() ([]*Article, error)
	CreateArticle(userId string, article Article) error
	UpdateArticle(newArticle UpdatedArticle, userid string) error
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
	_, err := s.db.Exec(SQLstatement, user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	SQLStatment := `SELECT * FROM users WHERE email=$1`
	res := s.db.QueryRow(SQLStatment, email)

	user := &User{}
	err := res.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresStore) getAllArticles() ([]*Article, error) {

	SQLstatement := `SELECT * FROM article`

	rows, err := s.db.Query(SQLstatement)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	articles := []*Article{}

	for rows.Next() {
		article, err := ScanIntoArticle(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row into article: %w", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return articles, nil
}

func (s *PostgresStore) CreateArticle(userId string, article Article) error {
	sqlStatement := `INSERT INTO article(article_title,article_content,user_id) VALUES($1,$2,$3)`
	res, err := s.db.Exec(sqlStatement, &article.Title, &article.Content, &userId)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

func (s *PostgresStore) UpdateArticle(newArticle UpdatedArticle, userid string) error {

	sqlStatment := `UPDATE article SET article_title=$1,article_content = $2 WHERE user_id =$3`

	_, err := s.db.Exec(sqlStatment, &newArticle.Title, &newArticle.Content, &userid)

	if err != nil {
		return err
	}
	return nil
}
func ScanIntoArticle(rows *sql.Rows) (*Article, error) {
	article := new(Article)
	err := rows.Scan(&article.ID,
		&article.Title,
		&article.Content,
		&article.CreatedAt,
		&article.UserId)

	return article, err
}
