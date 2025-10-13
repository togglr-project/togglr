package algorithms

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"

	"github.com/togglr-project/togglr/internal/domain"
)

type algorithmModel struct {
	ID              string          `db:"id"`
	Name            string          `db:"name"`
	Description     string          `db:"description"`
	Slug            string          `db:"slug"`
	Kind            string          `db:"kind"`
	DefaultSettings json.RawMessage `db:"default_settings"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

func (m *algorithmModel) toDomain() domain.Algorithm {
	var settings map[string]decimal.Decimal
	if m.DefaultSettings != nil {
		if err := json.Unmarshal(m.DefaultSettings, &settings); err != nil {
			slog.Error("unmarshal algorithm settings", "settings", string(m.DefaultSettings), "error", err)
		}
	}

	return domain.Algorithm{
		ID:              domain.AlgorithmID(m.ID),
		Name:            m.Name,
		Slug:            m.Slug,
		Kind:            domain.AlgorithmKind(m.Kind),
		Description:     m.Description,
		DefaultSettings: settings,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
