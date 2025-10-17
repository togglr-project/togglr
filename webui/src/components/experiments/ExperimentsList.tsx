import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  Switch,
  IconButton,
  Tooltip,
  CircularProgress,
} from '@mui/material';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { FeatureAlgorithm } from '../../generated/api/client';

interface ExperimentsListProps {
  projectId: string;
  environmentKey: string;
  onEdit: (experiment: FeatureAlgorithm) => void;
  onDelete: (experiment: FeatureAlgorithm) => void;
  onToggle: (experiment: FeatureAlgorithm) => void;
  onView: (experiment: FeatureAlgorithm) => void;
  togglingExperimentId?: string | null;
}

const ExperimentsList: React.FC<ExperimentsListProps> = ({
  projectId,
  environmentKey,
  onEdit,
  onDelete,
  onToggle,
  onView,
  togglingExperimentId,
}) => {
  const { data: experimentsResp, isLoading, error } = useQuery({
    queryKey: ['feature-algorithms', projectId, environmentKey],
    queryFn: async () => {
      const res = await apiClient.listFeatureAlgorithms(projectId, environmentKey);
      return res.data;
    },
    enabled: !!projectId && !!environmentKey,
  });

  const experiments = experimentsResp?.feature_algorithms || [];

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Typography color="error">
        Failed to load experiments.
      </Typography>
    );
  }

  if (experiments.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 4 }}>
        <Typography variant="body1" color="text.secondary">
          No experiments yet.
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
      {experiments.map((experiment) => (
        <Card 
          key={experiment.id} 
          sx={{ 
            minHeight: 80,
            cursor: 'pointer',
            '&:hover': {
              boxShadow: 2,
            },
          }}
          onClick={() => onView(experiment)}
        >
          <CardContent sx={{ p: 2, '&:last-child': { pb: 2 } }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
              <Box sx={{ flexGrow: 1, minWidth: 0 }}>
                <Typography 
                  variant="h6" 
                  sx={{ 
                    fontWeight: 600,
                    fontSize: '1rem',
                    mb: 0.5,
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                  }}
                  title={experiment.feature.name}
                >
                  {experiment.feature.name}
                </Typography>
                <Typography 
                  variant="body2" 
                  color="text.secondary"
                  sx={{ 
                    fontSize: '0.8rem',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                  }}
                  title={experiment.feature.key}
                >
                  {experiment.feature.key}
                </Typography>
              </Box>

              <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, overflow: 'hidden' }}>
                <Chip
                  size="small"
                  label={experiment.feature.enabled ? 'enabled' : 'disabled'}
                  color={experiment.feature.enabled ? 'success' : 'default'}
                  sx={{ 
                    fontSize: '0.7rem',
                    height: 20,
                  }}
                />
                
                <Box sx={{ width: 1, height: 12, bgcolor: 'divider', opacity: 0.5 }} />
                
                <Chip
                  size="small"
                  label={`Algorithm: ${experiment.algorithm_slug}`}
                  color="info"
                  variant="outlined"
                  sx={{ 
                    fontSize: '0.7rem',
                    height: 20,
                  }}
                />
              </Box>

              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Box 
                  sx={{ position: 'relative' }}
                  onClick={(e) => {
                    e.stopPropagation();
                  }}
                >
                  <Switch
                    size="small"
                    checked={experiment.enabled}
                    onChange={(e) => {
                      e.stopPropagation();
                      onToggle(experiment);
                    }}
                    onClick={(e) => {
                      e.stopPropagation();
                    }}
                    disabled={togglingExperimentId === experiment.id}
                  />
                  {togglingExperimentId === experiment.id && (
                    <CircularProgress
                      size={16}
                      sx={{
                        position: 'absolute',
                        top: '50%',
                        left: '50%',
                        marginTop: '-8px',
                        marginLeft: '-8px',
                      }}
                    />
                  )}
                </Box>
                
                <Box sx={{ display: 'flex', gap: 0.5, ml: 1 }}>
                  <Tooltip title="Edit experiment">
                    <IconButton 
                      size="small" 
                      onClick={(e) => {
                        e.stopPropagation();
                        onEdit(experiment);
                      }}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Delete experiment">
                    <IconButton 
                      size="small" 
                      onClick={(e) => {
                        e.stopPropagation();
                        onDelete(experiment);
                      }}
                      sx={{ color: 'error.main' }}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                </Box>
              </Box>
            </Box>
          </CardContent>
        </Card>
      ))}
    </Box>
  );
};

export default ExperimentsList;
