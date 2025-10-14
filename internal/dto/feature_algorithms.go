package dto

import (
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func DomainFeatureAlgorithmToAPI(alg domain.FeatureAlgorithm) generatedapi.FeatureAlgorithm {
	settings := make(generatedapi.FeatureAlgorithmSettings, len(alg.Settings))
	for key, value := range alg.Settings {
		settings[key] = value.InexactFloat64()
	}

	return generatedapi.FeatureAlgorithm{
		FeatureID:     alg.FeatureID.String(),
		EnvironmentID: int64(alg.EnvironmentID),
		AlgorithmSlug: alg.AlgorithmSlug,
		Enabled:       alg.Enabled,
		Settings:      settings,
	}
}

func DomainFeatureAlgorithmsToAPI(algs []domain.FeatureAlgorithm) []generatedapi.FeatureAlgorithm {
	result := make([]generatedapi.FeatureAlgorithm, len(algs))
	for i, alg := range algs {
		result[i] = DomainFeatureAlgorithmToAPI(alg)
	}

	return result
}
