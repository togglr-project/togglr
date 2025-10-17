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
import type { FeatureAlgorithm, UpdateFeatureAlgorithmRequest, AuthCredentialsMethodEnum } from '../../generated/api/client';
import { useAuth } from '../../auth/AuthContext';
import GuardResponseHandler from '../pending-changes/GuardResponseHandler';
import { useApprovePendingChange } from '../../hooks/usePendingChanges';

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
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [enabled, setEnabled] = useState(false);
  const [settings, setSettings] = useState<{ [key: string]: number }>({});
  const [error, setError] = useState<string | null>(null);
  
  // Guard workflow state
  const [guardResponse, setGuardResponse] = useState<{
    pendingChange?: any;
    conflictError?: string;
    forbiddenError?: string;
  }>({});
  
  const approveMutation = useApprovePendingChange();

  useEffect(() => {
    if (experiment) {
      setEnabled(experiment.enabled);
      setSettings(experiment.settings);
      setError(null);
    }
  }, [experiment]);

  // Handle auto-approve for single-user projects
  const handleAutoApprove = (authMethod: AuthCredentialsMethodEnum, credential: string, sessionId?: string) => {
    if (!guardResponse.pendingChange?.id || !user) return;
    
    approveMutation.mutate(
      {
        id: guardResponse.pendingChange.id,
        request: {
          approver_user_id: user.id,
          approver_name: user.username,
          auth: {
            method: authMethod,
            credential,
            ...(sessionId && { session_id: sessionId }),
          },
        },
      },
      {
        onSuccess: () => {
          setGuardResponse({});
          queryClient.invalidateQueries({ queryKey: ['feature-algorithms'] });
          onClose();
        },
      }
    );
  };

  const updateMutation = useMutation({
    mutationFn: async (data: UpdateFeatureAlgorithmRequest) => {
      if (!experiment) throw new Error('No experiment selected');
      const environmentId = parseInt(localStorage.getItem('currentEnvId') || '0');
      const res = await apiClient.updateFeatureAlgorithm(experiment.feature_id, environmentId, data);
      return { data: res.data, status: res.status };
    },
    onSuccess: (result) => {
      if (result.status === 202) {
        // Pending change created - handle as guard workflow
        setGuardResponse({ pendingChange: result.data });
        return;
      }
      if (result.status === 409) {
        // Conflict - feature locked by another pending change
        setGuardResponse({ conflictError: 'Feature is already locked by another pending change' });
        return;
      }
      if (result.status === 403) {
        // Forbidden - user doesn't have permission to modify guarded feature
        setGuardResponse({ forbiddenError: 'You don\'t have permission to modify this guarded feature' });
        return;
      }
      // Normal success - update applied immediately
      queryClient.invalidateQueries({ queryKey: ['feature-algorithms'] });
      onClose();
    },
    onError: (error: any) => {
      // Handle guard workflow responses (pending change or conflict)
      if (error.response?.status === 202) {
        setGuardResponse({ pendingChange: error.response.data });
        return;
      }
      if (error.response?.status === 409) {
        setGuardResponse({ conflictError: error.response.data.message || 'Feature is already locked by another pending change' });
        return;
      }
      if (error.response?.status === 403) {
        setGuardResponse({ forbiddenError: 'You don\'t have permission to modify this guarded feature' });
        return;
      }
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

      {/* Guard Response Handler */}
      <GuardResponseHandler
        pendingChange={guardResponse.pendingChange}
        conflictError={guardResponse.conflictError}
        forbiddenError={guardResponse.forbiddenError}
        onClose={() => setGuardResponse({})}
        onParentClose={onClose}
        onApprove={handleAutoApprove}
        approveLoading={approveMutation.isPending}
      />
    </Dialog>
  );
};

export default EditExperimentDialog;
