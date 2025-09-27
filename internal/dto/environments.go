package dto

import (
	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func DomainEnvironmentToAPI(env domain.Environment) generatedapi.Environment {
	projectID, _ := uuid.Parse(string(env.ProjectID))
	apiKey, _ := uuid.Parse(env.APIKey)

	return generatedapi.Environment{
		ID:        int64(env.ID),
		ProjectID: projectID,
		Key:       env.Key,
		Name:      env.Name,
		APIKey:    apiKey,
		CreatedAt: env.CreatedAt,
	}
}

func DomainEnvironmentsToAPI(environments []domain.Environment) []generatedapi.Environment {
	result := make([]generatedapi.Environment, len(environments))
	for i, env := range environments {
		result[i] = DomainEnvironmentToAPI(env)
	}

	return result
}

func APIEnvironmentToDomain(env generatedapi.Environment) domain.Environment {
	return domain.Environment{
		ID:        domain.EnvironmentID(env.ID),
		ProjectID: domain.ProjectID(env.ProjectID.String()),
		Key:       env.Key,
		Name:      env.Name,
		APIKey:    env.APIKey.String(),
		CreatedAt: env.CreatedAt,
	}
}
