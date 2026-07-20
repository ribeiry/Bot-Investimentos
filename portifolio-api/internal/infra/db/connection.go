package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (*sql.DB, error) {
	var err error
	sqlLite, err := sql.Open("sqlite3", "./data/portfolio.db")

	if err != nil {
		return nil, err
	}
	err = sqlLite.Ping()
	if err != nil {
		log.Println("Erro connectionDB")
		return nil, err
	}
	return sqlLite, nil
}
