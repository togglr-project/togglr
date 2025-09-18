package auditlog

import (
	"context"
	"encoding/json"
	"fmt"

	appctx "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

// ActorFromContext returns an audit actor string. If a user is present in context, returns "user:<id>", else "system".
func ActorFromContext(ctx context.Context) string {
	actor := appctx.Actor(ctx)
	if actor == domain.AuditActorUser {
		if uid := appctx.UserID(ctx); uid != 0 {
			return fmt.Sprintf("user:%d", uid)
		}
	}

	return string(actor)
}

func RequestIDFromContext(ctx context.Context) string {
	if reqID := appctx.RequestID(ctx); reqID != "" {
		return reqID
	}

	return "00000000-0000-0000-0000-000000000000"
}

// Write inserts an audit log entry within the current transaction (if any).
// oldVal and newVal are marshaled to JSON. If nil, the corresponding JSON is NULL.
func Write(
	ctx context.Context,
	exec db.Tx,
	projectID domain.ProjectID,
	featureID domain.FeatureID,
	entity domain.EntityType,
	actor string,
	action domain.AuditAction,
	oldVal any,
	newVal any,
) error {
	var (
		oldJSON []byte
		newJSON []byte
		err     error
	)

	if oldVal != nil {
		oldJSON, err = json.Marshal(oldVal)
		if err != nil {
			return fmt.Errorf("marshal old value: %w", err)
		}
	}

	if newVal != nil {
		newJSON, err = json.Marshal(newVal)
		if err != nil {
			return fmt.Errorf("marshal new value: %w", err)
		}
	}

	const query = `
		INSERT INTO audit_log (project_id, feature_id, entity, actor, action, old_value, new_value, request_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if _, err := exec.Exec(
		ctx,
		query,
		projectID,
		featureID,
		entity,
		actor,
		action,
		oldJSON,
		newJSON,
		RequestIDFromContext(ctx),
	); err != nil {
		return fmt.Errorf("insert audit_log: %w", err)
	}

	return nil
}
