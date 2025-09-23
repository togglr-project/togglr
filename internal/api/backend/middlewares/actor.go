package middlewares

import (
	"net/http"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
)

func ActorMdw(next http.Handler, actor domain.AuditActor) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := appcontext.WithActor(request.Context(), actor)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
