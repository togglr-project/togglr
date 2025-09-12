import React, { useState } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Alert,
  Container,
  Paper,
  CircularProgress,
} from '@mui/material';
import { useLicense } from '../../auth/LicenseContext';
import apiClient from '../../api/apiClient';
import LogoImg from '../LogoImg';

const LicenseExpiredPage: React.FC = () => {
  const [newLicense, setNewLicense] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitSuccess, setSubmitSuccess] = useState<string | null>(null);
  const { licenseStatus, checkLicenseStatus } = useLicense();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newLicense.trim()) {
      setSubmitError('Please enter a license key');
      return;
    }

    try {
      setIsSubmitting(true);
      setSubmitError(null);
      setSubmitSuccess(null);

      // Submit new license
      await apiClient.updateLicense({
        license_text: newLicense,
      });

      setSubmitSuccess('License updated successfully');
      setNewLicense('');
      
      // Refresh license status
      await checkLicenseStatus();
    } catch (err: any) {
      console.error('Failed to update license:', err);
      setSubmitError(err.response?.data?.message || 'Failed to update license');
    } finally {
      setIsSubmitting(false);
    }
  };

  const formatDate = (dateString: string) => {
    try {
      return new Date(dateString).toLocaleDateString();
    } catch {
      return dateString;
    }
  };

  const getLicenseTypeLabel = (type: string) => {
    switch (type.toLowerCase()) {
      case 'trial':
        return 'Trial';
      case 'commercial':
        return 'Commercial';
      default:
        return type;
    }
  };

  return (
    <Container maxWidth="md">
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: '100vh',
          py: 4,
        }}
      >
        <Box sx={{ mb: 4 }}>
          <LogoImg size="large" />
        </Box>

        <Paper
          elevation={3}
          sx={{
            p: 4,
            width: '100%',
            maxWidth: 600,
          }}
        >
          <Typography variant="h4" component="h1" gutterBottom align="center" color="error">
            License Required
          </Typography>

          <Alert severity="error" sx={{ mb: 3 }}>
            Your license has expired or is invalid. Please enter a valid license key to continue using the application.
          </Alert>

          {licenseStatus?.license && (
            <Card variant="outlined" sx={{ mb: 3 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Current License Information
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  <strong>Type:</strong> {getLicenseTypeLabel(licenseStatus.license.type || 'Unknown')}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  <strong>Issued:</strong> {licenseStatus.license.issued_at ? formatDate(licenseStatus.license.issued_at) : 'Unknown'}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  <strong>Expires:</strong> {licenseStatus.license.expires_at ? formatDate(licenseStatus.license.expires_at) : 'Unknown'}
                </Typography>
                {licenseStatus.license.days_until_expiry !== undefined && (
                  <Typography variant="body2" color="text.secondary">
                    <strong>Days until expiry:</strong> {licenseStatus.license.days_until_expiry}
                  </Typography>
                )}
              </CardContent>
            </Card>
          )}

          <Box component="form" onSubmit={handleSubmit}>
            <Typography variant="h6" gutterBottom>
              Enter New License Key
            </Typography>
            
            {submitError && (
              <Alert severity="error" sx={{ mb: 2 }}>
                {submitError}
              </Alert>
            )}

            {submitSuccess && (
              <Alert severity="success" sx={{ mb: 2 }}>
                {submitSuccess}
              </Alert>
            )}

            <TextField
              fullWidth
              multiline
              rows={4}
              label="License Key"
              value={newLicense}
              onChange={(e) => setNewLicense(e.target.value)}
              margin="normal"
              variant="outlined"
              placeholder="Paste your license key here..."
              disabled={isSubmitting}
            />

            <Button
              type="submit"
              variant="contained"
              color="primary"
              fullWidth
              size="large"
              disabled={isSubmitting || !newLicense.trim()}
              sx={{ mt: 2 }}
            >
              {isSubmitting ? (
                <>
                  <CircularProgress size={20} sx={{ mr: 1 }} />
                  Updating License...
                </>
              ) : (
                'Update License'
              )}
            </Button>
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default LicenseExpiredPage; 