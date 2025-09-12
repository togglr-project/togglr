import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  TextField,
  Button,
  Switch,
  FormControlLabel,
  Grid,
  Alert,
  CircularProgress,
  Divider,
} from '@mui/material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../../api/apiClient';
import { useNotification } from '../../../App';

const LDAPConfigTab: React.FC = () => {
  const { showNotification } = useNotification();
  const queryClient = useQueryClient();

  // Form state
  const [formData, setFormData] = useState({
    enabled: false,
    url: '',
    bind_dn: '',
    bind_password: '',
    user_base_dn: '',
    user_filter: '',
    user_name_attr: '',
    user_email_attr: '',
    start_tls: false,
    insecure_tls: false,
    timeout: '30s',
    sync_interval: 3600,
    background_sync_enabled: true,
  });

  // Fetch current LDAP config
  const {
    data: currentConfig,
    isLoading: isLoadingConfig,
    error: configError
  } = useQuery({
    queryKey: ['ldapConfig'],
    queryFn: async () => {
      try {
        const response = await apiClient.getLDAPConfig();
        return response.data;
      } catch (error) {
        console.error('Error fetching LDAP config:', error);
        return null;
      }
    }
  });

  // Update form data when config is loaded
  React.useEffect(() => {
    if (currentConfig) {
      setFormData({
        enabled: currentConfig.enabled,
        url: currentConfig.url,
        bind_dn: currentConfig.bind_dn,
        bind_password: currentConfig.bind_password,
        user_base_dn: currentConfig.user_base_dn,
        user_filter: currentConfig.user_filter,
        user_name_attr: currentConfig.user_name_attr,
        user_email_attr: currentConfig.user_email_attr,
        start_tls: currentConfig.start_tls,
        insecure_tls: currentConfig.insecure_tls,
        timeout: currentConfig.timeout,
        sync_interval: currentConfig.sync_interval || 3600,
        background_sync_enabled: (currentConfig.sync_interval || 0) > 0,
      });
    }
  }, [currentConfig]);

  // Mutation for updating LDAP config
  const updateConfigMutation = useMutation({
    mutationFn: async (config: typeof formData) => {
      const response = await apiClient.updateLDAPConfig(config);
      return response.data;
    },
    onSuccess: () => {
      showNotification('LDAP configuration updated successfully', 'success');
      queryClient.invalidateQueries({ queryKey: ['ldapConfig'] });
    },
    onError: (error) => {
      showNotification(`Failed to update LDAP configuration: ${error}`, 'error');
    }
  });

  // Mutation for testing connection
  const testConnectionMutation = useMutation({
    mutationFn: async () => {
      const response = await apiClient.testLDAPConnection({
        url: formData.url,
        bind_dn: formData.bind_dn,
        bind_password: formData.bind_password,
        user_base_dn: formData.user_base_dn,
        user_filter: formData.user_filter,
        user_name_attr: formData.user_name_attr,
        start_tls: formData.start_tls,
        insecure_tls: formData.insecure_tls,
        timeout: formData.timeout
      });
      return response.data;
    },
    onSuccess: (data) => {
      if (data.success) {
        showNotification('LDAP connection test successful', 'success');
      } else {
        showNotification(`LDAP connection test failed: ${data.message}`, 'error');
      }
    },
    onError: (error) => {
      showNotification(`LDAP connection test failed: ${error}`, 'error');
    }
  });

  const handleInputChange = (field: keyof typeof formData) => (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    if (event.target.type === 'checkbox') {
      const value = event.target.checked;
      
      // Special handling for background_sync_enabled
      if (field === 'background_sync_enabled') {
        setFormData(prev => ({
          ...prev,
          background_sync_enabled: value,
          sync_interval: value ? (prev.sync_interval || 3600) : 0
        }));
        return;
      }
      
      setFormData(prev => ({
        ...prev,
        [field]: value
      }));
    } else if (event.target.type === 'number') {
      const value = parseInt(event.target.value) || 0;
      setFormData(prev => ({
        ...prev,
        [field]: value
      }));
    } else {
      const value = event.target.value;
      setFormData(prev => ({
        ...prev,
        [field]: value
      }));
    }
  };

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    
    // Prepare config data for API
    const configData = {
      ...formData,
      sync_interval: formData.background_sync_enabled ? formData.sync_interval : 0
    };
    
    updateConfigMutation.mutate(configData);
  };

  const handleTestConnection = () => {
    testConnectionMutation.mutate();
  };

  if (isLoadingConfig) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (configError) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        Error loading LDAP configuration. Please try again later.
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
          LDAP Configuration
        </Typography>
        <Typography 
          variant="body2" 
          color="text.secondary"
          sx={{ maxWidth: '600px' }}
        >
          Configure LDAP server connection settings for user and team synchronization.
        </Typography>
      </Box>

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
        <form onSubmit={handleSubmit}>
          <Typography variant="h6" sx={{ mb: 3, fontWeight: 600 }}>Connection Settings</Typography>

          <FormControlLabel
            control={
              <Switch
                checked={formData.enabled}
                onChange={handleInputChange('enabled')}
                color="primary"
              />
            }
            label="Enable LDAP Integration"
            sx={{ mb: 3 }}
          />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="LDAP Server URL"
                value={formData.url}
                onChange={handleInputChange('url')}
                placeholder="ldap://ldap.example.com:389"
                helperText="LDAP server URL with protocol and port"
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Bind DN"
                value={formData.bind_dn}
                onChange={handleInputChange('bind_dn')}
                placeholder="cn=admin,dc=example,dc=com"
                helperText="DN for binding to LDAP server"
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Bind Password"
                type="password"
                value={formData.bind_password}
                onChange={handleInputChange('bind_password')}
                placeholder="Password for bind DN"
                helperText="Password for the bind DN"
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Connection Timeout"
                value={formData.timeout}
                onChange={handleInputChange('timeout')}
                placeholder="30s"
                helperText="Connection timeout (e.g., 30s, 1m)"
              />
            </Grid>
          </Grid>

          <Box sx={{ mt: 3, mb: 3 }}>
            <FormControlLabel
              control={
                <Switch
                  checked={formData.background_sync_enabled}
                  onChange={handleInputChange('background_sync_enabled')}
                  color="primary"
                />
              }
              label="Enable background synchronization of users and groups"
            />
          </Box>

          {formData.background_sync_enabled && (
            <Grid container spacing={3} sx={{ mb: 3 }}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Sync Interval (seconds)"
                  type="number"
                  value={formData.sync_interval}
                  onChange={handleInputChange('sync_interval')}
                  placeholder="3600"
                  helperText="Background sync interval in seconds (default: 3600 = 1 hour)"
                  inputProps={{ min: 60, max: 86400 }}
                />
              </Grid>
            </Grid>
          )}

          <Box sx={{ mt: 3, mb: 3 }}>
            <FormControlLabel
              control={
                <Switch
                  checked={formData.start_tls}
                  onChange={handleInputChange('start_tls')}
                  color="primary"
                />
              }
              label="Use StartTLS"
            />
            <FormControlLabel
              control={
                <Switch
                  checked={formData.insecure_tls}
                  onChange={handleInputChange('insecure_tls')}
                  color="primary"
                />
              }
              label="Skip TLS Certificate Verification"
              sx={{ ml: 4 }}
            />
          </Box>

          <Divider sx={{ my: 3 }} />

          <Typography variant="h6" sx={{ mb: 3, fontWeight: 600 }}>User Configuration</Typography>

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="User Base DN"
                value={formData.user_base_dn}
                onChange={handleInputChange('user_base_dn')}
                placeholder="ou=users,dc=example,dc=com"
                helperText="Base DN for user search"
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="User Filter"
                value={formData.user_filter}
                onChange={handleInputChange('user_filter')}
                placeholder="(objectClass=person)"
                helperText="LDAP filter for user search"
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Username Attribute"
                value={formData.user_name_attr}
                onChange={handleInputChange('user_name_attr')}
                placeholder="uid"
                helperText="Attribute containing username"
                required
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Email Attribute"
                value={formData.user_email_attr}
                onChange={handleInputChange('user_email_attr')}
                placeholder="mail"
                helperText="Attribute containing email address"
                required
              />
            </Grid>
          </Grid>

          <Divider sx={{ my: 3 }} />

          <Typography variant="h6" sx={{ mb: 3, fontWeight: 600 }}>Group Configuration (Optional)</Typography>

          <Box sx={{ mt: 4, display: 'flex', gap: 2 }}>
            <Button
              type="submit"
              variant="contained"
              color="primary"
              disabled={updateConfigMutation.isPending}
              startIcon={updateConfigMutation.isPending ? <CircularProgress size={20} /> : null}
            >
              Save Configuration
            </Button>
            <Button
              variant="outlined"
              color="primary"
              onClick={handleTestConnection}
              disabled={testConnectionMutation.isPending || !formData.url || !formData.bind_dn}
              startIcon={testConnectionMutation.isPending ? <CircularProgress size={20} /> : null}
            >
              Test Connection
            </Button>
          </Box>
        </form>
      </Paper>
    </>
  );
};

export default LDAPConfigTab; 