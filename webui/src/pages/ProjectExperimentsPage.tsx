import React, { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import {
  Add as AddIcon,
} from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import apiClient from '../api/apiClient';
import type { Project } from '../generated/api/client';
import { useRBAC } from '../auth/permissions';

interface ProjectResponse { project: Project }

const ProjectExperimentsPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const [environmentKey, setEnvironmentKey] = useState<string>(() => {
    // Try to get from localStorage first, fallback to 'prod'
    return localStorage.getItem('currentEnvironmentKey') || 'prod';
  });
  
  // RBAC checks for current project
  const rbac = useRBAC(projectId);

  // Check project access
  if (!rbac.canManageProject()) {
    return (
      <AuthenticatedLayout showBackButton backTo="/dashboard">
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="error" gutterBottom>
            Access Denied
          </Typography>
          <Typography variant="body2" color="text.secondary">
            You don't have permission to view this project.
          </Typography>
        </Box>
      </AuthenticatedLayout>
    );
  }

  const { data: projectResp, isLoading: loadingProject } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const res = await apiClient.getProject(projectId);
      return res.data as ProjectResponse;
    },
    enabled: !!projectId,
  });

  // Get environments for the project
  const { data: environmentsResp, isLoading: loadingEnvironments } = useQuery({
    queryKey: ['project-environments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectEnvironments(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  const environments = environmentsResp?.items ?? [];

  // Initialize environment ID in localStorage when environments are loaded
  React.useEffect(() => {
    if (environments.length > 0 && environmentKey) {
      const selectedEnv = environments.find(env => env.key === environmentKey);
      if (selectedEnv) {
        localStorage.setItem('currentEnvId', selectedEnv.id.toString());
        console.log('[ProjectExperimentsPage] Initialized environment in localStorage:', { id: selectedEnv.id, key: selectedEnv.key });
      }
    }
  }, [environments, environmentKey]);

  const project = projectResp?.project;

  return (
    <AuthenticatedLayout showBackButton backTo="/dashboard">
      <Paper sx={{ p: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" sx={{ color: 'primary.light' }}>Experiments</Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            {/* Environment selector */}
            <FormControl size="small" sx={{ minWidth: 200 }}>
              <InputLabel>Environment</InputLabel>
              <Select
                value={environmentKey}
                label="Environment"
                size="small"
                onChange={(e) => {
                  setEnvironmentKey(e.target.value);
                  // Find the environment ID and save it to localStorage
                  const selectedEnv = environments.find(env => env.key === e.target.value);
                  if (selectedEnv) {
                    localStorage.setItem('currentEnvId', selectedEnv.id.toString());
                    localStorage.setItem('currentEnvironmentKey', selectedEnv.key);
                    console.log('[ProjectExperimentsPage] Saved environment to localStorage:', { id: selectedEnv.id, key: selectedEnv.key });
                  }
                }}
                disabled={loadingEnvironments}
              >
                {environments.map((env) => (
                  <MenuItem key={env.id} value={env.key} data-env-id={env.id}>
                    {env.name} ({env.key})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
            {rbac.canManageFeature() && (
              <Button variant="contained" startIcon={<AddIcon />} size="small">
                Add Experiment
              </Button>
            )}
          </Box>
        </Box>

        {(loadingProject || loadingEnvironments) && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}

        {!loadingProject && !loadingEnvironments && (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography variant="body1" color="text.secondary">
              No experiments yet. API methods for experiments will be implemented soon.
            </Typography>
          </Box>
        )}
      </Paper>
    </AuthenticatedLayout>
  );
};

export default ProjectExperimentsPage;
