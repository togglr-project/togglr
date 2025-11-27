package bandit

import (
	"math"
)

// evalContextualThompson implements Thompson Sampling with context features.
// It uses a linear model with Bayesian updates for each variant.
func (m *BanditManager) evalContextualThompson(state *AlgorithmState, ctx map[string]any) string {
	cState := state.ContextualState
	if cState == nil {
		return ""
	}

	cState.mu.Lock()
	defer cState.mu.Unlock()

	features := ContextToFeatures(ctx, cState.FeatureDim)
	dim := cState.FeatureDim

	// Prior variance for sampling
	priorVar := getSettingAsFloat64(state.Settings, "prior_variance", 1.0)

	var bestVariant string
	bestSample := math.Inf(-1)

	for variant, vs := range cState.Variants {
		// Compute posterior mean: theta = A^-1 * b
		theta := solveLinear(vs.A, vs.B, dim)

		// Sample from posterior: theta_sample ~ N(theta, A^-1 * prior_var)
		// Simplified: add noise proportional to uncertainty
		thetaSample := make([]float64, dim)
		for i := range dim {
			// Approximate posterior variance from diagonal of A^-1
			variance := priorVar / (vs.A[i*dim+i] + 1e-6)
			stdDev := math.Sqrt(variance)
			thetaSample[i] = theta[i] + m.randSrc.NormFloat64()*stdDev
		}

		// Compute sampled reward: x^T * theta_sample
		sample := dotProduct(features, thetaSample)

		if sample > bestSample {
			bestSample = sample
			bestVariant = variant
		}
	}

	if bestVariant != "" {
		cState.Variants[bestVariant].Pulls++
	}

	return bestVariant
}

// updateContextualThompson updates the Thompson Sampling model with observed reward.
func (m *BanditManager) updateContextualThompson(
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

	// Update A = A + x * x^T (same as LinUCB)
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
