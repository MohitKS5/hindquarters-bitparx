package routing

import (
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	ap "github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	"log"
	"encoding/json"
	"net/url"
)

func GetRecentTrades(query url.Values) (res *[]rt.Trade, err error) {
	req, err := ap.NewRequestWithHeader(BASE_URL+"/api/v1/trades", http.MethodGet, query)
	if err != nil {
		log.Println("error at request: ", err)
		return
	}

	resp, err := DialBnb(req)

	if err != nil {
		return
	}

	res = new([]rt.Trade)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}

func GetAggregateTrades(query url.Values) (res *[]rt.Trade, err error) {
	req, err := ap.NewRequestWithHeader(BASE_URL+"/api/v1/aggTrades", http.MethodGet, query)
	if err != nil {
		log.Println("error at request: ", err)
		return
	}

	resp, err := DialBnb(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	res = new([]rt.Trade)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
