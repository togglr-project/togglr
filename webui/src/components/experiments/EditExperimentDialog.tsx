import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Switch,
  FormControlLabel,
  Box,
  Typography,
  CircularProgress,
  Alert,
} from '@mui/material';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { FeatureAlgorithm, UpdateFeatureAlgorithmRequest } from '../../generated/api/client';

interface EditExperimentDialogProps {
  open: boolean;
  onClose: () => void;
  experiment: FeatureAlgorithm | null;
}

const EditExperimentDialog: React.FC<EditExperimentDialogProps> = ({
  open,
  onClose,
  experiment,
}) => {
  const queryClient = useQueryClient();
  const [enabled, setEnabled] = useState(false);
  const [settings, setSettings] = useState<{ [key: string]: number }>({});
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (experiment) {
      setEnabled(experiment.enabled);
      setSettings(experiment.settings);
      setError(null);
    }
  }, [experiment]);

  const updateMutation = useMutation({
    mutationFn: async (data: UpdateFeatureAlgorithmRequest) => {
      if (!experiment) throw new Error('No experiment selected');
      const environmentId = parseInt(localStorage.getItem('currentEnvId') || '0');
      await apiClient.updateFeatureAlgorithm(experiment.feature_id, environmentId, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feature-algorithms'] });
      onClose();
    },
    onError: (error: any) => {
      setError('Failed to update experiment. Please try again.');
    },
  });

  const handleClose = () => {
    onClose();
  };

  const handleSubmit = () => {
    if (!experiment) return;
    
    updateMutation.mutate({
      enabled,
      settings,
    });
  };

  const handleSettingsChange = (key: string, value: number) => {
    setSettings(prev => ({
      ...prev,
      [key]: value
    }));
  };

  if (!experiment) return null;

  return (
    <Dialog open={open} onClose={handleClose} fullWidth maxWidth="sm">
      <DialogTitle>Edit Experiment</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, mt: 1 }}>
          {error && (
            <Alert severity="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}
          
          <Box>
            <Typography variant="subtitle2" color="text.secondary">
              Feature
            </Typography>
            <Typography variant="h6">{experiment.feature.name}</Typography>
            <Typography variant="body2" color="text.secondary" sx={{ fontFamily: 'monospace' }}>
              {experiment.feature.key}
            </Typography>
          </Box>

          <Box>
            <Typography variant="subtitle2" color="text.secondary">
              Algorithm
            </Typography>
            <Typography variant="body1">{experiment.algorithm_slug}</Typography>
          </Box>

          <Box>
            <FormControlLabel
              control={
                <Switch
                  checked={enabled}
                  onChange={(e) => setEnabled(e.target.checked)}
                />
              }
              label="Enabled"
            />
          </Box>

          <Box>
            <Typography variant="subtitle1" sx={{ mb: 2 }}>
              Algorithm Settings
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              {Object.entries(settings).map(([key, value]) => (
                <TextField
                  key={key}
                  label={key}
                  type="number"
                  value={value}
                  onChange={(e) => handleSettingsChange(key, parseFloat(e.target.value) || 0)}
                  size="small"
                  fullWidth
                />
              ))}
            </Box>
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={updateMutation.isPending}
        >
          {updateMutation.isPending ? <CircularProgress size={20} /> : 'Update'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default EditExperimentDialog;
