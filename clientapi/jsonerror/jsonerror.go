package jsonerror

import (
	"fmt"
	"net/http"

	"github.com/bitparx/util"
)

// Error represents the "standard error response"
type ParxError struct {
	ErrCode string `json:"errcode"`
	Err     string `json:"error"`
}

func (e *ParxError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrCode, e.Err)
}

// InternalServerError returns a 500 Internal Server Error
// format.
func InternalServerError() util.JSONResponse {
	return util.JSONResponse{
		Code: http.StatusInternalServerError,
		JSON: Unknown("Internal Server Error"),
	}
}

// Unknown is an unexpected error
func Unknown(msg string) *ParxError {
	return &ParxError{"M_UNKNOWN", msg}
}

// Forbidden is an error when the client tries to access a resource
// they are not allowed to access.
func Forbidden(msg string) *ParxError {
	return &ParxError{"M_FORBIDDEN", msg}
}

// BadJSON is an error when the client supplies malformed JSON.
func BadJSON(msg string) *ParxError {
	return &ParxError{"M_BAD_JSON", msg}
}

// NotJSON is an error when the client supplies something that is not JSON
// to a JSON endpoint.
func NotJSON(msg string) *ParxError {
	return &ParxError{"M_NOT_JSON", msg}
}

// NotFound is an error when the client tries to access an unknown resource.
func NotFound(msg string) *ParxError {
	return &ParxError{"M_NOT_FOUND", msg}
}

// MissingArgument is an error when the client tries to access a resource
// without providing an argument that is required.
func MissingArgument(msg string) *ParxError {
	return &ParxError{"M_MISSING_ARGUMENT", msg}
}

// InvalidArgumentValue is an error when the client tries to provide an
// invalid value for a valid argument
func InvalidArgumentValue(msg string) *ParxError {
	return &ParxError{"M_INVALID_ARGUMENT_VALUE", msg}
}

// MissingToken is an error when the client tries to access a resource which
// requires authentication without supplying credentials.
func MissingToken(msg string) *ParxError {
	return &ParxError{"M_MISSING_TOKEN", msg}
}

// UnknownToken is an error when the client tries to access a resource which
// requires authentication and supplies an unrecognized token
func UnknownToken(msg string) *ParxError {
	return &ParxError{"M_UNKNOWN_TOKEN", msg}
}

// WeakPassword is an error which is returned when the client tries to register
func WeakPassword(msg string) *ParxError {
	return &ParxError{"M_WEAK_PASSWORD", msg}
}

// InvalidUsername is an error returned when the client tries to register an
// invalid username
func InvalidUsername(msg string) *ParxError {
	return &ParxError{"M_INVALID_USERNAME", msg}
}

// UserInUse is an error returned when the client tries to register an
// username that already exists
func UserInUse(msg string) *ParxError {
	return &ParxError{"M_USER_IN_USE", msg}
}

// GuestAccessForbidden is an error which is returned when the client is
// forbidden from accessing a resource as a guest.
func GuestAccessForbidden(msg string) *ParxError {
	return &ParxError{"M_GUEST_ACCESS_FORBIDDEN", msg}
}

// LimitExceededError is a rate-limiting error.
type LimitExceededError struct {
	ParxError
	RetryAfterMS int64 `json:"retry_after_ms,omitempty"`
}

// LimitExceeded is an error when the client tries to send events too quickly.
func LimitExceeded(msg string, retryAfterMS int64) *LimitExceededError {
	return &LimitExceededError{
		ParxError:  ParxError{"M_LIMIT_EXCEEDED", msg},
		RetryAfterMS: retryAfterMS,
	}
}

// NotTrusted is an error which is returned when the client asks the server to
// proxy a request to a server that isn't trusted
func NotTrusted(serverName string) *ParxError {
	return &ParxError{
		ErrCode: "M_SERVER_NOT_TRUSTED",
		Err:     fmt.Sprintf("Untrusted server '%s'", serverName),
	}
}
