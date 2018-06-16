package routing

import (
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"net/http"
	"github.com/bitparx/binanceapi/auth/authParams"
	"encoding/json"
	"log"
	"io/ioutil"
	"net/url"
)

// Do send request
func Depth(query url.Values) (res *rt.DepthResponse, err error) {
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
	res = new(rt.DepthResponse)
	res.LastUpdateID = jsondata["lastUpdateId"].(float64)
	bids := jsondata["bids"].([]interface{})
	bidsLen := len(bids)
	res.Bids = make([]rt.Bid, bidsLen)
	for i := 0; i < bidsLen; i++ {
		item := bids[i].([]interface{})
		res.Bids[i] = rt.Bid{
			Price:    item[0].(string),
			Quantity: item[1].(string),
		}
	}
	asks := jsondata["asks"].([]interface{})
	asksLen := len(asks)
	res.Asks = make([]rt.Ask, asksLen)
	for i := 0; i < asksLen; i++ {
		item := asks[i].([]interface{})
		res.Asks[i] = rt.Ask{
			Price:    item[0].(string),
			Quantity: item[1].(string),
		}
	}

	return res, nil
}
