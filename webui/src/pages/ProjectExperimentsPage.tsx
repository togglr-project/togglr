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
  Tabs,
  Tab,
} from '@mui/material';
import {
  Add as AddIcon,
} from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import apiClient from '../api/apiClient';
import type { Project, AuthCredentialsMethodEnum } from '../generated/api/client';
import { useRBAC } from '../auth/permissions';
import { useAuth } from '../auth/AuthContext';
import GuardResponseHandler from '../components/pending-changes/GuardResponseHandler';
import { useApprovePendingChange } from '../hooks/usePendingChanges';
import { 
  ExperimentsList, 
  AlgorithmsList, 
  CreateExperimentDialog, 
  EditExperimentDialog,
  ExperimentDetailsDialog
} from '../components/experiments';
import type { FeatureAlgorithm } from '../generated/api/client';

interface ProjectResponse { project: Project }

const ProjectExperimentsPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [environmentKey, setEnvironmentKey] = useState<string>(() => {
    return localStorage.getItem('currentEnvironmentKey') || 'prod';
  });
  const [activeTab, setActiveTab] = useState(0);
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [detailsDialogOpen, setDetailsDialogOpen] = useState(false);
  const [selectedExperiment, setSelectedExperiment] = useState<FeatureAlgorithm | null>(null);
  const [togglingExperimentId, setTogglingExperimentId] = useState<string | null>(null);
  
  // Guard workflow state
  const [guardResponse, setGuardResponse] = useState<{
    pendingChange?: any;
    conflictError?: string;
    forbiddenError?: string;
  }>({});
  
  const rbac = useRBAC(projectId);
  const approveMutation = useApprovePendingChange();

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

  // Check project access after all hooks
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
          queryClient.invalidateQueries({ queryKey: ['feature-algorithms', projectId, environmentKey] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
        },
      }
    );
  };

  const handleEdit = (experiment: FeatureAlgorithm) => {
    setSelectedExperiment(experiment);
    setEditDialogOpen(true);
  };

  const handleView = (experiment: FeatureAlgorithm) => {
    setSelectedExperiment(experiment);
    setDetailsDialogOpen(true);
  };

  const handleDelete = async (experiment: FeatureAlgorithm) => {
    if (window.confirm('Are you sure you want to delete this experiment?')) {
      try {
        const environmentId = parseInt(localStorage.getItem('currentEnvId') || '0');
        const res = await apiClient.deleteFeatureAlgorithm(experiment.feature_id, environmentId);
        
        // Handle guard workflow responses
        if (res.status === 202) {
          setGuardResponse({ pendingChange: res.data });
          return;
        }
        if (res.status === 409) {
          setGuardResponse({ conflictError: 'Feature is already locked by another pending change' });
          return;
        }
        
        // Normal success - invalidate queries
        queryClient.invalidateQueries({ queryKey: ['feature-algorithms', projectId, environmentKey] });
      } catch (error: any) {
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
        console.error('Failed to delete experiment:', error);
      }
    }
  };

  const handleToggle = async (experiment: FeatureAlgorithm) => {
    setTogglingExperimentId(experiment.id);
    try {
      const environmentId = parseInt(localStorage.getItem('currentEnvId') || '0');
      const res = await apiClient.updateFeatureAlgorithm(experiment.feature_id, environmentId, {
        enabled: !experiment.enabled,
        settings: experiment.settings,
      });
      
      // Handle guard workflow responses
      if (res.status === 202) {
        setGuardResponse({ pendingChange: res.data });
        return;
      }
      if (res.status === 409) {
        setGuardResponse({ conflictError: 'Feature is already locked by another pending change' });
        return;
      }
      
      // Normal success - invalidate queries
      queryClient.invalidateQueries({ queryKey: ['feature-algorithms', projectId, environmentKey] });
    } catch (error: any) {
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
      console.error('Failed to toggle experiment:', error);
    } finally {
      setTogglingExperimentId(null);
    }
  };

  return (
    <AuthenticatedLayout showBackButton backTo="/dashboard">
      <Paper sx={{ p: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" sx={{ color: 'primary.light' }}>Experiments</Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <FormControl size="small" sx={{ minWidth: 200 }}>
              <InputLabel>Environment</InputLabel>
              <Select
                value={environmentKey}
                label="Environment"
                size="small"
                onChange={(e) => {
                  setEnvironmentKey(e.target.value);
                  const selectedEnv = environments.find(env => env.key === e.target.value);
                  if (selectedEnv) {
                    localStorage.setItem('currentEnvId', selectedEnv.id.toString());
                    localStorage.setItem('currentEnvironmentKey', selectedEnv.key);
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
              <Button 
                variant="contained" 
                startIcon={<AddIcon />} 
                size="small"
                onClick={() => setCreateDialogOpen(true)}
              >
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
          <>
            <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
              <Tabs value={activeTab} onChange={(e, newValue) => setActiveTab(newValue)}>
                <Tab label="Experiments" />
                <Tab label="Algorithms" />
              </Tabs>
            </Box>
            
            {activeTab === 0 && (
              <ExperimentsList 
                projectId={projectId}
                environmentKey={environmentKey}
                onEdit={handleEdit}
                onDelete={handleDelete}
                onToggle={handleToggle}
                onView={handleView}
                togglingExperimentId={togglingExperimentId}
              />
            )}
            
            {activeTab === 1 && (
              <AlgorithmsList />
            )}
          </>
        )}
      </Paper>

      <CreateExperimentDialog
        open={createDialogOpen}
        onClose={() => setCreateDialogOpen(false)}
        projectId={projectId}
        environmentKey={environmentKey}
      />

      <EditExperimentDialog
        open={editDialogOpen}
        onClose={() => setEditDialogOpen(false)}
        experiment={selectedExperiment}
      />

      <ExperimentDetailsDialog
        open={detailsDialogOpen}
        onClose={() => setDetailsDialogOpen(false)}
        experiment={selectedExperiment}
      />

      {/* Guard Response Handler */}
      <GuardResponseHandler
        pendingChange={guardResponse.pendingChange}
        conflictError={guardResponse.conflictError}
        forbiddenError={guardResponse.forbiddenError}
        onClose={() => setGuardResponse({})}
        onApprove={handleAutoApprove}
        approveLoading={approveMutation.isPending}
      />
    </AuthenticatedLayout>
  );
};

export default ProjectExperimentsPage;
