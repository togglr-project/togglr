package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// UpdateFeature handles PUT /api/v1/features/{feature_id}.
func (r *RestAPI) UpdateFeature(
	ctx context.Context,
	req *generatedapi.CreateFeatureRequest,
	params generatedapi.UpdateFeatureParams,
) (generatedapi.UpdateFeatureRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Get environment key from query parameters
	environmentKey := params.EnvironmentKey

	// Load feature to get a project and check permissions
	existing, err := r.featuresUseCase.GetByIDWithEnv(ctx, featureID, environmentKey)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature failed", "error", err)

		return nil, err
	}

	if err := r.permissionsService.CanManageProject(ctx, existing.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", existing.ProjectID)

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

	feature := domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          featureID,
			ProjectID:   existing.ProjectID,
			Key:         req.Key,
			Name:        req.Name,
			Description: req.Description.Or(""),
			Kind:        domain.FeatureKind(req.Kind),
			RolloutKey:  domain.RuleAttribute(req.RolloutKey.Or("")),
		},
		DefaultValue: req.DefaultValue,
		Enabled:      req.Enabled,
	}

	// Build variants
	variants := make([]domain.FlagVariant, 0, len(req.Variants))
	for _, variant := range req.Variants {
		variants = append(variants, domain.FlagVariant{
			ID:             domain.FlagVariantID(variant.ID.String()),
			ProjectID:      existing.ProjectID,
			FeatureID:      featureID,
			Name:           variant.Name,
			RolloutPercent: uint8(variant.RolloutPercent),
		})
	}

	// Build rules with structured conditions
	rules := make([]domain.Rule, 0, len(req.Rules))

	for _, rr := range req.Rules {
		expr, err := exprFromAPI(rr.Conditions)
		if err != nil {
			slog.Error("parse rule conditions", "error", err)

			return nil, err
		}

		var segmentIDRef *domain.SegmentID

		if rr.SegmentID.IsSet() {
			segmentID := domain.SegmentID(rr.SegmentID.Value.String())
			segmentIDRef = &segmentID
		}

		rules = append(rules, domain.Rule{
			ID:            domain.RuleID(rr.ID.String()),
			ProjectID:     existing.ProjectID,
			FeatureID:     featureID,
			Conditions:    expr,
			SegmentID:     segmentIDRef,
			IsCustomized:  rr.IsCustomized,
			Action:        domain.RuleAction(rr.Action),
			FlagVariantID: optString2FlagVariantIDRef(rr.FlagVariantID),
			Priority:      uint8(rr.Priority.Or(0)),
		})
	}

	updated, guardResult, err := r.featuresUseCase.UpdateWithChildren(ctx, environmentKey, feature, variants, rules)
	if err != nil {
		slog.Error("update feature with children failed", "error", err)

		return nil, err
	}

	// Handle guard result
	if guardResult.Error != nil {
		slog.Error("guard check failed", "error", guardResult.Error)

		return nil, guardResult.Error
	}

	if guardResult.ChangeConflict {
		return &generatedapi.ErrorConflict{Error: generatedapi.ErrorConflictError{
			Message: generatedapi.NewOptString("Feature is already locked by another pending change"),
		}}, nil
	}

	if guardResult.Pending {
		// Convert pending change to response
		pendingChangeResp := convertPendingChangeToResponse(guardResult.PendingChange)

		return &pendingChangeResp, nil
	}

	// Map to response DTO
	respVariants := make([]generatedapi.FlagVariant, 0, len(updated.FlagVariants))
	for _, variant := range updated.FlagVariants {
		respVariants = append(respVariants, generatedapi.FlagVariant{
			ID:             variant.ID.String(),
			FeatureID:      variant.FeatureID.String(),
			Name:           variant.Name,
			RolloutPercent: int(variant.RolloutPercent),
		})
	}

	respRules := make([]generatedapi.Rule, 0, len(updated.Rules))

	for _, it := range updated.Rules {
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
	nextStateEnabled, nextStateTime := r.featureProcessor.NextState(updated)

	var nextState generatedapi.OptNilBool

	var nextStateTimeOpt generatedapi.OptNilDateTime

	if !nextStateTime.IsZero() {
		nextState = generatedapi.NewOptNilBool(nextStateEnabled)
		nextStateTimeOpt = generatedapi.NewOptNilDateTime(nextStateTime)
	}

	resp := &generatedapi.FeatureDetailsResponse{
		Feature: generatedapi.FeatureExtended{
			ID:            updated.ID.String(),
			ProjectID:     updated.ProjectID.String(),
			Key:           updated.Key,
			RolloutKey:    ruleAttribute2OptString(updated.RolloutKey),
			Name:          updated.Name,
			Description:   generatedapi.NewOptNilString(updated.Description),
			Kind:          generatedapi.FeatureKind(updated.Kind),
			CreatedAt:     updated.CreatedAt,
			UpdatedAt:     updated.UpdatedAt,
			IsActive:      r.featureProcessor.IsFeatureActive(updated),
			NextState:     nextState,
			NextStateTime: nextStateTimeOpt,
		},
		Variants: respVariants,
		Rules:    respRules,
		Tags:     respTags,
	}

	return resp, nil
}
