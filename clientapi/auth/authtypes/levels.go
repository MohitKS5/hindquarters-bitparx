
package authtypes

// Levels represents the relationship between a user and actions he can perform
type Membership struct {
	Username	string
	Admin    	bool
	Moderator   bool
}
