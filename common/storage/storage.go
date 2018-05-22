package storage

import (
"database/sql"
"fmt"
_ "github.com/lib/pq"
)

const (
	DB_USER = "postgres"
	DB_PASSWORD = "bitparx"
	DB_NAME = "bitparx"
)

func PostgresConnect() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s port=**** sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	return db
}
