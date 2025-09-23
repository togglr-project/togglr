package apibackend

import (
	"context"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) GetSSOProviders(ctx context.Context) (generatedapi.GetSSOProvidersRes, error) {
	providers, err := r.usersUseCase.GetSSOProviders(ctx)
	if err != nil {
		return nil, err
	}

	apiProviders := make([]generatedapi.SSOProvider, 0, len(providers))
	for i := range providers {
		provider := providers[i]
		apiProviders = append(apiProviders, generatedapi.SSOProvider{
			Name:        provider.GetName(),
			DisplayName: provider.GetDisplayName(),
			Type:        generatedapi.SSOProviderTypeSaml,
			IconURL:     generatedapi.NewOptString(provider.GetIconURL()),
			Enabled:     provider.IsEnabled(),
		})
	}

	return &generatedapi.SSOProvidersResponse{Providers: apiProviders}, nil
}
