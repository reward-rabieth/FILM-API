package main

import (
	"database/sql"
	"fmt"

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

func NewPostgresStore() (*PostgresStore, error) {

	constr := "user=postgres dbname=postgres password=2981 sslmode=disable"
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

	query := `insert into
	 comment (text,ip_address,created_at)
     values($1,$2,$3)`
	resp, err := s.db.Query(query, com.Text, com.IpAddress, com.CreatedAt)
	if err != nil {
		return err
	}

	fmt.Printf("%v/n", resp)

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
