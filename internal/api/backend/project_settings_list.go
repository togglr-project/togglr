package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ListProjectSettings handles GET /api/v1/projects/{project_id}/settings
func (r *RestAPI) ListProjectSettings(
	ctx context.Context,
	params generatedapi.ListProjectSettingsParams,
) (generatedapi.ListProjectSettingsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check access permissions for the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", projectID)

		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			}}, nil
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		}

		return nil, err
	}

	// Parse pagination parameters
	page := 1
	perPage := 20
	if params.Page.Set {
		page = int(params.Page.Value)
	}
	if params.PerPage.Set {
		perPage = int(params.PerPage.Value)
	}

	// Get project settings
	settings, total, err := r.projectSettingsUseCase.List(ctx, projectID, page, perPage)
	if err != nil {
		slog.Error("list project settings failed", "error", err)
		return nil, err
	}

	// Convert to API response
	itemsResp := make([]generatedapi.ProjectSetting, 0, len(settings))
	for _, setting := range settings {
		apiSetting := dto.DomainProjectSettingToAPI(*setting)
		itemsResp = append(itemsResp, apiSetting)
	}

	resp := generatedapi.ListProjectSettingsResponse{
		Data: itemsResp,
		Pagination: generatedapi.OptPagination{
			Value: generatedapi.Pagination{
				Total:   uint(total),
				Page:    uint(page),
				PerPage: uint(perPage),
			},
			Set: true,
		},
	}

	return &resp, nil
}
