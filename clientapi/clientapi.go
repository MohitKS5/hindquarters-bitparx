package clientapi

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"github.com/bitparx/clientapi/auth/storage/levels"
	"github.com/bitparx/common/storage"
)

func PostgresConnect() (accountDB *accounts.Database, deviceDB *devices.Database, levelsDB *levels.Database) {
	dbinfo := storage.PostgresConnectCredentials()
	serverName := storage.SERVER_NAME
	//connect to accounts database
	accountDB, err := accounts.NewDatabase(dbinfo, serverName)
	if err != nil {
		panic(err)
	}

	// connect to devices database
	deviceDB, err = devices.NewDatabase(dbinfo, serverName)
	if err != nil {
		panic(err)
	}

	// connect to levels database
	levelsDB, err = levels.NewDatabase(dbinfo, serverName)
	if err != nil {
		fmt.Println("LevelsDB Error")
		panic(err)
	}

	//defer db.db.Close()
	return
}
