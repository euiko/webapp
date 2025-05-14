//go:build !embed

package static

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/euiko/webapp/core"
	"github.com/euiko/webapp/pkg/log"
)

func createStaticRoutes(r core.Router, s *Settings) {
	url, err := url.Parse(s.Proxy.Upstream)
	if err != nil {
		log.Fatal("invalid target", log.WithField("target", s.Proxy.Upstream))
	}

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
	})
}
