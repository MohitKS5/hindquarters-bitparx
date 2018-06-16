package routing

import (
	"encoding/json"
	rt "github.com/bitparx/binanceapi/rest_api/response_types"
	"github.com/bitparx/binanceapi/auth/authParams"
	"net/http"
	"github.com/bitparx/util"
	"log"
	"io/ioutil"
	"net/url"
	"errors"
)

// Do send request
func getAccount(query url.Values) (res *rt.Account, err error) {
	req, err := authParams.NewRequestWithSignature(BASE_URL+"/api/v3/account", http.MethodGet, query)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(resp.Body)
		return nil, util.LogThenError(err, "request")
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return nil, errors.New(bodyString)
	}
	defer resp.Body.Close()
	res = new(rt.Account)
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, util.LogThenError(err, "decoder")
	}
	return res, nil
}
