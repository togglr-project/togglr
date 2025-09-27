package domain

import (
	"time"
)

type ProjectID string

type Project struct {
	ID          ProjectID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ArchivedAt  *time.Time
}

type ProjectDTO struct {
	Name        string
	Description string
}

func (id ProjectID) String() string {
	return string(id)
}

// GetAPIKeyForEnvironment returns the API key for a specific environment
// This method will be implemented in the repository layer.
func (p *Project) GetAPIKeyForEnvironment(envID EnvironmentID, environments []Environment) string {
	for _, env := range environments {
		if env.ID == envID && env.ProjectID == p.ID {
			return env.APIKey
		}
	}

	return ""
}
