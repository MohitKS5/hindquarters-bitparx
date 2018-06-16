package routing

import (
	"net/http"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"encoding/json"
	"net/url"
)

func GetExchangeInfo(query url.Values) (exInfo *rt.ExchangeInfo, err error) {
	resp, err := http.Get(BASE_URL + "/api/v1/exchangeInfo")
	exInfo = new(rt.ExchangeInfo)
	err = json.NewDecoder(resp.Body).Decode(&exInfo)
	if err != nil {
		return nil, err
	}
	return
}
