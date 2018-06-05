package levels

import (
	"context"
	"database/sql"
	"github.com/bitparx/common"
	"github.com/bitparx/clientapi/auth/authtypes"
	"fmt"
	"github.com/bitparx/common/config"
)

const levelsSchema = `
-- Stores data about accounts.
CREATE TABLE IF NOT EXISTS account_levels (
    -- The user ID username for this device. This is preferable to storing the full user_id
    -- as it is smaller, makes it clearer that we only manage accounts for our own users, and may make
    -- migration to different domain names easier.
    username TEXT NOT NULL PRIMARY KEY,
    -- When this level was first assigned on the network, as a unix timestamp (ms resolution).
    -- assigned_ts BIGINT NOT NULL,
    -- Who approved the level
    -- todo assignee TEXT NOT NULL,
    -- the levels start here 
    -- is admin
    admin BOOLEAN,
    -- is moderator
    moderator BOOLEAN
    -- TODO: device keys, device display names, last used ts and IP address?, token restrictions (if 3rd-party OAuth app)
)
`

const insertLevelSQL = "" +
	"INSERT INTO account_levels(username, admin, moderator) VALUES ($1, $2, $3)"

const selectAccountByLocalpartSQL = "" +
	"SELECT * FROM account_levels WHERE username = $1"

const selectAccountsByAdminSQL = "" +
	"SELECT * FROM account_levels WHERE admin = $1"

const selectLocalpartsByAdminSQL = "" +
	"SELECT username FROM account_levels WHERE admin = $1"

const updateLevelAdminSQL = "" +
	"UPDATE account_levels SET admin = $1 WHERE username = $2"

const updateLevelModeratorSQL = "" +
	"UPDATE account_levels SET moderator = $1 WHERE username = $2"

const deleteAccountSQL = "" +
	"DELETE FROM account_levels WHERE username = $1"

const selectAllAccountsStmt = "" +
	"SELECT * FROM account_levels"

type levelsStatements struct {
	insertLevelStmt              *sql.Stmt
	selectAccountByLocalpartstmt *sql.Stmt
	selectAccountByAdminstmt     *sql.Stmt
	selectLocalpartsByAdminstmt  *sql.Stmt
	updateLevelAdminStmt         *sql.Stmt
	updateLevelModeratorStmt     *sql.Stmt
	deleteAccountStmt            *sql.Stmt
	selectAllAccountsStmt        *sql.Stmt
	serverName                   string
}

func (s *levelsStatements) prepare(db *sql.DB, server string) (err error) {
	_, err = db.Exec(levelsSchema)
	if err != nil {
		return
	}

	// make default user admin
	_, err = db.Exec(insertLevelSQL, config.DEFAULT_ADMIN_USERNAME, true, true)
	if err != nil && !common.IsUniqueConstraintViolationErr(err) {
		return
	}

	if s.insertLevelStmt, err = db.Prepare(insertLevelSQL); err != nil {
		return
	}
	if s.selectAccountByLocalpartstmt, err = db.Prepare(selectAccountByLocalpartSQL); err != nil {
		return
	}
	if s.selectAccountByAdminstmt, err = db.Prepare(selectAccountsByAdminSQL); err != nil {
		return
	}
	if s.selectLocalpartsByAdminstmt, err = db.Prepare(selectLocalpartsByAdminSQL); err != nil {
		return
	}
	if s.updateLevelAdminStmt, err = db.Prepare(updateLevelAdminSQL); err != nil {
		return
	}

	if s.updateLevelModeratorStmt, err = db.Prepare(updateLevelModeratorSQL); err != nil {
		return
	}

	if s.deleteAccountStmt, err = db.Prepare(deleteAccountSQL); err != nil {
		return
	}

	if s.selectAllAccountsStmt, err = db.Prepare(selectAllAccountsStmt); err != nil {
		return
	}
	s.serverName = server
	return
}

// insertLevels creates a new device. Returns an error if any account with the same access token already exists.
// Returns the nil on success.
func (s *levelsStatements) insertAccount(
	ctx context.Context, txn *sql.Tx, localpart string) error {
	stmt := common.TxStmt(txn, s.insertLevelStmt)
	if _, err := stmt.ExecContext(ctx, localpart, false, false); err != nil {
		return err
	}
	return nil
}

func (s *levelsStatements) selectAccountByLocalpart(
	ctx context.Context, localpart string,
) (*authtypes.LevelsData, error) {
	var acc authtypes.LevelsData
	stmt := s.selectAccountByLocalpartstmt
	err := stmt.QueryRowContext(ctx, localpart).Scan(&acc.Username, &acc.Access.Admin, &acc.Access.Moderator)
	if err != nil {
		return nil, err
	}
	return &acc, err
}

func (s *levelsStatements) deleteDevicesByLocalpart(
	ctx context.Context, txn *sql.Tx, localpart string,
) error {
	stmt := common.TxStmt(txn, s.deleteAccountStmt)
	_, err := stmt.ExecContext(ctx, localpart)
	return err
}

func (s *levelsStatements) updateLevelAdmin(
	ctx context.Context, txn *sql.Tx, admin sql.NullBool, localpart string) error {
	stmt := common.TxStmt(txn, s.updateLevelAdminStmt)
	res, err := stmt.ExecContext(ctx, admin, localpart)
	fmt.Println("reached database" + localpart)
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (s *levelsStatements) updateLevelModerator(
	ctx context.Context, txn *sql.Tx, moderator sql.NullBool, localpart string) error {
	stmt := common.TxStmt(txn, s.updateLevelModeratorStmt)
	res, err := stmt.ExecContext(ctx, moderator, localpart)
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (s *levelsStatements) selectLocalpartsByAdmin(
	ctx context.Context, admin bool,
) ([]authtypes.LevelsData, error) {
	accounts := []authtypes.LevelsData{}

	rows, err := s.selectLocalpartsByAdminstmt.QueryContext(ctx, admin)

	if err != nil {
		return accounts, err
	}

	for rows.Next() {
		var acc authtypes.LevelsData
		err = rows.Scan(&acc.Username)
		if err != nil {
			return accounts, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (s *levelsStatements) selectAllAccounts(
	ctx context.Context) ([]authtypes.LevelsData, error) {
	accounts := []authtypes.LevelsData{}

	rows, err := s.selectAllAccountsStmt.QueryContext(ctx)

	if err != nil {
		return accounts, err
	}

	for rows.Next() {
		var acc authtypes.LevelsData
		err = rows.Scan(&acc.Username, &acc.Access.Admin, &acc.Access.Moderator)
		if err != nil {
			return accounts, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}
