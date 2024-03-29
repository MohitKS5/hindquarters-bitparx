package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/common/storage"
	"golang.org/x/crypto/bcrypt"
	"github.com/bitparx/common/config"
	"github.com/bitparx/common"
)

const accountsSchema = `
-- Stores data about accounts.
CREATE TABLE IF NOT EXISTS account_accounts (
    -- The Bitparx_Server user ID username for this account
    username TEXT NOT NULL PRIMARY KEY,
    -- When this account was first created, as a unix timestamp (ms resolution).
    created_ts BIGINT NOT NULL,
    -- The password hash for this account. Can be NULL if this is a passwordless account.
    password_hash TEXT
    -- TODO:
    -- is_guest, is_admin, upgraded_ts, devices, any email reset stuff?
);
`

const insertAccountSQL = "" +
	"INSERT INTO account_accounts(username, created_ts, password_hash) VALUES ($1, $2, $3)"

const selectAccountByUsernameSQL = "" +
	"SELECT username FROM account_accounts WHERE username = $1"

const selectPasswordHashSQL = "" +
	"SELECT password_hash FROM account_accounts WHERE username = $1"

const selectAllAccountsSQL = "" +
	"SELECT A.username, created_ts, display_name, avatar_url FROM account_accounts A " +
	"INNER JOIN account_profiles B ON A.username = B.username"

// TODO: Update password

type accountsStatements struct {
	insertAccountStmt           *sql.Stmt
	selectAccountByUsernameStmt *sql.Stmt
	selectPasswordHashStmt      *sql.Stmt
	selectAllAccountsStmt       *sql.Stmt
	serverName                  string
}

func (s *accountsStatements) prepare(db *sql.DB, server string) (err error) {
	_, err = db.Exec(accountsSchema)
	if err != nil {
		return
	}

	// create default user (Administrator, bakaitBitcoin)
	createdTimeMS := time.Now().UnixNano() / 1000000
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(config.DEFAULT_ADMIN_PASSWORD), bcrypt.DefaultCost)
	_, err = db.Exec(insertAccountSQL, config.DEFAULT_ADMIN_USERNAME, createdTimeMS, string(hashBytes))
	if err != nil && !common.IsUniqueConstraintViolationErr(err) {
		return
	}

	if s.insertAccountStmt, err = db.Prepare(insertAccountSQL); err != nil {
		return
	}
	if s.selectAccountByUsernameStmt, err = db.Prepare(selectAccountByUsernameSQL); err != nil {
		return
	}
	if s.selectPasswordHashStmt, err = db.Prepare(selectPasswordHashSQL); err != nil {
		return
	}
	if s.selectAllAccountsStmt, err = db.Prepare(selectAllAccountsSQL); err != nil {
		return
	}
	s.serverName = server
	return
}

// insertAccount creates a new account. 'hash' should be the password hash for this account.
// Returns an error if this account already exists. Returns the account
// on success.
func (s *accountsStatements) insertAccount(
	ctx context.Context, username, hash string) (*authtypes.Account, error) {
	createdTimeMS := time.Now().UnixNano() / 1000000
	stmt := s.insertAccountStmt

	var err error
	_, err = stmt.ExecContext(ctx, username, createdTimeMS, hash)
	if err != nil {
		return nil, err
	}

	return &authtypes.Account{
		Username:   username,
		UserID:     makeUserID(username, s.serverName),
		ServerName: s.serverName,
	}, nil
}

// returns hash of username provided. used for login
func (s *accountsStatements) selectPasswordHash(
	ctx context.Context, username string,
) (hash string, err error) {
	err = s.selectPasswordHashStmt.QueryRowContext(ctx, username).Scan(&hash)
	return
}

// returns the account object corresponding to the username
func (s *accountsStatements) selectAccountByUsername(
	ctx context.Context, username string,
) (*authtypes.Account, error) {
	var acc authtypes.Account
	stmt := s.selectAccountByUsernameStmt
	err := stmt.QueryRowContext(ctx, username).Scan(&acc.Username)
	acc.UserID = makeUserID(username, s.serverName)
	acc.ServerName = s.serverName
	fmt.Println(acc)
	return &acc, err
}

// returns all accounts
func (s *accountsStatements) selectAllAccounts(ctx context.Context) (accounts []authtypes.Account, err error) {
	accounts = []authtypes.Account{}
	rows, err := s.selectAllAccountsStmt.QueryContext(ctx)
	if err != nil {
		return
	}
	for rows.Next() {
		acc := authtypes.Account{}
		err = rows.Scan(&acc.Username, &acc.Created, &acc.Profile.DisplayName, &acc.Profile.AvatarURL)
		if err != nil {
			return
		}
		acc.UserID = makeUserID(acc.Username, storage.SERVER_NAME)
		accounts = append(accounts, acc)
	}
	return
}

// make userID by concatenating username:servername
func makeUserID(username string, server string) string {
	return fmt.Sprintf("@%s:%s", username, string(server))
}
