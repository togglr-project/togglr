package project_approvers

import (
	"context"
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

type projectApproverModel struct {
	ProjectID string    `db:"project_id"`
	UserID    int       `db:"user_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

func (m *projectApproverModel) toDomain() domain.ProjectApprover {
	return domain.ProjectApprover{
		ProjectID: domain.ProjectID(m.ProjectID),
		UserID:    m.UserID,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
	}
}

// Create adds a new approver to a project.
func (r *Repository) Create(ctx context.Context, approver domain.ProjectApprover) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO project_approvers (project_id, user_id, role)
VALUES ($1, $2, $3)
ON CONFLICT (project_id, user_id) DO UPDATE SET role = EXCLUDED.role`

	_, err := executor.Exec(ctx, query,
		approver.ProjectID,
		approver.UserID,
		approver.Role,
	)
	if err != nil {
		return fmt.Errorf("insert project approver: %w", err)
	}

	return nil
}

// GetByProjectID retrieves all approvers for a project.
func (r *Repository) GetByProjectID(ctx context.Context, projectID domain.ProjectID) ([]domain.ProjectApprover, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT project_id, user_id, role, created_at
FROM project_approvers
WHERE project_id = $1::uuid
ORDER BY created_at`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project approvers: %w", err)
	}
	defer rows.Close()

	var approvers []domain.ProjectApprover

	for rows.Next() {
		var model projectApproverModel

		err := rows.Scan(
			&model.ProjectID,
			&model.UserID,
			&model.Role,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan project approver: %w", err)
		}

		approvers = append(approvers, model.toDomain())
	}

	return approvers, nil
}

// Delete removes an approver from a project.
func (r *Repository) Delete(ctx context.Context, projectID domain.ProjectID, userID int) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM project_approvers
WHERE project_id = $1::uuid AND user_id = $2`

	_, err := executor.Exec(ctx, query, projectID, userID)
	if err != nil {
		return fmt.Errorf("delete project approver: %w", err)
	}

	return nil
}

// IsUserApprover checks if a user is an approver for a project.
func (r *Repository) IsUserApprover(ctx context.Context, projectID domain.ProjectID, userID int) (bool, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT 1 FROM project_approvers
WHERE project_id = $1::uuid AND user_id = $2`

	var exists int

	err := executor.QueryRow(ctx, query, projectID, userID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("check user approver: %w", err)
	}

	return true, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
