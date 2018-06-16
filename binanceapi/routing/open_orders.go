package routing

import (
	"encoding/json"
	authp "github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"net/url"
)

// Do send request
func ListOpenOrders(query url.Values) (res *[]rt.Order, err error) {
	req, err := authp.NewRequestWithSignature(BASE_URL+"/api/v3/openOrders", http.MethodGet, query)
	resp, err := DialBnb(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	res = new([]rt.Order)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
