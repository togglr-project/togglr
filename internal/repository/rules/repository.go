package rules

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

// Create inserts a new rule and returns the created entity.
//
//nolint:lll // long query string is acceptable
func (r *Repository) Create(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO rules (feature_id, condition, flag_variant_id, priority)
VALUES ($1, $2, $3, $4)
RETURNING id, feature_id, condition, flag_variant_id, priority, created_at`

	var model ruleModel
	if err := executor.QueryRow(ctx, query,
		rule.FeatureID,
		[]byte(rule.Condition),
		rule.FlagVariantID,
		int(rule.Priority),
	).Scan(
		&model.ID,
		&model.FeatureID,
		&model.Condition,
		&model.FlagVariantID,
		&model.Priority,
		&model.CreatedAt,
	); err != nil {
		return domain.Rule{}, fmt.Errorf("insert rule: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.RuleID) (domain.Rule, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM rules WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Rule{}, fmt.Errorf("query rule by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ruleModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Rule{}, domain.ErrEntityNotFound
		}
		return domain.Rule{}, fmt.Errorf("collect rule row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.Rule, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM rules ORDER BY priority DESC, created_at DESC`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query rules: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[ruleModel])
	if err != nil {
		return nil, fmt.Errorf("collect rule rows: %w", err)
	}

	items := make([]domain.Rule, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.Rule, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM rules WHERE feature_id = $1 ORDER BY priority DESC, created_at DESC`

	rows, err := executor.Query(ctx, query, featureID)
	if err != nil {
		return nil, fmt.Errorf("query rules by feature_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[ruleModel])
	if err != nil {
		return nil, fmt.Errorf("collect rule rows: %w", err)
	}

	items := make([]domain.Rule, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

// Update updates existing rule by ID and returns the updated entity.
//
//nolint:lll // long query string is acceptable
func (r *Repository) Update(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE rules
SET feature_id = $1, condition = $2, flag_variant_id = $3, priority = $4
WHERE id = $5
RETURNING id, feature_id, condition, flag_variant_id, priority, created_at`

	var model ruleModel
	if err := executor.QueryRow(ctx, query,
		rule.FeatureID,
		[]byte(rule.Condition),
		rule.FlagVariantID,
		int(rule.Priority),
		rule.ID,
	).Scan(
		&model.ID,
		&model.FeatureID,
		&model.Condition,
		&model.FlagVariantID,
		&model.Priority,
		&model.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Rule{}, domain.ErrEntityNotFound
		}
		return domain.Rule{}, fmt.Errorf("update rule: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id domain.RuleID) error {
	executor := r.getExecutor(ctx)

	const query = `DELETE FROM rules WHERE id = $1`

	ct, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
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
