RBAC usage in webui

Overview
- Frontend RBAC is driven by GET /api/v1/users/me returning:
  - is_superuser: boolean
  - project_permissions: { [project_id: string]: string[] }

Helpers
- src/auth/permissions.ts provides:
  - hasPermission(projectId, perm)
  - useRBAC(projectId) hook returning:
    - isSuperuser, has(), and guards:
      - canViewProject, canManageProject
      - canViewFeature, canToggleFeature, canManageFeature
      - canManageSegment, canManageSchedule
      - canViewAudit, canManageMembership
  - Guard static helpers for non-React modules

Examples
- In a component bound to a project:
  import { useRBAC } from './auth/permissions';
  const rbac = useRBAC(projectId);
  if (!rbac.canViewProject()) return null;
  <Button disabled={!rbac.canToggleFeature()}>Toggle</Button>

- Outside React (utility/service):
  import { Guard, PERMISSIONS } from './auth/permissions';
  const allowed = Guard.canManageProject(projectId, user.is_superuser, user.project_permissions);

Guidelines
- Always use canXxx() functions in UI components instead of inline checks.
- Superuser bypass is built in (all guards return true when is_superuser is true).
- When project_permissions lacks an entry for a project, all guards return false for that project.

Notes
- User data is stored in AuthContext (src/auth/AuthContext.tsx) and fetched via getCurrentUser().
- Stick to these guards for visibility and enabled/disabled states across the UI.
