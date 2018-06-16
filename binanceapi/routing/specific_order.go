package routing

import (
	ap "github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"log"
	"encoding/json"
	"net/url"
)

func GetOrderById(query url.Values) (res *rt.Order,err error) {
	req,err := ap.NewRequestWithSignature(BASE_URL+"/api/v3/order",http.MethodGet, query)
	if err!=nil{
		log.Println(err)
		return
	}

	client:=&http.Client{}
	resp,err:=client.Do(req)
	if err!=nil || resp.StatusCode != http.StatusOK{
		log.Println("error: ",err, "\nresp",resp.Status)
		return
	}

	res = new(rt.Order)
	err = json.NewDecoder(resp.Body).Decode(res)
	return
}
