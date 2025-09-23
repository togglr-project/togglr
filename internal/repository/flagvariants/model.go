package flagvariants

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type flagVariantModel struct {
	ID             string `db:"id"`
	ProjectID      string `db:"project_id"`
	FeatureID      string `db:"feature_id"`
	Name           string `db:"name"`
	RolloutPercent int    `db:"rollout_percent"`
	// no timestamps in table; keeping struct minimal
	_createdAt time.Time // placeholder to avoid empty import for time in case of future changes
}

func (m *flagVariantModel) toDomain() domain.FlagVariant {
	return domain.FlagVariant{
		ID:             domain.FlagVariantID(m.ID),
		ProjectID:      domain.ProjectID(m.ProjectID),
		FeatureID:      domain.FeatureID(m.FeatureID),
		Name:           m.Name,
		RolloutPercent: uint8(m.RolloutPercent),
	}
}
