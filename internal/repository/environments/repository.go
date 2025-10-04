package environments

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/repository/auditlog"
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

func (r *Repository) Create(ctx context.Context, env domain.Environment) (domain.Environment, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO environments (project_id, key, name, api_key, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, project_id, key, name, api_key, created_at`

	var model environmentModel

	err := executor.QueryRow(ctx, query,
		env.ProjectID,
		env.Key,
		env.Name,
		env.APIKey,
		env.CreatedAt,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Key,
		&model.Name,
		&model.APIKey,
		&model.CreatedAt,
	)
	if err != nil {
		return domain.Environment{}, fmt.Errorf("insert environment: %w", err)
	}

	newEnv := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newEnv.ProjectID,
		domain.FeatureID(""), // No specific feature for environment creation
		domain.EntityEnvironment,
		newEnv.ID.String(),
		domain.AuditActionCreate,
		nil,
		newEnv,
		newEnv.ID,
	); err != nil {
		return domain.Environment{}, fmt.Errorf("audit environment create: %w", err)
	}

	return newEnv, nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.EnvironmentID) (domain.Environment, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM environments WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Environment{}, fmt.Errorf("query environment by id: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[environmentModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Environment{}, domain.ErrEntityNotFound
		}

		return domain.Environment{}, fmt.Errorf("collect environment row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByProjectIDAndKey(ctx context.Context, projectID domain.ProjectID, key string) (domain.Environment, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM environments WHERE project_id = $1::uuid AND key = $2 LIMIT 1`

	rows, err := executor.Query(ctx, query, projectID, key)
	if err != nil {
		return domain.Environment{}, fmt.Errorf("query environment by project_id and key: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[environmentModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Environment{}, domain.ErrEntityNotFound
		}

		return domain.Environment{}, fmt.Errorf("collect environment row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) ListByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.Environment, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM environments WHERE project_id = $1::uuid ORDER BY created_at ASC`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query environments by project_id: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[environmentModel])
	if err != nil {
		return nil, fmt.Errorf("collect environment rows: %w", err)
	}

	environments := make([]domain.Environment, 0, len(models))
	for _, m := range models {
		environments = append(environments, m.toDomain())
	}

	return environments, nil
}

func (r *Repository) GetByAPIKey(ctx context.Context, apiKey string) (domain.Environment, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM environments WHERE api_key = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, apiKey)
	if err != nil {
		return domain.Environment{}, fmt.Errorf("query environment by api_key: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[environmentModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Environment{}, domain.ErrEntityNotFound
		}

		return domain.Environment{}, fmt.Errorf("collect environment row: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, env domain.Environment) (domain.Environment, error) {
	executor := r.getExecutor(ctx)

	// Read old state for audit purposes within the same transaction.
	oldEnv, err := r.GetByID(ctx, env.ID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return domain.Environment{}, err
		}

		return domain.Environment{}, fmt.Errorf("get environment before update: %w", err)
	}

	const query = `
UPDATE environments
SET project_id = $1, key = $2, name = $3, api_key = $4
WHERE id = $5
RETURNING id, project_id, key, name, api_key, created_at`

	var model environmentModel

	err = executor.QueryRow(ctx, query,
		env.ProjectID,
		env.Key,
		env.Name,
		env.APIKey,
		env.ID,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.Key,
		&model.Name,
		&model.APIKey,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Environment{}, domain.ErrEntityNotFound
		}

		return domain.Environment{}, fmt.Errorf("update environment: %w", err)
	}

	newEnv := model.toDomain()
	if err := auditlog.Write(
		ctx,
		executor,
		newEnv.ProjectID,
		domain.FeatureID(""), // No specific feature for environment update
		domain.EntityEnvironment,
		newEnv.ID.String(),
		domain.AuditActionUpdate,
		oldEnv,
		newEnv,
		newEnv.ID,
	); err != nil {
		return domain.Environment{}, fmt.Errorf("audit environment update: %w", err)
	}

	return newEnv, nil
}

func (r *Repository) Delete(ctx context.Context, id domain.EnvironmentID) error {
	executor := r.getExecutor(ctx)

	// Read old state and write audit log before deletion.
	oldEnv, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := auditlog.Write(
		ctx,
		executor,
		oldEnv.ProjectID,
		domain.FeatureID(""), // No specific feature for environment deletion
		domain.EntityEnvironment,
		oldEnv.ID.String(),
		domain.AuditActionDelete,
		oldEnv,
		nil,
		oldEnv.ID,
	); err != nil {
		return fmt.Errorf("audit environment delete: %w", err)
	}

	const query = `DELETE FROM environments WHERE id = $1`

	ct, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete environment: %w", err)
	}

	// If nothing was deleted, return not found
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
