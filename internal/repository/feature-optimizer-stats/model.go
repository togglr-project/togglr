package feature_optimizer_stats

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type featureOptimizerStatsModel struct {
	ProjectID      string          `db:"project_id"`
	EnvironmentID  int64           `db:"environment_id"`
	FeatureID      string          `db:"feature_id"`
	AlgorithmSlug  string          `db:"algorithm_slug"`
	FeatureKey     string          `db:"feature_key"`
	EnvironmentKey string          `db:"environment_key"`
	Iteration      uint64          `db:"iteration"`
	CurrentValue   decimal.Decimal `db:"current_value"`
	BestValue      decimal.Decimal `db:"best_value"`
	BestReward     decimal.Decimal `db:"best_reward"`
	MetricSum      decimal.Decimal `db:"metric_sum"`
	LastError      decimal.Decimal `db:"last_error"`
	Integral       decimal.Decimal `db:"integral"`
	StepSize       decimal.Decimal `db:"step_size"`
	Temperature    decimal.Decimal `db:"temperature"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

func (m *featureOptimizerStatsModel) toDomain() domain.FeatureOptimizerStats {
	return domain.FeatureOptimizerStats{
		ProjectID:      domain.ProjectID(m.ProjectID),
		EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
		FeatureID:      domain.FeatureID(m.FeatureID),
		AlgorithmSlug:  m.AlgorithmSlug,
		FeatureKey:     m.FeatureKey,
		EnvironmentKey: m.EnvironmentKey,
		Iteration:      m.Iteration,
		CurrentValue:   m.CurrentValue,
		BestValue:      m.BestValue,
		BestReward:     m.BestReward,
		MetricSum:      m.MetricSum,
		LastError:      m.LastError,
		Integral:       m.Integral,
		StepSize:       m.StepSize,
		Temperature:    m.Temperature,
		UpdatedAt:      m.UpdatedAt,
	}
}
