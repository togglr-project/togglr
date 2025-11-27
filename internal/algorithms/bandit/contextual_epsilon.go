package bandit

import (
	"math"
)

// evalContextualEpsilon implements Epsilon-Greedy with context features.
// With probability epsilon, it explores randomly; otherwise exploits the best predicted variant.
func (m *BanditManager) evalContextualEpsilon(state *AlgorithmState, ctx map[string]any) string {
	cState := state.ContextualState
	if cState == nil {
		return ""
	}

	cState.mu.Lock()
	defer cState.mu.Unlock()

	epsilon := getSettingAsFloat64(state.Settings, "epsilon", 0.1)

	// Exploration: random variant
	if m.randSrc.Float64() < epsilon {
		variants := make([]string, 0, len(cState.Variants))
		for v := range cState.Variants {
			variants = append(variants, v)
		}

		if len(variants) == 0 {
			return ""
		}

		chosen := variants[m.randSrc.Intn(len(variants))]
		cState.Variants[chosen].Pulls++

		return chosen
	}

	// Exploitation: best predicted variant
	features := ContextToFeatures(ctx, cState.FeatureDim)
	dim := cState.FeatureDim

	var bestVariant string
	bestScore := math.Inf(-1)

	for variant, vs := range cState.Variants {
		// Compute predicted reward: theta = A^-1 * b, score = x^T * theta
		theta := solveLinear(vs.A, vs.B, dim)
		score := dotProduct(features, theta)

		if score > bestScore {
			bestScore = score
			bestVariant = variant
		}
	}

	// If no clear winner, pick random
	if bestVariant == "" {
		variants := make([]string, 0, len(cState.Variants))
		for v := range cState.Variants {
			variants = append(variants, v)
		}

		if len(variants) == 0 {
			return ""
		}

		bestVariant = variants[m.randSrc.Intn(len(variants))]
	}

	cState.Variants[bestVariant].Pulls++

	return bestVariant
}

// updateContextualEpsilon updates the model with observed reward.
func (m *BanditManager) updateContextualEpsilon(
	state *AlgorithmState,
	variantKey string,
	reward float64,
	ctx map[string]any,
) {
	cState := state.ContextualState
	if cState == nil {
		return
	}

	cState.mu.Lock()
	defer cState.mu.Unlock()

	vs, ok := cState.Variants[variantKey]
	if !ok {
		return
	}

	features := ContextToFeatures(ctx, cState.FeatureDim)
	dim := cState.FeatureDim

	// Update A = A + x * x^T
	for i := range dim {
		for j := range dim {
			vs.A[i*dim+j] += features[i] * features[j]
		}
	}

	// Update b = b + reward * x
	for i := range dim {
		vs.B[i] += reward * features[i]
	}

	vs.TotalRew += reward
	if reward > 0 {
		vs.Successes++
	} else {
		vs.Failures++
	}
}
