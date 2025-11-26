package apibackend

import (
	"context"
	"log/slog"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListAlgorithms(ctx context.Context) (generatedapi.ListAlgorithmsRes, error) {
	list, err := r.algorithmsRepo.List(ctx)
	if err != nil {
		slog.Error("error listing algorithms", "error", err)

		return nil, r.NewError(ctx, err)
	}

	result := make([]generatedapi.Algorithm, 0, len(list))
	for _, alg := range list {
		if alg.AlgorithmType() == domain.AlgorithmTypeBayesOpt ||
			alg.AlgorithmType() == domain.AlgorithmTypeCEM {
			continue
		}

		settings := make(generatedapi.AlgorithmDefaultSettings, len(alg.DefaultSettings))
		for key, value := range alg.DefaultSettings {
			settings[key] = value.InexactFloat64()
		}

		result = append(result, generatedapi.Algorithm{
			Slug:            alg.Slug,
			Name:            alg.Name,
			Description:     alg.Description,
			Kind:            generatedapi.AlgorithmKind(alg.Kind),
			DefaultSettings: settings,
		})
	}

	return &generatedapi.ListAlgorithmsResponse{Algorithms: result}, nil
}
