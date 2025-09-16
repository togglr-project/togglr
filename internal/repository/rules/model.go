package rules

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/rom8726/etoggle/internal/domain"
)

type ruleModel struct {
	ID            string         `db:"id"`
	ProjectID     string         `db:"project_id"`
	FeatureID     string         `db:"feature_id"`
	Condition     []byte         `db:"condition"`
	Action        string         `db:"action"`
	FlagVariantID sql.NullString `db:"flag_variant_id"`
	Priority      int            `db:"priority"`
	CreatedAt     time.Time      `db:"created_at"`
}

func (m *ruleModel) toDomain() domain.Rule {
	var conditions domain.Conditions
	err := json.Unmarshal(m.Condition, &conditions)
	if err != nil {
		slog.Error("unmarshal rule condition", "conditions", string(m.Condition), "error", err)
	}

	var flagVariantID *domain.FlagVariantID
	if m.FlagVariantID.Valid {
		variantID := domain.FlagVariantID(m.FlagVariantID.String)
		flagVariantID = &variantID
	}

	return domain.Rule{
		ID:            domain.RuleID(m.ID),
		ProjectID:     domain.ProjectID(m.ProjectID),
		FeatureID:     domain.FeatureID(m.FeatureID),
		Conditions:    conditions,
		Action:        domain.RuleAction(m.Action),
		FlagVariantID: flagVariantID,
		Priority:      uint8(m.Priority),
		CreatedAt:     m.CreatedAt,
	}
}
