package apibackend

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/go-faster/jx"
	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:nilerr // it's ok here
func (r *RestAPI) GetAuditLogEntry(
	ctx context.Context,
	params generatedapi.GetAuditLogEntryParams,
) (generatedapi.GetAuditLogEntryRes, error) {
	id := domain.AuditLogID(params.ID)

	entry, err := r.auditLogRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString("audit log not found"),
			}}, nil
		}

		slog.Error("get audit log by id failed", "error", err, "id", id)

		return nil, err
	}

	// Permission check: user must be allowed to view audit logs for the project
	if err := r.permissionsService.CanViewAudit(ctx, entry.ProjectID); err != nil {
		return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
			Message: generatedapi.NewOptString("permission denied"),
		}}, nil
	}

	resp := convertDomainAuditLog(entry)

	return &resp, nil
}

//nolint:gocognit // complex conversion function, refactor later
func convertDomainAuditLog(a domain.AuditLog) generatedapi.AuditLog {
	var (
		envKey   generatedapi.OptString
		userOpt  generatedapi.OptString
		featOpt  generatedapi.OptUUID
		reqIDOpt generatedapi.OptUUID
		oldOpt   generatedapi.OptNilAuditLogOldValue
		newOpt   generatedapi.OptNilAuditLogNewValue
	)

	if a.EnvKey != "" {
		envKey = generatedapi.NewOptString(a.EnvKey)
	}
	if a.Username != "" {
		userOpt = generatedapi.NewOptString(a.Username)
	}
	if s := a.FeatureID.String(); s != "" {
		if u, err := uuid.Parse(s); err == nil {
			featOpt = generatedapi.NewOptUUID(u)
		}
	}
	if a.RequestID != "" {
		if u, err := uuid.Parse(a.RequestID); err == nil {
			reqIDOpt = generatedapi.NewOptUUID(u)
		}
	}
	if len(a.OldValue) > 0 {
		// Check if it's valid JSON and convert to map[string]jx.Raw
		var oldValue map[string]interface{}
		if err := json.Unmarshal(a.OldValue, &oldValue); err == nil {
			tempOld := make(generatedapi.AuditLogOldValue)
			for k, v := range oldValue {
				if vBytes, err := json.Marshal(v); err == nil {
					tempOld[k] = jx.Raw(vBytes)
				}
			}
			oldOpt = generatedapi.NewOptNilAuditLogOldValue(tempOld)
		}
	}
	if len(a.NewValue) > 0 {
		// Check if it's valid JSON and convert to map[string]jx.Raw
		var newValue map[string]interface{}
		if err := json.Unmarshal(a.NewValue, &newValue); err == nil {
			tempNew := make(generatedapi.AuditLogNewValue)
			for k, v := range newValue {
				if vBytes, err := json.Marshal(v); err == nil {
					tempNew[k] = jx.Raw(vBytes)
				}
			}
			newOpt = generatedapi.NewOptNilAuditLogNewValue(tempNew)
		}
	}

	projUUID, _ := uuid.Parse(string(a.ProjectID))
	entityIDUUID := uuid.Nil
	if a.EntityID != "" {
		if u, err := uuid.Parse(a.EntityID); err == nil {
			entityIDUUID = u
		}
	}

	return generatedapi.AuditLog{
		ID:             int64(a.ID.Uint64()),
		ProjectID:      projUUID,
		EnvironmentID:  int(a.EnvironmentID),
		EnvironmentKey: envKey,
		Entity:         string(a.Entity),
		EntityID:       entityIDUUID,
		FeatureID:      featOpt,
		Action:         string(a.Action),
		Actor:          a.Actor,
		Username:       userOpt,
		RequestID:      reqIDOpt,
		OldValue:       oldOpt,
		NewValue:       newOpt,
		CreatedAt:      a.CreatedAt,
	}
}
