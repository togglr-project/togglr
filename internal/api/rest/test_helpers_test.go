package rest

import (
	"time"

	"github.com/rom8726/etoggl/internal/domain"
)

// domainUserIDPtr создает указатель на domain.UserID.
func domainUserIDPtr(id uint) *domain.UserID {
	domainID := domain.UserID(id)

	return &domainID
}

// domainProjectIDPtr создает указатель на domain.ProjectID.
func domainProjectIDPtr(id string) *domain.ProjectID {
	domainID := domain.ProjectID(id)

	return &domainID
}

// timePtr создает указатель на time.Time.
func timePtr(t time.Time) *time.Time {
	return &t
}

// stringPtr создает указатель на string.
func stringPtr(s string) *string {
	return &s
}

// intPtr создает указатель на int.
func intPtr(i int) *int {
	return &i
}

// uintPtr создает указатель на uint.
func uintPtr(u uint) *uint {
	return &u
}

// boolPtr создает указатель на bool.
func boolPtr(b bool) *bool {
	return &b
}
