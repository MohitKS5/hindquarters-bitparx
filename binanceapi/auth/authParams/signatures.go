package authParams

import (
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	cfg "github.com/bitparx/common/config"
	"bytes"
	"fmt"
	"time"
	"encoding/hex"
)

// Returns a request with given url, method, querystring along with
// api key header
func NewRequestWithHeader(url, method string, query map[string]string) (req *http.Request, err error) {
	body := generateQueryString(query)
	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, url+"?"+body, nil)
		break
	default:
		req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-MBX-APIKEY", cfg.API_KEY)

	return req, nil
}

// RequestWithWithHeaders returns a request with given url, querystring, method
// along with generated signature and api-key header
func NewRequestWithSignature(url, method string, query map[string]string) (*http.Request, error) {

	querystring := generateQueryString(query)

	// add timestamp parameter
	createdTimeMS := time.Now().UnixNano() / 1000000
	querystring = fmt.Sprintf("%s&timestamp=%d", querystring, createdTimeMS)

	// generate signature
	mac := hmac.New(sha256.New, []byte(cfg.SECRET_KEY))
	mac.Write([]byte(querystring))
	generatedMAC := hex.EncodeToString(mac.Sum(nil))

	// generate body
	body := fmt.Sprintf("%s&signature=%s", querystring, string(generatedMAC))

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-MBX-APIKEY", cfg.API_KEY)

	return req, nil
}

// generate query string from query map and adds timestamp to it
func generateQueryString(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}
	querystring := query["query_string"]
	if querystring == "" {
		for key := range query {
			querystring = fmt.Sprintf("%s%s=%s&", querystring, key, query[key])
		}
	}
	if last := len(querystring) - 1; querystring[last] == '&' {
		querystring = querystring[:last]
	}
	return querystring
}
