package apibackend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
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

	// Build filter from query params
	filter := contract.SegmentsListFilter{SortDesc: true}
	if params.TextSelector.Set {
		ts := params.TextSelector.Value
		filter.TextSelector = &ts
	}
	page := uint(1)
	perPage := uint(20)
	if params.Page.Set {
		page = params.Page.Value
	}
	if params.PerPage.Set {
		perPage = params.PerPage.Value
	}
	filter.Page = page
	filter.PerPage = perPage
	if params.SortBy.Set {
		filter.SortBy = string(params.SortBy.Value)
	}
	if params.SortOrder.Set {
		filter.SortDesc = params.SortOrder.Value == generatedapi.SortOrderDesc
	}

	items, total, err := r.segmentsUseCase.ListByProjectIDFiltered(ctx, projectID, filter)
	if err != nil {
		slog.Error("list project segments failed", "error", err)
		return nil, err
	}

	itemsResp := make([]generatedapi.Segment, 0, len(items))
	for _, it := range items {
		expr, err := exprToAPI(it.Conditions)
		if err != nil {
			slog.Error("build segment conditions response", "error", err)
			return nil, err
		}
		itemsResp = append(itemsResp, generatedapi.Segment{
			ID:          it.ID.String(),
			ProjectID:   it.ProjectID.String(),
			Name:        it.Name,
			Description: generatedapi.NewOptNilString(it.Description),
			Conditions:  expr,
			CreatedAt:   it.CreatedAt,
			UpdatedAt:   it.UpdatedAt,
		})
	}

	resp := generatedapi.ListSegmentsResponse{
		Items: itemsResp,
		Pagination: generatedapi.Pagination{
			Total:   uint(total),
			Page:    page,
			PerPage: perPage,
		},
	}

	return &resp, nil
}
