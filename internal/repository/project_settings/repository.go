package project_settings

import (
	"context"
	"encoding/json"
	"errors"
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
	return &Repository{
		db: pool,
	}
}

type projectSettingModel struct {
	ID        int             `db:"id"`
	ProjectID string          `db:"project_id"`
	Name      string          `db:"name"`
	Value     json.RawMessage `db:"value"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (m *projectSettingModel) toDomain() (domain.ProjectSetting, error) {
	var value interface{}
	if err := json.Unmarshal(m.Value, &value); err != nil {
		return domain.ProjectSetting{}, fmt.Errorf("unmarshal setting value: %w", err)
	}

	return domain.ProjectSetting{
		ID:        m.ID,
		ProjectID: domain.ProjectID(m.ProjectID),
		Name:      m.Name,
		Value:     value,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

// Set sets a project setting
func (r *Repository) Set(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
	value interface{},
) error {
	executor := r.getExecutor(ctx)

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal setting value: %w", err)
	}

	const query = `
INSERT INTO project_settings (project_id, name, value)
VALUES ($1, $2, $3)
ON CONFLICT (project_id, name) DO UPDATE SET value = EXCLUDED.value, updated_at = now()`

	_, err = executor.Exec(ctx, query, projectID, name, valueJSON)
	if err != nil {
		return fmt.Errorf("set project setting: %w", err)
	}

	return nil
}

// Get retrieves a project setting
func (r *Repository) Get(
	ctx context.Context,
	projectID domain.ProjectID,
	name string,
) (domain.ProjectSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT id, project_id, name, value, created_at, updated_at
FROM project_settings
WHERE project_id = $1 AND name = $2`

	var model projectSettingModel
	err := executor.QueryRow(ctx, query, projectID, name).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Name,
		&model.Value,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ProjectSetting{}, domain.ErrEntityNotFound
		}
		return domain.ProjectSetting{}, fmt.Errorf("get project setting: %w", err)
	}

	return model.toDomain()
}

// GetAll retrieves all settings for a project
func (r *Repository) GetAll(ctx context.Context, projectID domain.ProjectID) ([]domain.ProjectSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT id, project_id, name, value, created_at, updated_at
FROM project_settings
WHERE project_id = $1
ORDER BY name`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project settings: %w", err)
	}
	defer rows.Close()

	var settings []domain.ProjectSetting
	for rows.Next() {
		var model projectSettingModel
		err := rows.Scan(
			&model.ID,
			&model.ProjectID,
			&model.Name,
			&model.Value,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan project setting: %w", err)
		}

		setting, err := model.toDomain()
		if err != nil {
			return nil, err
		}

		settings = append(settings, setting)
	}

	return settings, nil
}

// Delete removes a project setting
func (r *Repository) Delete(ctx context.Context, projectID domain.ProjectID, name string) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM project_settings
WHERE project_id = $1 AND name = $2`

	_, err := executor.Exec(ctx, query, projectID, name)
	if err != nil {
		return fmt.Errorf("delete project setting: %w", err)
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
