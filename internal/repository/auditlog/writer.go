package auditlog

import (
	"context"
	"encoding/json"
	"fmt"

	appctx "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
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
	featureID domain.FeatureID,
	actor string,
	action string,
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
		INSERT INTO audit_log (feature_id, actor, action, old_value, new_value)
		VALUES ($1, $2, $3, $4, $5)
	`

	if _, err := exec.Exec(ctx, query, featureID, actor, action, oldJSON, newJSON); err != nil {
		return fmt.Errorf("insert audit_log: %w", err)
	}

	return nil
}
