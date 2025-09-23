package domain

import (
	"time"
)

type TagID string

type Tag struct {
	ID          TagID
	ProjectID   ProjectID
	CategoryID  *CategoryID
	Name        string
	Slug        string
	Description *string
	Color       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Category    *Category
}

type TagDTO struct {
	ProjectID   ProjectID
	CategoryID  *CategoryID
	Name        string
	Slug        string
	Description *string
	Color       *string
}

func (id TagID) String() string {
	return string(id)
}
