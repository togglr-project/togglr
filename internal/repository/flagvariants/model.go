package flagvariants

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type flagVariantModel struct {
	ID             string    `db:"id"`
	ProjectID      string    `db:"project_id"`
	FeatureID      string    `db:"feature_id"`
	EnvironmentID  int64     `db:"environment_id"`
	Name           string    `db:"name"`
	RolloutPercent int       `db:"rollout_percent"`
	_createdAt     time.Time //nolint:unused // needed for pgx.CollectRows
}

func (m *flagVariantModel) toDomain() domain.FlagVariant {
	return domain.FlagVariant{
		ID:             domain.FlagVariantID(m.ID),
		ProjectID:      domain.ProjectID(m.ProjectID),
		FeatureID:      domain.FeatureID(m.FeatureID),
		EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
		Name:           m.Name,
		RolloutPercent: uint8(m.RolloutPercent),
	}
}

type flagVariantExtModel struct {
	ID             string `db:"id"`
	ProjectID      string `db:"project_id"`
	FeatureID      string `db:"feature_id"`
	EnvironmentID  int64  `db:"environment_id"`
	Name           string `db:"name"`
	RolloutPercent int    `db:"rollout_percent"`
	FeatureKey     string `db:"feature_key"`
	EnvKey         string `db:"env_key"`
}

func (m *flagVariantExtModel) toDomain() domain.FlagVariantExtended {
	return domain.FlagVariantExtended{
		FlagVariant: domain.FlagVariant{
			ID:             domain.FlagVariantID(m.ID),
			ProjectID:      domain.ProjectID(m.ProjectID),
			FeatureID:      domain.FeatureID(m.FeatureID),
			EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
			Name:           m.Name,
			RolloutPercent: uint8(m.RolloutPercent),
		},
		FeatureKey: m.FeatureKey,
		EnvKey:     m.EnvKey,
	}
}
