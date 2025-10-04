import React from 'react';
import { userPermissions } from '../hooks/userPermissions.ts';

interface PermissionGuardProps {
  children: React.ReactNode;
  fallback?: React.ReactNode;
  // Project permissions
  projectId?: number;
  canRead?: boolean;
  canWrite?: boolean;
  canDelete?: boolean;
  canManage?: boolean;
  // Team permissions
  teamId?: number;
  requireOwner?: boolean;
  requireAdmin?: boolean;
  requireMember?: boolean;
  // General permissions
  requireCreateProjects?: boolean;
  requireCreateTeams?: boolean;
  requireManageUsers?: boolean;
  // Superuser
  requireSuperuser?: boolean;
}

const PermissionGuard: React.FC<PermissionGuardProps> = ({
  children,
  fallback = null,
  projectId,
  canRead,
  canWrite,
  canDelete,
  canManage,
  teamId,
  requireOwner,
  requireAdmin,
  requireMember,
  requireCreateProjects,
  requireCreateTeams,
  requireManageUsers,
  requireSuperuser,
}) => {
  const {
    canReadProject,
    canWriteProject,
    canDeleteProject,
    canManageProject,
    isTeamOwner,
    isTeamAdmin,
    isTeamMember,
    canCreateProjects,
    canCreateTeams,
    canManageUsers,
    isSuperuser,
  } = userPermissions();

  // Check superuser
  if (requireSuperuser && !isSuperuser()) {
    return <>{fallback}</>;
  }

  // Check general permissions
  if (requireCreateProjects && !canCreateProjects()) {
    return <>{fallback}</>;
  }

  if (requireCreateTeams && !canCreateTeams()) {
    return <>{fallback}</>;
  }

  if (requireManageUsers && !canManageUsers()) {
    return <>{fallback}</>;
  }

  // Check project permissions
  if (projectId) {
    if (canRead && !canReadProject(projectId)) {
      return <>{fallback}</>;
    }

    if (canWrite && !canWriteProject(projectId)) {
      return <>{fallback}</>;
    }

    if (canDelete && !canDeleteProject(projectId)) {
      return <>{fallback}</>;
    }

    if (canManage && !canManageProject(projectId)) {
      return <>{fallback}</>;
    }
  }

  // Check team permissions
  if (teamId) {
    if (requireOwner && !isTeamOwner(teamId)) {
      return <>{fallback}</>;
    }

    if (requireAdmin && !isTeamAdmin(teamId)) {
      return <>{fallback}</>;
    }

    if (requireMember && !isTeamMember(teamId)) {
      return <>{fallback}</>;
    }
  }

  return <>{children}</>;
};

export default PermissionGuard;
