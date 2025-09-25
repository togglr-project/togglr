package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// GetSegment handles GET /api/v1/segments/{segment_id}.
func (r *RestAPI) GetSegment(
	ctx context.Context,
	params generatedapi.GetSegmentParams,
) (generatedapi.GetSegmentRes, error) {
	id := domain.SegmentID(params.SegmentID)

	seg, err := r.segmentsUseCase.GetByID(ctx, id)
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
	if err := r.permissionsService.CanAccessProject(ctx, seg.ProjectID); err != nil {
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

	respExpr, err := exprToAPI(seg.Conditions)
	if err != nil {
		slog.Error("build segment conditions response", "error", err)

		return nil, err
	}

	resp := &generatedapi.SegmentResponse{Segment: generatedapi.Segment{
		ID:          seg.ID.String(),
		ProjectID:   seg.ProjectID.String(),
		Name:        seg.Name,
		Description: generatedapi.NewOptNilString(seg.Description),
		Conditions:  respExpr,
		CreatedAt:   seg.CreatedAt,
		UpdatedAt:   seg.UpdatedAt,
	}}

	return resp, nil
}
