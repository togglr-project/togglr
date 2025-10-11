package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateNotificationSetting(
	ctx context.Context,
	req *generatedapi.CreateNotificationSettingRequest,
	params generatedapi.CreateNotificationSettingParams,
) (generatedapi.CreateNotificationSettingRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())
	envKey := params.EnvironmentKey

	env, err := r.environmentsUseCase.GetByProjectIDAndKey(ctx, projectID, envKey)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		default:
			slog.Error("get environment failed", "error", err, "project_id", projectID, "env_key", envKey)

			return nil, err
		}
	}

	// Convert request to domain DTO
	settingDTO := dto.MakeNotificationSettingDTO(req, projectID, env.ID)

	// Call the service
	setting, err := r.featureNotificationsUseCase.CreateNotificationSetting(ctx, settingDTO)
	if err != nil {
		slog.Error("create notification setting failed", "error", err, "project_id", projectID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		default:
			slog.Error("create notification setting failed", "error", err, "project_id", projectID)

			return nil, err
		}
	}

	// Convert a domain model to an API model
	apiSetting := dto.DomainNotificationSettingToAPI(setting, envKey)

	return &apiSetting, nil
}
