import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Button,
  Paper,
  Grid,
  CircularProgress,
  Alert,
  Card,
  CardContent,
  Chip,
  Stack,
} from '@mui/material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../../api/apiClient';
import { useNotification } from '../../../App';
import LDAPSyncProgressBar from '../../LDAPSyncProgressBar';

// LDAP Sync Tab component
const LDAPSyncTab: React.FC = () => {
  const { showNotification } = useNotification();
  const queryClient = useQueryClient();
  const [syncStarted, setSyncStarted] = useState(false);

  // Fetch LDAP sync status
  const {
    data: syncStatus,
    isLoading: isLoadingStatus,
    error: statusError
  } = useQuery({
    queryKey: ['ldapSyncStatus'],
    queryFn: async () => {
      try {
        const response = await apiClient.getLDAPSyncStatus();
        return response.data;
      } catch (error) {
        console.error('Error fetching LDAP sync status:', error);
        return {
          is_running: false,
          last_sync_time: null,
          status: '',
          total_users: 0,
          synced_users: 0,
          errors: 0,
          warnings: 0,
          last_sync_duration: ''
        };
      }
    },
    refetchInterval: 5000 // Refresh every 5 seconds
  });

  // Fetch LDAP sync progress when sync is running or just started
  const {
    data: syncProgress,
    isLoading: isLoadingProgress
  } = useQuery({
    queryKey: ['ldapSyncProgress'],
    queryFn: async () => {
      try {
        const response = await apiClient.getLDAPSyncProgress();
        return response.data;
      } catch (error) {
        console.error('Error fetching LDAP sync progress:', error);
        return {
          is_running: false,
          progress: 0,
          current_step: '',
          processed_items: 0,
          total_items: 0,
          estimated_time: '',
          start_time: null
        };
      }
    },
    refetchInterval: (syncStatus?.is_running || syncStarted) ? 1000 : false, // Refresh every second when sync is running or just started
    enabled: !!(syncStatus?.is_running || syncStarted)
  });

  // Mutation for starting user sync
  const syncUsersMutation = useMutation({
    mutationFn: async () => {
      const response = await apiClient.syncLDAPUsers();
      return response.data;
    },
    onSuccess: () => {
      showNotification('User synchronization started successfully', 'success');
      setSyncStarted(true);
      // Add delay to allow sync to start in the backend
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey: ['ldapSyncStatus'] });
        queryClient.invalidateQueries({ queryKey: ['ldapSyncProgress'] });
      }, 1000);
    },
    onError: (error) => {
      showNotification(`Failed to start user synchronization: ${error}`, 'error');
    }
  });

  // Mutation for canceling sync
  const cancelSyncMutation = useMutation({
    mutationFn: async () => {
      const response = await apiClient.cancelLDAPSync();
      return response.data;
    },
    onSuccess: () => {
      showNotification('Synchronization cancelled successfully', 'success');
      queryClient.invalidateQueries({ queryKey: ['ldapSyncStatus'] });
      queryClient.invalidateQueries({ queryKey: ['ldapSyncProgress'] });
    },
    onError: (error) => {
      showNotification(`Failed to cancel synchronization: ${error}`, 'error');
    }
  });

  // Mutation for testing connection
  const testConnectionMutation = useMutation({
    mutationFn: async () => {
      const response = await apiClient.testLDAPConnection({
        url: '',
        bind_dn: '',
        bind_password: '',
        user_base_dn: '',
        user_filter: '',
        user_name_attr: '',
        start_tls: false,
        insecure_tls: false,
        timeout: ''
      });
      return response.data;
    },
    onSuccess: () => {
      showNotification('LDAP connection test successful', 'success');
    },
    onError: (error) => {
      showNotification(`LDAP connection test failed: ${error}`, 'error');
    }
  });

  // Format date for display
  const formatDate = (dateString: string | null | undefined) => {
    if (!dateString) return 'Never';
    
    try {
      const date = new Date(dateString);
      // Check if date is valid and not the zero date
      if (isNaN(date.getTime()) || date.getFullYear() === 1) {
        return 'Never';
      }
      return date.toLocaleString();
    } catch {
      return 'Never';
    }
  };

  // Function to get status color
  const getStatusColor = (status: string | undefined) => {
    if (!status) return 'default';
    
    switch (status.toLowerCase()) {
      case 'completed':
        return 'success';
      case 'failed':
        return 'error';
      case 'cancelled':
        return 'warning';
      case 'running':
        return 'primary';
      default:
        return 'default';
    }
  };

  // Function to get status label
  const getStatusLabel = (status: string | undefined, isRunning: boolean | undefined) => {
    if (isRunning) return 'Running';
    if (!status) return 'Idle';
    
    return status.charAt(0).toUpperCase() + status.slice(1);
  };

  // Handle sync button clicks
  const handleSyncUsers = () => {
    syncUsersMutation.mutate();
  };

  const handleCancelSync = () => {
    cancelSyncMutation.mutate();
  };

  const handleTestConnection = () => {
    testConnectionMutation.mutate();
  };

  // Reset syncStarted when sync is no longer running
  useEffect(() => {
    if (!syncStatus?.is_running && syncStarted) {
      setSyncStarted(false);
    }
  }, [syncStatus?.is_running, syncStarted]);

  // Loading state
  if (isLoadingStatus) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  // Error state
  if (statusError) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        Error loading LDAP sync status. Please try again later.
      </Alert>
    );
  }

  return (
    <>
      <Box 
        sx={{ 
          display: 'flex', 
          flexDirection: 'column',
          mb: 4,
          pb: 2,
          borderBottom: (theme) => `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)'}`
        }}
      >
        <Typography 
          variant="h6"
          sx={{ 
            fontWeight: 600,
            mb: 0.5
          }}
        >
          LDAP Synchronization
        </Typography>
        <Typography 
          variant="body2" 
          color="text.secondary"
          sx={{ maxWidth: '600px' }}
        >
          Manage LDAP synchronization settings and manually trigger synchronization of users and teams.
        </Typography>
      </Box>

      {/* Sync Controls */}
      <Paper 
        sx={{ 
          p: 3, 
          mb: 3,
          borderRadius: 2,
          background: (theme) => theme.palette.mode === 'dark'
            ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
            : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
          backdropFilter: 'blur(8px)',
          boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)'
        }}
      >
        <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>Synchronization Controls</Typography>

        <Grid container spacing={2} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={4}>
            <Button 
              variant="contained" 
              color="primary" 
              fullWidth
              onClick={handleSyncUsers}
              disabled={syncStatus?.is_running || syncUsersMutation.isPending}
              startIcon={syncUsersMutation.isPending ? <CircularProgress size={20} /> : null}
            >
              Sync Users
            </Button>
          </Grid>
        </Grid>

        {syncStatus?.is_running && (
          <Box sx={{ mt: 2 }}>
            <Button 
              variant="outlined" 
              color="error" 
              onClick={handleCancelSync}
              disabled={cancelSyncMutation.isPending}
              startIcon={cancelSyncMutation.isPending ? <CircularProgress size={20} /> : null}
            >
              Cancel Sync
            </Button>
          </Box>
        )}
      </Paper>

      {/* Sync Progress */}
      {syncStatus?.is_running && (
        isLoadingProgress ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', my: 2 }}>
            <CircularProgress size={32} />
          </Box>
        ) : (
          <LDAPSyncProgressBar
            isRunning={!!syncProgress?.is_running}
            progress={typeof syncProgress?.progress === 'number' ? syncProgress.progress : 0}
            currentStep={syncProgress?.current_step}
            processedItems={syncProgress?.processed_items}
            totalItems={syncProgress?.total_items}
            estimatedTime={syncProgress?.estimated_time}
            startTime={syncProgress?.start_time}
          />
        )
      )}

      {/* Sync Status */}
      <Paper 
        sx={{ 
          p: 3, 
          mb: 3,
          borderRadius: 2,
          background: (theme) => theme.palette.mode === 'dark'
            ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
            : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
          backdropFilter: 'blur(8px)',
          boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)'
        }}
      >
        <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>Sync Status</Typography>

        <Grid container spacing={3}>
          <Grid item xs={12} sm={6}>
            <Card variant="outlined" sx={{ height: '100%' }}>
              <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                  <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>Last Sync</Typography>
                  <Chip 
                    label={getStatusLabel(syncStatus?.status, syncStatus?.is_running)} 
                    color={getStatusColor(syncStatus?.status)}
                    size="small"
                    variant="outlined"
                  />
                </Box>
                <Stack spacing={1}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2" color="text.secondary">Last Run:</Typography>
                    <Typography variant="body2">{formatDate(syncStatus?.last_sync_time)}</Typography>
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2" color="text.secondary">Duration:</Typography>
                    <Typography variant="body2">{syncStatus?.last_sync_duration || 'N/A'}</Typography>
                  </Box>
                </Stack>
              </CardContent>
            </Card>
          </Grid>

          <Grid item xs={12} sm={6}>
            <Card variant="outlined" sx={{ height: '100%' }}>
              <CardContent>
                <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 600 }}>Statistics</Typography>
                <Stack spacing={1}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2" color="text.secondary">Total Users:</Typography>
                    <Typography variant="body2">{syncStatus?.total_users || 0}</Typography>
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Typography variant="body2" color="text.secondary">Synced Users:</Typography>
                    <Typography variant="body2">{syncStatus?.synced_users || 0}</Typography>
                  </Box>
                </Stack>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      </Paper>

      {/* Connection Test */}
      <Paper 
        sx={{ 
          p: 3, 
          mb: 3,
          borderRadius: 2,
          background: (theme) => theme.palette.mode === 'dark'
            ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
            : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
          backdropFilter: 'blur(8px)',
          boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)'
        }}
      >
        <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>Connection Test</Typography>

        <Box sx={{ mb: 2 }}>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Test the connection to your LDAP server with the current configuration.
          </Typography>

          <Button 
            variant="outlined" 
            color="primary" 
            onClick={handleTestConnection}
            disabled={testConnectionMutation.isPending}
            startIcon={testConnectionMutation.isPending ? <CircularProgress size={20} /> : null}
          >
            Test Connection
          </Button>
        </Box>
      </Paper>
    </>
  );
};

export default LDAPSyncTab;
