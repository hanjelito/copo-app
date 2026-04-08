package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func To(target string) http.HandlerFunc {
	u, _ := url.Parse(target)
	p := httputil.NewSingleHostReverseProxy(u)
	return func(w http.ResponseWriter, r *http.Request) {
		p.ServeHTTP(w, r)
	}
}
