package routing

import (
	"net/http"

	"github.com/bitparx/clientapi/auth"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/httputil"
	"github.com/bitparx/clientapi/jsonerror"
	"github.com/bitparx/util"
	"fmt"
	"log"
	"encoding/json"
	"github.com/bitparx/clientapi/auth/storage/devices"
)

type loginFlows struct {
	Flows []flow `json:"flows"`
}

type flow struct {
	Type   string   `json:"type"`
	Stages []string `json:"stages"`
}

type passwordRequest struct {
	User               string  `json:"user"`
	Password           string  `json:"password"`
	InitialDisplayName *string `json:"initial_device_display_name"`
}

type loginResponse struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	Server  string `json:"server"`
	DeviceID    string `json:"device_id"`
}

// Login implements GET and POST /login
func Login(
	req *http.Request, accountDB *accounts.Database, deviceDB *devices.Database) util.JSONResponse {
	if req.Method == http.MethodGet { // TODO: support other forms of login other than password, depending on config options
		return util.JSONResponse{
			Code: http.StatusForbidden,
			JSON: util.JSONResponse{
				Code: 403,
				JSON: "please login",
			},
		}
	} else if req.Method == http.MethodPost {
		var r passwordRequest
		resErr := httputil.UnmarshalJSONRequest(req, &r)
		if resErr != nil {
			return *resErr
		}
		if r.User == "" {
			return util.JSONResponse{
				Code: http.StatusBadRequest,
				JSON: jsonerror.BadJSON("'user' must be supplied."),
			}
		}

		fmt.Println("Processing login request")

		username, err := ParseUsernameParam(r.User)
		if err != nil {
			return util.JSONResponse{
				Code: http.StatusBadRequest,
				JSON: jsonerror.InvalidUsername(err.Error()),
			}
		}

		acc, err := accountDB.GetAccountByPassword(req.Context(), username, r.Password)
		if err != nil {
			// Technically we could tell them if the user does not exist by checking if err == sql.ErrNoRows
			// but that would leak the existence of the user.
			return util.JSONResponse{
				Code: http.StatusForbidden,
				JSON: jsonerror.Forbidden("username or password was incorrect, or the account does not exist"),
			}
		}

		token, err := auth.GenerateAccessToken()
		if err != nil {
			httputil.LogThenError(req, err)
		}

		//// TODO: Use the device ID in the request
		dev, err := deviceDB.CreateDevice(
			req.Context(), acc.Username, nil, token, r.InitialDisplayName,
		)
		if err != nil {
			return util.JSONResponse{
				Code: http.StatusInternalServerError,
				JSON: jsonerror.Unknown("failed to create device: " + err.Error()),
			}
		}

		return util.JSONResponse{
			Code: http.StatusOK,
			JSON: loginResponse{
				UserID:      acc.UserID,
				AccessToken: token,
				DeviceID:    dev.ID,
			},
		}
	}
	return util.JSONResponse{
		Code: http.StatusMethodNotAllowed,
		JSON: jsonerror.NotFound("Bad method"),
	}
}

func ParseUsernameParam(username string) (string, error) {
	// Todo regex to check invalid characters in username
	return username, nil
}

func LoginHandler( accountDB *accounts.Database, deviceDB *devices.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		res := Login(r,accountDB, deviceDB)
		json.NewEncoder(w).Encode(res);
	}
}
