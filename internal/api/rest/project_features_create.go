package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) CreateProjectFeature(
	ctx context.Context,
	req *generatedapi.CreateFeatureRequest,
	params generatedapi.CreateProjectFeatureParams,
) (generatedapi.CreateProjectFeatureRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check if the user can manage the project
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

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("project not found"),
			}}, nil
		}

		return nil, err
	}

	feature := domain.Feature{
		ProjectID:      projectID,
		Key:            req.Key,
		Name:           req.Name,
		Description:    req.Description.Or(""),
		Kind:           domain.FeatureKind(req.Kind),
		DefaultVariant: req.DefaultVariant,
		Enabled:        req.Enabled.Or(true),
	}

	created, err := r.featuresUseCase.Create(ctx, feature)
	if err != nil {
		slog.Error("create project feature failed", "error", err)
		return nil, err
	}

	resp := &generatedapi.FeatureResponse{Feature: generatedapi.Feature{
		ID:             created.ID.String(),
		ProjectID:      created.ProjectID.String(),
		Key:            created.Key,
		Name:           created.Name,
		Description:    generatedapi.NewOptNilString(created.Description),
		Kind:           generatedapi.FeatureKind(created.Kind),
		DefaultVariant: created.DefaultVariant,
		Enabled:        created.Enabled,
		CreatedAt:      created.CreatedAt,
		UpdatedAt:      created.UpdatedAt,
	}}

	return resp, nil
}
