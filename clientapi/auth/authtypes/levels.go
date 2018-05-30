
package authtypes

import "database/sql"

// LevelsData represents the relationship between a user and actions he can perform
type LevelsData struct {
	Username	string
	Access		Levels
}

type Levels struct {
	Admin sql.NullBool
	Moderator sql.NullBool
}
