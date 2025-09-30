package realtime

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) FetchAfter(ctx context.Context, after time.Time) ([]domain.RealtimeEvent, error) {
	exec := r.getExecutor(ctx)

	const query = `
SELECT source,
       event_id,
       project_id,
       environment_id,
       environment_key,
       entity,
       entity_id,
       action,
       created_at
FROM v_realtime_events
WHERE created_at > $1
ORDER BY created_at ASC`

	rows, err := exec.Query(ctx, query, after)
	if err != nil {
		return nil, fmt.Errorf("query realtime events: %w", err)
	}
	defer rows.Close()

	type rowModel struct {
		Source         string    `db:"source"`
		EventID        string    `db:"event_id"`
		ProjectID      string    `db:"project_id"`
		EnvironmentID  int64     `db:"environment_id"`
		EnvironmentKey string    `db:"environment_key"`
		Entity         string    `db:"entity"`
		EntityID       string    `db:"entity_id"`
		Action         string    `db:"action"`
		CreatedAt      time.Time `db:"created_at"`
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[rowModel])
	if err != nil {
		return nil, fmt.Errorf("collect realtime events: %w", err)
	}

	res := make([]domain.RealtimeEvent, 0, len(items))
	for i := range items {
		m := items[i]
		res = append(res, domain.RealtimeEvent{
			Source:         m.Source,
			EventID:        m.EventID,
			ProjectID:      domain.ProjectID(m.ProjectID),
			EnvironmentID:  m.EnvironmentID,
			EnvironmentKey: m.EnvironmentKey,
			Entity:         m.Entity,
			EntityID:       m.EntityID,
			Action:         m.Action,
			CreatedAt:      m.CreatedAt,
		})
	}

	return res, nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}
	return r.db
}
