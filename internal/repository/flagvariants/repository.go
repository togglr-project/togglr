package flagvariants

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/etoggle/internal/domain"
	"github.com/rom8726/etoggle/internal/repository/auditlog"
	"github.com/rom8726/etoggle/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) Create(ctx context.Context, v domain.FlagVariant) (domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	var (
		query string
		args  []any
	)

	if v.ID != "" {
		// Use client-provided ID
		query = `
INSERT INTO flag_variants (id, feature_id, name, rollout_percent)
VALUES ($1, $2, $3, $4)
RETURNING id, feature_id, name, rollout_percent`
		args = []any{v.ID, v.FeatureID, v.Name, int(v.RolloutPercent)}
	} else {
		query = `
INSERT INTO flag_variants (feature_id, name, rollout_percent)
VALUES ($1, $2, $3)
RETURNING id, feature_id, name, rollout_percent`
		args = []any{v.FeatureID, v.Name, int(v.RolloutPercent)}
	}

	var model flagVariantModel
	if err := executor.QueryRow(ctx, query, args...).Scan(
		&model.ID,
		&model.FeatureID,
		&model.Name,
		&model.RolloutPercent,
	); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("insert flag_variant: %w", err)
	}

	newVariant := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newVariant.FeatureID,
		auditlog.EntityFlagVariant,
		auditlog.ActorFromContext(ctx),
		auditlog.ActionCreate,
		nil,
		newVariant,
	); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("audit flag_variant create: %w", err)
	}

	return newVariant, nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.FlagVariantID) (domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM flag_variants WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.FlagVariant{}, fmt.Errorf("query flag_variant by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[flagVariantModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FlagVariant{}, domain.ErrEntityNotFound
		}
		return domain.FlagVariant{}, fmt.Errorf("collect flag_variant row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) List(ctx context.Context) ([]domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM flag_variants ORDER BY name ASC`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query flag_variants: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[flagVariantModel])
	if err != nil {
		return nil, fmt.Errorf("collect flag_variant rows: %w", err)
	}

	items := make([]domain.FlagVariant, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) ListByFeatureID(ctx context.Context, featureID domain.FeatureID) ([]domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM flag_variants WHERE feature_id = $1 ORDER BY name ASC`

	rows, err := executor.Query(ctx, query, featureID)
	if err != nil {
		return nil, fmt.Errorf("query flag_variants by feature_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[flagVariantModel])
	if err != nil {
		return nil, fmt.Errorf("collect flag_variant rows: %w", err)
	}

	items := make([]domain.FlagVariant, 0, len(models))
	for _, m := range models {
		items = append(items, m.toDomain())
	}

	return items, nil
}

func (r *Repository) Update(ctx context.Context, v domain.FlagVariant) (domain.FlagVariant, error) {
	executor := r.getExecutor(ctx)

	// Read old state for audit within the same transaction.
	oldVariant, err := r.GetByID(ctx, v.ID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.FlagVariant{}, err
		}
		return domain.FlagVariant{}, fmt.Errorf("get flag_variant before update: %w", err)
	}

	const query = `
UPDATE flag_variants
SET feature_id = $1, name = $2, rollout_percent = $3
WHERE id = $4
RETURNING id, feature_id, name, rollout_percent`

	var model flagVariantModel
	if err := executor.QueryRow(ctx, query, v.FeatureID, v.Name, int(v.RolloutPercent), v.ID).Scan(
		&model.ID,
		&model.FeatureID,
		&model.Name,
		&model.RolloutPercent,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FlagVariant{}, domain.ErrEntityNotFound
		}
		return domain.FlagVariant{}, fmt.Errorf("update flag_variant: %w", err)
	}

	newVariant := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newVariant.FeatureID,
		auditlog.EntityFlagVariant,
		auditlog.ActorFromContext(ctx),
		auditlog.ActionUpdate,
		oldVariant,
		newVariant,
	); err != nil {
		return domain.FlagVariant{}, fmt.Errorf("audit flag_variant update: %w", err)
	}

	return newVariant, nil
}

func (r *Repository) Delete(ctx context.Context, id domain.FlagVariantID) error {
	executor := r.getExecutor(ctx)

	oldVariant, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		oldVariant.FeatureID,
		auditlog.EntityFlagVariant,
		auditlog.ActorFromContext(ctx),
		auditlog.ActionDelete,
		oldVariant,
		nil,
	); err != nil {
		return fmt.Errorf("audit flag_variant delete: %w", err)
	}

	const query = `DELETE FROM flag_variants WHERE id = $1`

	ct, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete flag_variant: %w", err)
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
