import React from 'react';
import { Box, Typography, Paper, Grid, Avatar, Chip, Divider } from '@mui/material';
import { Person as PersonIcon, Email as EmailIcon, AdminPanelSettings as AdminIcon } from '@mui/icons-material';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import TwoFactorAuthSection from '../components/TwoFactorAuthSection';
import { useAuth } from '../auth/AuthContext';

const AccountPage: React.FC = () => {
  const { user, isLoading, error } = useAuth();

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
        </Grid>
      </Box>
    </AuthenticatedLayout>
  );
};

export default AccountPage;
