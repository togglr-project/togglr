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

// ListProjectSegments handles GET /api/v1/projects/{project_id}/segments
func (r *RestAPI) ListProjectSegments(
	ctx context.Context,
	params generatedapi.ListProjectSegmentsParams,
) (generatedapi.ListProjectSegmentsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check access permissions for the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
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

	items, err := r.segmentsUseCase.ListByProjectID(ctx, projectID)
	if err != nil {
		slog.Error("list project segments failed", "error", err)
		return nil, err
	}

	resp := make(generatedapi.ListSegmentsResponse, 0, len(items))
	for _, it := range items {
		conds := make([]generatedapi.RuleCondition, 0, len(it.Conditions))
		for _, c := range it.Conditions {
			var raw jx.Raw
			if c.Value != nil {
				bytes, mErr := json.Marshal(c.Value)
				if mErr != nil {
					slog.Error("marshal condition value", "error", mErr)
					return nil, mErr
				}
				raw = bytes
			}
			conds = append(conds, generatedapi.RuleCondition{
				Attribute: generatedapi.RuleAttribute(c.Attribute),
				Operator:  generatedapi.RuleOperator(c.Operator),
				Value:     raw,
			})
		}

		resp = append(resp, generatedapi.Segment{
			ID:          it.ID.String(),
			ProjectID:   it.ProjectID.String(),
			Name:        it.Name,
			Description: generatedapi.NewOptNilString(it.Description),
			Conditions:  conds,
			CreatedAt:   it.CreatedAt,
			UpdatedAt:   it.UpdatedAt,
		})
	}

	return &resp, nil
}
