package apibackend

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// UpdateProjectSetting handles PUT /api/v1/projects/{project_id}/settings/{setting_name}.
func (r *RestAPI) UpdateProjectSetting(
	ctx context.Context,
	req *generatedapi.UpdateProjectSettingRequest,
	params generatedapi.UpdateProjectSettingParams,
) (generatedapi.UpdateProjectSettingRes, error) {
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

	var value any
	_ = json.Unmarshal([]byte(req.Value), &value)

	// Update project setting
	setting, err := r.projectSettingsUseCase.Update(ctx, projectID, params.SettingName, value)
	if err != nil {
		slog.Error("update project setting failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("setting not found"),
			}}, nil
		}

		return nil, err
	}

	// Convert to API response
	apiSetting := dto.DomainProjectSettingToAPI(*setting)

	return &apiSetting, nil
}
