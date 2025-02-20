//go:build !embed

package webapp

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/settings"
	"github.com/go-chi/chi/v5"
)

func createStaticRoutes(r chi.Router, s *settings.StaticServer) {
	url, err := url.Parse(s.Proxy.Upstream)
	if err != nil {
		log.Fatal("invalid target", log.WithField("target", s.Proxy.Upstream))
	}

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		httputil.NewSingleHostReverseProxy(url).ServeHTTP(w, r)
	})
}
