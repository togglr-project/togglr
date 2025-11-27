package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEvalCEM_SelectsEliteMean(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"elite_fraction":  decimal.NewFromFloat(0.5),
			"population_size": decimal.NewFromFloat(4),
		},
	}

	// 4 samples, elite_fraction=0.5 means top 2 will be selected
	samples := []decimal.Decimal{
		decimal.NewFromFloat(1.0),
		decimal.NewFromFloat(2.0),
		decimal.NewFromFloat(3.0),
		decimal.NewFromFloat(4.0),
	}
	rewards := []decimal.Decimal{
		decimal.NewFromFloat(0.1),
		decimal.NewFromFloat(0.9), // Elite
		decimal.NewFromFloat(0.2),
		decimal.NewFromFloat(0.8), // Elite
	}

	result := m.evalCEM(state, samples, rewards)

	// Elite samples are 2.0 (reward 0.9) and 4.0 (reward 0.8)
	// Mean of elites = (2.0 + 4.0) / 2 = 3.0
	assert.True(t, result.Equal(decimal.NewFromFloat(3.0)),
		"Should return mean of elite samples")
}

func TestEvalCEM_DefaultSettings(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings:     make(map[string]decimal.Decimal),
	}

	// 20 samples (default population_size)
	samples := make([]decimal.Decimal, 20)
	rewards := make([]decimal.Decimal, 20)
	for i := 0; i < 20; i++ {
		samples[i] = decimal.NewFromFloat(float64(i))
		rewards[i] = decimal.NewFromFloat(float64(i)) // Higher sample = higher reward
	}

	result := m.evalCEM(state, samples, rewards)

	// With elite_fraction=0.2 and pop_size=20, top 4 samples are elites
	// Top 4 by reward: 19, 18, 17, 16
	// Mean = (19 + 18 + 17 + 16) / 4 = 17.5
	assert.True(t, result.Equal(decimal.NewFromFloat(17.5)),
		"Should return mean of top 20% samples")
}
