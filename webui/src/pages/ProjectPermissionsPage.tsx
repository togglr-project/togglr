import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  CircularProgress,
  Button,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  Card,
  CardContent,
  CardActions,
  Tooltip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Tabs,
  Tab,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Divider,
  Avatar,
  Badge,
} from '@mui/material';
import { 
  People as PeopleIcon, 
  Add as AddIcon, 
  Edit as EditIcon,
  Delete as DeleteIcon,
  PersonAdd as PersonAddIcon,
  AdminPanelSettings as AdminIcon,
  Security as SecurityIcon,
  ExpandMore as ExpandMoreIcon,
  Check as CheckIcon,
  Close as CloseIcon
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams, Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import { useRBAC } from '../auth/permissions';
import type { 
  Membership, 
  Role, 
  Permission, 
  User,
  CreateMembershipRequest,
  UpdateMembershipRequest,
  ListRolePermissions200ResponseInner
} from '../generated/api/client';

const ProjectPermissionsPage: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const { user } = useAuth();
  const rbac = useRBAC(projectId);
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState(0);
  const [addMemberOpen, setAddMemberOpen] = useState(false);
  const [editMemberOpen, setEditMemberOpen] = useState(false);
  const [deleteMemberOpen, setDeleteMemberOpen] = useState(false);
  const [selectedMembership, setSelectedMembership] = useState<Membership | null>(null);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [selectedRole, setSelectedRole] = useState<string>('');
  const [error, setError] = useState<string | null>(null);

  // Check project access
  if (!rbac.canViewProject()) {
    return (
      <AuthenticatedLayout showBackButton backTo="/dashboard">
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="error" gutterBottom>
            Access Denied
          </Typography>
          <Typography variant="body2" color="text.secondary">
            You don't have permission to view this project.
          </Typography>
        </Box>
      </AuthenticatedLayout>
    );
  }

  // Check membership management permissions
  if (!rbac.canManageMembership()) {
    return (
      <AuthenticatedLayout showBackButton backTo="/dashboard">
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="error" gutterBottom>
            Access Denied
          </Typography>
          <Typography variant="body2" color="text.secondary">
            You don't have permission to manage project memberships.
          </Typography>
        </Box>
      </AuthenticatedLayout>
    );
  }

  // Fetch project memberships
  const { data: memberships, isLoading: membershipsLoading } = useQuery<Membership[]>({
    queryKey: ['project-memberships', projectId],
    queryFn: async () => {
      if (!projectId) throw new Error('Project ID is required');
      const res = await apiClient.listProjectMemberships(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  // Fetch all users (for adding new members)
  const { data: allUsers, isLoading: usersLoading } = useQuery<User[]>({
    queryKey: ['users'],
    queryFn: async () => {
      const res = await apiClient.listUsers();
      return res.data;
    },
  });

  // Fetch roles
  const { data: roles, isLoading: rolesLoading } = useQuery<Role[]>({
    queryKey: ['roles'],
    queryFn: async () => {
      const res = await apiClient.listRoles();
      return res.data;
    },
  });

  // Fetch permissions
  const { data: permissions, isLoading: permissionsLoading } = useQuery<Permission[]>({
    queryKey: ['permissions'],
    queryFn: async () => {
      const res = await apiClient.listPermissions();
      return res.data;
    },
  });

  // Fetch role permissions for all roles
  const { data: rolePermissions, isLoading: rolePermissionsLoading } = useQuery<Record<string, string[]>>({
    queryKey: ['role-permissions'],
    queryFn: async () => {
      const res = await apiClient.listRolePermissions();
      // console.log('listRolePermissions response:', res.data);
      
      // Group permissions by role_id
      const rolePerms: Record<string, string[]> = {};
      res.data.forEach((item: ListRolePermissions200ResponseInner) => {
        if (item.role && item.permissions) {
          const roleId = item.role.id;
          if (!rolePerms[roleId]) {
            rolePerms[roleId] = [];
          }
          item.permissions.forEach(permission => {
            if (permission.key) {
              rolePerms[roleId].push(permission.key);
            }
          });
        }
      });
      
      // console.log('Processed role permissions:', rolePerms);
      return rolePerms;
    },
  });

  // Get users not yet in project
  const availableUsers = allUsers?.filter(user => 
    !memberships?.some(membership => membership.user_id === user.id)
  ) || [];

  // Create membership mutation
  const createMembershipMutation = useMutation({
    mutationFn: async (data: CreateMembershipRequest) => {
      if (!projectId) throw new Error('Project ID is required');
      const res = await apiClient.createProjectMembership(projectId, data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project-memberships', projectId] });
      setAddMemberOpen(false);
      setSelectedUser(null);
      setSelectedRole('');
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to add member');
    },
  });

  // Update membership mutation
  const updateMembershipMutation = useMutation({
    mutationFn: async (data: { membershipId: string; updateData: UpdateMembershipRequest }) => {
      if (!projectId) throw new Error('Project ID is required');
      const res = await apiClient.updateProjectMembership(projectId, data.membershipId, data.updateData);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project-memberships', projectId] });
      setEditMemberOpen(false);
      setSelectedMembership(null);
      setSelectedRole('');
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to update member');
    },
  });

  // Delete membership mutation
  const deleteMembershipMutation = useMutation({
    mutationFn: async (membershipId: string) => {
      if (!projectId) throw new Error('Project ID is required');
      await apiClient.deleteProjectMembership(projectId, membershipId);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project-memberships', projectId] });
      setDeleteMemberOpen(false);
      setSelectedMembership(null);
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to remove member');
    },
  });

  const handleAddMember = () => {
    if (!selectedUser || !selectedRole) {
      setError('Please select a user and role');
      return;
    }

    createMembershipMutation.mutate({
      user_id: selectedUser.id,
      role_id: selectedRole,
    });
  };

  const handleEditMember = () => {
    if (!selectedMembership || !selectedRole) {
      setError('Please select a role');
      return;
    }

    updateMembershipMutation.mutate({
      membershipId: selectedMembership.id,
      updateData: {
        role_id: selectedRole,
      },
    });
  };

  const handleDeleteMember = () => {
    if (!selectedMembership) return;

    deleteMembershipMutation.mutate(selectedMembership.id);
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  const getRolePermissions = (roleId: string) => {
    const perms = rolePermissions?.[roleId] || [];
    // console.log(`getRolePermissions for role ${roleId}:`, perms);
    
    // Ensure we return an array of strings
    if (Array.isArray(perms)) {
      const filteredPerms = perms.filter(perm => typeof perm === 'string');
      // console.log(`Filtered permissions for role ${roleId}:`, filteredPerms);
      return filteredPerms;
    }
    return [];
  };

  const getUserInitials = (user: User) => {
    if (user.username) {
      return user.username.charAt(0).toUpperCase();
    }
    return 'U';
  };

  if (!projectId) {
    return <Navigate to="/projects" replace />;
  }

  return (
    <AuthenticatedLayout>
      <Box sx={{ p: 3 }}>
        <PageHeader
          title="Project Permissions"
          subtitle="Manage team members and their roles"
          icon={<PeopleIcon />}
        >
          <Button
            variant="contained"
            startIcon={<PersonAddIcon />}
            onClick={() => setAddMemberOpen(true)}
            disabled={availableUsers.length === 0}
          >
            Add Member
          </Button>
        </PageHeader>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {/* Tabs */}
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={activeTab} onChange={handleTabChange} aria-label="permissions tabs">
            <Tab label="Team Members" />
            <Tab label="Roles & Permissions" />
          </Tabs>
        </Box>

        {/* Team Members Tab */}
        {activeTab === 0 && (
          <Box>
            {membershipsLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                <CircularProgress />
              </Box>
            ) : memberships && memberships.length > 0 ? (
              <TableContainer component={Paper}>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>User</TableCell>
                      <TableCell>Role</TableCell>
                      <TableCell>Permissions</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {memberships.map((membership) => {
                      const role = roles?.find(r => r.id === membership.role_id);
                      const user = allUsers?.find(u => u.id === membership.user_id);
                      
                      return (
                        <TableRow key={membership.id}>
                          <TableCell>
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                              <Avatar sx={{ bgcolor: 'primary.main' }}>
                                {user ? getUserInitials(user) : 'U'}
                              </Avatar>
                              <Box>
                                <Typography variant="subtitle2">
                                  {user?.username || 'Unknown User'}
                                </Typography>
                                <Typography variant="body2" color="text.secondary">
                                  {user?.email || 'No email'}
                                </Typography>
                              </Box>
                            </Box>
                          </TableCell>
                          <TableCell>
                            <Chip
                              label={role?.name || 'Unknown Role'}
                              color={role?.key === 'admin' ? 'warning' : 'default'}
                              size="small"
                              icon={role?.key === 'admin' ? <AdminIcon /> : undefined}
                            />
                          </TableCell>
                          <TableCell>
                            <Typography variant="body2" color="text.secondary">
                              {role ? `${getRolePermissions(role.id).length} permissions` : 'Loading...'}
                            </Typography>
                            {role && getRolePermissions(role.id).length > 0 && (
                              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 1 }}>
                                {getRolePermissions(role.id).slice(0, 3).map((permission) => (
                                  <Chip
                                    key={permission}
                                    label={permission}
                                    size="small"
                                    variant="outlined"
                                    sx={{ fontSize: '0.7rem', height: 20 }}
                                  />
                                ))}
                                {getRolePermissions(role.id).length > 3 && (
                                  <Chip
                                    label={`+${getRolePermissions(role.id).length - 3} more`}
                                    size="small"
                                    variant="outlined"
                                    sx={{ fontSize: '0.7rem', height: 20 }}
                                  />
                                )}
                              </Box>
                            )}
                          </TableCell>
                          <TableCell>
                            <Box sx={{ display: 'flex', gap: 1 }}>
                              <Tooltip title="Edit Role">
                                <IconButton
                                  size="small"
                                  onClick={() => {
                                    setSelectedMembership(membership);
                                    setSelectedRole(membership.role_id);
                                    setEditMemberOpen(true);
                                  }}
                                >
                                  <EditIcon />
                                </IconButton>
                              </Tooltip>
                              <Tooltip title="Remove Member">
                                <IconButton
                                  size="small"
                                  color="error"
                                  onClick={() => {
                                    setSelectedMembership(membership);
                                    setDeleteMemberOpen(true);
                                  }}
                                >
                                  <DeleteIcon />
                                </IconButton>
                              </Tooltip>
                            </Box>
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </TableContainer>
            ) : (
              <Paper sx={{ p: 4, textAlign: 'center' }}>
                <PeopleIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
                <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
                  No team members found
                </Typography>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                  Add team members to start collaborating on this project.
                </Typography>
                <Button
                  variant="contained"
                  startIcon={<PersonAddIcon />}
                  onClick={() => setAddMemberOpen(true)}
                  disabled={availableUsers.length === 0}
                >
                  Add First Member
                </Button>
              </Paper>
            )}
          </Box>
        )}

        {/* Roles & Permissions Tab */}
        {activeTab === 1 && (
          <Box>
            {rolesLoading || permissionsLoading || rolePermissionsLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                <CircularProgress />
              </Box>
            ) : (
              <Grid container spacing={3}>
                {roles?.map((role) => (
                  <Grid item xs={12} md={6} key={role.id}>
                    <Card>
                      <CardContent>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
                          <SecurityIcon color="primary" />
                          <Typography variant="h6" sx={{ flexGrow: 1 }}>
                            {role.name}
                          </Typography>
                          <Chip
                            label={role.key}
                            size="small"
                            color={role.key === 'admin' ? 'warning' : 'default'}
                          />
                        </Box>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                          {role.description || 'No description'}
                        </Typography>
                        <Typography variant="subtitle2" gutterBottom>
                          Permissions:
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {getRolePermissions(role.id).map((permission) => (
                            <Chip
                              key={permission}
                              label={permission}
                              size="small"
                              variant="outlined"
                            />
                          ))}
                        </Box>
                      </CardContent>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            )}
          </Box>
        )}

        {/* Add Member Dialog */}
        <Dialog open={addMemberOpen} onClose={() => setAddMemberOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Add Team Member</DialogTitle>
          <DialogContent>
            <FormControl fullWidth margin="normal">
              <InputLabel>User</InputLabel>
              <Select
                value={selectedUser?.id || ''}
                onChange={(e) => {
                  const user = allUsers?.find(u => u.id === e.target.value);
                  setSelectedUser(user || null);
                }}
                label="User"
              >
                {availableUsers.map((user) => (
                  <MenuItem key={user.id} value={user.id}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Avatar sx={{ width: 24, height: 24, bgcolor: 'primary.main' }}>
                        {getUserInitials(user)}
                      </Avatar>
                      <Box>
                        <Typography variant="body2">{user.username}</Typography>
                        <Typography variant="caption" color="text.secondary">
                          {user.email}
                        </Typography>
                      </Box>
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            <FormControl fullWidth margin="normal">
              <InputLabel>Role</InputLabel>
              <Select
                value={selectedRole}
                onChange={(e) => setSelectedRole(e.target.value)}
                label="Role"
              >
                {roles?.map((role) => (
                  <MenuItem key={role.id} value={role.id}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <SecurityIcon />
                      <Box>
                        <Typography variant="body2">{role.name}</Typography>
                        <Typography variant="caption" color="text.secondary">
                          {role.description}
                        </Typography>
                      </Box>
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setAddMemberOpen(false)}>Cancel</Button>
            <Button
              onClick={handleAddMember}
              variant="contained"
              disabled={!selectedUser || !selectedRole || createMembershipMutation.isPending}
            >
              Add Member
            </Button>
          </DialogActions>
        </Dialog>

        {/* Edit Member Dialog */}
        <Dialog open={editMemberOpen} onClose={() => setEditMemberOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Edit Team Member Role</DialogTitle>
          <DialogContent>
            <FormControl fullWidth margin="normal">
              <InputLabel>Role</InputLabel>
              <Select
                value={selectedRole}
                onChange={(e) => setSelectedRole(e.target.value)}
                label="Role"
              >
                {roles?.map((role) => (
                  <MenuItem key={role.id} value={role.id}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <SecurityIcon />
                      <Box>
                        <Typography variant="body2">{role.name}</Typography>
                        <Typography variant="caption" color="text.secondary">
                          {role.description}
                        </Typography>
                      </Box>
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setEditMemberOpen(false)}>Cancel</Button>
            <Button
              onClick={handleEditMember}
              variant="contained"
              disabled={!selectedRole || updateMembershipMutation.isPending}
            >
              Update Role
            </Button>
          </DialogActions>
        </Dialog>

        {/* Delete Member Dialog */}
        <Dialog open={deleteMemberOpen} onClose={() => setDeleteMemberOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Remove Team Member</DialogTitle>
          <DialogContent>
            <Typography>
              Are you sure you want to remove this team member from the project? 
              This action cannot be undone.
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDeleteMemberOpen(false)}>Cancel</Button>
            <Button
              onClick={handleDeleteMember}
              variant="contained"
              color="error"
              disabled={deleteMembershipMutation.isPending}
            >
              Remove Member
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </AuthenticatedLayout>
  );
};

export default ProjectPermissionsPage;
