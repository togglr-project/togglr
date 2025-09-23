package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// GetFeature handles GET /api/v1/features/{feature_id}
func (r *RestAPI) GetFeature(
	ctx context.Context,
	params generatedapi.GetFeatureParams,
) (generatedapi.GetFeatureRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Load feature and check access to its project
	feature, err := r.featuresUseCase.GetExtendedByID(ctx, featureID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}
		slog.Error("get feature failed", "error", err)
		return nil, err
	}

	if err := r.permissionsService.CanAccessProject(ctx, feature.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", feature.ProjectID)

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

	// Map to API DTOs
	respVariants := make([]generatedapi.FlagVariant, 0, len(feature.FlagVariants))
	for _, v := range feature.FlagVariants {
		respVariants = append(respVariants, generatedapi.FlagVariant{
			ID:             v.ID.String(),
			FeatureID:      v.FeatureID.String(),
			Name:           v.Name,
			RolloutPercent: int(v.RolloutPercent),
		})
	}

	respRules := make([]generatedapi.Rule, 0, len(feature.Rules))
	for _, it := range feature.Rules {
		expr, err := exprToAPI(it.Conditions)
		if err != nil {
			slog.Error("build rule conditions response", "error", err)
			return nil, err
		}

		var segmentID generatedapi.OptString
		if it.SegmentID != nil {
			segmentID = generatedapi.NewOptString(it.SegmentID.String())
		}

		respRules = append(respRules, generatedapi.Rule{
			ID:            it.ID.String(),
			FeatureID:     it.FeatureID.String(),
			Conditions:    expr,
			SegmentID:     segmentID,
			IsCustomized:  it.IsCustomized,
			Action:        generatedapi.RuleAction(it.Action),
			FlagVariantID: flagVariantRef2OptString(it.FlagVariantID),
			Priority:      int(it.Priority),
			CreatedAt:     it.CreatedAt,
		})
	}

	// Get tags
	tags, err := r.featureTagsUseCase.ListFeatureTags(ctx, featureID)
	if err != nil {
		slog.Error("list feature tags failed", "error", err, "feature_id", featureID)
		return nil, err
	}

	// Convert tags to response
	respTags := make([]generatedapi.ProjectTag, len(tags))
	for i, tag := range tags {
		item := generatedapi.ProjectTag{
			ID:        uuid.MustParse(tag.ID.String()),
			ProjectID: uuid.MustParse(tag.ProjectID.String()),
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}

		if tag.CategoryID != nil {
			item.CategoryID = generatedapi.NewOptNilUUID(uuid.MustParse(tag.CategoryID.String()))
		}
		if tag.Description != nil {
			item.Description = generatedapi.NewOptNilString(*tag.Description)
		}
		if tag.Color != nil {
			item.Color = generatedapi.NewOptNilString(*tag.Color)
		}

		// Convert category
		if tag.Category != nil {
			catItem := generatedapi.Category{
				ID:        uuid.MustParse(tag.Category.ID.String()),
				Name:      tag.Category.Name,
				Slug:      tag.Category.Slug,
				Kind:      generatedapi.CategoryKind(tag.Category.Kind),
				CreatedAt: tag.Category.CreatedAt,
				UpdatedAt: tag.Category.UpdatedAt,
			}

			if tag.Category.Description != nil {
				catItem.Description = generatedapi.NewOptNilString(*tag.Category.Description)
			}
			if tag.Category.Color != nil {
				catItem.Color = generatedapi.NewOptNilString(*tag.Category.Color)
			}

			item.Category = generatedapi.NewOptCategory(catItem)
		}

		respTags[i] = item
	}

	// Get next state information
	nextStateEnabled, nextStateTime := r.featureProcessor.NextState(feature)

	var nextState generatedapi.OptNilBool
	var nextStateTimeOpt generatedapi.OptNilDateTime
	if !nextStateTime.IsZero() {
		nextState = generatedapi.NewOptNilBool(nextStateEnabled)
		nextStateTimeOpt = generatedapi.NewOptNilDateTime(nextStateTime)
	}

	resp := &generatedapi.FeatureDetailsResponse{
		Feature: generatedapi.FeatureExtended{
			ID:             feature.ID.String(),
			ProjectID:      feature.ProjectID.String(),
			Key:            feature.Key,
			Name:           feature.Name,
			Description:    generatedapi.NewOptNilString(feature.Description),
			Kind:           generatedapi.FeatureKind(feature.Kind),
			DefaultVariant: feature.DefaultVariant,
			Enabled:        feature.Enabled,
			RolloutKey:     ruleAttribute2OptString(feature.RolloutKey),
			CreatedAt:      feature.CreatedAt,
			UpdatedAt:      feature.UpdatedAt,
			IsActive:       r.featureProcessor.IsFeatureActive(feature),
			NextState:      nextState,
			NextStateTime:  nextStateTimeOpt,
		},
		Variants: respVariants,
		Rules:    respRules,
		Tags:     respTags,
	}

	return resp, nil
}
