package routing

import (
	"github.com/bitparx/binanceapi/rest_api/response_types"
	"net/http"
	"github.com/bitparx/binanceapi/auth/authParams"
	"encoding/json"
	"log"
	"io/ioutil"
)

const BASE_URL = "https://api.binance.com"

// Do send request
func Depth() (res *response_types.DepthResponse, err error) {
	query := map[string]string{"symbol": "BNBBTC", "limit": "5"}
	req, err := authParams.NewRequestWithHeader(BASE_URL+"/api/v1/depth", http.MethodGet, query)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return
	}
	defer resp.Body.Close()

	var jsondata map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&jsondata)
	if err != nil {
		panic(err)
	}
	res = new(response_types.DepthResponse)
	res.LastUpdateID = jsondata["lastUpdateId"].(float64)
	bids := jsondata["bids"].([]interface{})
	bidsLen := len(bids)
	res.Bids = make([]response_types.Bid, bidsLen)
	for i := 0; i < bidsLen; i++ {
		item := bids[i].([]interface{})
		res.Bids[i] = response_types.Bid{
			Price:    item[0].(string),
			Quantity: item[1].(string),
		}
	}
	asks := jsondata["asks"].([]interface{})
	asksLen := len(asks)
	res.Asks = make([]response_types.Ask, asksLen)
	for i := 0; i < asksLen; i++ {
		item := asks[i].([]interface{})
		res.Asks[i] = response_types.Ask{
			Price:    item[0].(string),
			Quantity: item[1].(string),
		}
	}

	return res, nil
}
