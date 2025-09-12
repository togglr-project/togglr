import React from 'react';
import { Box, LinearProgress, Typography, Paper, useTheme, Fade } from '@mui/material';

interface LDAPSyncProgressBarProps {
  isRunning: boolean;
  progress: number; // 0..100 (percentage)
  currentStep?: string;
  processedItems?: number;
  totalItems?: number;
  estimatedTime?: string;
  startTime?: string;
}

const LDAPSyncProgressBar: React.FC<LDAPSyncProgressBarProps> = ({
  isRunning,
  progress,
  currentStep,
  processedItems,
  totalItems,
  estimatedTime,
  startTime,
}) => {
  const theme = useTheme();

  if (!isRunning) return null;

  const percent = Math.round(progress || 0);

  return (
    <Fade in={isRunning} unmountOnExit>
      <Paper elevation={3} sx={{
        p: 2,
        mb: 2,
        background: theme.palette.background.paper,
        border: `1.5px solid ${theme.palette.primary.main}`,
        maxWidth: 480,
        mx: 'auto',
      }}>
        <Typography variant="subtitle1" fontWeight={600} gutterBottom>
          LDAP Users and Groups Synchronization
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Box sx={{ flex: 1 }}>
            <LinearProgress 
              variant="determinate" 
              value={percent} 
              sx={{ height: 10, borderRadius: 5, background: theme.palette.action.hover }}
              color="primary"
            />
          </Box>
          <Typography variant="body2" color="text.secondary" minWidth={40}>
            {percent}%
          </Typography>
        </Box>
        <Box sx={{ mt: 1, display: 'flex', flexDirection: 'column', gap: 0.5 }}>
          {currentStep && (
            <Typography variant="body2" color="text.primary">
              Step: {currentStep}
            </Typography>
          )}
          {(typeof processedItems === 'number' && typeof totalItems === 'number') && (
            <Typography variant="body2" color="text.secondary">
              Processed: {processedItems} of {totalItems}
            </Typography>
          )}
          {estimatedTime && (
            <Typography variant="body2" color="text.secondary">
              Estimated time remaining: {estimatedTime}
            </Typography>
          )}
          {startTime && (
            <Typography variant="caption" color="text.disabled">
              Started: {new Date(startTime).toLocaleString()}
            </Typography>
          )}
        </Box>
      </Paper>
    </Fade>
  );
};

export default LDAPSyncProgressBar; 