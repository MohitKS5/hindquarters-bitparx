package routing

import (
	"encoding/json"
	authp "github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
)

// Do send request
func ListOpenOrders(symbol, recWindow string) (res *[]rt.Order, err error) {
	query := map[string]string{"recWindow": recWindow}

	// add symbol if provided
	if symbol != "" {
		query["symbol"] = symbol
	}
	req, err := authp.NewRequestWithSignature(BASE_URL+"/api/v3/openOrders", http.MethodGet, query)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	res = new([]rt.Order)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
