package bandit

import (
	"math/rand"

	"github.com/shopspring/decimal"
)

func newTestManager(seed int64) *BanditManager {
	return &BanditManager{
		randSrc: rand.New(rand.NewSource(seed)),
	}
}

func newTestState(variants []string) *AlgorithmState {
	variantsMap := make(map[string]*VariantStats, len(variants))
	for _, v := range variants {
		variantsMap[v] = &VariantStats{}
	}
	return &AlgorithmState{
		Variants:    variantsMap,
		VariantsArr: variants,
		Settings:    make(map[string]decimal.Decimal),
	}
}
