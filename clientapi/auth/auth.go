// Package auth implements authentication checks and storage.
package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/clientapi/httputil"
	"github.com/bitparx/common/jsonerror"
	"github.com/bitparx/util"
)

// OWASP recommends at least 128 bits of entropy for tokens: https://www.owasp.org/index.php/Insufficient_Session-ID_Length
// 32 bytes => 256 bits
var tokenByteLength = 32

// The length of generated device IDs
var deviceIDByteLength = 6

// DeviceDatabase represents a device database.
type DeviceDatabase interface {
	// Look up the device matching the given access token.
	GetDeviceByAccessToken(ctx context.Context, token string) (*authtypes.Device, error)
}

// GenerateDeviceID creates a new device id. Returns an error if failed to generate
// random bytes.
func GenerateDeviceID() (string, error) {
	b := make([]byte, deviceIDByteLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// url-safe no padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}


// VerifyAccessToken verifies that an access token was supplied in the given HTTP request
// and returns the device it corresponds to. Returns resErr (an error response which can be
// sent to the client) if the token is invalid or there was a problem querying the database.
func VerifyAccessToken(req *http.Request, deviceDB DeviceDatabase) (device *authtypes.Device, resErr *util.JSONResponse) {
	token, err := extractAccessToken(req)
	if err != nil {
		resErr = &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: jsonerror.MissingToken(err.Error()),
		}
		return
	}
	device, err = deviceDB.GetDeviceByAccessToken(req.Context(), token)
	if err != nil {
		if err == sql.ErrNoRows {
			resErr = &util.JSONResponse{
				Code: http.StatusUnauthorized,
				JSON: jsonerror.UnknownToken("Unknown token"),
			}
		} else {
			jsonErr := httputil.LogThenError(req, err)
			resErr = &jsonErr
		}
	}
	return
}

// GenerateAccessToken creates a new access token. Returns an error if failed to generate
// random bytes.
func GenerateAccessToken() (string, error) {
	b := make([]byte, tokenByteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// url-safe no padding
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// extractAccessToken from a request, or return an error detailing what went wrong. The
// error message MUST be human-readable and comprehensible to the client.
func extractAccessToken(req *http.Request) (string, error) {
	authBearer := req.Header.Get("Authorization")
	queryToken := req.URL.Query().Get("access_token")
	if authBearer != "" && queryToken != "" {
		return "", fmt.Errorf("mixing Authorization headers and access_token query parameters")
	}

	if queryToken != "" {
		return queryToken, nil
	}

	if authBearer != "" {
		parts := strings.SplitN(authBearer, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return "", fmt.Errorf("invalid Authorization header")
		}
		return parts[1], nil
	}

	return "", fmt.Errorf("missing access token")
}
