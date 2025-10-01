package rbac

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type roleModel struct {
	ID          string    `db:"id"`
	Key         string    `db:"key"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (role *roleModel) toDomain() domain.Role {
	return domain.Role{
		ID:          domain.RoleID(role.ID),
		Key:         role.Key,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
	}
}
