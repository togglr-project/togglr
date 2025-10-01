package rbac

import (
	"context"
	"fmt"

	appctx "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	"github.com/togglr-project/togglr/internal/repository/membershipaudit"
	"github.com/togglr-project/togglr/pkg/db"
)

var _ contract.MembershipsUseCase = (*Service)(nil)

type Service struct {
	rolesRepo       contract.RolesRepository
	permsRepo       contract.PermissionsRepository
	membershipsRepo contract.MembershipsRepository
	tx              db.TxManager
}

func New(
	rolesRepo contract.RolesRepository,
	permsRepo contract.PermissionsRepository,
	membershipsRepo contract.MembershipsRepository,
	tx db.TxManager,
) *Service {
	return &Service{rolesRepo: rolesRepo, permsRepo: permsRepo, membershipsRepo: membershipsRepo, tx: tx}
}

// Roles & permissions

func (s *Service) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return s.rolesRepo.List(ctx)
}

func (s *Service) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	return s.permsRepo.List(ctx)
}

func (s *Service) GetRolePermissions(ctx context.Context, roleID domain.RoleID) ([]domain.Permission, error) {
	return s.permsRepo.ListForRole(ctx, roleID)
}

func (s *Service) ListRolePermissions(ctx context.Context) (map[domain.Role][]domain.Permission, error) {
	return s.permsRepo.ListForAllRoles(ctx)
}

// Memberships

func (s *Service) ListProjectMemberships(ctx context.Context, projectID domain.ProjectID) ([]domain.ProjectMembership, error) {
	return s.membershipsRepo.ListForProject(ctx, projectID)
}

func (s *Service) CreateProjectMembership(
	ctx context.Context,
	projectID domain.ProjectID,
	userID int,
	roleID domain.RoleID,
) (domain.ProjectMembership, error) {
	var created domain.ProjectMembership
	if err := s.tx.ReadCommitted(ctx, func(ctx context.Context) error {
		m, err := s.membershipsRepo.Create(ctx, projectID, userID, roleID)
		if err != nil {
			return err
		}
		created = m

		actorID := int(appctx.UserID(ctx))
		exec := db.TxFromContext(ctx)
		if err := membershipaudit.Write(ctx, exec, string(m.ID), actorID, "create", nil, m); err != nil {
			return fmt.Errorf("write membership audit: %w", err)
		}

		return nil
	}); err != nil {
		return domain.ProjectMembership{}, err
	}

	return created, nil
}

func (s *Service) GetProjectMembership(
	ctx context.Context,
	projectID domain.ProjectID,
	membershipID domain.MembershipID,
) (domain.ProjectMembership, error) {
	return s.membershipsRepo.Get(ctx, projectID, membershipID)
}

func (s *Service) UpdateProjectMembership(
	ctx context.Context,
	projectID domain.ProjectID,
	membershipID domain.MembershipID,
	roleID domain.RoleID,
) (domain.ProjectMembership, error) {
	var updated domain.ProjectMembership
	if err := s.tx.ReadCommitted(ctx, func(ctx context.Context) error {
		old, err := s.membershipsRepo.Get(ctx, projectID, membershipID)
		if err != nil {
			return err
		}

		m, err := s.membershipsRepo.Update(ctx, projectID, membershipID, roleID)
		if err != nil {
			return err
		}
		updated = m

		actorID := int(appctx.UserID(ctx))
		exec := db.TxFromContext(ctx)
		if err := membershipaudit.Write(ctx, exec, string(m.ID), actorID, "update", old, m); err != nil {
			return fmt.Errorf("write membership audit: %w", err)
		}

		return nil
	}); err != nil {
		return domain.ProjectMembership{}, err
	}

	return updated, nil
}

func (s *Service) DeleteProjectMembership(
	ctx context.Context,
	projectID domain.ProjectID,
	membershipID domain.MembershipID,
) error {
	return s.tx.ReadCommitted(ctx, func(ctx context.Context) error {
		old, err := s.membershipsRepo.Get(ctx, projectID, membershipID)
		if err != nil {
			return err
		}

		if err := s.membershipsRepo.Delete(ctx, projectID, membershipID); err != nil {
			return err
		}

		actorID := int(appctx.UserID(ctx))
		exec := db.TxFromContext(ctx)
		if err := membershipaudit.Write(ctx, exec, string(old.ID), actorID, "delete", old, nil); err != nil {
			return fmt.Errorf("write membership audit: %w", err)
		}

		return nil
	})
}
