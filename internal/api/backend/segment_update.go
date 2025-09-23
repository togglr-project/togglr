package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// UpdateSegment handles PUT /api/v1/segments/{segment_id}
func (r *RestAPI) UpdateSegment(
	ctx context.Context,
	req *generatedapi.UpdateSegmentRequest,
	params generatedapi.UpdateSegmentParams,
) (generatedapi.UpdateSegmentRes, error) {
	id := domain.SegmentID(params.SegmentID)

	current, err := r.segmentsUseCase.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("segment not found"),
			}}, nil
		}
		slog.Error("get segment before update failed", "error", err)
		return nil, err
	}

	// Check permissions to manage the project
	if err := r.permissionsService.CanManageProject(ctx, current.ProjectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", current.ProjectID)

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

	// Convert API expression tree to domain
	expr, err := exprFromAPI(req.Conditions)
	if err != nil {
		slog.Error("parse segment conditions", "error", err)
		return nil, err
	}

	updated, err := r.segmentsUseCase.Update(ctx, domain.Segment{
		ID:          id,
		ProjectID:   current.ProjectID,
		Name:        req.Name,
		Description: req.Description.Or(""),
		Conditions:  expr,
	})
	if err != nil {
		slog.Error("update segment failed", "error", err)
		return nil, err
	}

	respExpr, err := exprToAPI(updated.Conditions)
	if err != nil {
		slog.Error("build segment conditions response", "error", err)
		return nil, err
	}

	resp := &generatedapi.SegmentResponse{Segment: generatedapi.Segment{
		ID:          updated.ID.String(),
		ProjectID:   updated.ProjectID.String(),
		Name:        updated.Name,
		Description: generatedapi.NewOptNilString(updated.Description),
		Conditions:  respExpr,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}}

	return resp, nil
}
