package apibackend

import (
	"context"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) CreateRuleAttribute(
	ctx context.Context,
	req *generatedapi.CreateRuleAttributeRequest,
) (generatedapi.CreateRuleAttributeRes, error) {
	// Only superuser can create attributes
	if !appcontext.IsSuper(ctx) {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("unauthorized"),
		}}, nil
	}

	if req.Name == "" {
		return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
			Message: generatedapi.NewOptString("name is required"),
		}}, nil
	}

	var desc *string
	if req.Description.IsSet() {
		v, _ := req.Description.Get()
		desc = &v
	}

	_, err := r.ruleAttributesUseCase.Create(ctx, domain.RuleAttribute(req.Name), desc)
	if err != nil {
		slog.Error("create rule attribute failed", "error", err)

		return nil, err
	}

	return &generatedapi.CreateRuleAttributeNoContent{}, nil
}
