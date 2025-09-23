package apibackend

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) AddFeatureTag(
	ctx context.Context,
	req *generatedapi.AddFeatureTagRequest,
	params generatedapi.AddFeatureTagParams,
) (generatedapi.AddFeatureTagRes, error) {
	userID := etogglcontext.UserID(ctx)
	featureID := domain.FeatureID(params.FeatureID.String())
	tagID := domain.TagID(req.TagID.String())

	// Add tag to feature
	err := r.featureTagsUseCase.AddFeatureTag(ctx, featureID, tagID)
	if err != nil {
		slog.Error("add feature tag failed", "error", err, "user_id", userID, "feature_id", featureID, "tag_id", tagID)

		if err == domain.ErrEntityNotFound {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		// Check for "already associated" error
		if err.Error() == "tag already associated with feature" {
			return &generatedapi.Error{Error: generatedapi.ErrorError{
				Message: generatedapi.NewOptString("tag already associated with feature"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.AddFeatureTagCreated{}, nil
}
