package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) ListSegmentDesyncFeatureIDs(
	ctx context.Context,
	params generatedapi.ListSegmentDesyncFeatureIDsParams,
) (generatedapi.ListSegmentDesyncFeatureIDsRes, error) {
	segmentID := domain.SegmentID(params.SegmentID)

	segment, err := r.segmentsUseCase.GetByID(ctx, segmentID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("segment not found"),
			}}, nil
		}

		slog.Error("get segment failed", "error", err)

		return nil, err
	}

	// Check access permissions for the owning project
	if err := r.permissionsService.CanAccessProject(ctx, segment.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", segment.ProjectID)

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

	list, err := r.segmentsUseCase.ListDesyncFeatureIDs(ctx, segmentID)
	if err != nil {
		slog.Error("list segment failed", "error", err)

		return nil, err
	}

	listResp := make(generatedapi.ListFeatureIDsResponse, 0, len(list))
	for _, item := range list {
		listResp = append(listResp, item.String())
	}

	return &listResp, nil
}
