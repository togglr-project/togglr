package apibackend

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/go-faster/jx"

	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

// GetSegment handles GET /api/v1/segments/{segment_id}
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

	conds := make([]generatedapi.RuleCondition, 0, len(seg.Conditions))
	for _, c := range seg.Conditions {
		var raw jx.Raw
		if c.Value != nil {
			b, mErr := json.Marshal(c.Value)
			if mErr != nil {
				slog.Error("marshal condition value", "error", mErr)
				return nil, mErr
			}
			raw = b
		}
		conds = append(conds, generatedapi.RuleCondition{
			Attribute: generatedapi.RuleAttribute(c.Attribute),
			Operator:  generatedapi.RuleOperator(c.Operator),
			Value:     raw,
		})
	}

	resp := &generatedapi.SegmentResponse{Segment: generatedapi.Segment{
		ID:          seg.ID.String(),
		ProjectID:   seg.ProjectID.String(),
		Name:        seg.Name,
		Description: generatedapi.NewOptNilString(seg.Description),
		Conditions:  conds,
		CreatedAt:   seg.CreatedAt,
		UpdatedAt:   seg.UpdatedAt,
	}}

	return resp, nil
}
