package bandit

import (
	"math"

	"github.com/shopspring/decimal"
)

// evalLinUCB implements the Linear Upper Confidence Bound algorithm.
// It chooses the variant with the highest UCB score based on context features.
func (m *BanditManager) evalLinUCB(state *AlgorithmState, ctx map[string]any) string {
	cState := state.ContextualState
	if cState == nil {
		return ""
	}

	cState.mu.Lock()
	defer cState.mu.Unlock()

	// Get alpha parameter (exploration factor)
	alpha := getSettingAsFloat64(state.Settings, "alpha", 1.0)

	// Convert context to feature vector
	features := ContextToFeatures(ctx, cState.FeatureDim)
	dim := cState.FeatureDim

	var bestVariant string
	bestScore := math.Inf(-1)

	for variant, vs := range cState.Variants {
		// Compute theta = A^-1 * b (using simplified approach)
		theta := solveLinear(vs.A, vs.B, dim)

		// Compute expected reward: x^T * theta
		expectedReward := dotProduct(features, theta)

		// Compute UCB bonus: alpha * sqrt(x^T * A^-1 * x)
		aInvX := solveLinear(vs.A, features, dim)
		ucbBonus := alpha * math.Sqrt(dotProduct(features, aInvX))

		// Total score
		score := expectedReward + ucbBonus

		if score > bestScore {
			bestScore = score
			bestVariant = variant
		}
	}

	if bestVariant != "" {
		cState.Variants[bestVariant].Pulls++
	}

	return bestVariant
}

// updateLinUCB updates the LinUCB model with observed reward
func (m *BanditManager) updateLinUCB(state *AlgorithmState, variantKey string, reward float64, ctx map[string]any) {
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
	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			vs.A[i*dim+j] += features[i] * features[j]
		}
	}

	// Update b = b + reward * x
	for i := 0; i < dim; i++ {
		vs.B[i] += reward * features[i]
	}

	vs.TotalRew += reward
	if reward > 0 {
		vs.Successes++
	} else {
		vs.Failures++
	}
}

// solveLinear solves Ax = b using Gaussian elimination (simplified)
func solveLinear(A, b []float64, dim int) []float64 {
	// Create working copies
	a := make([]float64, dim*dim)
	copy(a, A)
	x := make([]float64, dim)
	copy(x, b)

	// Forward elimination with partial pivoting
	for k := 0; k < dim; k++ {
		// Find pivot
		maxVal := math.Abs(a[k*dim+k])
		maxRow := k
		for i := k + 1; i < dim; i++ {
			if math.Abs(a[i*dim+k]) > maxVal {
				maxVal = math.Abs(a[i*dim+k])
				maxRow = i
			}
		}

		// Swap rows
		if maxRow != k {
			for j := k; j < dim; j++ {
				a[k*dim+j], a[maxRow*dim+j] = a[maxRow*dim+j], a[k*dim+j]
			}
			x[k], x[maxRow] = x[maxRow], x[k]
		}

		// Check for singular matrix
		if math.Abs(a[k*dim+k]) < 1e-10 {
			continue
		}

		// Eliminate column
		for i := k + 1; i < dim; i++ {
			factor := a[i*dim+k] / a[k*dim+k]
			for j := k; j < dim; j++ {
				a[i*dim+j] -= factor * a[k*dim+j]
			}
			x[i] -= factor * x[k]
		}
	}

	// Back substitution
	for i := dim - 1; i >= 0; i-- {
		if math.Abs(a[i*dim+i]) < 1e-10 {
			x[i] = 0
			continue
		}
		for j := i + 1; j < dim; j++ {
			x[i] -= a[i*dim+j] * x[j]
		}
		x[i] /= a[i*dim+i]
	}

	return x
}

func dotProduct(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// getLinUCBSettings returns default settings for LinUCB
func getLinUCBSettings() map[string]decimal.Decimal {
	return map[string]decimal.Decimal{
		"alpha":       decimal.NewFromFloat(1.0),  // Exploration parameter
		"feature_dim": decimal.NewFromFloat(32.0), // Feature dimension
	}
}
