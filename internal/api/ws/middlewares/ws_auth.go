package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// WSAuthMiddleware extracts the user ID from WebSocket subprotocol and sets it in the context.
// It looks for subprotocol in the format "bearer,<token>" or falls back to the Authorization header.
//
//nolint:gocognit // fix later
func WSAuthMiddleware(tokenizer contract.Tokenizer, usersSrv contract.UsersUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			var token string

			// First, try to get a token from WebSocket subprotocol
			if subprotocols := request.Header.Get("Sec-WebSocket-Protocol"); subprotocols != "" {
				slog.Debug("ws protocol", slog.String("protocol", subprotocols))

				protocols := strings.Split(subprotocols, ",")
				for _, protocol := range protocols {
					protocol = strings.TrimSpace(protocol)
					if strings.HasPrefix(protocol, "token.") {
						// Extract token from subprotocol
						token = strings.TrimPrefix(protocol, "token.")
						slog.Debug("extracted token from subprotocol", slog.String("token", token[:10]+"..."))

						break
					}
				}
			}

			// Fallback to the query parameter if no subprotocol token found
			if token == "" {
				token = request.URL.Query().Get("token")
				if token != "" {
					slog.Debug("extracted token from query parameter", slog.String("token", token[:10]+"..."))
				}
			}

			// Fallback to the Authorization header if no subprotocol or query token found
			if token == "" {
				authHeader := request.Header.Get("Authorization")
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			// Set Authorization header for downstream handlers
			if token != "" {
				request.Header.Set("Authorization", "Bearer "+token)
			}

			if token == "" {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}

			// Verify the token and get the user ID
			claims, err := tokenizer.VerifyToken(token, domain.TokenTypeAccess)
			if err != nil {
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}

			slog.Debug("token verified successfully", slog.Any("user_id", claims.UserID))

			// Get the user
			user, err := usersSrv.GetByID(request.Context(), domain.UserID(claims.UserID))
			if err != nil {
				slog.Info("user not found",
					slog.Any("user_id", claims.UserID), slog.String("error", err.Error()))
				writer.WriteHeader(http.StatusUnauthorized)

				return
			}

			slog.Debug("user found",
				slog.Int("user_id", int(user.ID)), slog.String("username", user.Username))

			// Set the user ID and superuser flag in the context
			ctx := appcontext.WithUserID(request.Context(), user.ID)
			ctx = appcontext.WithUsername(ctx, user.Username)
			ctx = appcontext.WithIsSuper(ctx, user.IsSuperuser)

			// Continue with the modified context
			slog.Debug("proceeding with authenticated context")
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
