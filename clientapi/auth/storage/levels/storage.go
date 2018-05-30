package levels

import (
	"database/sql"
	"context"
	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/common"
	"fmt"
)

// Database represents a device database.
type Database struct {
	db       *sql.DB
	accounts levelsStatements
}

// NewDatabase creates a new device database
func NewDatabase(dataSourceName string, serverName string) (*Database, error) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("postgres", dataSourceName); err != nil {
		return nil, err
	}
	d := levelsStatements{}
	if err = d.prepare(db, serverName); err != nil {
		fmt.Println("error at prepare")
		return nil, err
	}
	return &Database{db, d}, nil
}

// GetDeviceByLocalpart returns the account matching the given localpart.
// Returns sql.ErrNoRows if no matching device was found.
func (d *Database) GetAccountByLocalpart(
	ctx context.Context, localpart string,
) (*authtypes.LevelsData, error) {
	return d.accounts.selectAccountByLocalpart(ctx, localpart)
}

// CreateLevel makes a new device associated with the given localpart and set all levels to false.
// Returns nil on success.
func (d *Database) CreateLevel(
	ctx context.Context, localpart string) error {
	return common.WithTransaction(d.db, func(txn *sql.Tx) error {
		return d.accounts.insertAccount(ctx, txn, localpart)
	})
}

// GetAllAccounts proves with all acoount levels
func (d *Database) GetAllAccounts(
	ctx context.Context) ([]authtypes.LevelsData, error) {
	return d.accounts.selectAllAccounts(ctx)
}

// handle PUT /levels/admin/${localpart}
func (d *Database) UpdateLevelAdmin(ctx context.Context, admin bool, localpart string) error {
	return common.WithTransaction(d.db, func(txn *sql.Tx) error {
		return d.accounts.updateLevelAdmin(ctx, txn, admin, localpart)
	})
}

// handle PUT /levels/moderator/${localpart}
func (d *Database) UpdateLevelModerator(ctx context.Context, moderator bool, localpart string) error {
	return common.WithTransaction(d.db, func(txn *sql.Tx) error {
		return d.accounts.updateLevelModerator(ctx, txn, moderator, localpart)
	})
}