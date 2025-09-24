import React from 'react';
import { useLicense } from '../../auth/LicenseContext';
import { useAuth } from '../../auth/AuthContext';
import LicenseExpiredPage from './LicenseExpiredPage';
import { Box, CircularProgress, Typography } from '@mui/material';

interface LicenseGuardProps {
  children: React.ReactNode;
}

const LicenseGuard: React.FC<LicenseGuardProps> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const { isLicenseValid, isLoading } = useLicense();

  // If not authenticated, don't check license
  if (!isAuthenticated) {
    return <>{children}</>;
  }

  // Show loading while checking license
  if (isLoading) {
    return (
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: '100vh',
        }}
      >
        <CircularProgress size={60} />
        <Typography variant="h6" sx={{ mt: 2 }}>
          Checking license status...
        </Typography>
      </Box>
    );
  }

  // If license is not valid, show license expired page
  if (!isLicenseValid) {
    return <LicenseExpiredPage />;
  }

  // License is valid, show children
  return <>{children}</>;
};

export default LicenseGuard; 