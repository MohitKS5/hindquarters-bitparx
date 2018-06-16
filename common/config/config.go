package config

import (
	"github.com/bitparx/clientapi/auth/authtypes"
)

type Bitparx struct {
	// The configuration required for a server.
	Bitparx_Server struct {
		// The name of the server. This is usually the domain name, e.g 'bitpsrx.org', 'localhost'.
		ServerName string `yaml:"server_name"`
		// If set, allows registration by anyone who also has the shared
		// secret, even if registration is otherwise disabled.
		RegistrationSharedSecret string `yaml:"registration_shared_secret"`
		// This Home Server's ReCAPTCHA public key.
		//RecaptchaPublicKey string `yaml:"recaptcha_public_key"`
		//// This Home Server's ReCAPTCHA private key.
		//RecaptchaPrivateKey string `yaml:"recaptcha_private_key"`
		//// Boolean stating whether catpcha registration is enabled
		//// and required
		//RecaptchaEnabled bool `yaml:"enable_registration_captcha"`
		//// Secret used to bypass the captcha registration entirely
		//RecaptchaBypassSecret string `yaml:"captcha_bypass_secret"`
		//// HTTP API endpoint used to verify whether the captcha response
		//// was successful
		//RecaptchaSiteVerifyAPI string `yaml:"recaptcha_siteverify_api"`
		//// If set disables new users from registering (except via shared
		//// secrets)
		RegistrationDisabled bool `yaml:"registration_disabled"`
	} `yaml:"bitparx"`

	// Any information derived from the configuration options for later use.
	Derived struct {
		Registration struct {
			// Flows is a slice of flows, which represent one possible way that the client can authenticate a request.
			// As long as the generated flows only rely on config file options,
			// we can generate them on startup and store them until needed
			Flows []authtypes.Flow `json:"flows"`

			// Params that need to be returned to the client during
			// registration in order to complete registration stages.
			Params map[string]interface{} `json:"params"`
		}
	}
}

const (
	BINANCE_REST_URL = "https://api.binance.com"
)

func LoadConfig() *Bitparx {
	var allowedFlows = []authtypes.Flow{
		{
			[]authtypes.LoginType{authtypes.LoginTypeDummy},
		},
	}
	var bitparx Bitparx
	bitparx.Derived.Registration.Flows = allowedFlows
	bitparx.Bitparx_Server.RegistrationDisabled = false
	return &bitparx
}
