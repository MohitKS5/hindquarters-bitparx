package proxy_handles

import (
	"net/url"
	"net/http/httputil"
	"net/http"
)

type Prox struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxy(target string) *Prox {
	targetUrl, _ := url.Parse(target)
	return &Prox{target: targetUrl, proxy: httputil.NewSingleHostReverseProxy(targetUrl)}
}

func (p *Prox) Handle(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
