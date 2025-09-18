import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
  Alert,
  CircularProgress,
} from '@mui/material';
import {
  People as PeopleIcon,
  AdminPanelSettings as AdminPanelSettingsIcon,
  VpnKey as LicenseIcon,
  Sync as SyncIcon,
  FolderOutlined as ProjectsIcon,
} from '@mui/icons-material';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import { userPermissions } from '../hooks/userPermissions';
import Layout from '../components/Layout';
import TabPanel from '../components/admin/TabPanel';
import UsersTab from '../components/admin/UsersTab';
import ProjectsTab from '../components/admin/ProjectsTab';
import LicenseTab from '../components/admin/LicenseTab';
import ExternalAuthTab from '../components/admin/ExternalAuthTab';
import CreateUserDialog from '../components/admin/CreateUserDialog';
import CreateProjectDialog from '../components/admin/CreateProjectDialog';
import ConfirmationDialog from '../components/admin/ConfirmationDialog';
import { useNotification } from '../App';

const AdminPage: React.FC = () => {
  const { user } = useAuth();
  const { isSuperuser } = userPermissions();
  const { showNotification } = useNotification();
  const [tabValue, setTabValue] = useState(0);
  const [createUserDialogOpen, setCreateUserDialogOpen] = useState(false);
  const [createProjectDialogOpen, setCreateProjectDialogOpen] = useState(false);
  const [deleteUserDialogOpen, setDeleteUserDialogOpen] = useState(false);
  const [deleteProjectDialogOpen, setDeleteProjectDialogOpen] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null);
  const [selectedProjectId, setSelectedProjectId] = useState<string | null>(null);
  const [snackbar, setSnackbar] = useState<{
    open: boolean;
    message: string;
    severity: 'success' | 'error' | 'info' | 'warning';
  }>({ open: false, message: '', severity: 'info' });

  // Check if user is superuser
  if (!isSuperuser()) {
    return (
      <Layout>
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '50vh' }}>
          <Alert severity="error" sx={{ maxWidth: 600 }}>
            <Typography variant="h6" gutterBottom>
              Access Denied
            </Typography>
            <Typography>
              You don't have permission to access the admin panel. Only superusers can access this page.
            </Typography>
          </Alert>
        </Box>
      </Layout>
    );
  }

  // Fetch users data
  const {
    data: users,
    isLoading: isLoadingUsers,
    error: usersError,
    refetch: refetchUsers,
  } = useQuery({
    queryKey: ['users'],
    queryFn: async () => {
      const response = await apiClient.listUsers();
      return response.data;
    },
  });

  // Fetch projects data
  const {
    data: projects,
    isLoading: isLoadingProjects,
    error: projectsError,
    refetch: refetchProjects,
  } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const response = await apiClient.listProjects();
      return response.data;
    },
  });

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleCreateUser = async (username: string, email: string, password: string, isSuperuser: boolean) => {
    try {
      await apiClient.createUser({
        username,
        email,
        password,
        is_superuser: isSuperuser,
      });
      setCreateUserDialogOpen(false);
      showNotification('User created successfully', 'success');
      refetchUsers();
    } catch (error: any) {
      const message = error?.response?.data?.error?.message || 'Failed to create user';
      showNotification(message, 'error');
    }
  };

  const handleCreateProject = async (name: string, description: string) => {
    try {
      await apiClient.addProject({ name, description });
      setCreateProjectDialogOpen(false);
      showNotification('Project created successfully', 'success');
      refetchProjects();
    } catch (error: any) {
      const message = error?.response?.data?.error?.message || 'Failed to create project';
      showNotification(message, 'error');
    }
  };

  const handleToggleUserStatus = async (userId: number, isActive: boolean) => {
    try {
      await apiClient.setUserActiveStatus(userId, { is_active: isActive });
      showNotification(`User ${isActive ? 'activated' : 'deactivated'} successfully`, 'success');
      refetchUsers();
    } catch (error) {
      showNotification(`Failed to update user status: ${error}`, 'error');
    }
  };

  const handleToggleSuperuserStatus = async (userId: number, isSuperuser: boolean) => {
    try {
      await apiClient.setSuperuserStatus(userId, { is_superuser: isSuperuser });
      showNotification(`User ${isSuperuser ? 'granted' : 'revoked'} superuser status`, 'success');
      refetchUsers();
    } catch (error) {
      showNotification(`Failed to update superuser status: ${error}`, 'error');
    }
  };

  const handleDeleteUser = async (userId: number) => {
    try {
      await apiClient.deleteUser(userId);
      showNotification('User deleted successfully', 'success');
      setDeleteUserDialogOpen(false);
      setSelectedUserId(null);
      refetchUsers();
    } catch (error) {
      showNotification(`Failed to delete user: ${error}`, 'error');
    }
  };

  const handleArchiveProject = async (projectId: string) => {
    // This will be implemented with the actual API call
    console.log('Archiving project:', projectId);
    showNotification('Project archived successfully', 'success');
    setDeleteProjectDialogOpen(false);
    setSelectedProjectId(null);
    refetchProjects();
  };

  const openDeleteUserDialog = (userId: number) => {
    setSelectedUserId(userId);
    setDeleteUserDialogOpen(true);
  };

  const openDeleteProjectDialog = (projectId: string) => {
    setSelectedProjectId(projectId);
    setDeleteProjectDialogOpen(true);
  };

  return (
    <Layout>
      <Box sx={{ width: '100%' }}>
        <Box sx={{ mb: 3 }}>
          <Typography
            variant="h4"
            component="h1"
            gutterBottom
            sx={{
              fontWeight: 700,
              background: (theme) => theme.palette.mode === 'dark'
                ? 'linear-gradient(45deg, #8352ff 10%, #5e72e4 90%)'
                : 'linear-gradient(45deg, #5e72e4 30%, #8352ff 90%)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              mb: 1
            }}
          >
            Admin Panel
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'text.secondary',
              maxWidth: '800px',
              fontSize: '1.1rem'
            }}
          >
            Manage users, projects, licenses, and system configuration.
          </Typography>
        </Box>

        <Paper
          sx={{
            width: '100%',
            background: (theme) => theme.palette.mode === 'dark'
              ? 'linear-gradient(to bottom, rgba(65, 68, 74, 0.3), rgba(55, 58, 64, 0.3))'
              : 'linear-gradient(to bottom, rgba(255, 255, 255, 0.7), rgba(245, 245, 245, 0.7))',
            backdropFilter: 'blur(10px)',
            boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.03)',
            borderRadius: 2
          }}
        >
          <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Tabs
              value={tabValue}
              onChange={handleTabChange}
              aria-label="admin tabs"
              sx={{
                '& .MuiTab-root': {
                  fontWeight: 500,
                  transition: 'all 0.2s ease-in-out',
                  '&:hover': {
                    color: 'primary.main',
                    opacity: 0.8
                  }
                },
                '& .Mui-selected': {
                  fontWeight: 600
                }
              }}
            >
              <Tab 
                label="Users" 
                icon={<PeopleIcon />} 
                iconPosition="start"
              />
              <Tab 
                label="Projects" 
                icon={<ProjectsIcon />} 
                iconPosition="start"
              />
              <Tab
                label="External Auth"
                icon={<SyncIcon />}
                iconPosition="start"
              />
              <Tab 
                label="License" 
                icon={<LicenseIcon />} 
                iconPosition="start"
              />
            </Tabs>
          </Box>

          {/* Users Tab */}
          <TabPanel value={tabValue} index={0}>
            <UsersTab
              users={users}
              isLoading={isLoadingUsers}
              error={usersError}
              onCreateUser={() => setCreateUserDialogOpen(true)}
              onToggleUserStatus={handleToggleUserStatus}
              onToggleSuperuserStatus={handleToggleSuperuserStatus}
              onDeleteUser={openDeleteUserDialog}
            />
          </TabPanel>

          {/* Projects Tab */}
          <TabPanel value={tabValue} index={1}>
            <ProjectsTab
              projects={projects}
              isLoading={isLoadingProjects}
              error={projectsError}
              onCreateProject={() => setCreateProjectDialogOpen(true)}
              onArchiveProject={openDeleteProjectDialog}
              setSnackbar={setSnackbar}
            />
          </TabPanel>

          {/* External Auth Tab */}
          <TabPanel value={tabValue} index={2}>
            <ExternalAuthTab />
          </TabPanel>

          {/* License Tab */}
          <TabPanel value={tabValue} index={3}>
            <LicenseTab />
          </TabPanel>
        </Paper>

        {/* Dialogs */}
        <CreateUserDialog
          open={createUserDialogOpen}
          onClose={() => setCreateUserDialogOpen(false)}
          onCreateUser={handleCreateUser}
        />

        <CreateProjectDialog
          open={createProjectDialogOpen}
          onClose={() => setCreateProjectDialogOpen(false)}
          onCreateProject={handleCreateProject}
          isLoadingTeams={false}
        />

        <ConfirmationDialog
          open={deleteUserDialogOpen}
          onCancel={() => setDeleteUserDialogOpen(false)}
          onConfirm={() => selectedUserId && handleDeleteUser(selectedUserId)}
          title="Delete User"
          message="Are you sure you want to delete this user? This action cannot be undone."
          confirmButtonText="Delete"
          cancelButtonText="Cancel"
        />

        <ConfirmationDialog
          open={deleteProjectDialogOpen}
          onCancel={() => setDeleteProjectDialogOpen(false)}
          onConfirm={() => selectedProjectId && handleArchiveProject(selectedProjectId)}
          title="Archive Project"
          message="Are you sure you want to archive this project? This action cannot be undone."
          confirmButtonText="Archive"
          cancelButtonText="Cancel"
        />
      </Box>
    </Layout>
  );
};

export default AdminPage;
