package main

import (
	"github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	"fmt"
	"io/ioutil"
	"flag"
	"log"
	"github.com/bitparx/binanceapi/bit_sockets/bnb_socket"
)

func test() {
	query := map[string]string{
		"query_string": "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=50000000&",
	}
	req, err := authParams.NewRequestWithSignature("https://api.binance.com/api/v3/order/test", http.MethodPost, query)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body: ", string(body))
}

func testing() {
	address := flag.String("address", "localhost:4200", "http service address")
	http.HandleFunc("/", testhandle)
	http.HandleFunc("/close", closeTestHandle)
	log.Fatal(http.ListenAndServe(*address, nil))
	//routing.Depth()
}

func testhandle(w http.ResponseWriter, r *http.Request) {
	bnb_socket.Userconnect()
}

func closeTestHandle(w http.ResponseWriter, r *http.Request) {
	bnb_socket.Disconnectuser()
}
