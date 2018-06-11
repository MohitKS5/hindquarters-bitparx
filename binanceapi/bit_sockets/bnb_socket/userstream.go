package bnb_socket

import (
	"github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	"encoding/json"
)

const (
	ACCOUNT_INFO = "outboundAccountInfo"
	ORDER_UPDATE = "executionReport"
)

type key struct {
	ListenKey string `json:"listenKey"`
}

// connects to the userDataStream
func Userconnect() error {
	// get listen key
	req, err := authParams.NewRequestWithHeader("https://api.binance.com/api/v1/userDataStream", http.MethodPost, map[string]string{})
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	keys := &key{}
	err = json.NewDecoder(resp.Body).Decode(keys)
	if err != nil {
		panic(err)
	}

	// connect to the userStream
	GetMessages("/ws/" + keys.ListenKey)
	return nil
}

func Disconnectuser() (err error) {
	// delete listen key
	req, err := authParams.NewRequestWithHeader("https://api.binance.com/api/v1/userDataStream", http.MethodDelete, map[string]string{})
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	// stop go func
	disconnect()
	return
}
