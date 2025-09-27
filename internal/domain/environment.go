package domain

import (
	"fmt"
	"time"
)

type EnvironmentID int64

type Environment struct {
	ID        EnvironmentID
	ProjectID ProjectID
	Key       string // dev, stage, prod
	Name      string // Development, Staging, Production
	APIKey    string
	CreatedAt time.Time
}

type EnvironmentDTO struct {
	Key  string `json:"key" validate:"required,min=1,max=20"`
	Name string `json:"name" validate:"required,min=1,max=50"`
}

func (id EnvironmentID) String() string {
	return fmt.Sprintf("%d", id)
}
