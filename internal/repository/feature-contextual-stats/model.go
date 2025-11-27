package feature_contextual_stats

import (
	"encoding/json"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type rawModel struct {
	ProjectID      string          `db:"project_id"`
	EnvironmentID  int64           `db:"environment_id"`
	FeatureID      string          `db:"feature_id"`
	AlgorithmSlug  string          `db:"algorithm_slug"`
	VariantKey     string          `db:"variant_key"`
	FeatureKey     string          `db:"feature_key"`
	EnvironmentKey string          `db:"environment_key"`
	FeatureDim     int             `db:"feature_dim"`
	MatrixA        []byte          `db:"matrix_a"`
	VectorB        []byte          `db:"vector_b"`
	Pulls          uint64          `db:"pulls"`
	TotalReward    decimal.Decimal `db:"total_reward"`
	Successes      uint64          `db:"successes"`
	Failures       uint64          `db:"failures"`
}

func (m *rawModel) toDomain() domain.FeatureContextualStats {
	var matrixA []float64
	var vectorB []float64

	if len(m.MatrixA) > 0 {
		_ = json.Unmarshal(m.MatrixA, &matrixA)
	}

	if len(m.VectorB) > 0 {
		_ = json.Unmarshal(m.VectorB, &vectorB)
	}

	return domain.FeatureContextualStats{
		ProjectID:      domain.ProjectID(m.ProjectID),
		EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
		FeatureID:      domain.FeatureID(m.FeatureID),
		AlgorithmSlug:  m.AlgorithmSlug,
		VariantKey:     m.VariantKey,
		FeatureKey:     m.FeatureKey,
		EnvironmentKey: m.EnvironmentKey,
		FeatureDim:     m.FeatureDim,
		MatrixA:        matrixA,
		VectorB:        vectorB,
		Pulls:          m.Pulls,
		TotalReward:    m.TotalReward,
		Successes:      m.Successes,
		Failures:       m.Failures,
	}
}
