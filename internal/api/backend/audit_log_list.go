package apibackend

import (
	"context"
	"log/slog"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:nilerr // it's ok here
func (r *RestAPI) ListProjectAuditLogs(
	ctx context.Context,
	params generatedapi.ListProjectAuditLogsParams,
) (generatedapi.ListProjectAuditLogsRes, error) {
	projectID := domain.ProjectID(params.ProjectID.String())

	// Permission check
	if err := r.permissionsService.CanViewAudit(ctx, projectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	// Build filter
	filter := contract.AuditLogListFilter{
		ProjectID: projectID,
		SortBy:    "created_at",
		SortDesc:  true,
		Page:      1,
		PerPage:   20,
	}

	if params.EnvironmentKey.Set {
		filter.EnvironmentKey = &params.EnvironmentKey.Value
	}
	if params.Entity.Set {
		et := domain.EntityType(params.Entity.Value)
		filter.Entity = &et
	}
	if params.EntityID.Set {
		v := params.EntityID.Value.String()
		filter.EntityID = &v
	}
	if params.Actor.Set {
		filter.Actor = &params.Actor.Value
	}
	if params.From.Set {
		from := params.From.Value
		filter.From = &from
	}
	if params.To.Set {
		to := params.To.Value
		filter.To = &to
	}
	if params.SortBy.Set {
		filter.SortBy = string(params.SortBy.Value)
	}
	if params.SortOrder.Set {
		// OptSortOrder carries "asc" or "desc"; DESC means SortDesc=true
		filter.SortDesc = params.SortOrder.Value == generatedapi.SortOrderDesc
	}
	if params.Page.Set {
		filter.Page = int(params.Page.Value)
	}
	if params.PerPage.Set {
		filter.PerPage = int(params.PerPage.Value)
	}

	items, total, err := r.auditLogRepo.ListByProjectIDFiltered(ctx, filter)
	if err != nil {
		slog.Error("list audit logs failed", "error", err, "project_id", projectID)

		return nil, err
	}

	respItems := make([]generatedapi.AuditLog, 0, len(items))
	for _, it := range items {
		respItems = append(respItems, convertDomainAuditLog(it))
	}

	return &generatedapi.ListProjectAuditLogsOK{
		Items: respItems,
		Pagination: generatedapi.Pagination{
			Total:   uint(total),
			Page:    uint(filter.Page),
			PerPage: uint(filter.PerPage),
		},
	}, nil
}
