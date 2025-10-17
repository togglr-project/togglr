import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Switch,
  FormControlLabel,
  Box,
  Typography,
  CircularProgress,
  Chip,
  Alert,
} from '@mui/material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { FeatureExtended, Algorithm, CreateFeatureAlgorithmRequest } from '../../generated/api/client';
import SearchPanel from '../SearchPanel';

interface CreateExperimentDialogProps {
  open: boolean;
  onClose: () => void;
  projectId: string;
  environmentKey: string;
}

const getDefaultSettings = (algorithmSlug: string): { [key: string]: number } => {
  switch (algorithmSlug) {
    case 'epsilon-greedy':
      return { epsilon: 0.1 };
    case 'thompson-sampling':
      return { prior_beta: 1, prior_alpha: 1 };
    case 'ucb':
      return { confidence: 2.0 };
    default:
      return {};
  }
};

const CreateExperimentDialog: React.FC<CreateExperimentDialogProps> = ({
  open,
  onClose,
  projectId,
  environmentKey,
}) => {
  const queryClient = useQueryClient();
  const [selectedFeature, setSelectedFeature] = useState<FeatureExtended | null>(null);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState<string>('');
  const [enabled, setEnabled] = useState(true);
  const [settings, setSettings] = useState<{ [key: string]: number }>({});
  const [searchValue, setSearchValue] = useState('');
  const [error, setError] = useState<string | null>(null);

  const { data: featuresResp } = useQuery({
    queryKey: ['project-features', projectId, environmentKey],
    queryFn: async () => {
      const res = await apiClient.listProjectFeatures(projectId, environmentKey);
      return res.data;
    },
    enabled: open && !!projectId,
  });

  const { data: algorithmsResp } = useQuery({
    queryKey: ['algorithms'],
    queryFn: async () => {
      const res = await apiClient.listAlgorithms();
      return res.data;
    },
    enabled: open,
  });

  const features = featuresResp?.items || [];
  const algorithms = algorithmsResp?.algorithms || [];

  const filteredFeatures = features.filter(feature =>
    feature.name.toLowerCase().includes(searchValue.toLowerCase()) ||
    feature.key.toLowerCase().includes(searchValue.toLowerCase())
  );

  const selectedAlgorithmData = algorithms.find(alg => alg.slug === selectedAlgorithm);

  useEffect(() => {
    if (selectedAlgorithm) {
      const defaultSettings = getDefaultSettings(selectedAlgorithm);
      setSettings(defaultSettings);
    }
  }, [selectedAlgorithm]);

  const createMutation = useMutation({
    mutationFn: async (data: CreateFeatureAlgorithmRequest) => {
      if (!selectedFeature) throw new Error('No feature selected');
      const environmentId = parseInt(localStorage.getItem('currentEnvId') || '0');
      await apiClient.createFeatureAlgorithm(selectedFeature.id, environmentId, data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feature-algorithms'] });
      onClose();
      resetForm();
    },
    onError: (error: any) => {
      if (error.response?.status === 409) {
        setError('An algorithm already exists for this feature. Please select a different feature or edit the existing experiment.');
      } else {
        setError('Failed to create experiment. Please try again.');
      }
    },
  });

  const resetForm = () => {
    setSelectedFeature(null);
    setSelectedAlgorithm('');
    setEnabled(true);
    setSettings({});
    setSearchValue('');
    setError(null);
  };

  const handleClose = () => {
    resetForm();
    onClose();
  };

  const handleSubmit = () => {
    if (!selectedFeature || !selectedAlgorithm) return;
    
    createMutation.mutate({
      algorithm_slug: selectedAlgorithm,
      environment_id: parseInt(localStorage.getItem('currentEnvId') || '0'),
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

  return (
    <Dialog open={open} onClose={handleClose} fullWidth maxWidth="md">
      <DialogTitle>Create Experiment</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, mt: 1 }}>
          {error && (
            <Alert severity="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}
          
          <Box>
            <Typography variant="subtitle1" sx={{ mb: 2 }}>
              Select Feature
            </Typography>
            <SearchPanel
              searchValue={searchValue}
              onSearchChange={setSearchValue}
              placeholder="Search features..."
              projectId={projectId}
              showTagFilter={false}
            />
            <Box sx={{ maxHeight: 200, overflow: 'auto', mt: 1 }}>
              {filteredFeatures.map((feature) => (
                <Box
                  key={feature.id}
                  onClick={() => {
                    setSelectedFeature(feature);
                    setError(null);
                  }}
                  sx={{
                    p: 2,
                    border: '1px solid',
                    borderColor: selectedFeature?.id === feature.id ? 'primary.main' : 'divider',
                    borderRadius: 1,
                    mb: 1,
                    cursor: 'pointer',
                    bgcolor: selectedFeature?.id === feature.id ? 'primary.50' : 'background.paper',
                    '&:hover': {
                      bgcolor: selectedFeature?.id === feature.id ? 'primary.50' : 'action.hover',
                    },
                  }}
                >
                  <Typography variant="subtitle2">{feature.name}</Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ fontFamily: 'monospace' }}>
                    {feature.key}
                  </Typography>
                  <Box sx={{ display: 'flex', gap: 1, mt: 1 }}>
                    <Chip
                      size="small"
                      label={feature.kind}
                      color="primary"
                      variant="outlined"
                    />
                    <Chip
                      size="small"
                      label={feature.enabled ? 'enabled' : 'disabled'}
                      color={feature.enabled ? 'success' : 'default'}
                    />
                  </Box>
                </Box>
              ))}
            </Box>
          </Box>

          <Box>
            <Typography variant="subtitle1" sx={{ mb: 2 }}>
              Select Algorithm
            </Typography>
            <FormControl fullWidth>
              <InputLabel>Algorithm</InputLabel>
              <Select
                value={selectedAlgorithm}
                label="Algorithm"
                onChange={(e) => {
                  setSelectedAlgorithm(e.target.value);
                  setError(null);
                }}
              >
                {algorithms.map((algorithm) => (
                  <MenuItem key={algorithm.slug} value={algorithm.slug}>
                    <Box>
                      <Typography variant="body1">{algorithm.name}</Typography>
                      <Typography variant="body2" color="text.secondary">
                        {algorithm.description}
                      </Typography>
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
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

          {selectedAlgorithmData && Object.keys(settings).length > 0 && (
            <Box>
              <Typography variant="subtitle1" sx={{ mb: 2 }}>
                Algorithm Settings
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                {Object.entries(selectedAlgorithmData.default_settings).map(([key, defaultValue]) => (
                  <TextField
                    key={key}
                    label={key}
                    type="number"
                    value={settings[key] || defaultValue}
                    onChange={(e) => handleSettingsChange(key, parseFloat(e.target.value) || 0)}
                    size="small"
                    fullWidth
                  />
                ))}
              </Box>
            </Box>
          )}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={!selectedFeature || !selectedAlgorithm || createMutation.isPending}
        >
          {createMutation.isPending ? <CircularProgress size={20} /> : 'Create'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateExperimentDialog;
