package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) SyncCustomizedFeatureRule(
	ctx context.Context,
	params generatedapi.SyncCustomizedFeatureRuleParams,
) (generatedapi.SyncCustomizedFeatureRuleRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Load feature and check access to its project
	// Get environment key from query parameters
	environmentKey := params.EnvironmentKey

	feature, err := r.featuresUseCase.GetByIDWithEnv(ctx, featureID, environmentKey)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature failed", "error", err)

		return nil, err
	}

	if err := r.permissionsService.CanManageFeature(ctx, feature.ProjectID); err != nil {
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

	rule, err := r.rulesUseCase.SyncCustomized(ctx, domain.RuleID(params.RuleID))
	if err != nil {
		slog.Error("sync customized rule failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	// Build response Rule with expression tree
	respExpr, err := exprToAPI(rule.Conditions)
	if err != nil {
		slog.Error("build rule conditions response", "error", err)

		return nil, err
	}

	var segmentID generatedapi.OptString
	if rule.SegmentID != nil {
		segmentID = generatedapi.NewOptString(rule.SegmentID.String())
	}

	resp := &generatedapi.RuleResponse{Rule: generatedapi.Rule{
		ID:            rule.ID.String(),
		FeatureID:     rule.FeatureID.String(),
		Conditions:    respExpr,
		SegmentID:     segmentID,
		IsCustomized:  rule.IsCustomized,
		Action:        generatedapi.RuleAction(rule.Action),
		FlagVariantID: flagVariantRef2OptString(rule.FlagVariantID),
		Priority:      int(rule.Priority),
		CreatedAt:     rule.CreatedAt,
	}}

	return resp, nil
}
