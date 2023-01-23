package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateComment(*Comment) error
	//CreateCommentByid(*SavedComment) error
	Getcomments() ([]*Comment, error)
	//GetCommentsByid() ([]*SavedComment, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(cfg DBconfig) (*PostgresStore, error) {

	//constr := "dbuser=postgres dbname=postgres password=2981 sslmode=disable"

	constr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", cfg.User, cfg.Name, cfg.Password, cfg.SSLMode)

	db, err := sql.Open("postgres", constr)
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

func (s *PostgresStore) Init() error {
	return s.CreateCommentTable()

}

func (s *PostgresStore) CreateCommentTable() error {
	query := `create table if not exists comment(
		text varchar(50),
		ip_address varchar(50),
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateComment(com *Comment) error {

	now := time.Now().Local()
	com.CreatedAt = now

	query := `insert into
	 comment (text,ip_address,created_at)
     values($1,$2,$3)`
	_, err := s.db.Query(query, com.Text, com.IpAddress, com.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) Getcomments() ([]*Comment, error) {
	rows, err := s.db.Query("select * from comment")
	if err != nil {
		return nil, err
	}

	Datacomments := []*Comment{}
	for rows.Next() {
		comment := new(Comment)
		err := rows.Scan(
			&comment.Text,
			&comment.IpAddress,
			&comment.CreatedAt)
		if err != nil {
			return nil, err
		}

		Datacomments = append(Datacomments, comment)

	}
	return Datacomments, nil

}
