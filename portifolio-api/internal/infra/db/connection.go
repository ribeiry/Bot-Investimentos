package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL não definida")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("falha ao abrir conexão com o banco: %v", err)
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		log.Printf("falha ao pingar o banco: %v", err)
		return nil, err
	}
	return db, nil
}
