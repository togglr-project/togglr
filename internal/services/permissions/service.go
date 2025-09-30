package permissions

import (
	"context"

	etx "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
)

// Service handles permission checks for various operations.
type Service struct {
	projects contract.ProjectsRepository
	roles    contract.RolesRepository
	perms    contract.PermissionsRepository
	member   contract.MembershipsRepository
}

// New creates a new permissions service.
func New(
	projects contract.ProjectsRepository,
	roles contract.RolesRepository,
	perms contract.PermissionsRepository,
	member contract.MembershipsRepository,
) *Service {
	return &Service{projects: projects, roles: roles, perms: perms, member: member}
}

func (s *Service) isSuper(ctx context.Context) bool { return etx.IsSuper(ctx) }

// HasGlobalPermission checks global (non-project) permissions.
// For now, only superuser has global permissions.
func (s *Service) HasGlobalPermission(
	ctx context.Context,
	permKey domain.PermKey,
) (bool, error) {
	_ = permKey // reserved for future global-permissions storage

	return s.isSuper(ctx), nil
}

// HasProjectPermission checks if the user has a specific permission in the scope of the project.
func (s *Service) HasProjectPermission(
	ctx context.Context,
	projectID domain.ProjectID,
	permKey domain.PermKey,
) (bool, error) {
	if s.isSuper(ctx) {
		return true, nil
	}

	// Verify project exists (preserve current behavior and error mapping)
	if _, err := s.projects.GetByID(ctx, projectID); err != nil {
		return false, err
	}

	userID := etx.UserID(ctx)
	if userID == 0 {
		return false, domain.ErrUserNotFound
	}

	roleID, err := s.member.GetForUserProject(ctx, int(userID), projectID)
	if err != nil || roleID == "" {
		return false, err
	}

	return s.perms.RoleHasPermission(ctx, roleID, permKey)
}

// CanAccessProject checks if a user can access a project.
func (s *Service) CanAccessProject(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermProjectView)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanManageProject checks if a user can manage a project (create, update, delete).
func (s *Service) CanManageProject(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermProjectManage)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// GetAccessibleProjects returns all projects that a user can access.
func (s *Service) GetAccessibleProjects(
	ctx context.Context,
	projects []domain.Project,
) ([]domain.Project, error) {
	if s.isSuper(ctx) {
		return projects, nil
	}

	out := make([]domain.Project, 0, len(projects))

	for i := range projects {
		project := projects[i]

		ok, err := s.HasProjectPermission(ctx, project.ID, domain.PermProjectView)
		if err != nil {
			return nil, err
		}

		if ok {
			out = append(out, project)
		}
	}

	return out, nil
}

// GetMyProjectPermissions returns permissions for projects where the user has a membership.
func (s *Service) GetMyProjectPermissions(
	ctx context.Context,
) (map[domain.ProjectID][]domain.PermKey, error) {
	userID := etx.UserID(ctx)
	if userID == 0 {
		return nil, domain.ErrUserNotFound
	}

	all, err := s.projects.List(ctx)
	if err != nil {
		return nil, err
	}

	// Define the set of permission keys we expose via this endpoint
	permKeys := []domain.PermKey{
		domain.PermProjectView,
		domain.PermProjectManage,
		domain.PermProjectCreate,
		domain.PermFeatureView,
		domain.PermFeatureToggle,
		domain.PermFeatureManage,
		domain.PermAuditView,
		domain.PermMembershipManage,
		domain.PermSegmentManage,
		domain.PermScheduleManage,
	}

	result := make(map[domain.ProjectID][]domain.PermKey)

	for i := range all {
		p := all[i]

		// Check membership directly, do not use superuser bypass here
		roleID, mErr := s.member.GetForUserProject(ctx, int(userID), p.ID)
		if mErr != nil {
			return nil, mErr
		}

		if roleID == "" {
			continue // no membership — skip this project
		}

		// Collect granted permissions for the role
		var granted []domain.PermKey

		for _, key := range permKeys {
			has, perr := s.perms.RoleHasPermission(ctx, roleID, key)
			if perr != nil {
				return nil, perr
			}

			if has {
				granted = append(granted, key)
			}
		}

		if len(granted) > 0 {
			result[p.ID] = granted
		}
	}

	return result, nil
}
