package feature_algorithm_stats

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type featureAlgorithmStatsModel struct {
	FeatureID     string          `db:"feature_id"`
	EnvironmentID int64           `db:"environment_id"`
	AlgorithmSlug string          `db:"algorithm_slug"`
	VariantKey    string          `db:"variant_key"`
	Evaluations   uint64          `db:"evaluations"`
	Successes     uint64          `db:"successes"`
	Failures      uint64          `db:"failures"`
	MetricSum     decimal.Decimal `db:"metric_sum"`
	UpdatedAt     time.Time       `db:"updated_at"`
}

func (m *featureAlgorithmStatsModel) toDomain() domain.FeatureAlgorithmStats {
	return domain.FeatureAlgorithmStats{
		FeatureID:     domain.FeatureID(m.FeatureID),
		EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
		AlgorithmSlug: m.AlgorithmSlug,
		VariantKey:    m.VariantKey,
		Evaluations:   m.Evaluations,
		Successes:     m.Successes,
		Failures:      m.Failures,
		MetricSum:     m.MetricSum,
		UpdatedAt:     m.UpdatedAt,
	}
}
