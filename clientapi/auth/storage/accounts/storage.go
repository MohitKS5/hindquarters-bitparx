
package accounts

import (
	"context"
	"database/sql"

	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/common"
	"golang.org/x/crypto/bcrypt"
	// Import the postgres database driver.
	_ "github.com/lib/pq"
)

// Database represents an account database
type Database struct {
	db *sql.DB
	accounts     accountsStatements
	profiles     profilesStatements
	serverName   string
}

// NewDatabase creates a new accounts and profiles database
func NewDatabase(dataSourceName string, serverName string) (*Database, error) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("postgres", dataSourceName); err != nil {
		return nil, err
	}

	a := accountsStatements{}
	if err = a.prepare(db, serverName); err != nil {
		return nil, err
	}
	p := profilesStatements{}
	if err = p.prepare(db); err != nil {
		return nil, err
	}
	return &Database{db, a, p,serverName}, nil

}

// GetAccountByPassword returns the account associated with the given username and password.
// Returns sql.ErrNoRows if no account exists which matches the given username.
func (d *Database) GetAccountByPassword(
	ctx context.Context, username, plaintextPassword string,
) (*authtypes.Account, error) {
	hash, err := d.accounts.selectPasswordHash(ctx, username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintextPassword)); err != nil {
		return nil, err
	}
	return d.accounts.selectAccountByUsername(ctx, username)
}

// GetProfileByUsername returns the profile associated with the given username.
// Returns sql.ErrNoRows if no profile exists which matches the given username.
func (d *Database) GetProfileByUsername(
	ctx context.Context, username string,
) (*authtypes.Profile, error) {
	return d.profiles.selectProfileByUsername(ctx, username)
}

// SetAvatarURL updates the avatar URL of the profile associated with the given
// username. Returns an error if something went wrong with the SQL query
func (d *Database) SetAvatarURL(
	ctx context.Context, username string, avatarURL string,
) error {
	return d.profiles.setAvatarURL(ctx, username, avatarURL)
}

// SetDisplayName updates the display name of the profile associated with the given
// username. Returns an error if something went wrong with the SQL query
func (d *Database) SetDisplayName(
	ctx context.Context, username string, displayName string,
) error {
	return d.profiles.setDisplayName(ctx, username, displayName)
}

// CreateAccount makes a new account with the given login name and password, and creates an empty profile
// for this account. If the
// account already exists, it will return nil, nil.
func (d *Database) CreateAccount(
	ctx context.Context, username, plaintextPassword string,
) (*authtypes.Account, error) {
	var err error

	// Generate a password hash if this is not a password-less user
	hash := ""
	if plaintextPassword != "" {
		hash, err = hashPassword(plaintextPassword)
		if err != nil {
			return nil, err
		}
	}
	if err := d.profiles.insertProfile(ctx, username); err != nil {
		if common.IsUniqueConstraintViolationErr(err) {
			return nil, nil
		}
		return nil, err
	}
	return d.accounts.insertAccount(ctx, username, hash)
}

func hashPassword(plaintext string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	return string(hashBytes), err
}

// CheckAccountAvailability checks if the username/localpart is already present
// in the database.
// If the DB returns sql.ErrNoRows the Username isn't taken.
func (d *Database) CheckAccountAvailability(ctx context.Context, username string) (bool, error) {
	_, err := d.accounts.selectAccountByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return true, nil
	}
	return false, err
}
