package routing

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/bitparx/common/config"

	"github.com/bitparx/clientapi/auth"
	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/httputils"
	"github.com/bitparx/common/jsonerror"
	"github.com/bitparx/util"
	"log"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"github.com/bitparx/clientapi/auth/storage/levels"
	"database/sql"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 512
	maxUsernameLength = 254
	sessionIDLength   = 24
)

// sessionsDict keeps track of completed auth stages for each session.
type sessionsDict struct {
	sessions map[string][]authtypes.LoginType
}

// GetCompletedStages returns the completed stages for a session.
func (d sessionsDict) GetCompletedStages(sessionID string) []authtypes.LoginType {
	if completedStages, ok := d.sessions[sessionID]; ok {
		return completedStages
	}
	// Ensure that a empty slice is returned and not nil. See #399.
	return make([]authtypes.LoginType, 0)
}

// AAddCompletedStage records that a session has completed an auth stage.
func (d *sessionsDict) AddCompletedStage(sessionID string, stage authtypes.LoginType) {
	d.sessions[sessionID] = append(d.GetCompletedStages(sessionID), stage)
}

func newSessionsDict() *sessionsDict {
	return &sessionsDict{
		sessions: make(map[string][]authtypes.LoginType),
	}
}

var (
	// TODO: Remove old sessions. Need to do so on a session-specific timeout.
	// sessions stores the completed flow stages for all sessions. Referenced using their sessionID.
	sessions           = newSessionsDict()
	validUsernameRegex = regexp.MustCompile(`^[0-9a-z_\-./]+$`)
)

// registerRequest represents the submitted registration request.
// It can be broken down into 2 sections: the auth dictionary and registration parameters.
// Registration parameters vary depending on the request, and will need to remembered across
// sessions. If no parameters are supplied, the server should use the parameters previously
// remembered. If ANY parameters are supplied, the server should REPLACE all knowledge of
// previous parameters with the ones supplied. This mean you cannot "build up" request params.
type registerRequest struct {
	// registration parameters
	Password string `json:"password"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
	// user-interactive auth params
	Auth authDict `json:"auth"`

	InitialDisplayName *string `json:"initial_device_display_name"`
}

type authDict struct {
	Type    authtypes.LoginType `json:"type"`
	Session string              `json:"session"`
	Mac     []byte              `json:"mac"`

	// Recaptcha
	Response string `json:"response"`
	// TODO: Lots of custom keys depending on the type
}

type userInteractiveResponse struct {
	Flows     []authtypes.Flow       `json:"flows"`
	Completed []authtypes.LoginType  `json:"completed"`
	Params    map[string]interface{} `json:"params"`
	Session   string                 `json:"session"`
}

// newUserInteractiveResponse will return a struct to be sent back to the client
// during registration.
func newUserInteractiveResponse(
	sessionID string,
	fs []authtypes.Flow,
	params map[string]interface{},
) userInteractiveResponse {
	return userInteractiveResponse{
		fs, sessions.GetCompletedStages(sessionID), params, sessionID,
	}
}

type registerResponse struct {
	UserID      string           `json:"user_id"`
	AccessToken string           `json:"access_token"`
	ServerName  string           `json:"trade_server"`
	DeviceID    string           `json:"device_id"`
	Levels      authtypes.Levels `json:"accountlevels"`
}

// recaptchaResponse represents the HTTP response from a Google Recaptcha server
type recaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []int     `json:"error-codes"`
}

// validateUserName returns an error response if the username is invalid
func validateUserName(username string) *util.JSONResponse {
	if len(username) > maxUsernameLength {
		return &util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.BadJSON(fmt.Sprintf("'username' >%d characters", maxUsernameLength)),
		}
	} else if !validUsernameRegex.MatchString(username) {
		return &util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.InvalidUsername("User ID can only contain characters a-z, 0-9, or '_-./'"),
		}
	} else if username[0] == '_' { // Regex checks its not a zero length string
		return &util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.InvalidUsername("User ID can't start with a '_'"),
		}
	}
	return nil
}

// validatePassword returns an error response if the password is invalid
func validatePassword(password string) *util.JSONResponse {
	if len(password) > maxPasswordLength {
		return &util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.BadJSON(fmt.Sprintf("'password' >%d characters", maxPasswordLength)),
		}
	} else if len(password) > 0 && len(password) < minPasswordLength {
		return &util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.WeakPassword(fmt.Sprintf("password too weak: min %d chars", minPasswordLength)),
		}
	}
	return nil
}

// todo setup google recapcha api keys
//// validateRecaptcha returns an error response if the captcha response is invalid
//func validateRecaptcha(
//	cfg *config.Bitparx,
//	response string,
//	clientip string,
//) *util.JSONResponse {
//	if !cfg.Bitparx_Server.RecaptchaEnabled {
//		return &util.JSONResponse{
//			Code: http.StatusBadRequest,
//			JSON: jsonerror.BadJSON("Captcha registration is disabled"),
//		}
//	}
//
//	if response == "" {
//		return &util.JSONResponse{
//			Code: http.StatusBadRequest,
//			JSON: jsonerror.BadJSON("Captcha response is required"),
//		}
//	}
//
//	// Make a POST request to Google's API to check the captcha response
//	resp, err := http.PostForm(cfg.Bitparx_Server.RecaptchaSiteVerifyAPI,
//		url.Values{
//			"secret":   {cfg.Bitparx_Server.RecaptchaPrivateKey},
//			"response": {response},
//			"remoteip": {clientip},
//		},
//	)
//
//	if err != nil {
//		return &util.JSONResponse{
//			Code: http.StatusInternalServerError,
//			JSON: jsonerror.BadJSON("Error in requesting validation of captcha response"),
//		}
//	}
//
//	// Close the request once we're finishing reading from it
//	defer resp.Body.Close() // nolint: errcheck
//
//	// Grab the body of the response from the captcha server
//	var r recaptchaResponse
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return &util.JSONResponse{
//			Code: http.StatusInternalServerError,
//			JSON: jsonerror.BadJSON("Error in contacting captcha server" + err.Error()),
//		}
//	}
//	err = json.Unmarshal(body, &r)
//	if err != nil {
//		return &util.JSONResponse{
//			Code: http.StatusInternalServerError,
//			JSON: jsonerror.BadJSON("Error in unmarshaling captcha server's response: " + err.Error()),
//		}
//	}
//
//	// Check that we received a "success"
//	if !r.Success {
//		return &util.JSONResponse{
//			Code: http.StatusUnauthorized,
//			JSON: jsonerror.BadJSON("Invalid captcha response. Please try again."),
//		}
//	}
//	return nil
//}

// Register processes a /register request.
func Register(
	req *http.Request,
	accountDB *accounts.Database,
	deviceDB *devices.Database,
	levelDB *levels.Database,
	cfg *config.Bitparx,
) util.JSONResponse {

	if cfg.Bitparx_Server.RegistrationDisabled {
		return util.JSONResponse{
			Code: http.StatusForbidden,
			JSON: jsonerror.Forbidden("Signup has been disabled by admin"),
		}
	}

	var r registerRequest
	resErr := httputils.UnmarshalJSONRequest(req, &r)
	if resErr != nil {
		return *resErr
	}

	// Retrieve or generate the sessionID
	sessionID := r.Auth.Session
	if sessionID == "" {
		// Generate a new, random session ID
		sessionID = util.RandomString(sessionIDLength)
	}

	// If no auth type is specified by the client, send back the list of available flows
	if r.Auth.Type == "" {
		return util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: newUserInteractiveResponse(sessionID,
				cfg.Derived.Registration.Flows, cfg.Derived.Registration.Params),
		}
	}

	// Squash username to all lowercase letters
	r.Username = strings.ToLower(r.Username)

	if resErr = validateUserName(r.Username); resErr != nil {
		return *resErr
	}
	if resErr = validatePassword(r.Password); resErr != nil {
		return *resErr
	}

	// todo implement logrus
	//logger := util.GetLogger(req.Context())
	//logger.WithFields(log.Fields{
	//	"username":   r.Username,
	//	"auth.type":  r.Auth.Type,
	//	"session_id": r.Auth.Session,
	//}).Info("Processing registration request")

	return handleRegistrationFlow(req, r, sessionID, cfg, accountDB, deviceDB, levelDB)
}

// handleRegistrationFlow will direct and complete registration flow stages
// that the client has requested.
func handleRegistrationFlow(
	req *http.Request,
	r registerRequest,
	sessionID string,
	cfg *config.Bitparx,
	accountDB *accounts.Database,
	deviceDB *devices.Database,
	levelDB *levels.Database,
) util.JSONResponse {
	// TODO: Shared secret registration (create user scripts)
	// TODO: Enable registration config flag
	// TODO: Guest account upgrading

	// TODO: Handle loading of previous session parameters from database.
	// TODO: Handle mapping registrationRequest parameters into session parameters

	// TODO: email / msisdn auth types.

	switch r.Auth.Type {
	//case authtypes.LoginTypeRecaptcha:
	//	// Check given captcha response
	//	resErr := validateRecaptcha(cfg, r.Auth.Response, req.RemoteAddr)
	//	if resErr != nil {
	//		return *resErr
	//	}
	//
	//	// Add Recaptcha to the list of completed registration stages
	//	sessions.AddCompletedStage(sessionID, authtypes.LoginTypeRecaptcha)

	case authtypes.LoginTypeSharedSecret:
		// Check shared secret against config
		valid, err := isValidMacLogin(cfg, r.Username, r.Password, r.Admin, r.Auth.Mac)

		if err != nil {
			return httputils.LogThenError(req, err)
		} else if !valid {
			return util.MessageResponse(http.StatusForbidden, "HMAC incorrect")
		}

		// Add SharedSecret to the list of completed registration stages
		sessions.AddCompletedStage(sessionID, authtypes.LoginTypeSharedSecret)

	case authtypes.LoginTypeDummy:
		// there is nothing to do
		// Add Dummy to the list of completed registration stages
		sessions.AddCompletedStage(sessionID, authtypes.LoginTypeDummy)
		println("reached dummy")

	default:
		return util.JSONResponse{
			Code: http.StatusNotImplemented,
			JSON: jsonerror.Unknown("unknown/unimplemented auth type"),
		}
	}

	// Check if the user's registration flow has been completed successfully
	// A response with current registration flow and remaining available methods
	// will be returned if a flow has not been successfully completed yet
	return checkAndCompleteFlow(sessions.GetCompletedStages(sessionID),
		req, r, sessionID, cfg, accountDB, deviceDB, levelDB)
}

// checkAndCompleteFlow checks if a given registration flow is completed given
// a set of allowed flows. If so, registration is completed, otherwise a
// response with
func checkAndCompleteFlow(
	flow []authtypes.LoginType,
	req *http.Request,
	r registerRequest,
	sessionID string,
	cfg *config.Bitparx,
	accountDB *accounts.Database,
	deviceDB *devices.Database,
	levelDB *levels.Database,
) util.JSONResponse {
	if checkFlowCompleted(flow, cfg.Derived.Registration.Flows) {
		// This flow was completed, registration can continue
		println("flow completed")
		return completeRegistration(req.Context(), accountDB, deviceDB, levelDB,
			r.Username, r.Password, r.InitialDisplayName)
	}

	// There are still more stages to complete.
	// Return the flows and those that have been completed.
	return util.JSONResponse{
		Code: http.StatusUnauthorized,
		JSON: newUserInteractiveResponse(sessionID,
			cfg.Derived.Registration.Flows, cfg.Derived.Registration.Params),
	}
}

func completeRegistration(
	ctx context.Context,
	accountDB *accounts.Database,
	deviceDB *devices.Database,
	levelDB *levels.Database,
	username, password string,
	displayName *string,
) util.JSONResponse {
	if username == "" {
		return util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.BadJSON("missing username"),
		}
	}
	// Blank passwords are only allowed by registered application services
	if password == "" {
		return util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.BadJSON("missing password"),
		}
	}

	// todo handle errors properly, dont send sql errors to server
	err := levelDB.CreateLevel(ctx, username)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: jsonerror.Unknown("failed to create levels: " + err.Error()),
		}
	}

	acc, err := accountDB.CreateAccount(ctx, username, password)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: jsonerror.Unknown("failed to create account: " + err.Error()),
		}
	} else if acc == nil {
		return util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.UserInUse("Desired user ID is already taken."),
		}
	}

	token, err := auth.GenerateAccessToken()
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: jsonerror.Unknown("Failed to generate access token"),
		}
	}

	// // TODO: Use the device ID in the request.
	dev, err := deviceDB.CreateDevice(ctx, username, nil, token, displayName)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: jsonerror.Unknown("failed to create device: " + err.Error()),
		}
	}

	falseNullBool := sql.NullBool{false, true}
	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: registerResponse{
			UserID:      acc.UserID,
			AccessToken: token,
			ServerName:  acc.ServerName,
			Levels:      authtypes.Levels{falseNullBool, falseNullBool},
			DeviceID:    dev.ID,
		},
	}
}

// Used for shared secret registration.
// todo proper implementation
// Checks if the username, password and isAdmin flag matches the given mac.
func isValidMacLogin(
	cfg *config.Bitparx,
	username, password string,
	isAdmin bool,
	givenMac []byte,
) (bool, error) {
	sharedSecret := cfg.Bitparx_Server.RegistrationSharedSecret

	// Check that shared secret registration isn't disabled.
	if cfg.Bitparx_Server.RegistrationSharedSecret == "" {
		return false, errors.New("Shared secret registration is disabled")
	}

	// Double check that username/password don't contain the HMAC delimiters. We should have
	// already checked this.
	if strings.Contains(username, "\x00") {
		return false, errors.New("Username contains invalid character")
	}
	if strings.Contains(password, "\x00") {
		return false, errors.New("Password contains invalid character")
	}
	if sharedSecret == "" {
		return false, errors.New("Shared secret registration is disabled")
	}

	adminString := "notadmin"
	if isAdmin {
		adminString = "admin"
	}
	joined := strings.Join([]string{username, password, adminString}, "\x00")

	mac := hmac.New(sha1.New, []byte(sharedSecret))
	_, err := mac.Write([]byte(joined))
	if err != nil {
		return false, err
	}
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(givenMac, expectedMAC), nil
}

// checkFlows checks a single completed flow against another required one. If
// one contains at least all of the stages that the other does, checkFlows
// returns true.
func checkFlows(
	completedStages []authtypes.LoginType,
	requiredStages []authtypes.LoginType,
) bool {
	// Create temporary slices so they originals will not be modified on sorting
	completed := make([]authtypes.LoginType, len(completedStages))
	required := make([]authtypes.LoginType, len(requiredStages))
	copy(completed, completedStages)
	copy(required, requiredStages)

	// Sort the slices for simple comparison
	sort.Slice(completed, func(i, j int) bool { return completed[i] < completed[j] })
	sort.Slice(required, func(i, j int) bool { return required[i] < required[j] })

	// Iterate through each slice, going to the next required slice only once
	// we've found a match.
	i, j := 0, 0
	for j < len(required) {
		// Exit if we've reached the end of our input without being able to
		// match all of the required stages.
		if i >= len(completed) {
			return false
		}

		// If we've found a stage we want, move on to the next required stage.
		if completed[i] == required[j] {
			j++
		}
		i++
	}
	return true
}

// checkFlowCompleted checks if a registration flow complies with any allowed flow
// dictated by the server. Order of stages does not matter. A user may complete
// extra stages as long as the required stages of at least one flow is met.
func checkFlowCompleted(
	flow []authtypes.LoginType,
	allowedFlows []authtypes.Flow,
) bool {
	// Iterate through possible flows to check whether any have been fully completed.
	for _, allowedFlow := range allowedFlows {
		if checkFlows(flow, allowedFlow.Stages) {
			return true
		}
	}
	return false
}

type availableResponse struct {
	Available bool `json:"available"`
}

// RegisterAvailable checks if the username is already taken or invalid.
func RegisterAvailable(
	req *http.Request,
	accountDB *accounts.Database,
) util.JSONResponse {
	username := req.URL.Query().Get("username")

	// Squash username to all lowercase letters
	username = strings.ToLower(username)

	if err := validateUserName(username); err != nil {
		return *err
	}

	availability, availabilityErr := accountDB.CheckAccountAvailability(req.Context(), username)
	if availabilityErr != nil {
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: jsonerror.Unknown("failed to check availability: " + availabilityErr.Error()),
		}
	}
	if !availability {
		return util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: jsonerror.InvalidUsername("A different user ID has already been registered for this session"),
		}
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: availableResponse{
			Available: true,
		},
	}
}

// control registration
var cfg = config.LoadConfig()

func DisableRegistration(r *http.Request, set bool, levelDB *levels.Database) util.JSONResponse {
	cfg.Bitparx_Server.RegistrationDisabled = set
	return util.JSONResponse{
		Code: http.StatusOK,
	}
}

func RegisterHandler(accountDB *accounts.Database, deviceDB *devices.Database, levelDB *levels.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		Register(r, accountDB, deviceDB, levelDB, cfg).Encode(&w)
	}
}

func RegistrationHandler(levelDB *levels.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.Path)
		if !CheckAdmin(request, levelDB) {
			util.JSONResponse{
				Code: http.StatusUnauthorized,
			}.Encode(&writer)
		}
		switch request.Method {
		case http.MethodGet:
			util.JSONResponse{
				Code: http.StatusOK,
				JSON: cfg.Bitparx_Server.RegistrationDisabled,
			}.Encode(&writer)
			break
		case http.MethodPut:
			DisableRegistration(request, false, levelDB).Encode(&writer)
			break
		case http.MethodDelete:
			DisableRegistration(request, true, levelDB).Encode(&writer)
			break
		default:
			http.Error(writer, "Bad Method", http.StatusBadRequest)
		}
	}
}
