package notification_settings

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

var _ contract.NotificationSettingsRepository = (*Repository)(nil)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) CreateSetting(
	ctx context.Context,
	settingDTO domain.NotificationSettingDTO,
) (domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	model := settingFromDTO(settingDTO)

	const query = `
INSERT INTO notification_settings (project_id, environment_id, type, config, enabled, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *`

	rows, err := executor.Query(ctx, query,
		model.ProjectID,
		model.EnvironmentID,
		model.Type,
		model.Config,
		model.Enabled,
		model.CreatedAt,
		model.UpdatedAt,
	)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("insert notification setting: %w", err)
	}
	defer rows.Close()

	setting, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("collect notification setting: %w", err)
	}

	return setting.toDomain(), nil
}

func (r *Repository) GetSettingByID(
	ctx context.Context,
	id domain.NotificationSettingID,
) (domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM notification_settings WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("query notification setting by ID: %w", err)
	}
	defer rows.Close()

	setting, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.NotificationSetting{}, domain.ErrEntityNotFound
		}

		return domain.NotificationSetting{}, fmt.Errorf("collect notification setting: %w", err)
	}

	return setting.toDomain(), nil
}

func (r *Repository) UpdateSetting(ctx context.Context, setting domain.NotificationSetting) error {
	executor := r.getExecutor(ctx)

	model := settingFromDomain(setting)
	model.UpdatedAt = time.Now()

	const query = `
UPDATE notification_settings
SET config = $1, enabled = $2, updated_at = $3
WHERE id = $4`

	_, err := executor.Exec(ctx, query,
		model.Config,
		model.Enabled,
		model.UpdatedAt,
		model.ID,
	)
	if err != nil {
		return fmt.Errorf("update notification setting: %w", err)
	}

	return nil
}

func (r *Repository) DeleteSetting(ctx context.Context, id domain.NotificationSettingID) error {
	executor := r.getExecutor(ctx)

	// Then delete the setting
	const deleteSetting = `
DELETE FROM notification_settings
WHERE id = $1`

	_, err := executor.Exec(ctx, deleteSetting, id)
	if err != nil {
		return fmt.Errorf("delete notification setting: %w", err)
	}

	return nil
}

func (r *Repository) ListSettings(
	ctx context.Context,
	projectID domain.ProjectID,
	environmentID domain.EnvironmentID,
) ([]domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT *
FROM notification_settings
WHERE project_id = $1 AND environment_id = $2
ORDER BY id`

	rows, err := executor.Query(ctx, query, projectID, environmentID)
	if err != nil {
		return nil, fmt.Errorf("query notification settings: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		return nil, fmt.Errorf("collect notification settings: %w", err)
	}

	settings := make([]domain.NotificationSetting, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		settings = append(settings, model.toDomain())
	}

	return settings, nil
}

func (r *Repository) ListSettingsAll(
	ctx context.Context,
	projectID domain.ProjectID,
) ([]domain.NotificationSetting, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT *
FROM notification_settings
WHERE project_id = $1
ORDER BY id`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query notification settings: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[notificationSettingModel])
	if err != nil {
		return nil, fmt.Errorf("collect notification settings: %w", err)
	}

	settings := make([]domain.NotificationSetting, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		settings = append(settings, model.toDomain())
	}

	return settings, nil
}

func (r *Repository) CountSettings(
	ctx context.Context,
	projectID domain.ProjectID,
	envID domain.EnvironmentID,
) (uint, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT COUNT(*)
FROM notification_settings
WHERE project_id = $1 AND environment_id = $2`

	var count uint
	err := executor.QueryRow(ctx, query, projectID, envID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count notification settings: %w", err)
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
