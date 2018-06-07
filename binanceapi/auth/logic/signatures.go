package logic

import (
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	cfg "github.com/bitparx/common/config"
	"bytes"
	"fmt"
	"time"
)

// RequestWithWithHeaders returns a post request with given url, querystring, generated signature and api-key header
func RequestWithHeaders(url string, query map[string]string) (*http.Request, error) {
	// generate query string
	querystring := query["query_string"]
	if querystring == "" {
		for key := range query {
			querystring = fmt.Sprintf("%s%s=%s&", querystring, key, query[key])
		}
	}
	createdTimeMS := time.Now().UnixNano() / 1000000
	querystring = fmt.Sprintf("%stimestamp=%d", querystring, createdTimeMS)

	// generate signature
	mac := hmac.New(sha256.New, []byte(cfg.SECRET_KEY))
	mac.Write([]byte(querystring))
	generatedMAC := mac.Sum(nil)

	// generate body
	body := fmt.Sprintf("%s&signature=%s", querystring, generatedMAC)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-MBX-APIKEY", cfg.API_KEY)

	return req, nil
}
