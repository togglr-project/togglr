package pending_changes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/contract"
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

type pendingChangeModel struct {
	ID              string          `db:"id"`
	ProjectID       string          `db:"project_id"`
	RequestedBy     string          `db:"requested_by"`
	RequestUserID   *int            `db:"request_user_id"`
	Change          json.RawMessage `db:"change"`
	Status          string          `db:"status"`
	CreatedAt       time.Time       `db:"created_at"`
	ApprovedBy      *string         `db:"approved_by"`
	ApprovedUserID  *int            `db:"approved_user_id"`
	ApprovedAt      *time.Time      `db:"approved_at"`
	RejectedBy      *string         `db:"rejected_by"`
	RejectedAt      *time.Time      `db:"rejected_at"`
	RejectionReason *string         `db:"rejection_reason"`
}

func (m *pendingChangeModel) toDomain() (domain.PendingChange, error) {
	var change domain.PendingChangePayload
	if err := json.Unmarshal(m.Change, &change); err != nil {
		return domain.PendingChange{}, fmt.Errorf("unmarshal change: %w", err)
	}

	return domain.PendingChange{
		ID:              domain.PendingChangeID(m.ID),
		ProjectID:       domain.ProjectID(m.ProjectID),
		RequestedBy:     m.RequestedBy,
		RequestUserID:   m.RequestUserID,
		Change:          change,
		Status:          domain.PendingChangeStatus(m.Status),
		CreatedAt:       m.CreatedAt,
		ApprovedBy:      m.ApprovedBy,
		ApprovedUserID:  m.ApprovedUserID,
		ApprovedAt:      m.ApprovedAt,
		RejectedBy:      m.RejectedBy,
		RejectedAt:      m.RejectedAt,
		RejectionReason: m.RejectionReason,
	}, nil
}

type pendingChangeEntityModel struct {
	ID              string    `db:"id"`
	PendingChangeID string    `db:"pending_change_id"`
	Entity          string    `db:"entity"`
	EntityID        string    `db:"entity_id"`
	CreatedAt       time.Time `db:"created_at"`
}

func (m *pendingChangeEntityModel) toDomain() domain.PendingChangeEntity {
	return domain.PendingChangeEntity{
		ID:              m.ID,
		PendingChangeID: domain.PendingChangeID(m.PendingChangeID),
		Entity:          m.Entity,
		EntityID:        m.EntityID,
		CreatedAt:       m.CreatedAt,
	}
}

// Create creates a new pending change and its entities
func (r *Repository) Create(
	ctx context.Context,
	projectID domain.ProjectID,
	requestedBy string,
	requestUserID *int,
	change domain.PendingChangePayload,
) (domain.PendingChange, error) {
	executor := r.getExecutor(ctx)

	// First check for conflicts
	if err := r.CheckEntityConflict(ctx, change.Entities); err != nil {
		return domain.PendingChange{}, err
	}

	// Marshal change to JSON
	changeJSON, err := json.Marshal(change)
	if err != nil {
		return domain.PendingChange{}, fmt.Errorf("marshal change: %w", err)
	}

	// Insert pending change
	const insertQuery = `
INSERT INTO pending_changes (project_id, requested_by, request_user_id, change, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, project_id, requested_by, request_user_id, change, status, created_at, 
          approved_by, approved_user_id, approved_at, rejected_by, rejected_at, rejection_reason`

	var model pendingChangeModel
	err = executor.QueryRow(ctx, insertQuery,
		projectID,
		requestedBy,
		requestUserID,
		changeJSON,
		domain.PendingChangeStatusPending,
	).Scan(
		&model.ID,
		&model.ProjectID,
		&model.RequestedBy,
		&model.RequestUserID,
		&model.Change,
		&model.Status,
		&model.CreatedAt,
		&model.ApprovedBy,
		&model.ApprovedUserID,
		&model.ApprovedAt,
		&model.RejectedBy,
		&model.RejectedAt,
		&model.RejectionReason,
	)
	if err != nil {
		return domain.PendingChange{}, fmt.Errorf("insert pending change: %w", err)
	}

	// Insert entities
	for _, entity := range change.Entities {
		const entityQuery = `
INSERT INTO pending_change_entities (pending_change_id, entity, entity_id)
VALUES ($1, $2, $3)`

		_, err = executor.Exec(ctx, entityQuery,
			model.ID,
			entity.Entity,
			entity.EntityID,
		)
		if err != nil {
			return domain.PendingChange{}, fmt.Errorf("insert pending change entity: %w", err)
		}
	}

	return model.toDomain()
}

// GetByID retrieves a pending change by ID
func (r *Repository) GetByID(ctx context.Context, id domain.PendingChangeID) (domain.PendingChange, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT id, project_id, requested_by, request_user_id, change, status, created_at,
       approved_by, approved_user_id, approved_at, rejected_by, rejected_at, rejection_reason
FROM pending_changes
WHERE id = $1`

	var model pendingChangeModel
	err := executor.QueryRow(ctx, query, id).Scan(
		&model.ID,
		&model.ProjectID,
		&model.RequestedBy,
		&model.RequestUserID,
		&model.Change,
		&model.Status,
		&model.CreatedAt,
		&model.ApprovedBy,
		&model.ApprovedUserID,
		&model.ApprovedAt,
		&model.RejectedBy,
		&model.RejectedAt,
		&model.RejectionReason,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.PendingChange{}, domain.ErrEntityNotFound
		}
		return domain.PendingChange{}, fmt.Errorf("get pending change: %w", err)
	}

	return model.toDomain()
}

// List retrieves pending changes with filtering
func (r *Repository) List(
	ctx context.Context,
	filter contract.PendingChangesListFilter,
) ([]domain.PendingChange, int, error) {
	executor := r.getExecutor(ctx)

	// Build query
	query := `
SELECT id, project_id, requested_by, request_user_id, change, status, created_at,
       approved_by, approved_user_id, approved_at, rejected_by, rejected_at, rejection_reason
FROM pending_changes
WHERE 1=1`

	args := []interface{}{}
	argIndex := 1

	if filter.ProjectID != nil {
		query += fmt.Sprintf(" AND project_id = $%d", argIndex)
		args = append(args, *filter.ProjectID)
		argIndex++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND request_user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	// Add sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortDir := "ASC"
	if filter.SortDesc {
		sortDir = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)

	// Add pagination
	if filter.PerPage > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.PerPage)
		argIndex++

		if filter.Page > 0 {
			offset := (filter.Page - 1) * filter.PerPage
			query += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, offset)
		}
	}

	rows, err := executor.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list pending changes: %w", err)
	}
	defer rows.Close()

	var changes []domain.PendingChange
	for rows.Next() {
		var model pendingChangeModel
		err := rows.Scan(
			&model.ID,
			&model.ProjectID,
			&model.RequestedBy,
			&model.RequestUserID,
			&model.Change,
			&model.Status,
			&model.CreatedAt,
			&model.ApprovedBy,
			&model.ApprovedUserID,
			&model.ApprovedAt,
			&model.RejectedBy,
			&model.RejectedAt,
			&model.RejectionReason,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan pending change: %w", err)
		}

		change, err := model.toDomain()
		if err != nil {
			return nil, 0, err
		}

		changes = append(changes, change)
	}

	// Get total count
	countQuery := `
SELECT COUNT(*)
FROM pending_changes
WHERE 1=1`

	countArgs := []interface{}{}
	countArgIndex := 1

	if filter.ProjectID != nil {
		countQuery += fmt.Sprintf(" AND project_id = $%d", countArgIndex)
		countArgs = append(countArgs, *filter.ProjectID)
		countArgIndex++
	}

	if filter.Status != nil {
		countQuery += fmt.Sprintf(" AND status = $%d", countArgIndex)
		countArgs = append(countArgs, *filter.Status)
		countArgIndex++
	}

	if filter.UserID != nil {
		countQuery += fmt.Sprintf(" AND request_user_id = $%d", countArgIndex)
		countArgs = append(countArgs, *filter.UserID)
		countArgIndex++
	}

	var total int
	err = executor.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count pending changes: %w", err)
	}

	return changes, total, nil
}

// UpdateStatus updates the status of a pending change
func (r *Repository) UpdateStatus(
	ctx context.Context,
	id domain.PendingChangeID,
	status domain.PendingChangeStatus,
	approvedBy *string,
	approvedUserID *int,
	approvedAt *string,
	rejectedBy *string,
	rejectedAt *string,
	rejectionReason *string,
) error {
	executor := r.getExecutor(ctx)

	query := `
UPDATE pending_changes
SET status = $2, approved_by = $3, approved_user_id = $4, approved_at = $5,
    rejected_by = $6, rejected_at = $7, rejection_reason = $8
WHERE id = $1`

	var approvedAtTime *time.Time
	var rejectedAtTime *time.Time

	if approvedAt != nil {
		if t, err := time.Parse(time.RFC3339, *approvedAt); err == nil {
			approvedAtTime = &t
		}
	}

	if rejectedAt != nil {
		if t, err := time.Parse(time.RFC3339, *rejectedAt); err == nil {
			rejectedAtTime = &t
		}
	}

	_, err := executor.Exec(ctx, query,
		id,
		status,
		approvedBy,
		approvedUserID,
		approvedAtTime,
		rejectedBy,
		rejectedAtTime,
		rejectionReason,
	)
	if err != nil {
		return fmt.Errorf("update pending change status: %w", err)
	}

	return nil
}

// CheckEntityConflict checks if there are any pending changes for the given entities
func (r *Repository) CheckEntityConflict(
	ctx context.Context,
	entities []domain.EntityChange,
) error {
	executor := r.getExecutor(ctx)

	for _, entity := range entities {
		const query = `
SELECT 1 FROM pending_change_entities pce
JOIN pending_changes pc ON pc.id = pce.pending_change_id
WHERE pce.entity = $1 AND pce.entity_id = $2 AND pc.status = 'pending'`

		var exists int
		err := executor.QueryRow(ctx, query, entity.Entity, entity.EntityID).Scan(&exists)
		if err == nil {
			return fmt.Errorf("entity %s %s is already locked by another pending change", entity.Entity, entity.EntityID)
		}
		if err != pgx.ErrNoRows {
			return fmt.Errorf("check entity conflict: %w", err)
		}
	}

	return nil
}

// GetEntitiesByPendingChangeID retrieves all entities for a pending change
func (r *Repository) GetEntitiesByPendingChangeID(
	ctx context.Context,
	pendingChangeID domain.PendingChangeID,
) ([]domain.PendingChangeEntity, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT id, pending_change_id, entity, entity_id, created_at
FROM pending_change_entities
WHERE pending_change_id = $1`

	rows, err := executor.Query(ctx, query, pendingChangeID)
	if err != nil {
		return nil, fmt.Errorf("get pending change entities: %w", err)
	}
	defer rows.Close()

	var entities []domain.PendingChangeEntity
	for rows.Next() {
		var model pendingChangeEntityModel
		err := rows.Scan(
			&model.ID,
			&model.PendingChangeID,
			&model.Entity,
			&model.EntityID,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan pending change entity: %w", err)
		}

		entities = append(entities, model.toDomain())
	}

	return entities, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
