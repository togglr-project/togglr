package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateFeatureRule(
	ctx context.Context,
	req *generatedapi.CreateRuleRequest,
	params generatedapi.CreateFeatureRuleParams,
) (generatedapi.CreateFeatureRuleRes, error) {
	featureID := domain.FeatureID(params.FeatureID)
	environmentKey := params.EnvironmentKey

	feature, err := r.featuresUseCase.GetByIDWithEnvironment(ctx, featureID, environmentKey)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature for rule create failed", "error", err)

		return nil, err
	}

	// Check permissions on the owning project
	if err := r.permissionsService.CanManageProject(ctx, feature.ProjectID); err != nil {
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

	// Convert API expression tree to domain
	expr, err := exprFromAPI(req.Conditions)
	if err != nil {
		slog.Error("parse rule conditions", "error", err)

		return nil, err
	}

	var segmentIDRef *domain.SegmentID

	if req.SegmentID.IsSet() {
		segmentID := domain.SegmentID(req.SegmentID.Value.String())
		segmentIDRef = &segmentID
	}

	rule := domain.Rule{
		ProjectID:     feature.ProjectID,
		FeatureID:     featureID,
		Conditions:    expr,
		SegmentID:     segmentIDRef,
		IsCustomized:  req.IsCustomized,
		Action:        domain.RuleAction(req.Action),
		FlagVariantID: optString2FlagVariantIDRef(req.FlagVariantID),
		Priority:      uint8(req.Priority.Or(0)),
	}

	created, err := r.rulesUseCase.Create(ctx, rule)
	if err != nil {
		slog.Error("create rule failed", "error", err)

		return nil, err
	}

	// Build response Rule with an expression tree
	respExpr, err := exprToAPI(created.Conditions)
	if err != nil {
		slog.Error("build rule conditions response", "error", err)

		return nil, err
	}

	var segmentID generatedapi.OptString
	if rule.SegmentID != nil {
		segmentID = generatedapi.NewOptString(rule.SegmentID.String())
	}

	resp := &generatedapi.RuleResponse{Rule: generatedapi.Rule{
		ID:            created.ID.String(),
		FeatureID:     created.FeatureID.String(),
		Conditions:    respExpr,
		SegmentID:     segmentID,
		IsCustomized:  created.IsCustomized,
		Action:        generatedapi.RuleAction(created.Action),
		FlagVariantID: flagVariantRef2OptString(created.FlagVariantID),
		Priority:      int(created.Priority),
		CreatedAt:     created.CreatedAt,
	}}

	return resp, nil
}
