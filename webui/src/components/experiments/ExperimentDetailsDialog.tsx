import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography,
} from '@mui/material';
import type { FeatureAlgorithm } from '../../generated/api/client';

interface ExperimentDetailsDialogProps {
  open: boolean;
  onClose: () => void;
  experiment: FeatureAlgorithm | null;
}

const ExperimentDetailsDialog: React.FC<ExperimentDetailsDialogProps> = ({
  open,
  onClose,
  experiment,
}) => {
  if (!experiment) return null;

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle>
        Experiment {experiment.feature.name}: {experiment.algorithm_slug}
      </DialogTitle>
      <DialogContent>
        <Box sx={{ py: 2 }}>
          <Typography variant="body1" color="text.secondary">
            Experiment details will be implemented here.
          </Typography>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};

export default ExperimentDetailsDialog;
