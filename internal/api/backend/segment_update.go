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

	// Build domain conditions
	conds := make(domain.Conditions, 0, len(req.Conditions))
	for _, c := range req.Conditions {
		var val any
		if len(c.Value) > 0 {
			if err := json.Unmarshal(c.Value, &val); err != nil {
				slog.Error("unmarshal condition value", "error", err)
				return nil, err
			}
		}
		conds = append(conds, domain.Condition{
			Attribute: domain.RuleAttribute(c.Attribute),
			Operator:  domain.RuleOperator(c.Operator),
			Value:     val,
		})
	}

	updated, err := r.segmentsUseCase.Update(ctx, domain.Segment{
		ID:          id,
		ProjectID:   current.ProjectID,
		Name:        req.Name,
		Description: req.Description.Or(""),
		Conditions:  conds,
	})
	if err != nil {
		slog.Error("update segment failed", "error", err)
		return nil, err
	}

	respConds := make([]generatedapi.RuleCondition, 0, len(updated.Conditions))
	for _, c := range updated.Conditions {
		var raw jx.Raw
		if c.Value != nil {
			b, mErr := json.Marshal(c.Value)
			if mErr != nil {
				slog.Error("marshal condition value", "error", mErr)
				return nil, mErr
			}
			raw = b
		}
		respConds = append(respConds, generatedapi.RuleCondition{
			Attribute: generatedapi.RuleAttribute(c.Attribute),
			Operator:  generatedapi.RuleOperator(c.Operator),
			Value:     raw,
		})
	}

	resp := &generatedapi.SegmentResponse{Segment: generatedapi.Segment{
		ID:          updated.ID.String(),
		ProjectID:   updated.ProjectID.String(),
		Name:        updated.Name,
		Description: generatedapi.NewOptNilString(updated.Description),
		Conditions:  respConds,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}}

	return resp, nil
}
