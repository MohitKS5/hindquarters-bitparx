package proxy_handles

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
	"net/http/httputil"
	"log"
)

func (p *Prox) customHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Proxy", "GoProxy")
	p.proxy.Transport = &logTransport{}
}

type logTransport struct{}

func (t *logTransport) RoundTrip(request *http.Request) (*http.Response, error) {

	buf, _ := ioutil.ReadAll(request.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	fmt.Println("Request body : ", rdr1)
	request.Body = rdr2

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		print("\n\ncame in error resp here", err)
		return nil, err //Server is not reachable. Server not working
	}

	body, err := httputil.DumpResponse(response, true)
	if err != nil {
		print("\n\nerror in dumb response")
		// copying the response body did not work
		return nil, err
	}

	log.Println("Response Body : ", string(body))
	return response, err
}
