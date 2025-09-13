package domain

type FlagVariantID string

type FlagVariant struct {
	ID             FlagVariantID
	FeatureID      FeatureID
	Name           string // e.g. "A", "B"
	RolloutPercent uint8  // % of traffic (0..100)
}

func (id FlagVariantID) String() string {
	return string(id)
}
