package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// ProjectAccess middleware checks if the user has access to the project.
func ProjectAccess(permissionsService contract.PermissionsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodPost {
				next.ServeHTTP(writer, request)

				return
			}

			// Extract project ID from URL path
			// Expected format: /api/projects/{projectID}/...
			parts := strings.Split(request.URL.Path, "/")
			var projectIDStr string
			for i, part := range parts {
				if part == "projects" && i+1 < len(parts) {
					projectIDStr = parts[i+1]

					break
				}
			}

			if projectIDStr == "" || projectIDStr == "add" {
				next.ServeHTTP(writer, request)

				return
			}

			projectID := domain.ProjectID(projectIDStr)

			err := permissionsService.CanAccessProject(request.Context(), projectID)
			if err != nil {
				slog.Error("failed to check project access", "error", err, "projectID", projectID)

				switch {
				case errors.Is(err, domain.ErrEntityNotFound):
					errNotFound := generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
						Message: generatedapi.NewOptString(err.Error()),
					}}
					errNotFoundData, _ := errNotFound.MarshalJSON()
					http.Error(writer, string(errNotFoundData), http.StatusNotFound)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrPermissionDenied):
					errPermDenied := generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
						Message: generatedapi.NewOptString("permission denied"),
					}}
					errPermDeniedData, _ := errPermDenied.MarshalJSON()
					http.Error(writer, string(errPermDeniedData), http.StatusForbidden)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrUserNotFound):
					errUnauthorized := generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
						Message: generatedapi.NewOptString("unauthorized"),
					}}
					errUnauthorizedData, _ := errUnauthorized.MarshalJSON()
					http.Error(writer, string(errUnauthorizedData), http.StatusUnauthorized)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				default:
					errInternal := generatedapi.ErrorInternalServerError{
						Error: generatedapi.ErrorInternalServerErrorError{
							Message: generatedapi.NewOptString("internal server error"),
						},
					}
					errInternalData, _ := errInternal.MarshalJSON()
					http.Error(writer, string(errInternalData), http.StatusInternalServerError)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				}
			}

			// Store project ID in context for later use
			ctx := etogglcontext.WithProjectID(request.Context(), projectID)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

// ProjectManagement middleware checks if the user can manage the project.
//
//nolint:gocyclo // need refactoring
func ProjectManagement(permissionsService contract.PermissionsService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodGet {
				next.ServeHTTP(writer, request)

				return
			}

			// Extract project ID from URL path
			// Expected format: /api/projects/{projectID}/...
			parts := strings.Split(request.URL.Path, "/")
			var projectIDStr string
			for i, part := range parts {
				if part == "projects" && i+1 < len(parts) {
					projectIDStr = parts[i+1]

					break
				}
			}

			if projectIDStr == "" || projectIDStr == "add" {
				next.ServeHTTP(writer, request)

				return
			}

			projectID := domain.ProjectID(projectIDStr)

			err := permissionsService.CanManageProject(
				request.Context(),
				projectID,
			)
			if err != nil {
				slog.Error("failed to check project management permission", "error", err)

				switch {
				case errors.Is(err, domain.ErrEntityNotFound):
					errNotFound := generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
						Message: generatedapi.NewOptString(err.Error()),
					}}
					errNotFoundData, _ := errNotFound.MarshalJSON()
					http.Error(writer, string(errNotFoundData), http.StatusNotFound)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrPermissionDenied):
					errPermDenied := generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
						Message: generatedapi.NewOptString("permission denied"),
					}}
					errPermDeniedData, _ := errPermDenied.MarshalJSON()
					http.Error(writer, string(errPermDeniedData), http.StatusForbidden)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				case errors.Is(err, domain.ErrUserNotFound):
					errUnauthorized := generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
						Message: generatedapi.NewOptString("unauthorized"),
					}}
					errUnauthorizedData, _ := errUnauthorized.MarshalJSON()
					http.Error(writer, string(errUnauthorizedData), http.StatusUnauthorized)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				default:
					errInternal := generatedapi.ErrorInternalServerError{
						Error: generatedapi.ErrorInternalServerErrorError{
							Message: generatedapi.NewOptString("internal server error"),
						},
					}
					errInternalData, _ := errInternal.MarshalJSON()
					http.Error(writer, string(errInternalData), http.StatusInternalServerError)
					writer.Header().Set("Content-Type", "application/json; charset=utf-8")

					return
				}
			}

			// Store project ID in context for later use
			ctx := etogglcontext.WithProjectID(request.Context(), projectID)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}
