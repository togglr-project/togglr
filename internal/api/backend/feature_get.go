package apibackend

import (
	"context"
	"errors"
	"log/slog"

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
		},
		Variants: respVariants,
		Rules:    respRules,
	}

	return resp, nil
}
