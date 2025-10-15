package feature_algorithms

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type featureAlgorithmModel struct {
	EnvironmentID int64           `db:"environment_id"`
	FeatureID     string          `db:"feature_id"`
	AlgorithmSlug string          `db:"algorithm_slug"`
	Settings      json.RawMessage `db:"settings"`
	Enabled       bool            `db:"enabled"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}

func (m *featureAlgorithmModel) toDomain() domain.FeatureAlgorithm {
	var settings map[string]decimal.Decimal
	if m.Settings != nil {
		if err := json.Unmarshal(m.Settings, &settings); err != nil {
			slog.Error("unmarshal feature algorithm settings", "settings", string(m.Settings), "error", err)
		}
	}

	return domain.FeatureAlgorithm{
		EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
		FeatureID:     domain.FeatureID(m.FeatureID),
		AlgorithmSlug: m.AlgorithmSlug,
		Settings:      settings,
		Enabled:       m.Enabled,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

type featureAlgorithmExtModel struct {
	EnvironmentID int64           `db:"environment_id"`
	FeatureID     string          `db:"feature_id"`
	AlgorithmSlug string          `db:"algorithm_slug"`
	Settings      json.RawMessage `db:"settings"`
	Enabled       bool            `db:"enabled"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
	FeatureKey    string          `db:"feature_key"`
	EnvKey        string          `db:"env_key"`
}

func (m *featureAlgorithmExtModel) toDomain() domain.FeatureAlgorithmExtended {
	var settings map[string]decimal.Decimal
	if m.Settings != nil {
		if err := json.Unmarshal(m.Settings, &settings); err != nil {
			slog.Error("unmarshal feature algorithm settings", "settings", string(m.Settings), "error", err)
		}
	}

	return domain.FeatureAlgorithmExtended{
		FeatureAlgorithm: domain.FeatureAlgorithm{
			EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
			FeatureID:     domain.FeatureID(m.FeatureID),
			AlgorithmSlug: m.AlgorithmSlug,
			Settings:      settings,
			Enabled:       m.Enabled,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		},
		FeatureKey: m.FeatureKey,
		EnvKey:     m.EnvKey,
	}
}
