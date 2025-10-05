package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

//nolint:ireturn // repository uses db.Tx interface
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}

type projectHealthRow struct {
	ProjectID                  string `db:"project_id"`
	ProjectName                string `db:"project_name"`
	EnvironmentID              string `db:"environment_id"`
	EnvironmentKey             string `db:"environment_key"`
	TotalFeatures              int64  `db:"total_features"`
	EnabledFeatures            int64  `db:"enabled_features"`
	DisabledFeatures           int64  `db:"disabled_features"`
	AutoDisableManagedFeatures int64  `db:"auto_disable_managed_features"`
	UncategorizedFeatures      int64  `db:"uncategorized_features"`
	GuardedFeatures            int64  `db:"guarded_features"`
	PendingFeatures            int64  `db:"pending_features"`
	PendingGuardedFeatures     int64  `db:"pending_guarded_features"`
	HealthStatus               string `db:"health_status"`
}

func (r *Repository) ProjectHealth(
	ctx context.Context,
	envKey string,
	projectID *string,
) ([]domain.ProjectHealth, error) {
	exec := r.getExecutor(ctx)
	query := `SELECT project_id, project_name, environment_id, environment_key,
		total_features, enabled_features, disabled_features,
		auto_disable_managed_features, uncategorized_features,
		guarded_features, pending_features, pending_guarded_features, health_status
		FROM v_project_health WHERE environment_key = $1`
	args := []any{envKey}
	if projectID != nil {
		query += ` AND project_id = $2` //nolint:goconst // false positive
		args = append(args, *projectID)
	}
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query v_project_health: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[projectHealthRow])
	if err != nil {
		return nil, fmt.Errorf("collect v_project_health: %w", err)
	}
	out := make([]domain.ProjectHealth, 0, len(models))
	for i := range models {
		m := models[i]
		out = append(out, domain.ProjectHealth{
			ProjectID:                  m.ProjectID,
			ProjectName:                m.ProjectName,
			EnvironmentID:              m.EnvironmentID,
			EnvironmentKey:             m.EnvironmentKey,
			TotalFeatures:              uint(m.TotalFeatures),
			EnabledFeatures:            uint(m.EnabledFeatures),
			DisabledFeatures:           uint(m.DisabledFeatures),
			AutoDisableManagedFeatures: uint(m.AutoDisableManagedFeatures),
			UncategorizedFeatures:      uint(m.UncategorizedFeatures),
			GuardedFeatures:            uint(m.GuardedFeatures),
			PendingFeatures:            uint(m.PendingFeatures),
			PendingGuardedFeatures:     uint(m.PendingGuardedFeatures),
			HealthStatus:               domain.HealthStatus(m.HealthStatus),
		})
	}

	return out, nil
}

type categoryHealthRow struct {
	ProjectID                  string `db:"project_id"`
	ProjectName                string `db:"project_name"`
	EnvironmentID              string `db:"environment_id"`
	EnvironmentKey             string `db:"environment_key"`
	CategoryID                 string `db:"category_id"`
	CategoryName               string `db:"category_name"`
	CategorySlug               string `db:"category_slug"`
	TotalFeatures              int64  `db:"total_features"`
	EnabledFeatures            int64  `db:"enabled_features"`
	DisabledFeatures           int64  `db:"disabled_features"`
	PendingFeatures            int64  `db:"pending_features"`
	GuardedFeatures            int64  `db:"guarded_features"`
	AutoDisableManagedFeatures int64  `db:"auto_disable_managed_features"`
	PendingGuardedFeatures     int64  `db:"pending_guarded_features"`
	HealthStatus               string `db:"health_status"`
}

func (r *Repository) CategoryHealth(
	ctx context.Context,
	envKey string,
	projectID *string,
) ([]domain.CategoryHealth, error) {
	exec := r.getExecutor(ctx)
	query := `
SELECT project_id, project_name, environment_id, environment_key, category_id, category_name, 
	category_slug, total_features, enabled_features, disabled_features, pending_features, guarded_features,
	auto_disable_managed_features, pending_guarded_features, health_status
FROM v_project_category_health WHERE environment_key = $1`
	args := []any{envKey}
	if projectID != nil {
		query += ` AND project_id = $2`
		args = append(args, *projectID)
	}
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query v_project_category_health: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[categoryHealthRow])
	if err != nil {
		return nil, fmt.Errorf("collect v_project_category_health: %w", err)
	}
	out := make([]domain.CategoryHealth, 0, len(models))
	for i := range models {
		m := models[i]
		out = append(out, domain.CategoryHealth{
			ProjectID:                  m.ProjectID,
			ProjectName:                m.ProjectName,
			EnvironmentID:              m.EnvironmentID,
			EnvironmentKey:             m.EnvironmentKey,
			CategoryID:                 m.CategoryID,
			CategoryName:               m.CategoryName,
			CategorySlug:               m.CategorySlug,
			TotalFeatures:              uint(m.TotalFeatures),
			EnabledFeatures:            uint(m.EnabledFeatures),
			DisabledFeatures:           uint(m.DisabledFeatures),
			PendingFeatures:            uint(m.PendingFeatures),
			GuardedFeatures:            uint(m.GuardedFeatures),
			AutoDisableManagedFeatures: uint(m.AutoDisableManagedFeatures),
			PendingGuardedFeatures:     uint(m.PendingGuardedFeatures),
			HealthStatus:               domain.HealthStatus(m.HealthStatus),
		})
	}

	return out, nil
}

type recentActivityRow struct {
	ProjectID      string    `db:"project_id"`
	EnvironmentID  string    `db:"environment_id"`
	EnvironmentKey string    `db:"environment_key"`
	ProjectName    string    `db:"project_name"`
	RequestID      string    `db:"request_id"`
	Actor          string    `db:"actor"`
	CreatedAt      time.Time `db:"created_at"`
	Status         string    `db:"status"`
	Changes        []byte    `db:"changes"`
}

func (r *Repository) RecentActivity(
	ctx context.Context,
	envKey string,
	projectID *string,
	limit uint,
) ([]domain.RecentActivity, error) {
	exec := r.getExecutor(ctx)
	query := `
SELECT project_id, environment_id, environment_key, project_name, request_id, actor, created_at, status, changes
FROM v_project_recent_activity WHERE environment_key = $1`
	args := []any{envKey}
	if projectID != nil {
		query += ` AND project_id = $2`
		args = append(args, *projectID)
	}
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query v_project_recent_activity: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[recentActivityRow])
	if err != nil {
		return nil, fmt.Errorf("collect v_project_recent_activity: %w", err)
	}
	out := make([]domain.RecentActivity, 0, len(models))
	for i := range models {
		m := models[i]
		changes, err := parseRecentChanges(m.Changes)
		if err != nil {
			return nil, fmt.Errorf("parse changes json: %w", err)
		}
		out = append(out, domain.RecentActivity{
			ProjectID:      m.ProjectID,
			EnvironmentID:  m.EnvironmentID,
			EnvironmentKey: m.EnvironmentKey,
			ProjectName:    m.ProjectName,
			RequestID:      m.RequestID,
			Actor:          m.Actor,
			CreatedAt:      m.CreatedAt,
			Status:         m.Status,
			Changes:        changes,
		})
	}

	return out, nil
}

type riskyFeatureRow struct {
	ProjectID      string `db:"project_id"`
	ProjectName    string `db:"project_name"`
	EnvironmentID  string `db:"environment_id"`
	EnvironmentKey string `db:"environment_key"`
	FeatureID      string `db:"feature_id"`
	FeatureName    string `db:"feature_name"`
	Enabled        bool   `db:"enabled"`
	RiskyTags      string `db:"risky_tags"`
	HasPending     bool   `db:"has_pending"`
}

func (r *Repository) RiskyFeatures(
	ctx context.Context,
	envKey string,
	projectID *string,
	limit uint,
) ([]domain.RiskyFeature, error) {
	exec := r.getExecutor(ctx)
	query := `
SELECT project_id, project_name, environment_id, environment_key, feature_id,
       feature_name, enabled, risky_tags, has_pending
FROM v_project_top_risky_features WHERE environment_key = $1`
	args := []any{envKey}
	if projectID != nil {
		query += ` AND project_id = $2`
		args = append(args, *projectID)
	}
	query += ` ORDER BY feature_name ASC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query v_project_top_risky_features: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[riskyFeatureRow])
	if err != nil {
		return nil, fmt.Errorf("collect v_project_top_risky_features: %w", err)
	}
	out := make([]domain.RiskyFeature, 0, len(models))
	for i := range models {
		m := models[i]
		out = append(out, domain.RiskyFeature{
			ProjectID:      m.ProjectID,
			ProjectName:    m.ProjectName,
			EnvironmentID:  m.EnvironmentID,
			EnvironmentKey: m.EnvironmentKey,
			FeatureID:      m.FeatureID,
			FeatureName:    m.FeatureName,
			Enabled:        m.Enabled,
			HasPending:     m.HasPending,
			RiskyTags:      m.RiskyTags,
		})
	}

	return out, nil
}

type pendingSummaryRow struct {
	ProjectID             string     `db:"project_id"`
	ProjectName           string     `db:"project_name"`
	EnvironmentID         string     `db:"environment_id"`
	EnvironmentKey        string     `db:"environment_key"`
	TotalPending          int64      `db:"total_pending"`
	PendingFeatureChanges int64      `db:"pending_feature_changes"`
	PendingGuardedChanges int64      `db:"pending_guarded_changes"`
	OldestRequestAt       *time.Time `db:"oldest_request_at"`
}

func (r *Repository) PendingSummary(
	ctx context.Context,
	envKey string,
	projectID *string,
) ([]domain.PendingSummary, error) {
	exec := r.getExecutor(ctx)
	query := `
SELECT project_id, project_name, environment_id, environment_key, total_pending, pending_feature_changes,
pending_guarded_changes, oldest_request_at
FROM v_project_pending_summary WHERE environment_key = $1`
	args := []any{envKey}
	if projectID != nil {
		query += ` AND project_id = $2`
		args = append(args, *projectID)
	}
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query v_project_pending_summary: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[pendingSummaryRow])
	if err != nil {
		return nil, fmt.Errorf("collect v_project_pending_summary: %w", err)
	}
	out := make([]domain.PendingSummary, 0, len(models))
	for i := range models {
		m := models[i]
		out = append(out, domain.PendingSummary{
			ProjectID:             m.ProjectID,
			ProjectName:           m.ProjectName,
			EnvironmentID:         m.EnvironmentID,
			EnvironmentKey:        m.EnvironmentKey,
			TotalPending:          uint(m.TotalPending),
			PendingFeatureChanges: uint(m.PendingFeatureChanges),
			PendingGuardedChanges: uint(m.PendingGuardedChanges),
			OldestRequestAt:       m.OldestRequestAt,
		})
	}

	return out, nil
}

// parseRecentChanges decodes JSONB array of objects into []domain.RecentChange.
func parseRecentChanges(raw []byte) ([]domain.RecentChange, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var tmp []struct {
		Entity   string `json:"entity"`
		EntityID string `json:"entity_id"`
		Action   string `json:"action"`
	}
	if err := json.Unmarshal(raw, &tmp); err != nil {
		return nil, err
	}
	out := make([]domain.RecentChange, 0, len(tmp))
	for i := range tmp {
		out = append(out, domain.RecentChange{
			Entity:   tmp[i].Entity,
			EntityID: tmp[i].EntityID,
			Action:   tmp[i].Action,
		})
	}

	return out, nil
}

// TopActiveFeatureIDs returns feature IDs from v_top_active_features ordered by rank_score.
func (r *Repository) TopActiveFeatureIDs(
	ctx context.Context,
	envKey string,
	projectID *string,
	limit uint,
) ([]string, error) {
	exec := r.getExecutor(ctx)
	query := `SELECT feature_id FROM v_top_active_features WHERE environment_key = $1 AND enabled = true`
	args := []any{envKey}
	if projectID != nil {
		query += ` AND project_id = $2`
		args = append(args, *projectID)
	}
	query += ` ORDER BY rank_score DESC LIMIT $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query v_top_active_features: %w", err)
	}
	defer rows.Close()

	ids, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
		var id string
		if err := row.Scan(&id); err != nil {
			return "", err
		}

		return id, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect v_top_active_features ids: %w", err)
	}

	return ids, nil
}
