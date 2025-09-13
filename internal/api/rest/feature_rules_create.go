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

	// Build domain rule
	condBytes, _ := json.Marshal(req.Condition) // ogen type is map[string]jx.Raw; marshalling back to json
	rule := domain.Rule{
		FeatureID:     featureID,
		Condition:     json.RawMessage(condBytes),
		FlagVariantID: domain.FlagVariantID(req.FlagVariantID),
		Priority:      uint8(req.Priority.Or(0)),
	}

	created, err := r.rulesUseCase.Create(ctx, rule)
	if err != nil {
		slog.Error("create rule failed", "error", err)
		return nil, err
	}

	// Convert created.Condition (json) to generated RuleCondition
	var tmp map[string]json.RawMessage
	if err := json.Unmarshal(created.Condition, &tmp); err != nil {
		slog.Error("unmarshal created rule condition", "error", err)
		return nil, err
	}
	cond := make(generatedapi.RuleCondition, len(tmp))
	for k, v := range tmp {
		cond[k] = jx.Raw(v)
	}

	resp := &generatedapi.RuleResponse{Rule: generatedapi.Rule{
		ID:            created.ID.String(),
		FeatureID:     created.FeatureID.String(),
		Condition:     cond,
		FlagVariantID: created.FlagVariantID.String(),
		Priority:      int(created.Priority),
		CreatedAt:     created.CreatedAt,
	}}

	return resp, nil
}
