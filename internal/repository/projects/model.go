package projects

import (
	"database/sql"
	"time"

	"github.com/rom8726/etoggle/internal/domain"
)

type projectModel struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	APIKey      string         `db:"api_key"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	ArchivedAt  *time.Time     `db:"archived_at"`
}

func (m *projectModel) toDomain() domain.Project {
	return domain.Project{
		ID:          domain.ProjectID(m.ID),
		Name:        m.Name,
		Description: m.Description.String,
		APIKey:      m.APIKey,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		ArchivedAt:  m.ArchivedAt,
	}
}
