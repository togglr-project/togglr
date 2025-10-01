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

	// Verify the project exists (preserve current behavior and error mapping)
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

// CanToggleFeature checks if a user can toggle a feature.
func (s *Service) CanToggleFeature(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermFeatureToggle)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanManageFeature checks if a user can manage features (create, update, delete).
func (s *Service) CanManageFeature(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermFeatureManage)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanManageSegment checks if a user can manage segments.
func (s *Service) CanManageSegment(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermSegmentManage)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanManageSchedule checks if a user can manage feature schedules.
func (s *Service) CanManageSchedule(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermScheduleManage)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanViewAudit checks if a user can view audit logs.
func (s *Service) CanViewAudit(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermAuditView)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanManageMembership checks if a user can manage project memberships.
func (s *Service) CanManageMembership(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermMembershipManage)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanManageTags checks if a user can manage project tags.
func (s *Service) CanManageTags(ctx context.Context, projectID domain.ProjectID) error {
	ok, err := s.HasProjectPermission(ctx, projectID, domain.PermTagManage)
	if err != nil {
		return err
	}

	if !ok {
		return domain.ErrPermissionDenied
	}

	return nil
}

// CanManageCategories checks if a user can manage global categories.
// This is a global permission, so projectID is ignored.
func (s *Service) CanManageCategories(ctx context.Context) error {
	// For global permissions, check if user has category.manage on any project
	if s.isSuper(ctx) {
		return nil
	}

	userID := etx.UserID(ctx)
	if userID == 0 {
		return domain.ErrUserNotFound
	}

	// Get all projects where user has membership
	all, err := s.projects.List(ctx)
	if err != nil {
		return err
	}

	// Check if user has category.manage permission on any project
	for _, project := range all {
		roleID, err := s.member.GetForUserProject(ctx, int(userID), project.ID)
		if err != nil || roleID == "" {
			continue // no membership or error
		}

		has, err := s.perms.RoleHasPermission(ctx, roleID, domain.PermCategoryManage)
		if err != nil {
			continue // error checking permission
		}

		if has {
			return nil // user has category.manage on at least one project
		}
	}

	return domain.ErrPermissionDenied
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
		domain.PermCategoryManage,
		domain.PermTagManage,
	}

	result := make(map[domain.ProjectID][]domain.PermKey)

	for i := range all {
		project := all[i]

		// Check membership directly, do not use superuser bypass here
		roleID, mErr := s.member.GetForUserProject(ctx, int(userID), project.ID)
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
			result[project.ID] = granted
		}
	}

	return result, nil
}

func (s *Service) GetMyProjectRoles(ctx context.Context) (map[domain.ProjectID]domain.Role, error) {
	userID := etx.UserID(ctx)
	if userID == 0 {
		return nil, domain.ErrUserNotFound
	}

	all, err := s.projects.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make(map[domain.ProjectID]domain.Role)

	for i := range all {
		project := all[i]

		// Check membership directly, do not use superuser bypass here
		roleID, err := s.member.GetForUserProject(ctx, int(userID), project.ID)
		if err != nil {
			return nil, err
		}

		if roleID == "" {
			continue // no membership — skip this project
		}

		role, err := s.roles.GetByID(ctx, domain.RoleID(roleID))
		if err != nil {
			return nil, err
		}

		result[project.ID] = role
	}

	return result, nil
}
