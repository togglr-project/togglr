package customalgorithms

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type customAlgorithmModel struct {
	ID              string          `db:"id"`
	Slug            string          `db:"slug"`
	Name            string          `db:"name"`
	Description     *string         `db:"description"`
	Kind            string          `db:"kind"`
	WASMBinary      []byte          `db:"wasm_binary"`
	WASMHash        string          `db:"wasm_hash"`
	DefaultSettings json.RawMessage `db:"default_settings"`
	CreatedBy       *int64          `db:"created_by"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

func (m *customAlgorithmModel) toDomain() domain.CustomAlgorithm {
	var settings map[string]decimal.Decimal
	if m.DefaultSettings != nil {
		if err := json.Unmarshal(m.DefaultSettings, &settings); err != nil {
			slog.Error("unmarshal custom algorithm settings", "settings", string(m.DefaultSettings), "error", err)
		}
	}

	var description string
	if m.Description != nil {
		description = *m.Description
	}

	var createdBy *domain.UserID
	if m.CreatedBy != nil {
		uid := domain.UserID(*m.CreatedBy)
		createdBy = &uid
	}

	return domain.CustomAlgorithm{
		ID:              domain.CustomAlgorithmID(m.ID),
		Slug:            m.Slug,
		Name:            m.Name,
		Description:     description,
		Kind:            domain.AlgorithmKind(m.Kind),
		WASMBinary:      m.WASMBinary,
		WASMHash:        m.WASMHash,
		DefaultSettings: settings,
		CreatedBy:       createdBy,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

type customAlgorithmStatsModel struct {
	ProjectID      string          `db:"project_id"`
	FeatureID      string          `db:"feature_id"`
	EnvironmentID  int64           `db:"environment_id"`
	AlgorithmID    string          `db:"algorithm_id"`
	VariantKey     string          `db:"variant_key"`
	FeatureKey     string          `db:"feature_key"`
	EnvironmentKey string          `db:"environment_key"`
	State          json.RawMessage `db:"state"`
	Evaluations    int64           `db:"evaluations"`
	Successes      int64           `db:"successes"`
	Failures       int64           `db:"failures"`
	MetricSum      decimal.Decimal `db:"metric_sum"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

func (m *customAlgorithmStatsModel) toDomain() domain.CustomAlgorithmStats {
	return domain.CustomAlgorithmStats{
		ProjectID:      domain.ProjectID(m.ProjectID),
		FeatureID:      domain.FeatureID(m.FeatureID),
		EnvironmentID:  domain.EnvironmentID(m.EnvironmentID),
		AlgorithmID:    domain.CustomAlgorithmID(m.AlgorithmID),
		VariantKey:     m.VariantKey,
		FeatureKey:     m.FeatureKey,
		EnvironmentKey: m.EnvironmentKey,
		State:          m.State,
		Evaluations:    uint64(m.Evaluations),
		Successes:      uint64(m.Successes),
		Failures:       uint64(m.Failures),
		MetricSum:      m.MetricSum,
		UpdatedAt:      m.UpdatedAt,
	}
}
