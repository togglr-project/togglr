package middlewares

import (
	"net/http"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
)

func ActorMdw(next http.Handler, actor domain.AuditActor) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := etogglcontext.WithActor(request.Context(), actor)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
