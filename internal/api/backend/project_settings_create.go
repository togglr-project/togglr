package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// CreateProjectSetting handles POST /api/v1/projects/{project_id}/settings.
func (r *RestAPI) CreateProjectSetting(
	ctx context.Context,
	req *generatedapi.CreateProjectSettingRequest,
	params generatedapi.CreateProjectSettingParams,
) (generatedapi.CreateProjectSettingRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check access permissions for the project
	if err := r.permissionsService.CanManageProject(ctx, projectID); err != nil {
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

	// Create project setting
	setting, err := r.projectSettingsUseCase.Create(ctx, projectID, req.Name, req.Value)
	if err != nil {
		slog.Error("create project setting failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("project not found"),
			}}, nil
		}

		return nil, err
	}

	// Convert to API response
	apiSetting := dto.DomainProjectSettingToAPI(*setting)

	return &apiSetting, nil
}
