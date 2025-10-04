package apibackend

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/dto"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// GetFeature handles GET /api/v1/features/{feature_id}.
func (r *RestAPI) GetFeature(
	ctx context.Context,
	params generatedapi.GetFeatureParams,
) (generatedapi.GetFeatureRes, error) {
	featureID := domain.FeatureID(params.FeatureID)

	feature, err := r.featuresUseCase.GetExtendedByID(ctx, featureID, params.EnvironmentKey)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("feature not found"),
			}}, nil
		}

		slog.Error("get feature failed", "error", err)

		return nil, err
	}

	env, err := r.environmentsUseCase.GetByID(ctx, feature.EnvironmentID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("environment not found"),
			}}, nil
		}

		slog.Error("get environment failed", "error", err)

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
	respVariants := dto.DomainFlagVariantsToAPI(feature.FlagVariants)

	respRules, err := dto.DomainRulesToAPI(feature.Rules)
	if err != nil {
		slog.Error("build rule conditions response", "error", err)

		return nil, err
	}

	// Get tags
	tags, err := r.featureTagsUseCase.ListFeatureTags(ctx, featureID)
	if err != nil {
		slog.Error("list feature tags failed", "error", err, "feature_id", featureID)

		return nil, err
	}

	// Convert tags to response
	respTags := dto.DomainTagsToAPI(tags)

	// Get next state information
	nextStateEnabled, nextStateTime := r.featureProcessor.NextState(feature)

	// Get feature health
	health, err := r.errorReportsUseCase.GetFeatureHealth(ctx, feature.ProjectID, feature.Key, env.Key)
	if err != nil {
		slog.Error("get feature health failed", "error", err, "feature_id", featureID)

		return nil, err
	}

	// Create FeatureExtended with tags
	featureWithTags := feature
	featureWithTags.Tags = tags

	// Get next state information
	var nextStatePtr *bool

	var nextStateTimePtr *time.Time

	if !nextStateTime.IsZero() {
		nextStatePtr = &nextStateEnabled
		nextStateTimePtr = &nextStateTime
	}

	featureExtended := dto.DomainFeatureExtendedToAPI(
		featureWithTags,
		r.featureProcessor.IsFeatureActive(feature),
		nextStatePtr,
		nextStateTimePtr,
		health.Status,
	)

	resp := &generatedapi.FeatureDetailsResponse{
		Feature:  featureExtended,
		Variants: respVariants,
		Rules:    respRules,
		Tags:     respTags,
	}

	return resp, nil
}
