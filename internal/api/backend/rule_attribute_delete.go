package apibackend

import (
	"context"
	"errors"
	"log/slog"

	etogglcontext "github.com/rom8726/etoggle/internal/context"
	"github.com/rom8726/etoggle/internal/domain"
	generatedapi "github.com/rom8726/etoggle/internal/generated/server"
)

func (r *RestAPI) DeleteRuleAttribute(
	ctx context.Context,
	params generatedapi.DeleteRuleAttributeParams,
) (generatedapi.DeleteRuleAttributeRes, error) {
	// Only superuser can delete attributes
	if !etogglcontext.IsSuper(ctx) {
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
