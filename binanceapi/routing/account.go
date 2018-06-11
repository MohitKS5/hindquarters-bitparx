package routing

import (
	"encoding/json"
	"github.com/bitparx/binanceapi/rest_api/response_types"
	"github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
)

// Do send request
func GetAccount() (res *response_types.Account, err error) {
	req, err := authParams.NewRequestWithSignature("/api/v3/account", http.MethodGet, map[string]string{})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	acc := new(response_types.Account)
	err = json.NewDecoder(resp.Body).Decode(&acc)
	if err != nil {
		return nil, err
	}
	return res, nil
}
