package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB() *sql.DB {

	connStr := "host=localhost port=5432 user=postgres password=bekarys7 dbname=onlineshop sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
