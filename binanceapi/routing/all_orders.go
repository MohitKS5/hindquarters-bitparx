package routing

import (
	"net/http"
	"encoding/json"
	authp "github.com/bitparx/binanceapi/auth/authParams"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"net/url"
)

func ListAllOrders(query url.Values) (res *[]rt.Order, err error) {
	req, err := authp.NewRequestWithSignature(BASE_URL+"/api/v3/allOrders", http.MethodGet, query)
	resp, err := DialBnb(req)
	if err != nil {
		return nil, err
	}
	res = new([]rt.Order)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
