package environments

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

type environmentModel struct {
	ID        int64     `db:"id"`
	ProjectID string    `db:"project_id"`
	Key       string    `db:"key"`
	Name      string    `db:"name"`
	APIKey    string    `db:"api_key"`
	CreatedAt time.Time `db:"created_at"`
}

func (m environmentModel) toDomain() domain.Environment {
	return domain.Environment{
		ID:        domain.EnvironmentID(m.ID),
		ProjectID: domain.ProjectID(m.ProjectID),
		Key:       m.Key,
		Name:      m.Name,
		APIKey:    m.APIKey,
		CreatedAt: m.CreatedAt,
	}
}

func environmentFromDomain(env domain.Environment) environmentModel {
	return environmentModel{
		ID:        int64(env.ID),
		ProjectID: string(env.ProjectID),
		Key:       env.Key,
		Name:      env.Name,
		APIKey:    env.APIKey,
		CreatedAt: env.CreatedAt,
	}
}
