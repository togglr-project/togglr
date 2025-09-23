package domain

import (
	"time"
)

type CategoryID string

type CategoryKind string

const (
	CategoryKindSystem CategoryKind = "system"
	CategoryKindUser   CategoryKind = "user"
	CategoryKindNoCopy CategoryKind = "nocopy"
)

type Category struct {
	ID          CategoryID
	Name        string
	Slug        string
	Description *string
	Color       *string
	Kind        CategoryKind
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CategoryDTO struct {
	Name        string
	Slug        string
	Description *string
	Color       *string
	Kind        CategoryKind
}

func (id CategoryID) String() string {
	return string(id)
}

func (k CategoryKind) String() string {
	return string(k)
}

func (k CategoryKind) IsValid() bool {
	return k == CategoryKindSystem || k == CategoryKindUser || k == CategoryKindNoCopy
}
