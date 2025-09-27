package dto

import (
	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainProjectToAPI converts domain Project to generated API Project.
func DomainProjectToAPI(project domain.Project) generatedapi.Project {
	return generatedapi.Project{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
	}
}

// DomainProjectsToAPI converts slice of domain Projects to slice of generated API Projects.
func DomainProjectsToAPI(projects []domain.Project) []generatedapi.Project {
	resp := make([]generatedapi.Project, 0, len(projects))
	for _, project := range projects {
		resp = append(resp, DomainProjectToAPI(project))
	}

	return resp
}

// APIProjectToDomain converts generated API Project to domain Project.
func APIProjectToDomain(project generatedapi.Project) domain.Project {
	return domain.Project{
		ID:          domain.ProjectID(project.ID),
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
	}
}
