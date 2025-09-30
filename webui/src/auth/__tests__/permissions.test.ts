import { describe, it, expect } from 'vitest';
import { hasPermission, Guard, PERMISSIONS } from '../permissions';

describe('permissions helpers', () => {
  const pp: Record<string, string[]> = {
    '1': [
      PERMISSIONS.project.view,
      PERMISSIONS.feature.view,
      PERMISSIONS.feature.toggle,
      PERMISSIONS.audit.view,
    ],
    '2': [PERMISSIONS.project.view],
  };

  it('grants all when superuser', () => {
    expect(hasPermission('1', PERMISSIONS.feature.manage, { isSuperuser: true, projectPermissions: pp })).toBe(true);
    expect(Guard.canManageProject('2', true, pp)).toBe(true);
  });

  it('returns false when projectId missing', () => {
    expect(hasPermission(undefined, PERMISSIONS.feature.view, { projectPermissions: pp })).toBe(false);
  });

  it('checks project-specific permissions', () => {
    expect(Guard.canViewProject('1', false, pp)).toBe(true);
    expect(Guard.canToggleFeature('1', false, pp)).toBe(true);
    expect(Guard.canManageFeature('1', false, pp)).toBe(false);
  });

  it('returns false when project has no permissions entry', () => {
    expect(Guard.canViewProject('99', false, pp)).toBe(false);
  });
});
