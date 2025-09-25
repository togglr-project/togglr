package guard_service

import (
	"context"
	"errors"
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
	return &Repository{
		db: pool,
	}
}

// IsFeatureGuarded checks if a feature has the guarded tag.
func (r *Repository) IsFeatureGuarded(ctx context.Context, featureID domain.FeatureID) (bool, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT 1 FROM feature_tags ft
JOIN tags t ON ft.tag_id = t.id
JOIN categories c ON t.category_id = c.id
WHERE ft.feature_id = $1 AND c.slug = 'guarded'`

	var exists int

	err := executor.QueryRow(ctx, query, featureID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("check feature guarded: %w", err)
	}

	return true, nil
}

// IsEntityGuarded checks if any entity in the list is guarded.
func (r *Repository) IsEntityGuarded(ctx context.Context, entities []domain.EntityChange) (bool, error) {
	for _, entity := range entities {
		if entity.Entity == "feature" {
			featureID := domain.FeatureID(entity.EntityID)

			guarded, err := r.IsFeatureGuarded(ctx, featureID)
			if err != nil {
				return false, err
			}

			if guarded {
				return true, nil
			}
		}
		// Add checks for other entity types (rules, feature_schedules, etc.) as needed
	}

	return false, nil
}

// GetProjectActiveUserCount returns the number of active users in a project.
func (r *Repository) GetProjectActiveUserCount(ctx context.Context, projectID domain.ProjectID) (int, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT COUNT(*) FROM memberships
WHERE project_id = $1`

	var count int

	err := executor.QueryRow(ctx, query, projectID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get project active user count: %w", err)
	}

	return count, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
