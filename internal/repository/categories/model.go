package categories

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type categoryModel struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Description *string   `db:"description"`
	Color       *string   `db:"color"`
	Kind        string    `db:"kind"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m *categoryModel) toDomain() domain.Category {
	return domain.Category{
		ID:          domain.CategoryID(m.ID),
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
		Color:       m.Color,
		Kind:        domain.CategoryKind(m.Kind),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
