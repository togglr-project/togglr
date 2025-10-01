import { useMemo } from 'react';
import { useAuth } from './AuthContext';

// Permission keys as constants to avoid typos
export const PERMISSIONS = {
  project: {
    view: 'project.view',
    manage: 'project.manage',
  },
  feature: {
    view: 'feature.view',
    toggle: 'feature.toggle',
    manage: 'feature.manage',
  },
  segment: {
    manage: 'segment.manage',
  },
  schedule: {
    manage: 'schedule.manage',
  },
  audit: {
    view: 'audit.view',
  },
  membership: {
    manage: 'membership.manage',
  },
  tag: {
    manage: 'tag.manage',
  },
  category: {
    manage: 'category.manage',
  },
} as const;

export type PermissionKey =
  | typeof PERMISSIONS.project.view
  | typeof PERMISSIONS.project.manage
  | typeof PERMISSIONS.feature.view
  | typeof PERMISSIONS.feature.toggle
  | typeof PERMISSIONS.feature.manage
  | typeof PERMISSIONS.segment.manage
  | typeof PERMISSIONS.schedule.manage
  | typeof PERMISSIONS.audit.view
  | typeof PERMISSIONS.membership.manage
  | typeof PERMISSIONS.tag.manage
  | typeof PERMISSIONS.category.manage;

export function hasPermission(
  projectId: string | number | undefined,
  perm: PermissionKey,
  opts?: { isSuperuser?: boolean; projectPermissions?: Record<string, string[]> | undefined },
): boolean {
  const isSuperuser = opts?.isSuperuser ?? false;
  if (isSuperuser) return true;

  // For global permissions (like category.manage), check if user has it on any project
  if (perm === PERMISSIONS.category.manage) {
    const pp = opts?.projectPermissions;
    if (!pp) return false;
    
    // Check if user has category.manage permission on any project
    return Object.values(pp).some(permissions => 
      permissions.includes(PERMISSIONS.category.manage)
    );
  }

  // For project-specific permissions, we need projectId
  if (!projectId) return false;

  const pp = opts?.projectPermissions;
  if (!pp) return false;

  const perms = pp[String(projectId)];
  if (!perms || perms.length === 0) return false;

  return perms.includes(perm);
}

export function useRBAC(projectId?: string | number) {
  const { user } = useAuth();

  const guards = useMemo(() => {
    const superuser = Boolean(user?.is_superuser);
    const pp = user?.project_permissions;

    const check = (p: PermissionKey) => hasPermission(projectId, p, { isSuperuser: superuser, projectPermissions: pp });

    return {
      isSuperuser: superuser,
      has: check,
      canViewProject: () => check(PERMISSIONS.project.view),
      canManageProject: () => check(PERMISSIONS.project.manage),
      canViewFeature: () => check(PERMISSIONS.feature.view),
      canToggleFeature: () => check(PERMISSIONS.feature.toggle),
      canManageFeature: () => check(PERMISSIONS.feature.manage),
      canManageSegment: () => check(PERMISSIONS.segment.manage),
      canManageSchedule: () => check(PERMISSIONS.schedule.manage),
      canViewAudit: () => check(PERMISSIONS.audit.view),
      canManageMembership: () => check(PERMISSIONS.membership.manage),
      canManageTags: () => check(PERMISSIONS.tag.manage),
      canManageCategories: () => check(PERMISSIONS.category.manage),
    };
  }, [user, projectId]);

  return guards;
}

// Standalone helpers mirroring the spec (for modules that cannot use hooks)
export const Guard = {
  canViewProject(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.project.view, { isSuperuser, projectPermissions });
  },
  canManageProject(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.project.manage, { isSuperuser, projectPermissions });
  },
  canToggleFeature(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.feature.toggle, { isSuperuser, projectPermissions });
  },
  canManageFeature(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.feature.manage, { isSuperuser, projectPermissions });
  },
  canManageSegment(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.segment.manage, { isSuperuser, projectPermissions });
  },
  canManageSchedule(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.schedule.manage, { isSuperuser, projectPermissions });
  },
  canViewAudit(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.audit.view, { isSuperuser, projectPermissions });
  },
  canManageMembership(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.membership.manage, { isSuperuser, projectPermissions });
  },
  canManageTags(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    return hasPermission(projectId, PERMISSIONS.tag.manage, { isSuperuser, projectPermissions });
  },
  canManageCategories(projectId: string | number, isSuperuser?: boolean, projectPermissions?: Record<string, string[]>) {
    // For global permissions, we can pass any projectId (it will be ignored)
    return hasPermission(projectId, PERMISSIONS.category.manage, { isSuperuser, projectPermissions });
  },
};
