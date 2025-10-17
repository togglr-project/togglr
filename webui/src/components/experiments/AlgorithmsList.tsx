import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  CircularProgress,
  Grid,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { Algorithm } from '../../generated/api/client';

const AlgorithmsList: React.FC = () => {
  const { data: algorithmsResp, isLoading, error } = useQuery({
    queryKey: ['algorithms'],
    queryFn: async () => {
      const res = await apiClient.listAlgorithms();
      return res.data;
    },
  });

  const algorithms = algorithmsResp?.algorithms || [];

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
        Failed to load algorithms.
      </Typography>
    );
  }

  if (algorithms.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 4 }}>
        <Typography variant="body1" color="text.secondary">
          No algorithms available.
        </Typography>
      </Box>
    );
  }

  return (
    <Grid container spacing={2}>
      {algorithms.map((algorithm) => (
        <Grid item xs={12} md={6} lg={4} key={algorithm.slug}>
          <Card sx={{ height: '100%' }}>
            <CardContent>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  {algorithm.name}
                </Typography>
                <Chip 
                  label={algorithm.kind} 
                  size="small" 
                  color="primary"
                  sx={{ textTransform: 'capitalize' }}
                />
              </Box>
              
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                {algorithm.description}
              </Typography>
              
              <Typography variant="subtitle2" sx={{ mb: 1, fontWeight: 600 }}>
                Default Settings:
              </Typography>
              
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                {Object.entries(algorithm.default_settings).map(([key, value]) => (
                  <Chip
                    key={key}
                    label={`${key}: ${value}`}
                    size="small"
                    variant="outlined"
                    sx={{ fontSize: '0.7rem' }}
                  />
                ))}
              </Box>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
};

export default AlgorithmsList;
