package middlewares

import (
	"net/http"

	"github.com/google/uuid"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
)

const (
	RequestIDHeader = "X-Request-Id"
)

func RequestIDMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		reqID := request.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := etogglcontext.WithRequestID(request.Context(), reqID)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
