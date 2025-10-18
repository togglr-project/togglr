package feedback_events

import (
	"context"
	"encoding/json"
	"fmt"

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

func (r *Repository) AddEvent(ctx context.Context, event domain.FeedbackEventDTO) error {
	executor := r.getExecutor(ctx)

	contextData, err := json.Marshal(event.Context)
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	const query = `
INSERT INTO monitoring.feedback_events
	(feature_id, variant_key, event_type, reward, context, created_at)
VALUES ($1, $2, $3, $4, $5, NOW())`

	_, err = executor.Exec(ctx, query,
		event.FeatureID,
		event.VariantKey,
		event.EventType,
		event.Reward,
		contextData,
	)

	return err
}

func (r *Repository) AddEventsBatch(ctx context.Context, events []domain.FeedbackEventDTO) error {
	if len(events) == 0 {
		return nil
	}

	executor := r.getExecutor(ctx)
	batch := &pgx.Batch{}

	const query = `
INSERT INTO monitoring.feedback_events
    (project_id, environment_id, feature_id, feature_key, environment_key, variant_key, 
     event_type, algorithm_slug, reward, context, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())`

	for _, e := range events {
		contextData, err := json.Marshal(e.Context)
		if err != nil {
			return fmt.Errorf("failed to marshal context for feature_id=%s: %w", e.FeatureID, err)
		}

		batch.Queue(query,
			e.ProjectID,
			e.EnvironmentID,
			e.FeatureID,
			e.FeatureKey,
			e.EnvKey,
			e.VariantKey,
			e.EventType,
			e.AlgorithmSlug,
			e.Reward,
			contextData,
		)
	}

	br := executor.SendBatch(ctx, batch)
	defer func() { _ = br.Close() }()

	for range events {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}
	}

	return nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
