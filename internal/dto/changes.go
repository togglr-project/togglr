package dto

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainChangesToAPI converts domain changes to the API response format
func DomainChangesToAPI(result domain.ChangesListResult) generatedapi.ListChangesResponse {
	// Convert change groups
	items := make([]generatedapi.ChangeGroup, 0, len(result.Items))
	for _, group := range result.Items {
		changes := make([]generatedapi.Change, 0, len(group.Changes))
		for _, change := range group.Changes {
			// Parse EntityID as UUID
			entityID, err := uuid.Parse(change.EntityID)
			if err != nil {
				// If parsing fails, use zero UUID
				entityID = uuid.Nil
			}

			apiChange := generatedapi.Change{
				ID:       int64(change.ID),
				Entity:   generatedapi.EntityType(change.Entity),
				EntityID: entityID,
				Action:   generatedapi.AuditAction(change.Action),
			}

			// Handle old_value - convert to ChangeOldValue
			if change.OldValue != nil {
				var oldValue any
				if err := json.Unmarshal(*change.OldValue, &oldValue); err == nil {
					// For now, we'll create an empty ChangeOldValue
					// In a real implementation; you might want to store the actual value
					apiChange.OldValue = &generatedapi.ChangeOldValue{}
				}
			}

			// Handle new_value - convert to ChangeNewValue
			if change.NewValue != nil {
				var newValue any
				if err := json.Unmarshal(*change.NewValue, &newValue); err == nil {
					// For now, we'll create an empty ChangeNewValue
					// In a real implementation; you might want to store the actual value
					apiChange.NewValue = &generatedapi.ChangeNewValue{}
				}
			}

			changes = append(changes, apiChange)
		}

		// Parse RequestID as UUID
		requestID, err := uuid.Parse(group.RequestID)
		if err != nil {
			// If parsing fails, use zero UUID
			requestID = uuid.Nil
		}

		items = append(items, generatedapi.ChangeGroup{
			RequestID: requestID,
			Actor:     group.Actor,
			Username:  group.Username,
			CreatedAt: group.CreatedAt,
			Changes:   changes,
		})
	}

	// Calculate pagination
	page := 1
	perPage := 20
	if len(result.Items) > 0 {
		// This is a simplified calculation - in real implementation,
		// you'd need to pass page/perPage from the request
		perPage = len(result.Items)
	}

	// Parse ProjectID as UUID
	projectID, err := uuid.Parse(string(result.ProjectID))
	if err != nil {
		// If parsing fails, use zero UUID
		projectID = uuid.Nil
	}

	return generatedapi.ListChangesResponse{
		ProjectID: projectID,
		Items:     items,
		Pagination: generatedapi.Pagination{
			Total:   uint(result.Total),
			Page:    uint(page),
			PerPage: uint(perPage),
		},
	}
}

// APIChangesFilterToDomain converts API filter parameters to domain filter
func APIChangesFilterToDomain(
	projectID domain.ProjectID,
	params generatedapi.ListProjectChangesParams,
) domain.ChangesListFilter {
	filter := domain.ChangesListFilter{
		ProjectID: projectID,
		Page:      1,
		PerPage:   20,
		SortBy:    "created_at",
		SortDesc:  true,
	}

	// Handle pagination
	if params.Page.IsSet() {
		filter.Page = int(params.Page.Value)
	}
	if params.PerPage.IsSet() {
		filter.PerPage = int(params.PerPage.Value)
	}

	// Handle sorting
	if params.SortBy.IsSet() {
		filter.SortBy = string(params.SortBy.Value)
	}
	if params.SortOrder.IsSet() {
		filter.SortDesc = params.SortOrder.Value == "desc"
	}

	// Handle filters
	if params.Actor.IsSet() {
		filter.Actor = &params.Actor.Value
	}
	if params.Entity.IsSet() {
		entity := domain.EntityType(params.Entity.Value)
		filter.Entity = &entity
	}
	if params.Action.IsSet() {
		action := domain.AuditAction(params.Action.Value)
		filter.Action = &action
	}
	if params.FeatureID.IsSet() {
		featureID := domain.FeatureID(params.FeatureID.Value.String())
		filter.FeatureID = &featureID
	}
	if params.From.IsSet() {
		filter.From = &params.From.Value
	}
	if params.To.IsSet() {
		filter.To = &params.To.Value
	}

	return filter
}
