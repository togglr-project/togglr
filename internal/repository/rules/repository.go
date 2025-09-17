package rules

import (
	"context"
	"encoding/json"
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

// Create inserts a new rule and returns the created entity.
//
//nolint:lll // long query string is acceptable
func (r *Repository) Create(ctx context.Context, rule domain.Rule) (domain.Rule, error) {
	executor := r.getExecutor(ctx)

	conditionsData, err := json.Marshal(rule.Conditions)
	if err != nil {
		return domain.Rule{}, fmt.Errorf("marshal conditions: %w", err)
	}

	var (
		query string
		args  []any
	)

	if rule.ID != "" {
		query = `
INSERT INTO rules (id, project_id, feature_id, condition, segment_id, is_customized, action, flag_variant_id, priority)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, project_id, feature_id, condition, segment_id, is_customized, action, flag_variant_id, priority, created_at`
		args = []any{rule.ID, rule.ProjectID, rule.FeatureID, conditionsData, rule.SegmentID,
			rule.IsCustomized, rule.Action, rule.FlagVariantID, int(rule.Priority)}
	} else {
		query = `
INSERT INTO rules (project_id, feature_id, condition, segment_id, is_customized, action, flag_variant_id, priority)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, project_id, feature_id, condition, segment_id, is_customized, action, flag_variant_id, priority, created_at`
		args = []any{rule.ProjectID, rule.FeatureID, conditionsData, rule.SegmentID,
			rule.IsCustomized, rule.Action, rule.FlagVariantID, int(rule.Priority)}
	}

	var model ruleModel
	if err := executor.QueryRow(ctx, query, args...).Scan(
		&model.ID,
		&model.ProjectID,
		&model.FeatureID,
		&model.Condition,
		&model.SegmentID,
		&model.IsCustomized,
		&model.Action,
		&model.FlagVariantID,
		&model.Priority,
		&model.CreatedAt,
	); err != nil {
		return domain.Rule{}, fmt.Errorf("insert rule: %w", err)
	}

	newRule := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newRule.ProjectID,
		newRule.FeatureID,
		domain.EntityRule,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionCreate,
		nil,
		newRule,
	); err != nil {
		return domain.Rule{}, fmt.Errorf("audit rule create: %w", err)
	}

	return newRule, nil
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

	const query = `SELECT * FROM rules ORDER BY feature_id, priority`

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

	const query = `SELECT * FROM rules WHERE feature_id = $1 ORDER BY priority`

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

func (r *Repository) ListCustomizedFeatureIDsBySegment(
	ctx context.Context,
	segmentID domain.SegmentID,
) ([]domain.FeatureID, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT DISTINCT feature_id FROM rules WHERE segment_id = $1 AND is_customized = TRUE`

	rows, err := executor.Query(ctx, query, segmentID)
	if err != nil {
		return nil, fmt.Errorf("query distinct feature_ids by segment_id: %w", err)
	}
	defer rows.Close()

	ids := make([]domain.FeatureID, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan feature_id: %w", err)
		}
		ids = append(ids, domain.FeatureID(id))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate feature_ids: %w", err)
	}

	return ids, nil
}

func (r *Repository) ListNotCustomizedRulesBySegment(
	ctx context.Context,
	segmentID domain.SegmentID,
) ([]domain.Rule, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM rules WHERE segment_id = $1 AND is_customized = FALSE`

	rows, err := executor.Query(ctx, query, segmentID)
	if err != nil {
		return nil, fmt.Errorf("query rules by segment_id: %w", err)
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

	// Read old state for audit within the same transaction.
	oldRule, err := r.GetByID(ctx, rule.ID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.Rule{}, err
		}
		return domain.Rule{}, fmt.Errorf("get rule before update: %w", err)
	}

	const query = `
UPDATE rules
SET feature_id = $1, condition = $2, flag_variant_id = $3, priority = $4, action = $5, segment_id = $6, is_customized = $7
WHERE id = $8
RETURNING id, project_id, feature_id, condition, action, flag_variant_id, priority, segment_id, is_customized, created_at`

	conditionsData, err := json.Marshal(rule.Conditions)
	if err != nil {
		return domain.Rule{}, fmt.Errorf("marshal conditions: %w", err)
	}

	var model ruleModel
	if err := executor.QueryRow(ctx, query,
		rule.FeatureID,
		conditionsData,
		rule.FlagVariantID,
		int(rule.Priority),
		rule.Action,
		rule.SegmentID,
		rule.IsCustomized,
		rule.ID,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.FeatureID,
		&model.Condition,
		&model.Action,
		&model.FlagVariantID,
		&model.Priority,
		&model.SegmentID,
		&model.IsCustomized,
		&model.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Rule{}, domain.ErrEntityNotFound
		}
		return domain.Rule{}, fmt.Errorf("update rule: %w", err)
	}

	newRule := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newRule.ProjectID,
		newRule.FeatureID,
		domain.EntityRule,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionUpdate,
		oldRule,
		newRule,
	); err != nil {
		return domain.Rule{}, fmt.Errorf("audit rule update: %w", err)
	}

	return newRule, nil
}

func (r *Repository) Delete(ctx context.Context, id domain.RuleID) error {
	executor := r.getExecutor(ctx)

	oldRule, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		oldRule.ProjectID,
		oldRule.FeatureID,
		domain.EntityRule,
		auditlog.ActorFromContext(ctx),
		domain.AuditActionDelete,
		oldRule,
		nil,
	); err != nil {
		return fmt.Errorf("audit rule delete: %w", err)
	}

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
