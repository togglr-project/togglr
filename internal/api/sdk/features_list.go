package apisdk

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/go-faster/jx"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	generatedapi "github.com/rom8726/etoggle/internal/generated/sdkserver"
)

func (s *SDKRestAPI) ListProjectFeatures(ctx context.Context) (generatedapi.ListProjectFeaturesRes, error) {
	projectID := etogglcontext.ProjectID(ctx)
	list, err := s.featuresUseCase.ListExtendedByProjectID(ctx, projectID)
	if err != nil {
		slog.Error("list features failed", "error", err)

		return nil, err
	}

	result := make(generatedapi.ListFeaturesResponse, 0, len(list))
	for _, feature := range list {
		variants := make([]generatedapi.FlagVariant, 0, len(feature.FlagVariants))
		for _, variant := range feature.FlagVariants {
			variants = append(variants, generatedapi.FlagVariant{
				ID:             variant.ID.String(),
				FeatureID:      variant.FeatureID.String(),
				Name:           variant.Name,
				RolloutPercent: uint(variant.RolloutPercent),
			})
		}

		rules := make([]generatedapi.Rule, 0, len(feature.Rules))
		for _, rule := range feature.Rules {
			conds := make([]generatedapi.RuleCondition, 0, len(rule.Conditions))
			for _, condition := range rule.Conditions {
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

			rules = append(rules, generatedapi.Rule{
				ID:            rule.ID.String(),
				FeatureID:     rule.FeatureID.String(),
				Conditions:    conds,
				FlagVariantID: rule.FlagVariantID.String(),
				Priority:      uint(rule.Priority),
			})
		}

		result = append(result, generatedapi.FeatureDetailsResponse{
			Feature: generatedapi.Feature{
				ID:             feature.ID.String(),
				Key:            feature.Key,
				Name:           feature.Name,
				Kind:           generatedapi.FeatureKind(feature.Kind),
				DefaultVariant: feature.DefaultVariant,
				Enabled:        feature.Enabled,
			},
			Variants: variants,
			Rules:    rules,
		})
	}

	return &result, nil
}
