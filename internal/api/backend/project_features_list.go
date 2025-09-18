package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/contract"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// ListProjectFeatures handles GET /api/v1/projects/{project_id}/features
func (r *RestAPI) ListProjectFeatures(
	ctx context.Context,
	params generatedapi.ListProjectFeaturesParams,
) (generatedapi.ListProjectFeaturesRes, error) {
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

	// Build filter from query params
	filter := contract.FeaturesListFilter{SortDesc: true}
	if params.Kind.Set {
		k := domain.FeatureKind(string(params.Kind.Value))
		filter.Kind = &k
	}
	if params.Enabled.Set {
		e := params.Enabled.Value
		filter.Enabled = &e
	}
	if params.SortBy.Set {
		filter.SortBy = string(params.SortBy.Value)
	}
	if params.SortOrder.Set {
		filter.SortDesc = params.SortOrder.Value == generatedapi.SortOrderDesc
	}
	if params.TextSelector.Set {
		ts := params.TextSelector.Value
		filter.TextSelector = &ts
	}
	page := uint(1)
	perPage := uint(20)
	if params.Page.Set {
		page = params.Page.Value
	}
	if params.PerPage.Set {
		perPage = params.PerPage.Value
	}
	filter.Page = page
	filter.PerPage = perPage

	items, total, err := r.featuresUseCase.ListByProjectIDFiltered(ctx, projectID, filter)
	if err != nil {
		slog.Error("list project features failed", "error", err)
		return nil, err
	}

	itemsResp := make([]generatedapi.Feature, 0, len(items))
	for _, it := range items {
		itemsResp = append(itemsResp, generatedapi.Feature{
			ID:             it.ID.String(),
			ProjectID:      it.ProjectID.String(),
			Key:            it.Key,
			Name:           it.Name,
			Description:    generatedapi.NewOptNilString(it.Description),
			Kind:           generatedapi.FeatureKind(it.Kind),
			DefaultVariant: it.DefaultVariant,
			Enabled:        it.Enabled,
			RolloutKey:     ruleAttribute2OptString(it.RolloutKey),
			CreatedAt:      it.CreatedAt,
			UpdatedAt:      it.UpdatedAt,
		})
	}

	resp := generatedapi.ListFeaturesResponse{
		Items: itemsResp,
		Pagination: generatedapi.Pagination{
			Total:   uint(total),
			Page:    page,
			PerPage: perPage,
		},
	}

	return &resp, nil
}
