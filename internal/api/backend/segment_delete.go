package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// DeleteSegment handles DELETE /api/v1/segments/{segment_id}
func (r *RestAPI) DeleteSegment(
	ctx context.Context,
	params generatedapi.DeleteSegmentParams,
) (generatedapi.DeleteSegmentRes, error) {
	id := domain.SegmentID(params.SegmentID)

	// Load segment to know project and to return 404 if it doesn't exist
	seg, err := r.segmentsUseCase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("segment not found"),
			}}, nil
		}
		slog.Error("get segment before delete failed", "error", err, "segment_id", id)
		return nil, err
	}

	// Check permissions to manage the project
	if err := r.permissionsService.CanManageProject(ctx, seg.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", seg.ProjectID)

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

	if err := r.segmentsUseCase.Delete(ctx, id); err != nil {
		slog.Error("delete segment failed", "error", err, "segment_id", id)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("segment not found"),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteSegmentNoContent{}, nil
}
