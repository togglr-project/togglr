package bandit

import (
	"hash/fnv"
	"sync"

	"github.com/shopspring/decimal"
)

const (
	// DefaultFeatureDim is the default dimension for feature hashing
	DefaultFeatureDim = 32
)

// ContextualVariantState holds the model state for one variant in contextual bandits
type ContextualVariantState struct {
	// For LinUCB: A matrix (d x d) stored as flat array, b vector (d)
	A []float64 // d*d matrix stored row-major
	B []float64 // d vector

	// For Thompson Sampling: mean and precision for each feature
	Mu        []float64 // mean vector
	Precision []float64 // precision (inverse variance) vector

	// Common stats
	Pulls     uint64
	TotalRew  float64
	Successes uint64
	Failures  uint64
}

// ContextualAlgorithmState holds state for contextual bandit algorithms
type ContextualAlgorithmState struct {
	FeatureDim int
	Variants   map[string]*ContextualVariantState
	Settings   map[string]decimal.Decimal
	mu         sync.RWMutex
}

// NewContextualAlgorithmState creates a new contextual state with given dimension
func NewContextualAlgorithmState(dim int, variants []string, settings map[string]decimal.Decimal) *ContextualAlgorithmState {
	if dim <= 0 {
		dim = DefaultFeatureDim
	}

	state := &ContextualAlgorithmState{
		FeatureDim: dim,
		Variants:   make(map[string]*ContextualVariantState, len(variants)),
		Settings:   settings,
	}

	for _, v := range variants {
		state.Variants[v] = newContextualVariantState(dim)
	}

	return state
}

func newContextualVariantState(dim int) *ContextualVariantState {
	// Initialize A as identity matrix
	a := make([]float64, dim*dim)
	for i := 0; i < dim; i++ {
		a[i*dim+i] = 1.0
	}

	return &ContextualVariantState{
		A:         a,
		B:         make([]float64, dim),
		Mu:        make([]float64, dim),
		Precision: make([]float64, dim),
	}
}

// ContextToFeatures converts context map to feature vector using feature hashing
func ContextToFeatures(ctx map[string]any, dim int) []float64 {
	features := make([]float64, dim)

	// Bias term
	features[0] = 1.0

	for key, val := range ctx {
		idx := hashToIndex(key, dim)

		switch v := val.(type) {
		case float64:
			features[idx] += v
		case float32:
			features[idx] += float64(v)
		case int:
			features[idx] += float64(v)
		case int64:
			features[idx] += float64(v)
		case bool:
			if v {
				features[idx] += 1.0
			}
		case string:
			// Hash string value to create a feature
			valIdx := hashToIndex(key+"="+v, dim)
			features[valIdx] += 1.0
		}
	}

	// Normalize features
	var norm float64
	for _, f := range features {
		norm += f * f
	}
	if norm > 0 {
		norm = 1.0 / sqrt(norm)
		for i := range features {
			features[i] *= norm
		}
	}

	return features
}

func hashToIndex(s string, dim int) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))

	return int(h.Sum32()) % dim
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}

	return z
}
