package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DeleteProjectSetting handles DELETE /api/v1/projects/{project_id}/settings/{setting_name}.
func (r *RestAPI) DeleteProjectSetting(
	ctx context.Context,
	params generatedapi.DeleteProjectSettingParams,
) (generatedapi.DeleteProjectSettingRes, error) {
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

	// Delete project setting
	err := r.projectSettingsUseCase.Delete(ctx, projectID, params.SettingName)
	if err != nil {
		slog.Error("delete project setting failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("setting not found"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteProjectSettingNoContent{}, nil
}
