package routing

import (
	"encoding/json"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	"github.com/bitparx/util"
	"net/url"
)

// Do send request
func getAccount(query url.Values) (res *rt.Account, err error) {
	req, err := authParams.NewRequestWithSignature(BASE_URL+"/api/v3/account", http.MethodGet, query)
	resp, err := DialBnb(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res = new(rt.Account)
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, util.LogThenError(err, "decoder")
	}
	return res, nil
}
