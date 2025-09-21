import React from 'react';
import { Box, Typography } from '@mui/material';
import { Dashboard as DashboardIcon } from '@mui/icons-material';
import { Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import { useAuth } from '../auth/AuthContext';

const DashboardPage: React.FC = () => {
  const { isAuthenticated } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return (
    <AuthenticatedLayout>
      <PageHeader
        title="Dashboard"
        subtitle=""
        icon={<DashboardIcon />}
      />
      <Box sx={{ p: 2 }}>
        <Typography variant="body2" color="text.secondary">
          {/* Empty dashboard for now */}
        </Typography>
      </Box>
    </AuthenticatedLayout>
  );
};

export default DashboardPage;
