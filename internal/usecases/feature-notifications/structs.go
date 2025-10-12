package feature_notifications

import (
	"github.com/togglr-project/togglr/internal/domain"
)

type projectIDWithEnvID struct {
	ProjectID domain.ProjectID
	EnvID     domain.EnvironmentID
}
