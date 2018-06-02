package routing

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bitparx/clientapi/auth/authtypes"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"github.com/bitparx/clientapi/httputil"
	"github.com/bitparx/common/jsonerror"
	"github.com/bitparx/util"
	"fmt"
	"strings"
)

type deviceJSON struct {
	DeviceID string `json:"device_id"`
	UserID   string `json:"user_id"`
}

type devicesJSON struct {
	Devices []deviceJSON `json:"devices"`
}

type deviceUpdateJSON struct {
	DisplayName *string `json:"display_name"`
}

// return all active sessions (devices)
func GetAllDevices(req *http.Request, deviceDB *devices.Database) *util.JSONResponse {
	dev,err := deviceDB.GetALlDevices(req.Context())
	if err!=nil {
		return &util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: httputil.LogThenError(req,err),
		}
	}

	return &util.JSONResponse{
		Code: http.StatusOK,
		JSON: dev,
	}
}

func RouteHandlerDevices(database *devices.Database) http.HandlerFunc  {
	return func(writer http.ResponseWriter, request *http.Request) {
		res := GetAllDevices(request, database)
		err, ok := res.JSON.(*jsonerror.ParxError)
		if  ok {
			http.Error(writer, err.Err, res.Code)
		} else {
			json.NewEncoder(writer).Encode(res)
		}
	}
}

// GetDeviceByID handles /device/{deviceID}
func GetDeviceByID(
	req *http.Request, deviceDB *devices.Database, device *authtypes.Device,
	deviceID string,
) util.JSONResponse {
	localpart, _, err := SplitID('@', device.UserID)
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	ctx := req.Context()
	dev, err := deviceDB.GetDeviceByID(ctx, localpart, deviceID)
	if err == sql.ErrNoRows {
		return util.JSONResponse{
			Code: http.StatusNotFound,
			JSON: jsonerror.NotFound("Unknown device"),
		}
	} else if err != nil {
		return httputil.LogThenError(req, err)
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: deviceJSON{
			DeviceID: dev.ID,
			UserID:   dev.UserID,
		},
	}
}

// GetDevicesByLocalpart handles /devices
func GetDevicesByLocalpart(
	req *http.Request, deviceDB *devices.Database, device *authtypes.Device,
) util.JSONResponse {
	localpart, _, err := SplitID('@', device.UserID)
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	ctx := req.Context()
	deviceList, err := deviceDB.GetDevicesByLocalpart(ctx, localpart)

	if err != nil {
		return httputil.LogThenError(req, err)
	}

	res := devicesJSON{}

	for _, dev := range deviceList {
		res.Devices = append(res.Devices, deviceJSON{
			DeviceID: dev.ID,
			UserID:   dev.UserID,
		})
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: res,
	}
}

// UpdateDeviceByID handles PUT on /devices/{deviceID}
func UpdateDeviceByID(
	req *http.Request, deviceDB *devices.Database, device *authtypes.Device,
	deviceID string,
) util.JSONResponse {
	if req.Method != http.MethodPut {
		return util.JSONResponse{
			Code: http.StatusMethodNotAllowed,
			JSON: jsonerror.NotFound("Bad Method"),
		}
	}

	localpart, _, err := SplitID('@', device.UserID)
	if err != nil {
		return httputil.LogThenError(req, err)
	}

	ctx := req.Context()
	dev, err := deviceDB.GetDeviceByID(ctx, localpart, deviceID)
	if err == sql.ErrNoRows {
		return util.JSONResponse{
			Code: http.StatusNotFound,
			JSON: jsonerror.NotFound("Unknown device"),
		}
	} else if err != nil {
		return httputil.LogThenError(req, err)
	}

	if dev.UserID != device.UserID {
		return util.JSONResponse{
			Code: http.StatusForbidden,
			JSON: jsonerror.Forbidden("device not owned by current user"),
		}
	}

	defer req.Body.Close() // nolint: errcheck

	payload := deviceUpdateJSON{}

	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return httputil.LogThenError(req, err)
	}

	if err := deviceDB.UpdateDevice(ctx, localpart, deviceID, payload.DisplayName); err != nil {
		return httputil.LogThenError(req, err)
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: struct{}{},
	}
}

// SplitID splits a ID into a local part and a server name.
func SplitID(sigil byte, id string) (local string, domain string, err error) {
	// IDs have the format: SIGIL LOCALPART ":" DOMAIN
	// Split on the first ":" character since the domain can contain ":"
	// characters.
	if len(id) == 0 || id[0] != sigil {
		return "", "", fmt.Errorf("invalid ID %q doesn't start with %q", id, sigil)
	}
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		// The ID must have a ":" character.
		return "", "", fmt.Errorf("invalid ID %q missing ':'", id)
	}
	return parts[0][1:], string(parts[1]), nil
}