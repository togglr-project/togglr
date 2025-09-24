package apibackend

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/go-faster/jx"
	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// ListPendingChanges handles GET /api/v1/pending_changes
func (r *RestAPI) ListPendingChanges(
	ctx context.Context,
	params generatedapi.ListPendingChangesParams,
) (generatedapi.ListPendingChangesRes, error) {
	// Build filter
	filter := contract.PendingChangesListFilter{
		Page:     1,
		PerPage:  20,
		SortBy:   "created_at",
		SortDesc: true,
	}

	if params.ProjectID.Set {
		projectID := domain.ProjectID(params.ProjectID.Value.String())
		filter.ProjectID = &projectID
	}

	if params.Status.Set {
		status := domain.PendingChangeStatus(params.Status.Value)
		filter.Status = &status
	}

	if params.UserID.Set {
		userID := int(params.UserID.Value)
		filter.UserID = &userID
	}

	if params.Page.Set {
		filter.Page = params.Page.Value
	}

	if params.PerPage.Set {
		filter.PerPage = params.PerPage.Value
	}

	if params.SortBy.Set {
		filter.SortBy = string(params.SortBy.Value)
	}

	if params.SortDesc.Set {
		filter.SortDesc = params.SortDesc.Value
	}

	// Get pending changes
	changes, total, err := r.pendingChangesUseCase.List(ctx, filter)
	if err != nil {
		slog.Error("list pending changes failed", "error", err)
		return nil, err
	}

	// Convert to response format
	var responseChanges []generatedapi.PendingChangeResponse
	for _, change := range changes {
		responseChange := convertPendingChangeToResponse(&change)
		responseChanges = append(responseChanges, responseChange)
	}

	return &generatedapi.PendingChangesListResponse{
		Data: responseChanges,
		Pagination: generatedapi.Pagination{
			Total:   uint(total),
			Page:    filter.Page,
			PerPage: filter.PerPage,
		},
	}, nil
}

// convertPendingChangeToResponse converts domain.PendingChange to generatedapi.PendingChangeResponse
func convertPendingChangeToResponse(change *domain.PendingChange) generatedapi.PendingChangeResponse {
	// Convert entities
	var entities []generatedapi.EntityChange
	for _, entity := range change.Change.Entities {
		// Convert changes
		changes := make(map[string]generatedapi.ChangeValue)
		for field, changeValue := range entity.Changes {
			// Convert to JSON for storage
			oldJSON, _ := json.Marshal(changeValue.Old)
			newJSON, _ := json.Marshal(changeValue.New)

			changes[field] = generatedapi.ChangeValue{
				Old: jx.Raw(oldJSON),
				New: jx.Raw(newJSON),
			}
		}

		entityUUID, _ := uuid.Parse(entity.EntityID)
		entities = append(entities, generatedapi.EntityChange{
			Entity:   generatedapi.EntityChangeEntity(entity.Entity),
			EntityID: entityUUID,
			Action:   generatedapi.EntityChangeAction(entity.Action),
			Changes:  changes,
		})
	}

	// Convert meta
	meta := generatedapi.PendingChangeMeta{
		Reason: change.Change.Meta.Reason,
		Client: change.Change.Meta.Client,
		Origin: change.Change.Meta.Origin,
		SingleUserProject: generatedapi.OptBool{
			Value: change.Change.Meta.SingleUserProject,
			Set:   true,
		},
	}

	// Convert payload
	payload := generatedapi.PendingChangePayload{
		Entities: entities,
		Meta:     meta,
	}

	responseID, _ := uuid.Parse(change.ID.String())
	projectUUID, _ := uuid.Parse(change.ProjectID.String())

	response := generatedapi.PendingChangeResponse{
		ID:          responseID,
		ProjectID:   projectUUID,
		RequestedBy: change.RequestedBy,
		Change:      payload,
		Status:      generatedapi.PendingChangeResponseStatus(change.Status),
		CreatedAt:   change.CreatedAt,
	}

	// Set optional fields
	if change.RequestUserID != nil {
		response.RequestUserID = generatedapi.OptNilUint{
			Value: uint(*change.RequestUserID),
			Set:   true,
		}
	}

	if change.ApprovedBy != nil {
		response.ApprovedBy = generatedapi.OptNilString{
			Value: *change.ApprovedBy,
			Set:   true,
		}
	}

	if change.ApprovedUserID != nil {
		response.ApprovedUserID = generatedapi.OptNilUint{
			Value: uint(*change.ApprovedUserID),
			Set:   true,
		}
	}

	if change.ApprovedAt != nil {
		response.ApprovedAt = generatedapi.OptNilDateTime{
			Value: *change.ApprovedAt,
			Set:   true,
		}
	}

	if change.RejectedBy != nil {
		response.RejectedBy = generatedapi.OptNilString{
			Value: *change.RejectedBy,
			Set:   true,
		}
	}

	if change.RejectedAt != nil {
		response.RejectedAt = generatedapi.OptNilDateTime{
			Value: *change.RejectedAt,
			Set:   true,
		}
	}

	if change.RejectionReason != nil {
		response.RejectionReason = generatedapi.OptNilString{
			Value: *change.RejectionReason,
			Set:   true,
		}
	}

	return response
}
