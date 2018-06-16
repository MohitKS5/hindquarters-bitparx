package routing

import (
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	"encoding/json"
	"net/url"
)

func getMyTrades(query url.Values) (res *[]rt.TradeV3, err error) {
	req, err := authParams.NewRequestWithSignature(BASE_URL+"/api/v3/myTrades", http.MethodGet, query)
	resp, err := DialBnb(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res = new([]rt.TradeV3)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
