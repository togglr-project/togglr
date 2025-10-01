package rbac

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/pkg/db"
)

// Roles repository implementation.
// It implements contract.RolesRepository.
type Roles struct {
	db db.Tx
}

func NewRoles(pool *pgxpool.Pool) *Roles {
	return &Roles{db: pool}
}

func (r *Roles) GetByKey(ctx context.Context, key string) (domain.Role, error) {
	exec := getExecutor(ctx, r.db)

	const query = `select * from roles where key = $1 limit 1`

	rows, err := exec.Query(ctx, query, key)
	if err != nil {
		return domain.Role{}, fmt.Errorf("get role by key: %w", err)
	}
	defer rows.Close()

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roleModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Role{}, domain.ErrEntityNotFound
		}

		return domain.Role{}, fmt.Errorf("collect role: %w", err)
	}

	return role.toDomain(), nil
}

func (r *Roles) GetByID(ctx context.Context, id domain.RoleID) (domain.Role, error) {
	exec := getExecutor(ctx, r.db)

	const query = `select * from roles where id = $1 limit 1`

	rows, err := exec.Query(ctx, query, id)
	if err != nil {
		return domain.Role{}, fmt.Errorf("get role by key: %w", err)
	}
	defer rows.Close()

	role, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[roleModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Role{}, domain.ErrEntityNotFound
		}

		return domain.Role{}, fmt.Errorf("collect role: %w", err)
	}

	return role.toDomain(), nil
}

var _ contract.RolesRepository = (*Roles)(nil)

// Permissions repository implementation.
// It implements contract.PermissionsRepository.
type Permissions struct {
	db db.Tx
}

func NewPermissions(pool *pgxpool.Pool) *Permissions {
	return &Permissions{db: pool}
}

func (r *Permissions) RoleHasPermission(
	ctx context.Context,
	roleID string,
	key domain.PermKey,
) (bool, error) {
	exec := getExecutor(ctx, r.db)

	const query = `select exists(
		select 1 from role_permissions rp
		join permissions p on p.id = rp.permission_id
		where rp.role_id = $1 and p.key = $2
	)`

	var has bool
	if err := exec.QueryRow(ctx, query, roleID, string(key)).Scan(&has); err != nil {
		return false, fmt.Errorf("role has permission: %w", err)
	}

	return has, nil
}

var _ contract.PermissionsRepository = (*Permissions)(nil)

// Memberships repository implementation.
// It implements contract.MembershipsRepository.
type Memberships struct {
	db db.Tx
}

func NewMemberships(pool *pgxpool.Pool) *Memberships {
	return &Memberships{db: pool}
}

func (r *Memberships) GetForUserProject(
	ctx context.Context,
	userID int,
	projectID domain.ProjectID,
) (string, error) { // roleID
	exec := getExecutor(ctx, r.db)

	const query = `select role_id from memberships where project_id = $1::uuid and user_id = $2 limit 1`

	var roleID string
	if err := exec.QueryRow(ctx, query, projectID, userID).Scan(&roleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("get membership for user/project: %w", err)
	}

	return roleID, nil
}

var _ contract.MembershipsRepository = (*Memberships)(nil)

// helper to get tx from context
//
//nolint:ireturn // internal helper
func getExecutor(ctx context.Context, fallback db.Tx) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return fallback
}
