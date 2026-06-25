package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := "user=takehome password=takehome dbname=inventory sslmode=disable host=localhost port=5438"
	if host := os.Getenv("DB_HOST"); host != "" {
		connStr = fmt.Sprintf("user=takehome password=takehome dbname=inventory sslmode=disable host=%s port=5432", host)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// README requires connection pool size of 10
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
