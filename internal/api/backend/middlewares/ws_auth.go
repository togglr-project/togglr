package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// WSAuthMiddleware extracts the user ID from WebSocket subprotocol and sets it in the context.
// It looks for subprotocol in format "bearer,<token>" or falls back to Authorization header.
func WSAuthMiddleware(tokenizer contract.Tokenizer, usersSrv contract.UsersUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			var token string

			// First, try to get token from query parameter (most reliable)
			token = request.URL.Query().Get("token")

			// Fallback to WebSocket subprotocol if no query token found
			if token == "" {
				if subprotocols := request.Header.Get("Sec-WebSocket-Protocol"); subprotocols != "" {
					// Note: We can't use structured logging here as we don't have access to logger
					// This is just for debugging
					_ = subprotocols

					protocols := strings.Split(subprotocols, ",")
					for _, protocol := range protocols {
						protocol = strings.TrimSpace(protocol)
						if strings.HasPrefix(protocol, "bearer.") {
							// Token is base64 encoded
							encodedToken := strings.TrimPrefix(protocol, "bearer.")
							if decoded, err := base64.StdEncoding.DecodeString(encodedToken); err == nil {
								token = string(decoded)
								break
							}
						} else if strings.HasPrefix(protocol, "bearer,") {
							// Legacy format: bearer,token
							token = strings.TrimPrefix(protocol, "bearer,")
							break
						}
					}
				}
			}

			// Fallback to Authorization header if no subprotocol or query token found
			if token == "" {
				authHeader := request.Header.Get("Authorization")
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			// If no token found, pass through without authentication
			if token == "" {
				next.ServeHTTP(writer, request)
				return
			}

			// Verify the token and get the user ID
			claims, err := tokenizer.VerifyToken(token, domain.TokenTypeAccess)
			if err != nil {
				// Invalid token, pass through without authentication
				next.ServeHTTP(writer, request)
				return
			}

			// Get the user
			user, err := usersSrv.GetByID(request.Context(), domain.UserID(claims.UserID))
			if err != nil {
				// User isn't found, pass through without authentication
				next.ServeHTTP(writer, request)
				return
			}

			// Set the user ID and superuser flag in the context
			ctx := appcontext.WithUserID(request.Context(), user.ID)
			ctx = appcontext.WithUsername(ctx, user.Username)
			ctx = appcontext.WithIsSuper(ctx, user.IsSuperuser)

			// Continue with the modified context
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
