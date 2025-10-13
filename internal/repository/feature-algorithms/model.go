package feature_algorithms

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type featureAlgorithmModel struct {
	EnvironmentID int64           `db:"environment_id"`
	FeatureID     string          `db:"feature_id"`
	AlgorithmID   string          `db:"algorithm_id"`
	FlagVariantID sql.NullString  `db:"flag_variant_id"`
	Settings      json.RawMessage `db:"settings"`
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

	var flagVariantIDRef *domain.FlagVariantID
	if m.FlagVariantID.Valid {
		flagVariantID := domain.FlagVariantID(m.FlagVariantID.String)
		flagVariantIDRef = &flagVariantID
	}

	return domain.FeatureAlgorithm{
		EnvironmentID: domain.EnvironmentID(m.EnvironmentID),
		FeatureID:     domain.FeatureID(m.FeatureID),
		AlgorithmID:   domain.AlgorithmID(m.AlgorithmID),
		FlagVariantID: flagVariantIDRef,
		Settings:      settings,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
