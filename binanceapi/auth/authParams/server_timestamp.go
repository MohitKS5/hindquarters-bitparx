package authParams

import (
	"net/http"
	"encoding/json"
	"time"
	cfg "github.com/bitparx/common/config"
)

type server struct {
	ServerTime int64
}

var timeLag = CalcTimeLag(cfg.BINANCE_REST_URL)

func CalcTimeLag(url string) time.Duration {
	resp, err := http.Get(url + "/api/v1/time")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	res := new(server)
	err = json.NewDecoder(resp.Body).Decode(res)
	serverTime := time.Unix(res.ServerTime, 0)

	// add timestamp parameter
	myTimeMS := time.Unix(int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond), 0)
	return serverTime.Sub(myTimeMS)
}
