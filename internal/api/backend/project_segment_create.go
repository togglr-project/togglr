package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// CreateProjectSegment handles POST /api/v1/projects/{project_id}/segments.
func (r *RestAPI) CreateProjectSegment(
	ctx context.Context,
	req *generatedapi.CreateSegmentRequest,
	params generatedapi.CreateProjectSegmentParams,
) (generatedapi.CreateProjectSegmentRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check if the user can manage segments
	if err := r.permissionsService.CanManageSegment(ctx, projectID); err != nil {
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

	// Convert API expression tree to domain
	expr, err := exprFromAPI(req.Conditions)
	if err != nil {
		slog.Error("parse segment conditions", "error", err)

		return nil, err
	}

	segment := domain.Segment{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description.Or(""),
		Conditions:  expr,
	}

	created, err := r.segmentsUseCase.Create(ctx, segment)
	if err != nil {
		slog.Error("create segment failed", "error", err)

		return nil, err
	}

	// Build response
	exprOut, err := exprToAPI(created.Conditions)
	if err != nil {
		slog.Error("build segment conditions response", "error", err)

		return nil, err
	}

	resp := &generatedapi.SegmentResponse{Segment: generatedapi.Segment{
		ID:          created.ID.String(),
		ProjectID:   created.ProjectID.String(),
		Name:        created.Name,
		Description: generatedapi.NewOptNilString(created.Description),
		Conditions:  exprOut,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}}

	return resp, nil
}
