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

// UpdateFeature handles PUT /api/v1/features/{feature_id}
func (r *RestAPI) UpdateFeature(
	ctx context.Context,
	req *generatedapi.CreateFeatureRequest,
	params generatedapi.UpdateFeatureParams,
) (generatedapi.UpdateFeatureRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	// Load feature to get project and check permissions
	existing, err := r.featuresUseCase.GetByID(ctx, featureID)
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
		ID:             featureID,
		ProjectID:      existing.ProjectID,
		Key:            req.Key,
		Name:           req.Name,
		Description:    req.Description.Or(""),
		Kind:           domain.FeatureKind(req.Kind),
		DefaultVariant: req.DefaultVariant,
		Enabled:        req.Enabled.Or(existing.Enabled),
	}

	// Build variants
	variants := make([]domain.FlagVariant, 0, len(req.Variants))
	for _, variant := range req.Variants {
		variants = append(variants, domain.FlagVariant{
			ID:             domain.FlagVariantID(variant.ID),
			FeatureID:      featureID,
			Name:           variant.Name,
			RolloutPercent: uint8(variant.RolloutPercent),
		})
	}

	// Build rules with structured conditions
	rules := make([]domain.Rule, 0, len(req.Rules))
	for _, rr := range req.Rules {
		conds := make(domain.Conditions, 0, len(rr.Conditions))
		for _, condition := range rr.Conditions {
			var val any
			if len(condition.Value) > 0 {
				if err := json.Unmarshal(condition.Value, &val); err != nil {
					slog.Error("unmarshal condition value", "error", err)
					return nil, err
				}
			}
			conds = append(conds, domain.Condition{
				Attribute: domain.RuleAttribute(condition.Attribute),
				Operator:  domain.RuleOperator(condition.Operator),
				Value:     val,
			})
		}

		rules = append(rules, domain.Rule{
			ID:            domain.RuleID(rr.ID),
			FeatureID:     featureID,
			Conditions:    conds,
			FlagVariantID: domain.FlagVariantID(rr.FlagVariantID),
			Priority:      uint8(rr.Priority.Or(0)),
		})
	}

	updated, err := r.featuresUseCase.UpdateWithChildren(ctx, feature, variants, rules)
	if err != nil {
		slog.Error("update feature with children failed", "error", err)
		return nil, err
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
		conds := make([]generatedapi.RuleCondition, 0, len(it.Conditions))
		for _, condition := range it.Conditions {
			var raw jx.Raw
			if condition.Value != nil {
				b, mErr := json.Marshal(condition.Value)
				if mErr != nil {
					slog.Error("marshal condition value", "error", mErr)
					return nil, mErr
				}
				raw = b
			}
			conds = append(conds, generatedapi.RuleCondition{
				Attribute: generatedapi.RuleAttribute(condition.Attribute),
				Operator:  generatedapi.RuleOperator(condition.Operator),
				Value:     raw,
			})
		}

		respRules = append(respRules, generatedapi.Rule{
			ID:            it.ID.String(),
			FeatureID:     it.FeatureID.String(),
			Conditions:    conds,
			FlagVariantID: it.FlagVariantID.String(),
			Priority:      int(it.Priority),
			CreatedAt:     it.CreatedAt,
		})
	}

	resp := &generatedapi.FeatureDetailsResponse{
		Feature: generatedapi.Feature{
			ID:             updated.ID.String(),
			ProjectID:      updated.ProjectID.String(),
			Key:            updated.Key,
			Name:           updated.Name,
			Description:    generatedapi.NewOptNilString(updated.Description),
			Kind:           generatedapi.FeatureKind(updated.Kind),
			DefaultVariant: updated.DefaultVariant,
			Enabled:        updated.Enabled,
			CreatedAt:      updated.CreatedAt,
			UpdatedAt:      updated.UpdatedAt,
		},
		Variants: respVariants,
		Rules:    respRules,
	}

	return resp, nil
}
