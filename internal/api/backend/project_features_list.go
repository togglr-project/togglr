package apibackend

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
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
		k := domain.FeatureKind(params.Kind.Value)
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

	items, total, err := r.featuresUseCase.ListExtendedByProjectIDFiltered(ctx, projectID, filter)
	if err != nil {
		slog.Error("list project features failed", "error", err)
		return nil, err
	}

	itemsResp := make([]generatedapi.FeatureExtended, 0, len(items))
	for _, it := range items {
		// Get next state information
		nextStateEnabled, nextStateTime := r.featureProcessor.NextState(it)

		// Get feature tags
		tags, err := r.featureTagsUseCase.ListFeatureTags(ctx, it.ID)
		if err != nil {
			slog.Warn("failed to load feature tags", "error", err, "feature_id", it.ID)
			tags = []domain.Tag{} // Continue with empty tags
		}

		// Create FeatureExtended with tags
		featureWithTags := it
		featureWithTags.Tags = tags

		// Get next state information
		var nextStatePtr *bool
		var nextStateTimePtr *time.Time
		if !nextStateTime.IsZero() {
			nextStatePtr = &nextStateEnabled
			nextStateTimePtr = &nextStateTime
		}

		featureExtended := dto.DomainFeatureExtendedToAPI(featureWithTags, r.featureProcessor.IsFeatureActive(it), nextStatePtr, nextStateTimePtr)
		itemsResp = append(itemsResp, featureExtended)
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
