package main

import (
	"net/http"
	"github.com/bitparx/binanceapi/auth/logic"
	"fmt"
	"io/ioutil"
)

func main() {
	query := map[string]string{
		"query_string": "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=50000000&",
	}
	req, err := logic.NewRequestWithSignature("https://api.binance.com/api/v3/order/test", http.MethodPost, query)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body: ",string(body))
}
