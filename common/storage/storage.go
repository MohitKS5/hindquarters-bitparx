package storage

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/auth/storage/devices"
)

const (
	HOST        = "localhost"
	PORT        = 5433
	DB_USER     = "postgres"
	DB_PASSWORD = "bitparx"
	DB_NAME     = "bitparx"
	SERVER_NAME = "PostgreSQL 10 (x86)"
)

func PostgresConnect() (accountDB *accounts.Database, deviceDB *devices.Database) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, DB_USER, DB_PASSWORD, DB_NAME)
	accountDB, err := accounts.NewDatabase(dbinfo, SERVER_NAME)
	if err != nil {
		panic(err)
	}
	deviceDB, err = devices.NewDatabase(dbinfo, SERVER_NAME)
	if err != nil {
		panic(err)
	}

	//defer db.db.Close()
	return
}
