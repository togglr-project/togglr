import React from 'react';
import { Box, Typography } from '@mui/material';

const LicenseTab: React.FC = () => {
  return (
    <Box>
      <Typography variant="h5" component="h2" fontWeight={600} sx={{ color: 'primary.main', mb: 3 }}>
        License Management
      </Typography>
      <Typography variant="body1" color="text.secondary">
        License management is not available in the open source version.
      </Typography>
    </Box>
  );
};

export default LicenseTab;
