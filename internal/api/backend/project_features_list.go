package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

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

		var nextState generatedapi.OptNilBool
		var nextStateTimeOpt generatedapi.OptNilDateTime
		if !nextStateTime.IsZero() {
			nextState = generatedapi.NewOptNilBool(nextStateEnabled)
			nextStateTimeOpt = generatedapi.NewOptNilDateTime(nextStateTime)
		}

		// Get feature tags
		tags, err := r.featureTagsUseCase.ListFeatureTags(ctx, it.ID)
		if err != nil {
			slog.Warn("failed to load feature tags", "error", err, "feature_id", it.ID)
			tags = []domain.Tag{} // Continue with empty tags
		}

		// Convert tags to API format
		tagsResp := make([]generatedapi.ProjectTag, 0, len(tags))
		for _, tag := range tags {
			var description generatedapi.OptNilString
			if tag.Description != nil {
				description = generatedapi.NewOptNilString(*tag.Description)
			}

			var color generatedapi.OptNilString
			if tag.Color != nil {
				color = generatedapi.NewOptNilString(*tag.Color)
			}

			var categoryID generatedapi.OptNilUUID
			if tag.CategoryID != nil {
				catID, err := uuid.Parse(tag.CategoryID.String())
				if err == nil {
					categoryID = generatedapi.NewOptNilUUID(catID)
				}
			}

			tagID, err := uuid.Parse(tag.ID.String())
			if err != nil {
				slog.Warn("invalid tag ID", "error", err, "tag_id", tag.ID)
				continue
			}

			tagsResp = append(tagsResp, generatedapi.ProjectTag{
				ID:          tagID,
				Name:        tag.Name,
				Slug:        tag.Slug,
				Description: description,
				Color:       color,
				CategoryID:  categoryID,
			})
		}

		itemsResp = append(itemsResp, generatedapi.FeatureExtended{
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
			IsActive:       r.featureProcessor.IsFeatureActive(it),
			NextState:      nextState,
			NextStateTime:  nextStateTimeOpt,
			Tags:           tagsResp,
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
