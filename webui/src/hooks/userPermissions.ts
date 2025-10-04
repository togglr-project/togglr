import { useAuth } from '../auth/AuthContext';
import { jwtDecode } from 'jwt-decode';

interface ProjectPermission {
  can_read: boolean;
  can_write: boolean;
  can_delete: boolean;
  can_manage: boolean;
  team_role?: string;
}

interface UserPermissions {
  project_permissions?: Record<string, ProjectPermission>;
  team_roles?: Record<string, string>;
  can_create_projects: boolean;
  can_create_teams: boolean;
  can_manage_users: boolean;
}

interface TokenClaims {
  userId: number;
  username: string;
  isSuperuser: boolean;
  permissions?: UserPermissions;
}

export const userPermissions = () => {
  const { user } = useAuth();

  const getTokenClaims = (): TokenClaims | null => {
    const token = localStorage.getItem('accessToken');
    if (!token) return null;

    try {
      return jwtDecode<TokenClaims>(token);
    } catch (error) {
      console.error('Failed to decode JWT token:', error);
      return null;
    }
  };

  const claims = getTokenClaims();
  const permissions = claims?.permissions;

  const canReadProject = (projectId: number): boolean => {
    if (claims?.isSuperuser) return true;
    
    if (!permissions?.project_permissions) return false;
    return permissions.project_permissions[projectId]?.can_read || false;
  };

  const canWriteProject = (projectId: number): boolean => {
    if (claims?.isSuperuser) return true;
    
    if (!permissions?.project_permissions) return false;
    return permissions.project_permissions[projectId]?.can_write || false;
  };

  const canDeleteProject = (projectId: number): boolean => {
    if (claims?.isSuperuser) return true;
    
    if (!permissions?.project_permissions) return false;
    return permissions.project_permissions[projectId]?.can_delete || false;
  };

  const canManageProject = (projectId: number, teamId?: number): boolean => {
    if (claims?.isSuperuser) {
      console.log(`Superuser can manage project ${projectId}`);
      return true;
    }
    
    if (permissions?.project_permissions) {
      const projectPermission = permissions.project_permissions[projectId];
      if (projectPermission?.can_manage) {
        console.log(`User has explicit can_manage permission for project ${projectId}`);
        return true;
      }
      if (projectPermission?.team_role && (projectPermission.team_role === 'owner' || projectPermission.team_role === 'admin')) {
        console.log(`User has team role '${projectPermission.team_role}' for project ${projectId} from project_permissions`);
        return true;
      }
    }
    
    if (teamId && permissions?.team_roles) {
      const teamRole = permissions.team_roles[String(teamId)];
      if (teamRole === 'owner' || teamRole === 'admin') {
        console.log(`User has team role '${teamRole}' for team ${teamId}, can manage project ${projectId}`);
        return true;
      }
    }
    
    console.log(`User cannot manage project ${projectId}. Team ID: ${teamId}, Available team roles:`, permissions?.team_roles);
    return false;
  };

  const canManageProjectByTeamRole = (projectTeamId: number): boolean => {
    if (claims?.isSuperuser) return true;
    
    if (!permissions?.team_roles) return false;
    
    const teamRole = permissions.team_roles[projectTeamId];
    if (teamRole === 'owner' || teamRole === 'admin') {
      console.log(`User has team role '${teamRole}' for team ${projectTeamId}`);
      return true;
    }
    
    console.log(`User has team role '${teamRole}' for team ${projectTeamId}, cannot manage`);
    return false;
  };

  const getProjectTeamRole = (projectId: number): string | null => {
    if (!permissions?.project_permissions) return null;
    return permissions.project_permissions[projectId]?.team_role || null;
  };

  const getTeamRole = (teamId: number): string | null => {
    if (!permissions?.team_roles) return null;
    return permissions.team_roles[teamId] || null;
  };

  const isTeamOwner = (teamId: number): boolean => {
    return getTeamRole(teamId) === 'owner';
  };

  const isTeamAdmin = (teamId: number): boolean => {
    const role = getTeamRole(teamId);
    return role === 'owner' || role === 'admin';
  };

  const isTeamMember = (teamId: number): boolean => {
    const role = getTeamRole(teamId);
    return role === 'owner' || role === 'admin' || role === 'member';
  };

  const canCreateProjects = (): boolean => {
    const result = claims?.isSuperuser || permissions?.can_create_projects || false;
    return result;
  };

  const canCreateTeams = (): boolean => {
    return claims?.isSuperuser || permissions?.can_create_teams || false;
  };

  const canManageUsers = (): boolean => {
    return claims?.isSuperuser || permissions?.can_manage_users || false;
  };

  const isSuperuser = (): boolean => {
    return claims?.isSuperuser || false;
  };

  const getAccessibleProjectIds = (): number[] => {
    if (!permissions?.project_permissions) {
      console.log('No project permissions found');
      return [];
    }
    const ids = Object.keys(permissions.project_permissions).map(Number);
    console.log('Accessible project IDs:', ids);
    return ids;
  };

  const getUserTeamIds = (): number[] => {
    if (!permissions?.team_roles) return [];
    return Object.keys(permissions.team_roles).map(Number);
  };

  return {
    canReadProject,
    canWriteProject,
    canDeleteProject,
    canManageProject,
    getProjectTeamRole,

    getTeamRole,
    isTeamOwner,
    isTeamAdmin,
    isTeamMember,

    canCreateProjects,
    canCreateTeams,
    canManageUsers,
    isSuperuser,

    getAccessibleProjectIds,
    getUserTeamIds,

    permissions,
    claims,
  };
}; 