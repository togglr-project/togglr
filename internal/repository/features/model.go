package features

import (
	"database/sql"
	"time"

	"github.com/rom8726/etoggle/internal/domain"
)

type featureModel struct {
	ID             string         `db:"id"`
	ProjectID      string         `db:"project_id"`
	Key            string         `db:"key"`
	Name           string         `db:"name"`
	Description    sql.NullString `db:"description"`
	Kind           string         `db:"kind"`
	DefaultVariant string         `db:"default_variant"`
	RolloutKey     sql.NullString `db:"rollout_key"`
	Enabled        bool           `db:"enabled"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

func (f *featureModel) toDomain() domain.Feature {
	return domain.Feature{
		ID:             domain.FeatureID(f.ID),
		ProjectID:      domain.ProjectID(f.ProjectID),
		Key:            f.Key,
		Name:           f.Name,
		Description:    f.Description.String,
		Kind:           domain.FeatureKind(f.Kind),
		DefaultVariant: f.DefaultVariant,
		RolloutKey:     domain.RuleAttribute(f.RolloutKey.String),
		Enabled:        f.Enabled,
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
	}
}
