package apibackend

import (
	"context"
	"log/slog"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func (r *RestAPI) ListRuleAttributes(ctx context.Context) (generatedapi.ListRuleAttributesRes, error) {
	list, err := r.ruleAttributesUseCase.List(ctx)
	if err != nil {
		slog.Error("list rule attributes failed", "error", err)

		return nil, err
	}

	resp := make(generatedapi.ListRuleAttributesResponse, 0, len(list))
	for i := range list {
		item := list[i]
		var description generatedapi.OptString
		if item.Description != nil {
			description = generatedapi.NewOptString(*item.Description)
		}

		resp = append(resp, generatedapi.RuleAttributeEntity{
			Name:        item.Name.String(),
			Description: description,
		})
	}

	return &resp, nil
}
