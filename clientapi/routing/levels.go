package routing

import (
	"net/http"
	"github.com/bitparx/clientapi/auth/storage/levels"
	"database/sql"
	"github.com/bitparx/util"
	"github.com/bitparx/common/jsonerror"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
)

type levelsFunction func(*http.Request, *levels.Database, bool, string, string) *util.JSONResponse

func GetAllAccountLevels(
	req *http.Request,
	levelsDB *levels.Database,
	level bool,
	localpart string,
	levelToUpdate string,
	) *util.JSONResponse {
	accounts, err := levelsDB.GetAllAccounts(req.Context())
	if err == sql.ErrNoRows {
		return &util.JSONResponse{
			Code: http.StatusNotFound,
			JSON: jsonerror.NotFound("no users exist"),
		}
	}

	if err != nil {
		return &util.JSONResponse{
			Code: 500,
			JSON: jsonerror.Unknown(err.Error()),
		}
	}

	return &util.JSONResponse{
		Code: http.StatusOK,
		JSON: accounts,
	}
}

func UpdateLevelByLocalpart(r *http.Request, levelDB *levels.Database, level bool, localpart, levelToUpdate string) *util.JSONResponse {
	var err error = nil
	switch levelToUpdate {
	case "admin": err = levelDB.UpdateLevelAdmin(r.Context(), level, localpart)
	case "moderator":  err = levelDB.UpdateLevelModerator(r.Context(), level, localpart)
	default:
		return &util.JSONResponse{
			Code: http.StatusForbidden,
			JSON: jsonerror.Forbidden("not allowed to change this level"),
		}
	}
	if err != nil {
		return &util.JSONResponse{
			Code: http.StatusNetworkAuthenticationRequired,
			JSON: jsonerror.Unknown(err.Error()),
		}
	}

	return &util.JSONResponse{
		Code: http.StatusOK,
	}
}

func RouteLevelsHandler(
	levelsDB *levels.Database,
	level bool,
	taskfunction levelsFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars:= mux.Vars(r)
		log.Println(r.URL.Path)
		res := taskfunction(r, levelsDB, level,vars["localpart"], vars["levelname"])
		err, ok := res.JSON.(*jsonerror.ParxError)
		if ok {
			http.Error(w, err.Err, res.Code)
		} else {
			json.NewEncoder(w).Encode(res)
		}
	}
}
