package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ListFeatureRules handles GET /api/v1/features/{feature_id}/rules.
func (r *RestAPI) ListFeatureRules(
	ctx context.Context,
	params generatedapi.ListFeatureRulesParams,
) (generatedapi.ListFeatureRulesRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Ensure feature exists and get its project to check access rights
	feature, err := r.featuresUseCase.GetByID(ctx, featureID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature for list rules failed", "error", err)

		return nil, err
	}

	// Check access permissions for the owning project
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

	items, err := r.rulesUseCase.ListByFeatureID(ctx, featureID)
	if err != nil {
		slog.Error("list rules by feature failed", "error", err)

		return nil, err
	}

	resp := make(generatedapi.ListRulesResponse, 0, len(items))

	for _, it := range items {
		expr, err := exprToAPI(it.Conditions)
		if err != nil {
			slog.Error("build rule conditions response", "error", err)

			return nil, err
		}

		var segmentID generatedapi.OptString
		if it.SegmentID != nil {
			segmentID = generatedapi.NewOptString(it.SegmentID.String())
		}

		resp = append(resp, generatedapi.Rule{
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

	return &resp, nil
}
