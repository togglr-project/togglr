package apibackend

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/go-faster/jx"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) CreateFeatureRule(
	ctx context.Context,
	req *generatedapi.CreateRuleRequest,
	params generatedapi.CreateFeatureRuleParams,
) (generatedapi.CreateFeatureRuleRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Ensure feature exists and get its project
	feature, err := r.featuresUseCase.GetByID(ctx, featureID)
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

	// Build domain rule from structured conditions
	conds := make(domain.Conditions, 0, len(req.Conditions))
	for _, c := range req.Conditions {
		var val any
		if len(c.Value) > 0 {
			if err := json.Unmarshal(c.Value, &val); err != nil {
				slog.Error("unmarshal condition value", "error", err)
				return nil, err
			}
		}
		conds = append(conds, domain.Condition{
			Attribute: domain.RuleAttribute(c.Attribute),
			Operator:  domain.RuleOperator(c.Operator),
			Value:     val,
		})
	}

	rule := domain.Rule{
		ProjectID:     feature.ProjectID,
		FeatureID:     featureID,
		Conditions:    conds,
		Action:        domain.RuleAction(req.Action),
		FlagVariantID: optString2FlagVariantIDRef(req.FlagVariantID),
		Priority:      uint8(req.Priority.Or(0)),
	}

	created, err := r.rulesUseCase.Create(ctx, rule)
	if err != nil {
		slog.Error("create rule failed", "error", err)
		return nil, err
	}

	// Build response Rule with structured conditions
	respConds := make([]generatedapi.RuleCondition, 0, len(created.Conditions))
	for _, c := range created.Conditions {
		var raw jx.Raw
		if c.Value != nil {
			bytes, err := json.Marshal(c.Value)
			if err != nil {
				slog.Error("marshal condition value", "error", err)
				return nil, err
			}
			raw = bytes
		}
		respConds = append(respConds, generatedapi.RuleCondition{
			Attribute: generatedapi.RuleAttribute(c.Attribute),
			Operator:  generatedapi.RuleOperator(c.Operator),
			Value:     raw,
		})
	}

	resp := &generatedapi.RuleResponse{Rule: generatedapi.Rule{
		ID:            created.ID.String(),
		FeatureID:     created.FeatureID.String(),
		Conditions:    respConds,
		Action:        generatedapi.RuleAction(created.Action),
		FlagVariantID: flagVariantRef2OptString(created.FlagVariantID),
		Priority:      int(created.Priority),
		CreatedAt:     created.CreatedAt,
	}}

	return resp, nil
}
