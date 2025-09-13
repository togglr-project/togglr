package middlewares

import (
	"net/http"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
)

func WithRawRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixURL(r)
		ctx := etogglcontext.WithRawRequest(r.Context(), r)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func fixURL(req *http.Request) {
	if req.URL.Scheme == "" {
		if req.TLS != nil {
			req.URL.Scheme = "https"
		} else {
			req.URL.Scheme = "http"
		}
	}

	if req.URL.Host == "" {
		req.URL.Host = req.Host
	}
}
