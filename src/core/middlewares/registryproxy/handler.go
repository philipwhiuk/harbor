package registryproxy

import (
	"fmt"
	"github.com/goharbor/harbor/src/core/config"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type proxyHandler struct {
	handler http.Handler
}

func New(urls ...string) (http.Handler, error) {
	var registryURL string
	var err error
	if len(urls) > 1 {
		return nil, fmt.Errorf("the parm, urls should have only 0 or 1 elements")
	}
	if len(urls) == 0 {
		registryURL, err = config.RegistryURL()
		if err != nil {
			return nil, err
		}
	} else {
		registryURL = urls[0]
	}
	targetURL, err := url.Parse(registryURL)
	if err != nil {
		return nil, err
	}

	return &proxyHandler{
		handler: httputil.NewSingleHostReverseProxy(targetURL),
	}, nil

}

func (ph proxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ph.handler.ServeHTTP(rw, req)
}
