import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  TextField,
  Alert,
  Chip,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  CircularProgress,
  Divider,
  IconButton,
  Tooltip,
} from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '../../api/apiClient';
import { useNotification } from '../../App';
import { LicenseType } from '../../generated/api/client';

const LicenseTab: React.FC = () => {
  const { showNotification } = useNotification();
  const queryClient = useQueryClient();
  const [openUpdateDialog, setOpenUpdateDialog] = useState(false);
  const [newLicense, setNewLicense] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Fetch license status
  const {
    data: licenseStatus,
    isLoading,
    error,
  } = useQuery({
    queryKey: ['licenseStatus'],
    queryFn: async () => {
      const response = await apiClient.getLicenseStatus();
      return response.data;
    },
  });

  // Fetch product info
  const {
    data: productInfo,
    isLoading: isLoadingProductInfo,
    error: productInfoError,
  } = useQuery({
    queryKey: ['productInfo'],
    queryFn: async () => {
      const response = await apiClient.getProductInfo();
      return response.data;
    },
  });

  // Update license mutation
  const updateLicenseMutation = useMutation({
    mutationFn: async (licenseText: string) => {
      const response = await apiClient.updateLicense({
        license_text: licenseText,
      });
      return response.data;
    },
    onSuccess: () => {
      showNotification('License updated successfully!', 'success');
      setOpenUpdateDialog(false);
      setNewLicense('');
      queryClient.invalidateQueries({ queryKey: ['licenseStatus'] });
    },
    onError: (error: unknown) => {
      console.error('Failed to update license:', error);
      const errorMessage = (error as { response?: { data?: { message?: string } } })?.response?.data?.message || 'Failed to update license';
      showNotification(
        errorMessage,
        'error'
      );
    },
  });

  const handleOpenUpdateDialog = () => {
    setNewLicense('');
    setOpenUpdateDialog(true);
  };

  const handleCloseUpdateDialog = () => {
    setOpenUpdateDialog(false);
    setNewLicense('');
  };

  const handleUpdateLicense = async () => {
    if (!newLicense.trim()) {
      showNotification('Please enter a license key', 'error');
      return;
    }

    setIsSubmitting(true);
    try {
      await updateLicenseMutation.mutateAsync(newLicense);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCopyClientId = async () => {
    if (!productInfo?.client_id) {
      showNotification('No client ID available to copy', 'error');
      return;
    }

    try {
      await navigator.clipboard.writeText(productInfo.client_id);
      showNotification('Client ID copied to clipboard!', 'success');
    } catch (error) {
      console.error('Failed to copy to clipboard:', error);
      showNotification('Failed to copy Client ID to clipboard', 'error');
    }
  };

  const handleCopyLicenseText = async () => {
    if (!licenseStatus?.license?.license_text) {
      showNotification('No license text available to copy', 'error');
      return;
    }

    try {
      await navigator.clipboard.writeText(licenseStatus.license.license_text);
      showNotification('License text copied to clipboard!', 'success');
    } catch (error) {
      console.error('Failed to copy to clipboard:', error);
      showNotification('Failed to copy license text to clipboard', 'error');
    }
  };

  const formatDate = (dateString: string) => {
    try {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      });
    } catch {
      return dateString;
    }
  };

  const getLicenseTypeLabel = (type: LicenseType) => {
    switch (type) {
      case LicenseType.Trial:
      case LicenseType.TrialSelfSigned:
        return 'Trial';
      case LicenseType.Commercial:
        return 'Commercial';
      case LicenseType.Individual:
        return 'Individual';
      default:
        return type;
    }
  };

  const getStatusColor = (isValid: boolean) => {
    return isValid ? 'success' : 'error';
  };

  const getStatusText = (isValid: boolean) => {
    return isValid ? 'Valid' : 'Invalid';
  };

  if (isLoading || isLoadingProductInfo) {
    return (
      <Box display="flex" justifyContent="center" p={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        Failed to load license status. Please try again later.
      </Alert>
    );
  }

  if (productInfoError) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        Failed to load product information. Please try again later.
      </Alert>
    );
  }

  return (
    <Box>
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h5" component="h2" fontWeight={600} sx={{ color: 'primary.main' }}>
          License Management
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={handleOpenUpdateDialog}
          sx={{ fontWeight: 500 }}
        >
          Update License
        </Button>
      </Box>

      {/* Installation Info Card */}
      {productInfo && (
        <Card sx={{ mb: 3 }}>
          <CardContent sx={{ p: 3 }}>
            <Box sx={{ mb: 3 }}>
              <Typography variant="h6" gutterBottom fontWeight={600} sx={{ color: 'primary.light' }}>
                Installation Info
              </Typography>
            </Box>

            <Divider sx={{ my: 2 }} />

            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  Client ID
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>
                    {productInfo.client_id || 'N/A'}
                  </Typography>
                  {productInfo.client_id && (
                    <Tooltip title="Copy Client ID to clipboard">
                      <IconButton
                        size="small"
                        onClick={handleCopyClientId}
                        sx={{ color: 'primary.main' }}
                      >
                        <ContentCopyIcon fontSize="small" />
                      </IconButton>
                    </Tooltip>
                  )}
                </Box>
              </Grid>

              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  Created At
                </Typography>
                <Typography variant="body1">
                  {productInfo.created_at
                    ? formatDate(productInfo.created_at)
                    : 'N/A'}
                </Typography>
              </Grid>
            </Grid>
          </CardContent>
        </Card>
      )}

      {licenseStatus?.license && (
        <Card sx={{ mb: 3 }}>
          <CardContent sx={{ p: 3 }}>
            <Box sx={{ mb: 3 }}>
              <Typography variant="h6" gutterBottom fontWeight={600} sx={{ color: 'primary.light' }}>
                Current License Status
              </Typography>
              <Grid container spacing={2} alignItems="center">
                <Grid item>
                  <Chip
                    label={getStatusText(licenseStatus.license.is_valid || false)}
                    color={getStatusColor(licenseStatus.license.is_valid || false)}
                    variant="filled"
                    size="medium"
                  />
                </Grid>
                {licenseStatus.license.is_expired && (
                  <Grid item>
                    <Chip
                      label="Expired"
                      color="error"
                      variant="outlined"
                      size="medium"
                    />
                  </Grid>
                )}
              </Grid>
            </Box>

            <Divider sx={{ my: 2 }} />

            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  License ID
                </Typography>
                <Typography variant="body1" sx={{ fontFamily: 'monospace' }}>
                  {licenseStatus.license.id || 'N/A'}
                </Typography>
              </Grid>

              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  Type
                </Typography>
                <Typography variant="body1">
                  {licenseStatus.license.type
                    ? getLicenseTypeLabel(licenseStatus.license.type)
                    : 'N/A'}
                </Typography>
              </Grid>

              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  Issued At
                </Typography>
                <Typography variant="body1">
                  {licenseStatus.license.issued_at
                    ? formatDate(licenseStatus.license.issued_at)
                    : 'N/A'}
                </Typography>
              </Grid>

              <Grid item xs={12} md={6}>
                <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                  Expires At
                </Typography>
                <Typography variant="body1">
                  {licenseStatus.license.expires_at
                    ? formatDate(licenseStatus.license.expires_at)
                    : 'N/A'}
                </Typography>
              </Grid>

              {licenseStatus.license.days_until_expiry !== undefined && (
                <Grid item xs={12} md={6}>
                  <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                    Days Until Expiry
                  </Typography>
                  <Typography
                    variant="body1"
                    color={
                      licenseStatus.license.days_until_expiry <= 30
                        ? 'warning.main'
                        : 'text.primary'
                    }
                    fontWeight={
                      licenseStatus.license.days_until_expiry <= 30 ? 600 : 400
                    }
                  >
                    {licenseStatus.license.days_until_expiry} days
                    {licenseStatus.license.days_until_expiry <= 30 && (
                      <Chip
                        label="Expiring Soon"
                        color="warning"
                        size="small"
                        sx={{ ml: 1 }}
                      />
                    )}
                  </Typography>
                </Grid>
              )}
              
              {licenseStatus.license.license_text && (
                <Grid item xs={12}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
                    <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                      License Text
                    </Typography>
                    <Button
                      variant="outlined"
                      size="small"
                      startIcon={<ContentCopyIcon />}
                      onClick={handleCopyLicenseText}
                    >
                      Copy License Text
                    </Button>
                  </Box>
                  <TextField
                    fullWidth
                    multiline
                    rows={4}
                    value={licenseStatus.license.license_text}
                    InputProps={{
                      readOnly: true,
                      sx: { fontFamily: 'monospace', fontSize: '0.875rem' }
                    }}
                    variant="outlined"
                  />
                </Grid>
              )}
            </Grid>
          </CardContent>
        </Card>
      )}

      {/* Update License Dialog */}
      <Dialog
        open={openUpdateDialog}
        onClose={handleCloseUpdateDialog}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle sx={{ color: 'primary.main' }}>Update License</DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Enter your new license key to update the current license.
          </Typography>
          <TextField
            autoFocus
            margin="dense"
            label="License Key"
            type="text"
            fullWidth
            multiline
            rows={4}
            value={newLicense}
            onChange={(e) => setNewLicense(e.target.value)}
            placeholder="Paste your license key here..."
            variant="outlined"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseUpdateDialog} disabled={isSubmitting} size="small">
            Cancel
          </Button>
          <Button
            onClick={handleUpdateLicense}
            variant="contained"
            disabled={!newLicense.trim() || isSubmitting}
            startIcon={isSubmitting ? <CircularProgress size={16} /> : null}
            size="small"
          >
            {isSubmitting ? 'Updating...' : 'Update License'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default LicenseTab; 