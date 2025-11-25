package bandit

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestContextToFeatures_BasicTypes(t *testing.T) {
	ctx := map[string]any{
		"age":       30,
		"score":     0.85,
		"is_active": true,
		"country":   "US",
	}

	features := ContextToFeatures(ctx, 16)

	assert.Len(t, features, 16)
	assert.NotEqual(t, 0.0, features[0], "Bias term should be set")

	// Features should be normalized
	var norm float64
	for _, f := range features {
		norm += f * f
	}
	assert.InDelta(t, 1.0, norm, 0.01, "Features should be normalized")
}

func TestContextToFeatures_EmptyContext(t *testing.T) {
	ctx := map[string]any{}

	features := ContextToFeatures(ctx, 8)

	assert.Len(t, features, 8)
	// With empty context, only bias term is set
	assert.NotEqual(t, 0.0, features[0], "Bias term should be set")
}

func TestContextToFeatures_DifferentDimensions(t *testing.T) {
	ctx := map[string]any{
		"user.id": "123",
	}

	for _, dim := range []int{8, 16, 32, 64} {
		features := ContextToFeatures(ctx, dim)
		assert.Len(t, features, dim, "Feature dimension should match requested")
	}
}

func TestContextToFeatures_ConsistentHashing(t *testing.T) {
	ctx := map[string]any{
		"user.id":      "test_user",
		"country_code": "US",
	}

	// Same context should produce same features
	features1 := ContextToFeatures(ctx, 16)
	features2 := ContextToFeatures(ctx, 16)

	for i := range features1 {
		assert.Equal(t, features1[i], features2[i], "Same context should produce same features")
	}
}

func TestNewContextualAlgorithmState_InitializesCorrectly(t *testing.T) {
	variants := []string{"A", "B", "C"}
	settings := map[string]decimal.Decimal{
		"alpha": decimal.NewFromFloat(1.0),
	}

	state := NewContextualAlgorithmState(16, variants, settings)

	assert.Equal(t, 16, state.FeatureDim)
	assert.Len(t, state.Variants, 3)

	for _, v := range state.Variants {
		// A matrix should be identity (d*d)
		assert.Len(t, v.A, 16*16)
		// Check diagonal elements are 1
		for i := 0; i < 16; i++ {
			assert.Equal(t, 1.0, v.A[i*16+i], "Diagonal of A should be 1")
		}
		// B vector should be zeros
		assert.Len(t, v.B, 16)
		for _, b := range v.B {
			assert.Equal(t, 0.0, b, "B vector should be zeros")
		}
	}
}

func TestNewContextualAlgorithmState_DefaultDimension(t *testing.T) {
	state := NewContextualAlgorithmState(0, []string{"A"}, nil)

	assert.Equal(t, DefaultFeatureDim, state.FeatureDim)
}

func TestHashToIndex_Distribution(t *testing.T) {
	keys := []string{
		"user.id", "country_code", "device_type", "browser",
		"os", "platform", "age", "gender",
	}

	dim := 16
	indices := make(map[int]int)

	for _, key := range keys {
		idx := hashToIndex(key, dim)
		assert.GreaterOrEqual(t, idx, 0)
		assert.Less(t, idx, dim)
		indices[idx]++
	}

	// Should use multiple buckets (not all in one)
	assert.Greater(t, len(indices), 1, "Hash should distribute across buckets")
}

func TestSolveLinear_IdentityMatrix(t *testing.T) {
	// For identity matrix A, solution should be x = b
	dim := 4
	A := make([]float64, dim*dim)
	for i := 0; i < dim; i++ {
		A[i*dim+i] = 1.0
	}
	b := []float64{1.0, 2.0, 3.0, 4.0}

	x := solveLinear(A, b, dim)

	for i := range b {
		assert.InDelta(t, b[i], x[i], 0.001, "Solution should equal b for identity matrix")
	}
}

func TestDotProduct(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0}
	b := []float64{4.0, 5.0, 6.0}

	result := dotProduct(a, b)

	// 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
	assert.Equal(t, 32.0, result)
}
