package routing

import (
	"net/http"
	"github.com/bitparx/clientapi/auth/storage/levels"
	"database/sql"
	"github.com/bitparx/util"
	"github.com/bitparx/common/jsonerror"
	"log"
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

//SetLevelByLocalpart sets the level after checking if user sending the request is admin
func SetLevelByLocalpart(r *http.Request, levelDB *levels.Database, level sql.NullBool, localpart, levelToUpdate string) *util.JSONResponse {
	if !CheckAdmin(r, levelDB) {
		return &util.JSONResponse{
			Code: http.StatusUnauthorized,
			JSON: jsonerror.Forbidden("Not authorized for transaction"),
		}
	}
	return updateLevelByLocalpart(r, levelDB, level, localpart, levelToUpdate)
}

// sets level to null
func RequestLevelByLocalpart(r *http.Request, levelDB *levels.Database, level sql.NullBool, localpart, levelToUpdate string) *util.JSONResponse {
	ID := context.Get(r, "localpart").(string)
	localpart, _, _ = SplitID('@', ID)
	return updateLevelByLocalpart(r, levelDB, level, localpart, levelToUpdate)
}

// generic function to update levels
func updateLevelByLocalpart(r *http.Request, levelDB *levels.Database, level sql.NullBool, localpart, levelToUpdate string) *util.JSONResponse {
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
		taskfunction(r, levelsDB, level, vars["localpart"], vars["levelname"]).Encode(&w)
	}
}

func CheckAdmin(r *http.Request, levelsDB *levels.Database) bool {
	clientID := context.Get(r, "localpart").(string)
	client, _, err := SplitID('@', clientID)
	acc, err := levelsDB.GetAccountByLocalpart(r.Context(), client)
	if err != nil {
		return false
	}

	isAdmin := acc.Access.Admin
	return isAdmin.Bool && isAdmin.Valid
}
