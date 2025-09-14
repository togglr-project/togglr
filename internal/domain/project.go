package domain

import (
	"time"
)

type ProjectID string

type Project struct {
	ID          ProjectID
	Name        string
	Description string
	APIKey      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ArchivedAt  *time.Time
}

type ProjectDTO struct {
	Name        string
	Description string
	APIKey      string
}

func (id ProjectID) String() string {
	return string(id)
}
