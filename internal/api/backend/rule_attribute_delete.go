package apibackend

import (
	"context"
	"errors"
	"log/slog"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) DeleteRuleAttribute(
	ctx context.Context,
	params generatedapi.DeleteRuleAttributeParams,
) (generatedapi.DeleteRuleAttributeRes, error) {
	// Only superuser can delete attributes
	if !appcontext.IsSuper(ctx) {
		return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
			Message: generatedapi.NewOptString("unauthorized"),
		}}, nil
	}

	if params.Name == "" {
		return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
			Message: generatedapi.NewOptString("name is required"),
		}}, nil
	}

	if err := r.ruleAttributesUseCase.Delete(ctx, domain.RuleAttribute(params.Name)); err != nil {
		slog.Error("delete rule attribute failed", "error", err, "name", params.Name)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString("rule attribute not found"),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteRuleAttributeNoContent{}, nil
}
