
package accounts

import (
	"context"
	"database/sql"

	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/common/config"
	"github.com/bitparx/common"
)

const profilesSchema = `
-- Stores data about accounts profiles.
CREATE TABLE IF NOT EXISTS account_profiles (
    -- The bitparx user ID username for this account
    username TEXT NOT NULL PRIMARY KEY,
    -- The display name for this account
    display_name TEXT,
    -- The URL of the avatar for this account
    avatar_url TEXT
);
`

const insertProfileSQL = "" +
	"INSERT INTO account_profiles(username, display_name, avatar_url) VALUES ($1, $2, $3)"

const selectProfileByUsernameSQL = "" +
	"SELECT username, display_name, avatar_url FROM account_profiles WHERE username = $1"

const setAvatarURLSQL = "" +
	"UPDATE account_profiles SET avatar_url = $1 WHERE username = $2"

const setDisplayNameSQL = "" +
	"UPDATE account_profiles SET display_name = $1 WHERE username = $2"

type profilesStatements struct {
	insertProfileStmt            *sql.Stmt
	selectProfileByUsernameStmt *sql.Stmt
	setAvatarURLStmt             *sql.Stmt
	setDisplayNameStmt           *sql.Stmt
}

func (s *profilesStatements) prepare(db *sql.DB) (err error) {
	_, err = db.Exec(profilesSchema)
	if err != nil {
		return
	}

	_, err = db.Exec(insertProfileSQL,config.DEFAULT_ADMIN_USERNAME,"","")
	if err != nil && !common.IsUniqueConstraintViolationErr(err) {
		return
	}

	if s.insertProfileStmt, err = db.Prepare(insertProfileSQL); err != nil {
		return
	}
	if s.selectProfileByUsernameStmt, err = db.Prepare(selectProfileByUsernameSQL); err != nil {
		return
	}
	if s.setAvatarURLStmt, err = db.Prepare(setAvatarURLSQL); err != nil {
		return
	}
	if s.setDisplayNameStmt, err = db.Prepare(setDisplayNameSQL); err != nil {
		return
	}
	return
}

func (s *profilesStatements) insertProfile(
	ctx context.Context, username string,
) (err error) {
	_, err = s.insertProfileStmt.ExecContext(ctx, username, "", "")
	return
}

func (s *profilesStatements) selectProfileByUsername(
	ctx context.Context, username string,
) (*authtypes.Profile, error) {
	var profile authtypes.Profile
	err := s.selectProfileByUsernameStmt.QueryRowContext(ctx, username).Scan(
		&profile.Username, &profile.DisplayName, &profile.AvatarURL,
	)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *profilesStatements) setAvatarURL(
	ctx context.Context, username string, avatarURL string,
) (err error) {
	_, err = s.setAvatarURLStmt.ExecContext(ctx, avatarURL, username)
	return
}

func (s *profilesStatements) setDisplayName(
	ctx context.Context, username string, displayName string,
) (err error) {
	_, err = s.setDisplayNameStmt.ExecContext(ctx, displayName, username)
	return
}
