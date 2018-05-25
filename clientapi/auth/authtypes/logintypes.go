package authtypes

type LoginType string

// The relevant login types implemented
const (
	LoginTypeDummy              = "m.login.dummy"
	LoginTypeSharedSecret       = "login.shared_secret"
	LoginTypeRecaptcha          = "m.login.recaptcha"
)
