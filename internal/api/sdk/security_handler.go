package apisdk

import (
	"context"
	"sync"
	"time"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/sdkserver"
)

var _ generatedapi.SecurityHandler = (*SecurityHandler)(nil)

type cacheEntry struct {
	projectID domain.ProjectID
	envKey    string
	expiresAt time.Time
}

type SecurityHandler struct {
	projectsRepo contract.ProjectsRepository
	cache        sync.Map // key: apiKey(string) -> value: cacheEntry
	cacheTTL     time.Duration
}

func NewSecurityHandler(
	projectsRepo contract.ProjectsRepository,
) *SecurityHandler {
	return &SecurityHandler{
		projectsRepo: projectsRepo,
		cacheTTL:     time.Hour,
	}
}

func (r *SecurityHandler) HandleApiKeyAuth(
	ctx context.Context,
	_ generatedapi.OperationName,
	tokenHolder generatedapi.ApiKeyAuth,
) (context.Context, error) {
	if v, ok := r.cache.Load(tokenHolder.APIKey); ok {
		ce := v.(cacheEntry)
		if time.Now().Before(ce.expiresAt) {
			ctx = appcontext.WithProjectID(ctx, ce.projectID)
			ctx = appcontext.WithEnvKey(ctx, ce.envKey)

			return ctx, nil
		}
		// expired entry, remove
		r.cache.Delete(tokenHolder.APIKey)
	}

	project, err := r.projectsRepo.GetByAPIKey(ctx, tokenHolder.APIKey)
	if err != nil {
		return ctx, err
	}

	// cache the result
	r.cache.Store(tokenHolder.APIKey, cacheEntry{
		projectID: project.ID,
		envKey:    project.EnvKey,
		expiresAt: time.Now().Add(r.cacheTTL),
	})

	ctx = appcontext.WithProjectID(ctx, project.ID)
	ctx = appcontext.WithEnvKey(ctx, project.EnvKey)

	return ctx, nil
}
