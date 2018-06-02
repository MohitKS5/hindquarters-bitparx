package authtypes

// Account represents a an account on admin side
type Account struct {
	UserID     string
	Username   string
	ServerName string
	Profile    Profile
	Created    int64
	// TODO: Other flags like IsAdmin, IsGuest
	// TODO: Devices
	// TODO: Associations (e.g. with application services)
}
