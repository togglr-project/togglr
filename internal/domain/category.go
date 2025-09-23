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

type CategoryType string

const (
	CategoryTypeSafety CategoryType = "safety"
	CategoryTypeDomain CategoryType = "domain"
	CategoryTypeUser   CategoryType = "user"
)

type Category struct {
	ID          CategoryID
	Name        string
	Slug        string
	Description *string
	Color       *string
	Kind        CategoryKind
	Type        CategoryType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CategoryDTO struct {
	Name        string
	Slug        string
	Description *string
	Color       *string
	Kind        CategoryKind
	Type        CategoryType
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

func (typ CategoryType) String() string {
	return string(typ)
}

func (typ CategoryType) IsValid() bool {
	return typ == CategoryTypeUser || typ == CategoryTypeSafety || typ == CategoryTypeDomain
}
