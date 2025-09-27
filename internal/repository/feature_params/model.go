package feature_params

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type featureParamsModel struct {
	FeatureID     string    `db:"feature_id"`
	EnvironmentID int64     `db:"environment_id"`
	Enabled       bool      `db:"enabled"`
	DefaultValue  string    `db:"default_value"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (m featureParamsModel) toDomain() domain.FeatureParams {
	return domain.FeatureParams{
		FeatureID:     domain.FeatureID(m.FeatureID),
		EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
		Enabled:       m.Enabled,
		DefaultValue:  m.DefaultValue,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func featureParamsFromDomain(params domain.FeatureParams) featureParamsModel {
	return featureParamsModel{
		FeatureID:     string(params.FeatureID),
		EnvironmentID: int64(params.EnvironmentID),
		Enabled:       params.Enabled,
		DefaultValue:  params.DefaultValue,
		CreatedAt:     params.CreatedAt,
		UpdatedAt:     params.UpdatedAt,
	}
}
