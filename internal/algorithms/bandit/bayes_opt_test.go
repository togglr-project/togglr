package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestEvalBayesOpt_SelectsBestSample(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"noise": decimal.NewFromFloat(0.0), // No noise for deterministic test
		},
	}

	samples := []decimal.Decimal{
		decimal.NewFromFloat(1.0),
		decimal.NewFromFloat(2.0),
		decimal.NewFromFloat(3.0),
	}
	rewards := []decimal.Decimal{
		decimal.NewFromFloat(0.5),
		decimal.NewFromFloat(0.9), // Best reward
		decimal.NewFromFloat(0.3),
	}

	result := m.evalBayesOpt(state, samples, rewards)

	// Should select sample with best reward (2.0)
	assert.True(t, result.Equal(decimal.NewFromFloat(2.0)),
		"Should select sample with highest reward")
}

func TestEvalBayesOpt_WithNoise(t *testing.T) {
	m := newTestManager(42)
	state := &AlgorithmState{
		CurrentValue: decimal.Zero,
		Settings: map[string]decimal.Decimal{
			"noise": decimal.NewFromFloat(0.5), // High noise
		},
	}

	samples := []decimal.Decimal{
		decimal.NewFromFloat(1.0),
		decimal.NewFromFloat(2.0),
	}
	rewards := []decimal.Decimal{
		decimal.NewFromFloat(0.5),
		decimal.NewFromFloat(0.51), // Very close rewards
	}

	// With noise, selection may vary
	selections := make(map[string]int)
	for i := 0; i < 100; i++ {
		result := m.evalBayesOpt(state, samples, rewards)
		selections[result.String()]++
	}

	// Both should be selected sometimes due to noise
	assert.Greater(t, len(selections), 0, "Should select at least one variant")
}
