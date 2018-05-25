package routing

import (
	"net/http"

	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"github.com/bitparx/clientapi/httputil"
	"github.com/bitparx/clientapi/jsonerror"
	"github.com/bitparx/util"
	"github.com/bitparx/clientapi/auth"
	"encoding/json"
)

// Logout handles POST /logout, auth = true type endpoint
func Logout(
	req *http.Request, deviceDB *devices.Database, device *authtypes.Device,
) util.JSONResponse {
	if req.Method != http.MethodPost {
		return util.JSONResponse{
			Code: http.StatusMethodNotAllowed,
			JSON: jsonerror.NotFound("Bad method"),
		}
	}

	username, _, err := SplitID('@', device.UserID)
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	if err := deviceDB.RemoveDevice(req.Context(), device.ID, username); err != nil {
		return httputil.LogThenError(req, err)
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: struct{}{},
	}
}

// LogoutAll handles POST /logout/all
func LogoutAll(
	req *http.Request, deviceDB *devices.Database, device *authtypes.Device,
) util.JSONResponse {
	username, _, err := SplitID('@', device.UserID)
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	if err := deviceDB.RemoveAllDevices(req.Context(), username); err != nil {
		return httputil.LogThenError(req, err)
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: struct{}{},
	}
}

// LogoutHandler for router
func LogoutHandler(deviceDB *devices.Database) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		dev, err := auth.VerifyAccessToken(request, deviceDB)
		if err!=nil{
			json.NewEncoder(writer).Encode(err)
			return
		}
		Logout(request, deviceDB, dev)
		json.NewEncoder(writer).Encode(http.StatusOK)
	}
}
