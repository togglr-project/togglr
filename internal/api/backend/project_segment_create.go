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

// CreateProjectSegment handles POST /api/v1/projects/{project_id}/segments
func (r *RestAPI) CreateProjectSegment(
	ctx context.Context,
	req *generatedapi.CreateSegmentRequest,
	params generatedapi.CreateProjectSegmentParams,
) (generatedapi.CreateProjectSegmentRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check if the user can manage the project
	if err := r.permissionsService.CanManageProject(ctx, projectID); err != nil {
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

	segment := domain.Segment{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description.Or(""),
		Conditions:  conds,
	}

	created, err := r.segmentsUseCase.Create(ctx, segment)
	if err != nil {
		slog.Error("create segment failed", "error", err)
		return nil, err
	}

	// Build response
	respConds := make([]generatedapi.RuleCondition, 0, len(created.Conditions))
	for _, c := range created.Conditions {
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
		ID:          created.ID.String(),
		ProjectID:   created.ProjectID.String(),
		Name:        created.Name,
		Description: generatedapi.NewOptNilString(created.Description),
		Conditions:  respConds,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}}

	return resp, nil
}
