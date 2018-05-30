package storage

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"github.com/bitparx/clientapi/auth/storage/levels"
)

const (
	HOST        = "localhost"
	PORT        = 5433
	DB_USER     = "postgres"
	DB_PASSWORD = "bitparx"
	DB_NAME     = "bitparx"
	SERVER_NAME = "PostgreSQL 10 (x86)"
)

func PostgresConnect() (accountDB *accounts.Database, deviceDB *devices.Database, levelsDB *levels.Database) {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, DB_USER, DB_PASSWORD, DB_NAME)
	//connect to accounts database
	accountDB, err := accounts.NewDatabase(dbinfo, SERVER_NAME)
	if err != nil {
		panic(err)
	}

	// connect to devices database
	deviceDB, err = devices.NewDatabase(dbinfo, SERVER_NAME)
	if err != nil {
		panic(err)
	}

	// connect to levels database
	levelsDB, err = levels.NewDatabase(dbinfo, SERVER_NAME)
	if err != nil {
		fmt.Println("LevelsDB Error")
		panic(err)
	}

	//defer db.db.Close()
	return
}
