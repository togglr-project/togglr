package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) RemoveFeatureTag(
	ctx context.Context,
	params generatedapi.RemoveFeatureTagParams,
) (generatedapi.RemoveFeatureTagRes, error) {
	userID := appcontext.UserID(ctx)
	featureID := domain.FeatureID(params.FeatureID.String())
	tagID := domain.TagID(params.TagID.String())

	// Remove tag from feature
	err := r.featureTagsUseCase.RemoveFeatureTag(ctx, featureID, tagID)
	if err != nil {
		slog.Error("remove feature tag failed", "error", err, "user_id", userID, "feature_id", featureID, "tag_id", tagID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.RemoveFeatureTagNoContent{}, nil
}
