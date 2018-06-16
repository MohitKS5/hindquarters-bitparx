package authParams

import (
	"net/http"
	"crypto/hmac"
	"crypto/sha256"
	cfg "github.com/bitparx/common/config"
	"bytes"
	"time"
	"encoding/hex"
	"net/url"
)

// Returns a request with given url, method, querystring along with
// api key header
func NewRequestWithHeader(url, method string, query url.Values) (req *http.Request, err error) {
	//url := generateQueryString(query)
	body := query.Encode()
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
func NewRequestWithSignature(url, method string, query url.Values) (req *http.Request, err error) {

	// add timestamp parameter
	createdTimeMS := int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond)
	query.Add("timestamp", string(createdTimeMS+timeLag))

	// generate signature
	mac := hmac.New(sha256.New, []byte(cfg.SECRET_KEY))
	mac.Write([]byte(query.Encode()))
	generatedMAC := hex.EncodeToString(mac.Sum(nil))

	// generate body
	query.Add("signature", string(generatedMAC))
	body := query.Encode()

	// check method use as query param if get and body if post
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
