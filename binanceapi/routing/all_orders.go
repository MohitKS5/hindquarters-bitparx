package routing

import (
	"net/http"
	"encoding/json"
	authp "github.com/bitparx/binanceapi/auth/authParams"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
)

func ListAllOrders(symbol, recWindow string) (res *[]rt.Order, err error) {
	query := map[string]string{
		"symbol":    symbol,
		"recWindow": recWindow,
	}
	req, err := authp.NewRequestWithSignature(BASE_URL+"/api/v3/allOrders", http.MethodGet, query)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	res = new([]rt.Order)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
