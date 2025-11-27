package domain

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type CustomAlgorithmID string

func (id CustomAlgorithmID) String() string {
	return string(id)
}

// CustomAlgorithm represents a user-defined WASM algorithm.
type CustomAlgorithm struct {
	ID              CustomAlgorithmID          `db:"id"               pk:"true"`
	Slug            string                     `db:"slug"`
	Name            string                     `db:"name"`
	Description     string                     `db:"description"`
	Kind            AlgorithmKind              `db:"kind"`
	WASMBinary      []byte                     `db:"wasm_binary"`
	WASMHash        string                     `db:"wasm_hash"`
	DefaultSettings map[string]decimal.Decimal `db:"default_settings"`
	CreatedBy       *UserID                    `db:"created_by"`
	CreatedAt       time.Time                  `db:"created_at"`
	UpdatedAt       time.Time                  `db:"updated_at"`
}

// CustomAlgorithmDTO is used for creating/updating custom algorithms.
type CustomAlgorithmDTO struct {
	Slug            string
	Name            string
	Description     string
	Kind            AlgorithmKind
	WASMBinary      []byte
	DefaultSettings map[string]decimal.Decimal
	CreatedBy       *UserID
}

// CustomAlgorithmStats stores the state and statistics for a custom algorithm instance.
type CustomAlgorithmStats struct {
	ProjectID      ProjectID         `db:"project_id"`
	FeatureID      FeatureID         `db:"feature_id"`
	EnvironmentID  EnvironmentID     `db:"environment_id"`
	AlgorithmID    CustomAlgorithmID `db:"algorithm_id"`
	VariantKey     string            `db:"variant_key"`
	FeatureKey     string            `db:"feature_key"`
	EnvironmentKey string            `db:"environment_key"`
	State          json.RawMessage   `db:"state"`
	Evaluations    uint64            `db:"evaluations"`
	Successes      uint64            `db:"successes"`
	Failures       uint64            `db:"failures"`
	MetricSum      decimal.Decimal   `db:"metric_sum"`
	UpdatedAt      time.Time         `db:"updated_at"`
}

// WASMInput represents input data sent to WASM algorithms.
type WASMInput struct {
	// Common fields
	Settings map[string]float64 `json:"settings"`
	State    json.RawMessage    `json:"state"`

	// Bandit-specific fields
	Variants     []string                    `json:"variants,omitempty"`
	VariantStats map[string]WASMVariantStats `json:"variant_stats,omitempty"`

	// Contextual bandit fields
	Context map[string]any `json:"context,omitempty"`

	// Optimizer-specific fields
	CurrentValue float64 `json:"current_value,omitempty"`
	Iteration    uint64  `json:"iteration,omitempty"`
	MetricSum    float64 `json:"metric_sum,omitempty"`
	BestValue    float64 `json:"best_value,omitempty"`
	BestReward   float64 `json:"best_reward,omitempty"`
}

// WASMVariantStats represents per-variant statistics for WASM algorithms.
type WASMVariantStats struct {
	Evaluations uint64  `json:"evaluations"`
	Successes   uint64  `json:"successes"`
	Failures    uint64  `json:"failures"`
	MetricSum   float64 `json:"metric_sum"`
}

// WASMOutput represents output from WASM algorithms.
type WASMOutput struct {
	// For bandit algorithms
	SelectedVariant string `json:"selected_variant,omitempty"`

	// For optimizer algorithms
	OptimizedValue float64 `json:"optimized_value,omitempty"`

	// Updated state to persist
	NewState json.RawMessage `json:"new_state,omitempty"`

	// Error message if algorithm failed
	Error string `json:"error,omitempty"`
}

// WASMFeedbackInput represents feedback event data for WASM algorithms.
type WASMFeedbackInput struct {
	Settings   map[string]float64 `json:"settings"`
	State      json.RawMessage    `json:"state"`
	VariantKey string             `json:"variant_key"`
	EventType  string             `json:"event_type"` // "success", "failure", "evaluation"
	Reward     float64            `json:"reward"`
	Context    map[string]any     `json:"context,omitempty"`
}

// WASMFeedbackOutput represents output from WASM feedback handler.
type WASMFeedbackOutput struct {
	NewState json.RawMessage `json:"new_state,omitempty"`
	Error    string          `json:"error,omitempty"`
}
