package storage

import (
	"fmt"
	_ "github.com/lib/pq"
)

// todo take these as flags while running exe
const (
	HOST        = "localhost"
	PORT        = 5433
	DB_USER     = "postgres"
	DB_PASSWORD = "bitparx"
	DB_NAME     = "bitparx"
	SERVER_NAME = "PostgreSQL 10 (x86)"
)

func PostgresConnectCredentials() string {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, DB_USER, DB_PASSWORD, DB_NAME)
	return dbinfo
}
