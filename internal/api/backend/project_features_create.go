package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) CreateProjectFeature(
	ctx context.Context,
	req *generatedapi.CreateFeatureRequest,
	params generatedapi.CreateProjectFeatureParams,
) (generatedapi.CreateProjectFeatureRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check if the user can manage the project
	if err := r.permissionsService.CanManageProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", projectID)

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

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("project not found"),
			}}, nil
		}

		return nil, err
	}

	feature := domain.Feature{
		ProjectID:      projectID,
		Key:            req.Key,
		Name:           req.Name,
		Description:    req.Description.Or(""),
		Kind:           domain.FeatureKind(req.Kind),
		DefaultVariant: req.DefaultVariant,
		Enabled:        req.Enabled.Or(true),
		RolloutKey:     domain.RuleAttribute(req.RolloutKey.Or("")),
	}

	// Build inline flag variants
	variants := make([]domain.FlagVariant, 0, len(req.Variants))
	for _, v := range req.Variants {
		variants = append(variants, domain.FlagVariant{
			ID:             domain.FlagVariantID(v.ID),
			ProjectID:      projectID,
			Name:           v.Name,
			RolloutPercent: uint8(v.RolloutPercent),
		})
	}

	// Build inline rules with structured conditions
	rules := make([]domain.Rule, 0, len(req.Rules))
	for _, rr := range req.Rules {
		expr, err := exprFromAPI(rr.Conditions)
		if err != nil {
			slog.Error("build rule conditions response", "error", err)
			return nil, err
		}

		rules = append(rules, domain.Rule{
			ID:            domain.RuleID(rr.ID),
			ProjectID:     projectID,
			Conditions:    expr,
			IsCustomized:  rr.IsCustomized,
			Action:        domain.RuleAction(rr.Action),
			FlagVariantID: optString2FlagVariantIDRef(rr.FlagVariantID),
			Priority:      uint8(rr.Priority.Or(0)),
		})
	}

	created, err := r.featuresUseCase.CreateWithChildren(ctx, feature, variants, rules)
	if err != nil {
		slog.Error("create project feature with children failed", "error", err)
		return nil, err
	}

	resp := &generatedapi.FeatureResponse{Feature: generatedapi.Feature{
		ID:             created.ID.String(),
		ProjectID:      created.ProjectID.String(),
		Key:            created.Key,
		Name:           created.Name,
		Description:    generatedapi.NewOptNilString(created.Description),
		Kind:           generatedapi.FeatureKind(created.Kind),
		DefaultVariant: created.DefaultVariant,
		Enabled:        created.Enabled,
		CreatedAt:      created.CreatedAt,
		UpdatedAt:      created.UpdatedAt,
	}}

	return resp, nil
}
