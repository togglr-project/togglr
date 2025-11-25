package feature_contextual_stats

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type featureContextualStatsModel struct {
	ProjectID      string          `db:"project_id"`
	EnvironmentID  int64           `db:"environment_id"`
	FeatureID      string          `db:"feature_id"`
	AlgorithmSlug  string          `db:"algorithm_slug"`
	VariantKey     string          `db:"variant_key"`
	FeatureKey     string          `db:"feature_key"`
	EnvironmentKey string          `db:"environment_key"`
	FeatureDim     int             `db:"feature_dim"`
	MatrixA        []float64       `db:"matrix_a"`
	VectorB        []float64       `db:"vector_b"`
	Pulls          uint64          `db:"pulls"`
	TotalReward    decimal.Decimal `db:"total_reward"`
	Successes      uint64          `db:"successes"`
	Failures       uint64          `db:"failures"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

func (m *featureContextualStatsModel) toDomain() domain.FeatureContextualStats {
	return domain.FeatureContextualStats{
		ProjectID:      domain.ProjectID(m.ProjectID),
		EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
		FeatureID:      domain.FeatureID(m.FeatureID),
		AlgorithmSlug:  m.AlgorithmSlug,
		VariantKey:     m.VariantKey,
		FeatureKey:     m.FeatureKey,
		EnvironmentKey: m.EnvironmentKey,
		FeatureDim:     m.FeatureDim,
		MatrixA:        m.MatrixA,
		VectorB:        m.VectorB,
		Pulls:          m.Pulls,
		TotalReward:    m.TotalReward,
		Successes:      m.Successes,
		Failures:       m.Failures,
		UpdatedAt:      m.UpdatedAt,
	}
}
