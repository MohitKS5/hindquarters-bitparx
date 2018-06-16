package authParams


import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/bitparx/binanceapi/routing"
)

type server struct {
	ServerTime int64
}

var timeLag = CalcTimeLag(routing.BASE_URL)

func CalcTimeLag(url string) int64 {
	resp, err := http.Get(url + "/api/v1/time")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	res := new(server)
	err = json.NewDecoder(resp.Body).Decode(res)

	// add timestamp parameter
	myTimeMS := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
	return res.ServerTime - myTimeMS
}
