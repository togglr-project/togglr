package dto

import (
	"github.com/google/uuid"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

// DomainMembershipToAPI converts domain.ProjectMembership to generated API Membership.
func DomainMembershipToAPI(m domain.ProjectMembership) (generatedapi.Membership, error) {
	id, err := uuid.Parse(string(m.ID))
	if err != nil {
		return generatedapi.Membership{}, err
	}
	pid, err := uuid.Parse(string(m.ProjectID))
	if err != nil {
		return generatedapi.Membership{}, err
	}
	rid, err := uuid.Parse(string(m.RoleID))
	if err != nil {
		return generatedapi.Membership{}, err
	}

	return generatedapi.Membership{
		ID:        id,
		UserID:    int64(m.UserID),
		ProjectID: pid,
		RoleID:    rid,
		RoleKey:   m.RoleKey,
		RoleName:  m.RoleName,
		CreatedAt: m.CreatedAt,
	}, nil
}
