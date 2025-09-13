package rules

import (
	"time"

	"github.com/rom8726/etoggle/internal/domain"
)

type ruleModel struct {
	ID            string    `db:"id"`
	FeatureID     string    `db:"feature_id"`
	Condition     []byte    `db:"condition"`
	FlagVariantID string    `db:"flag_variant_id"`
	Priority      int       `db:"priority"`
	CreatedAt     time.Time `db:"created_at"`
}

func (m *ruleModel) toDomain() domain.Rule {
	return domain.Rule{
		ID:            domain.RuleID(m.ID),
		FeatureID:     domain.FeatureID(m.FeatureID),
		Condition:     m.Condition,
		FlagVariantID: domain.FlagVariantID(m.FlagVariantID),
		Priority:      uint8(m.Priority),
		CreatedAt:     m.CreatedAt,
	}
}
