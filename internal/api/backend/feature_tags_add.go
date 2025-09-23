package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) AddFeatureTag(
	ctx context.Context,
	req *generatedapi.AddFeatureTagRequest,
	params generatedapi.AddFeatureTagParams,
) (generatedapi.AddFeatureTagRes, error) {
	userID := appcontext.UserID(ctx)
	featureID := domain.FeatureID(params.FeatureID.String())
	tagID := domain.TagID(req.TagID.String())

	// Add tag to feature
	err := r.featureTagsUseCase.AddFeatureTag(ctx, featureID, tagID)
	if err != nil {
		slog.Error("add feature tag failed", "error", err, "user_id", userID, "feature_id", featureID, "tag_id", tagID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		// Check for the "already associated" error
		if err.Error() == "tag already associated with feature" {
			return &generatedapi.Error{Error: generatedapi.ErrorError{
				Message: generatedapi.NewOptString("tag already associated with feature"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.AddFeatureTagCreated{}, nil
}
