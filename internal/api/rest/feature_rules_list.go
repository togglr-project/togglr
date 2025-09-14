package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/go-faster/jx"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// ListFeatureRules handles GET /api/v1/features/{feature_id}/rules
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
		conds := make([]generatedapi.RuleCondition, 0, len(it.Conditions))
		for _, c := range it.Conditions {
			var raw jx.Raw
			if c.Value != nil {
				b, mErr := json.Marshal(c.Value)
				if mErr != nil {
					slog.Error("marshal condition value", "error", mErr)
					return nil, mErr
				}
				raw = b
			}
			conds = append(conds, generatedapi.RuleCondition{
				Attribute: generatedapi.RuleAttribute(c.Attribute),
				Operator:  generatedapi.RuleOperator(c.Operator),
				Value:     raw,
			})
		}

		resp = append(resp, generatedapi.Rule{
			ID:            it.ID.String(),
			FeatureID:     it.FeatureID.String(),
			Conditions:    conds,
			FlagVariantID: it.FlagVariantID.String(),
			Priority:      int(it.Priority),
			CreatedAt:     it.CreatedAt,
		})
	}

	return &resp, nil
}
