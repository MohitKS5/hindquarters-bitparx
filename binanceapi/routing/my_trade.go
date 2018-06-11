package routing

import (
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	"github.com/bitparx/binanceapi/utils"
	"encoding/json"
)
func getMyTrades(symbol string, recWindow int64, params ...utils.ReqParams) (res *[]rt.TradeV3, err error) {
	query := map[string]string{
		"symbol":    symbol,
		"recWindow": string(recWindow),
	}
	utils.MergeMaps(query,params[0])
	req, err := authParams.NewRequestWithSignature(BASE_URL+"/api/v3/myTrades", http.MethodGet, query)
	client := &http.Client{}
	resp,err := client.Do(req)
	defer resp.Body.Close()

	res = new([]rt.TradeV3)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}