package apibackend

import (
	"context"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) RemoveFeatureTag(
	ctx context.Context,
	params generatedapi.RemoveFeatureTagParams,
) (generatedapi.RemoveFeatureTagRes, error) {
	userID := etogglcontext.UserID(ctx)
	featureID := domain.FeatureID(params.FeatureID.String())
	tagID := domain.TagID(params.TagID.String())

	// Remove tag from feature
	err := r.featureTagsUseCase.RemoveFeatureTag(ctx, featureID, tagID)
	if err != nil {
		slog.Error("remove feature tag failed", "error", err, "user_id", userID, "feature_id", featureID, "tag_id", tagID)

		if err == domain.ErrEntityNotFound {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.RemoveFeatureTagNoContent{}, nil
}
