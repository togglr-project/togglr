import React from 'react';
import { Box, Typography, Paper, Grid, Avatar, Chip, Divider, Card, CardContent, List, ListItem, ListItemText, ListItemIcon, Accordion, AccordionSummary, AccordionDetails, Stack } from '@mui/material';
import { 
  Person as PersonIcon, 
  Email as EmailIcon, 
  AdminPanelSettings as AdminIcon,
  ExpandMore as ExpandMoreIcon,
  Folder as FolderIcon,
  Visibility as ViewIcon,
  Edit as EditIcon,
  ToggleOn as ToggleIcon,
  People as PeopleIcon,
  Schedule as ScheduleIcon,
  Assignment as AuditIcon,
  Settings as SettingsIcon
} from '@mui/icons-material';
import { useQuery } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import TwoFactorAuthSection from '../components/TwoFactorAuthSection';
import { useAuth } from '../auth/AuthContext';
import apiClient from '../api/apiClient';
import type { Project } from '../generated/api/client';

const AccountPage: React.FC = () => {
  const { user, isLoading, error } = useAuth();

  // Load projects for permissions display
  const { data: projects, isLoading: projectsLoading } = useQuery<Project[]>({
    queryKey: ['projects'],
    queryFn: async () => (await apiClient.listProjects()).data,
  });

  const getUserInitials = () => {
    if (!user?.username) return 'U';
    return user.username.charAt(0).toUpperCase();
  };

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // Get user's accessible projects with permissions and roles
  const getUserProjects = () => {
    if (!user || !projects) return [];
    
    // If user is superuser, show all projects with full permissions
    if (user.is_superuser) {
      return projects.map(project => ({
        ...project,
        permissions: ['project.view', 'project.manage', 'feature.view', 'feature.toggle', 'feature.manage', 'segment.manage', 'schedule.manage', 'audit.view', 'membership.manage', 'tag.manage', 'category.manage'],
        role: { name: 'Superuser', key: 'superuser', description: 'Full system access' }
      }));
    }
    
    // Otherwise, filter by project_permissions and include roles
    return projects
      .filter(project => {
        const permissions = user.project_permissions?.[project.id];
        return permissions && permissions.includes('project.view');
      })
      .map(project => ({
        ...project,
        permissions: user.project_permissions?.[project.id] || [],
        role: user.project_roles?.[project.id] || null
      }));
  };

  // Get permission icon and label
  const getPermissionInfo = (permission: string) => {
    const permissionMap: Record<string, { icon: React.ReactNode; label: string; color: string }> = {
      'project.view': { icon: <ViewIcon />, label: 'View Project', color: 'primary' },
      'project.manage': { icon: <SettingsIcon />, label: 'Manage Project', color: 'warning' },
      'feature.view': { icon: <ViewIcon />, label: 'View Features', color: 'info' },
      'feature.toggle': { icon: <ToggleIcon />, label: 'Toggle Features', color: 'success' },
      'feature.manage': { icon: <EditIcon />, label: 'Manage Features', color: 'success' },
      'segment.manage': { icon: <PeopleIcon />, label: 'Manage Segments', color: 'secondary' },
      'schedule.manage': { icon: <ScheduleIcon />, label: 'Manage Schedules', color: 'info' },
      'audit.view': { icon: <AuditIcon />, label: 'View Audit Logs', color: 'default' },
      'membership.manage': { icon: <PeopleIcon />, label: 'Manage Members', color: 'warning' },
      'tag.manage': { icon: <EditIcon />, label: 'Manage Tags', color: 'info' },
      'category.manage': { icon: <SettingsIcon />, label: 'Manage Categories', color: 'primary' }
    };
    
    return permissionMap[permission] || { icon: <ViewIcon />, label: permission, color: 'default' };
  };

  const userProjects = getUserProjects();

  return (
    <AuthenticatedLayout>
      <PageHeader
        title="My Account"
        subtitle="Manage account settings and security"
        icon={<PersonIcon />}
      />
      
      <Box sx={{ p: 3 }}>
        <Grid container spacing={3}>
          {/* User Information */}
          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 3, height: 'fit-content' }}>
              <Typography variant="h6" gutterBottom sx={{ color: 'primary.light', display: 'flex', alignItems: 'center', gap: 1 }}>
                <PersonIcon />
                User Information
              </Typography>
              
              {isLoading ? (
                <Typography>Loading...</Typography>
              ) : error ? (
                <Typography color="error">
                  Error loading user data. Please try again.
                </Typography>
              ) : user ? (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  {/* Avatar and main information */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
                    <Avatar 
                      sx={{ 
                        width: 64, 
                        height: 64, 
                        bgcolor: 'primary.main',
                        fontSize: '1.5rem',
                        fontWeight: 600
                      }}
                    >
                      {getUserInitials()}
                    </Avatar>
                    <Box>
                      <Typography variant="h5" sx={{ fontWeight: 600 }}>
                        {user.username}
                      </Typography>
                      <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
                        {user.is_superuser && (
                          <Chip 
                            icon={<AdminIcon />} 
                            label="Administrator" 
                            color="warning" 
                            size="small" 
                          />
                        )}
                        {user.is_active && (
                          <Chip 
                            label="Active" 
                            color="success" 
                            size="small" 
                          />
                        )}
                      </Box>
                    </Box>
                  </Box>

                  <Divider />

                  {/* Detailed information */}
                  <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <EmailIcon color="action" fontSize="small" />
                      <Typography variant="body2" color="text.secondary">
                        Email:
                      </Typography>
                      <Typography variant="body2">
                        {user.email || 'Not specified'}
                      </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <PersonIcon color="action" fontSize="small" />
                      <Typography variant="body2" color="text.secondary">
                        User ID:
                      </Typography>
                      <Typography variant="body2">
                        {user.id}
                      </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" color="text.secondary">
                        Created:
                      </Typography>
                      <Typography variant="body2">
                        {formatDate(user.created_at)}
                      </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" color="text.secondary">
                        Last login:
                      </Typography>
                      <Typography variant="body2">
                        {formatDate(user.last_login)}
                      </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" color="text.secondary">
                        External user:
                      </Typography>
                      <Chip 
                        label={user.is_external ? 'Yes' : 'No'} 
                        color={user.is_external ? 'warning' : 'default'} 
                        size="small" 
                      />
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" color="text.secondary">
                        Temporary password:
                      </Typography>
                      <Chip 
                        label={user.is_tmp_password ? 'Yes' : 'No'} 
                        color={user.is_tmp_password ? 'error' : 'default'} 
                        size="small" 
                      />
                    </Box>

                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" color="text.secondary">
                        User agreement accepted:
                      </Typography>
                      <Chip 
                        label={user.license_accepted ? 'Yes' : 'No'} 
                        color={user.license_accepted ? 'success' : 'error'} 
                        size="small" 
                      />
                    </Box>
                  </Box>
                </Box>
              ) : (
                <Typography>User not found</Typography>
              )}
            </Paper>
          </Grid>

          {/* Two-Factor Authentication Section */}
          <Grid item xs={12} md={6}>
            <TwoFactorAuthSection 
              userData={user}
              userLoading={isLoading}
              userError={error ? new Error(error) : null}
            />
          </Grid>

          {/* Project Permissions Section */}
          <Grid item xs={12}>
            <Paper sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom sx={{ color: 'primary.light', display: 'flex', alignItems: 'center', gap: 1, mb: 3 }}>
                <FolderIcon />
                Project Permissions
              </Typography>
              
              {projectsLoading ? (
                <Typography>Loading project permissions...</Typography>
              ) : userProjects.length === 0 ? (
                <Typography color="text.secondary">
                  No project access found.
                </Typography>
              ) : (
                <Box>
                  {userProjects.map((project) => (
                    <Accordion key={project.id} sx={{ mb: 2, boxShadow: 1 }}>
                      <AccordionSummary
                        expandIcon={<ExpandMoreIcon />}
                        sx={{
                          backgroundColor: 'primary.50',
                          '&:hover': { backgroundColor: 'primary.100' }
                        }}
                      >
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
                          <FolderIcon color="primary" />
                          <Box sx={{ flexGrow: 1 }}>
                            <Typography variant="h6" sx={{ fontWeight: 600 }}>
                              {project.name}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              {project.description || 'No description'}
                            </Typography>
                          </Box>
                          <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                            <Chip 
                              label={`${project.permissions.length} permissions`} 
                              size="small" 
                              color="primary" 
                              variant="outlined"
                            />
                            {project.role && (
                              <Chip 
                                label={project.role.name} 
                                size="small" 
                                color={user?.is_superuser ? "warning" : "secondary"} 
                                icon={user?.is_superuser ? <AdminIcon /> : undefined}
                                variant="filled"
                              />
                            )}
                          </Box>
                        </Box>
                      </AccordionSummary>
                      <AccordionDetails>
                        <Box>
                          {/* Role Information */}
                          {project.role && (
                            <Box sx={{ 
                              mb: 3, 
                              p: 2, 
                              backgroundColor: 'primary.50', 
                              borderRadius: 1, 
                              border: '1px solid', 
                              borderColor: 'primary.200' 
                            }}>
                              <Typography variant="subtitle2" gutterBottom sx={{ color: 'primary.dark', fontWeight: 600 }}>
                                Your Role
                              </Typography>
                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                                <Chip 
                                  label={project.role.name} 
                                  size="small" 
                                  color={user?.is_superuser ? "warning" : "secondary"} 
                                  icon={user?.is_superuser ? <AdminIcon /> : undefined}
                                  variant="filled"
                                />
                              </Box>
                              {project.role.description && (
                                <Typography variant="body2" sx={{ color: 'text.primary' }}>
                                  {project.role.description}
                                </Typography>
                              )}
                            </Box>
                          )}
                          
                          <Typography variant="subtitle2" gutterBottom sx={{ mb: 2, color: 'text.secondary' }}>
                            Available Permissions:
                          </Typography>
                          <Grid container spacing={1}>
                            {project.permissions.map((permission) => {
                              const { icon, label, color } = getPermissionInfo(permission);
                              return (
                                <Grid item xs={12} sm={6} md={4} key={permission}>
                                  <Card 
                                    variant="outlined" 
                                    sx={{ 
                                      p: 1.5, 
                                      height: '100%',
                                      borderColor: `${color}.light`,
                                      '&:hover': { 
                                        boxShadow: 2,
                                        transform: 'translateY(-1px)',
                                        transition: 'all 0.2s ease-in-out'
                                      }
                                    }}
                                  >
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                      <Box sx={{ color: `${color}.main` }}>
                                        {icon}
                                      </Box>
                                      <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                        {label}
                                      </Typography>
                                    </Box>
                                  </Card>
                                </Grid>
                              );
                            })}
                          </Grid>
                        </Box>
                      </AccordionDetails>
                    </Accordion>
                  ))}
                  
                  {/* Summary Stats */}
                  <Box sx={{ 
                    mt: 3, 
                    p: 2, 
                    backgroundColor: 'primary.50', 
                    borderRadius: 1,
                    border: '1px solid',
                    borderColor: 'primary.200'
                  }}>
                    <Typography variant="subtitle2" gutterBottom sx={{ color: 'primary.dark', fontWeight: 600 }}>
                      Summary
                    </Typography>
                    <Stack direction="row" spacing={3} flexWrap="wrap">
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <FolderIcon color="primary" fontSize="small" />
                        <Typography variant="body2" sx={{ color: 'text.primary' }}>
                          <strong>{userProjects.length}</strong> projects
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <ViewIcon color="info" fontSize="small" />
                        <Typography variant="body2" sx={{ color: 'text.primary' }}>
                          <strong>{userProjects.reduce((sum, p) => sum + p.permissions.length, 0)}</strong> total permissions
                        </Typography>
                      </Box>
                      {user?.is_superuser && (
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <AdminIcon color="warning" fontSize="small" />
                          <Typography variant="body2" color="warning.dark" sx={{ fontWeight: 500 }}>
                            <strong>Superuser</strong> - Full access
                          </Typography>
                        </Box>
                      )}
                    </Stack>
                  </Box>
                </Box>
              )}
            </Paper>
          </Grid>
        </Grid>
      </Box>
    </AuthenticatedLayout>
  );
};

export default AccountPage;
