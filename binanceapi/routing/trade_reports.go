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

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("error: ", err, "\nresp", resp.Status)
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

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("error: ", err, "\nresp", resp.Status)
		return
	}

	res = new([]rt.Trade)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
