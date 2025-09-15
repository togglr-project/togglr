package auditlog

import (
	"context"
	"encoding/json"
	"fmt"

	appctx "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

// ActorFromContext returns audit actor string. If a user is present in context, returns "user:<id>", else "system".
func ActorFromContext(ctx context.Context) string {
	if uid := appctx.UserID(ctx); uid != 0 {
		return fmt.Sprintf("user:%d", uid)
	}

	return "system"
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
		INSERT INTO audit_log (project_id, feature_id, entity, actor, action, old_value, new_value)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := exec.Exec(ctx, query, projectID, featureID, entity, actor, action, oldJSON, newJSON); err != nil {
		return fmt.Errorf("insert audit_log: %w", err)
	}

	return nil
}
