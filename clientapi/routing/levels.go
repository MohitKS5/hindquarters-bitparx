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
	"github.com/gorilla/context"
)

type levelsFunction func(*http.Request, *levels.Database, sql.NullBool, string, string) *util.JSONResponse

func GetAllAccountLevels(
	req *http.Request,
	levelsDB *levels.Database,
	level sql.NullBool,
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

func SetLevelByLocalpart(r *http.Request, levelDB *levels.Database, level sql.NullBool, localpart, levelToUpdate string) *util.JSONResponse {
	if !CheckPriviliges(r, levelDB) {
		return &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: jsonerror.Forbidden("Not authorized for transaction"),
		}
	}
	return UpdateLevelByLocalpart(r, levelDB, level, localpart, levelToUpdate)
}

// sets level to null
func RequestLevelByLocalpart(r *http.Request, levelDB *levels.Database, level sql.NullBool, localpart, levelToUpdate string) *util.JSONResponse {
	ID := context.Get(r, "localpart").(string)
	localpart, _, _ = SplitID('@', ID)
	return UpdateLevelByLocalpart(r, levelDB, level, localpart, levelToUpdate)
}

// sets first user an admin
func SetFirstUserLevel(r *http.Request, levelDB *levels.Database, level sql.NullBool, localpart, levelToUpdate string) *util.JSONResponse {
	accounts,_ := levelDB.GetAllAccounts(r.Context())
	if len(accounts) != 1 {
		return &util.JSONResponse{
			Code: http.StatusForbidden,
		}
	}
	return UpdateLevelByLocalpart(r, levelDB, level, accounts[0].Username, "admin")
}


func UpdateLevelByLocalpart(r *http.Request, levelDB *levels.Database, level sql.NullBool, localpart, levelToUpdate string) *util.JSONResponse {
	var err error = nil
	switch levelToUpdate {
	case "admin":
		err = levelDB.UpdateLevelAdmin(r.Context(), level, localpart)
	case "moderator":
		err = levelDB.UpdateLevelModerator(r.Context(), level, localpart)
	default:
		return &util.JSONResponse{
			Code: http.StatusForbidden,
			JSON: jsonerror.Forbidden("not allowed to change this level"),
		}
	}
	if err != nil {
		return &util.JSONResponse{
			Code: http.StatusNetworkAuthenticationRequired,
			JSON: jsonerror.Unknown("UNAUTHORIZED ACCESS"),
		}
	}

	return &util.JSONResponse{
		Code: http.StatusOK,
	}
}

func RouteLevelsHandler(
	levelsDB *levels.Database,
	level sql.NullBool,
	taskfunction levelsFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println(r.URL.Path)
		res := taskfunction(r, levelsDB, level, vars["localpart"], vars["levelname"])
		err, ok := res.JSON.(*jsonerror.ParxError)
		if ok {
			http.Error(w, err.Err, res.Code)
		} else {
			json.NewEncoder(w).Encode(res)
		}
	}
}

func CheckPriviliges(r *http.Request, levelsDB *levels.Database) bool {
	clientID := context.Get(r, "localpart").(string)
	client, _, err := SplitID('@', clientID)
	acc, err := levelsDB.GetAccountByLocalpart(r.Context(), client)
	if err != nil {
		return false
	}

	isAdmin := acc.Access.Admin
	return isAdmin.Bool && isAdmin.Valid

}
