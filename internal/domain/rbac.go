package domain

import "time"

type PermissionID string

type Permission struct {
	ID   PermissionID
	Key  PermKey
	Name string
}

type MembershipID string

type ProjectMembership struct {
	ID        MembershipID
	UserID    int
	ProjectID ProjectID
	RoleID    RoleID
	RoleKey   string
	RoleName  string
	CreatedAt time.Time
}
