package features

import (
	"database/sql"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type baseFeatureModel struct {
	ID          string         `db:"id"`
	ProjectID   string         `db:"project_id"`
	Key         string         `db:"key"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Kind        string         `db:"kind"`
	RolloutKey  sql.NullString `db:"rollout_key"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (f *baseFeatureModel) toDomain() domain.BasicFeature {
	return domain.BasicFeature{
		ID:          domain.FeatureID(f.ID),
		ProjectID:   domain.ProjectID(f.ProjectID),
		Key:         f.Key,
		Name:        f.Name,
		Description: f.Description.String,
		Kind:        domain.FeatureKind(f.Kind),
		RolloutKey:  domain.RuleAttribute(f.RolloutKey.String),
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

type featureFullModel struct {
	ID             string         `db:"id"`
	ProjectID      string         `db:"project_id"`
	EnvironmentID  int64          `db:"environment_id"`
	EnvironmentKey string         `db:"environment_key"`
	Key            string         `db:"key"`
	Name           string         `db:"name"`
	Description    sql.NullString `db:"description"`
	Kind           string         `db:"kind"`
	RolloutKey     sql.NullString `db:"rollout_key"`
	Enabled        bool           `db:"enabled"`
	DefaultValue   string         `db:"default_value"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

func (f *featureFullModel) toDomain() domain.Feature {
	return domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          domain.FeatureID(f.ID),
			ProjectID:   domain.ProjectID(f.ProjectID),
			Key:         f.Key,
			Name:        f.Name,
			Description: f.Description.String,
			Kind:        domain.FeatureKind(f.Kind),
			RolloutKey:  domain.RuleAttribute(f.RolloutKey.String),
			CreatedAt:   f.CreatedAt,
			UpdatedAt:   f.UpdatedAt,
		},
		EnvironmentID: domain.EnvironmentID(f.EnvironmentID),
		Enabled:       f.Enabled,
		DefaultValue:  f.DefaultValue,
	}
}
