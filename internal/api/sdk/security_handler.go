package apisdk

import (
	"context"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/contract"
	generatedapi "github.com/rom8726/etoggle/internal/generated/sdkserver"
)

var _ generatedapi.SecurityHandler = (*SecurityHandler)(nil)

type SecurityHandler struct {
	projectsRepo contract.ProjectsRepository
}

func NewSecurityHandler(
	projectsRepo contract.ProjectsRepository,
) *SecurityHandler {
	return &SecurityHandler{
		projectsRepo: projectsRepo,
	}
}

func (r *SecurityHandler) HandleApiKeyAuth(
	ctx context.Context,
	_ generatedapi.OperationName,
	tokenHolder generatedapi.ApiKeyAuth,
) (context.Context, error) {
	project, err := r.projectsRepo.GetByAPIKey(ctx, tokenHolder.APIKey)
	if err != nil {
		return ctx, err
	}

	ctx = etogglcontext.WithProjectID(ctx, project.ID)

	return ctx, nil
}
