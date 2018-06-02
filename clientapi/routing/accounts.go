package routing

import (
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"net/http"
	"github.com/bitparx/util"
	"github.com/bitparx/clientapi/httputil"
	"github.com/bitparx/common/jsonerror"
	"encoding/json"
)

// get all accounts
func GetAllAccounts(req *http.Request, accountDB *accounts.Database) *util.JSONResponse {
	acc, err := accountDB.GetAllAccounts(req.Context())
	if err != nil {
		return &util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: httputil.LogThenError(req, err),
		}
	}
	return &util.JSONResponse{
		Code: http.StatusOK,
		JSON: acc,
	}
}

func RouteHandlerAccounts(accountDB *accounts.Database)  http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {
		res := GetAllAccounts(request, accountDB)
		err, ok := res.JSON.(*jsonerror.ParxError)
		if  ok {
			http.Error(writer, err.Err, res.Code)
		} else {
			json.NewEncoder(writer).Encode(res)
		}
	}
}