import React, { useState } from 'react';
import { Box, Paper, Typography, Button, CircularProgress, Grid, Chip, Switch, Tooltip } from '@mui/material';
import { Add as AddIcon, Flag as FlagIcon } from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import type { Feature, Project } from '../generated/api/client';
import CreateFeatureDialog from '../components/features/CreateFeatureDialog';
import FeatureDetailsDialog from '../components/features/FeatureDetailsDialog';
import { useAuth } from '../auth/AuthContext';

interface ProjectResponse { project: Project }


const ProjectPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const { user } = useAuth();
  const queryClient = useQueryClient();

  const { data: projectResp, isLoading: loadingProject, error: projectError } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const res = await apiClient.getProject(projectId);
      return res.data as ProjectResponse;
    },
    enabled: !!projectId,
  });

  const { data: features, isLoading: loadingFeatures, error: featuresError } = useQuery<Feature[]>({
    queryKey: ['project-features', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectFeatures(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  // Create Feature Dialog state
  const [open, setOpen] = useState(false);

  const project = projectResp?.project;

  // Feature details dialog state & data
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [selectedFeature, setSelectedFeature] = useState<Feature | null>(null);

  // Permission to toggle features in this project (superuser can always toggle)
  const canToggleFeature = Boolean(user?.is_superuser || user?.project_permissions?.[projectId]?.includes('feature.toggle'));

  // Toggle mutation
  const toggleMutation = useMutation({
    mutationFn: async ({ featureId, enabled }: { featureId: string; enabled: boolean }) => {
      await apiClient.toggleFeature(featureId, { enabled });
    },
    onSuccess: (_data, variables) => {
      // Refresh lists and details after toggle
      queryClient.invalidateQueries({ queryKey: ['project-features', projectId] });
      queryClient.invalidateQueries({ queryKey: ['feature-details', variables.featureId] });
    },
  });

  const openFeatureDetails = (f: Feature) => {
    setSelectedFeature(f);
    setDetailsOpen(true);
  };

  return (
    <AuthenticatedLayout showBackButton backTo="/dashboard">
      <PageHeader
        title={project ? project.name : 'Project'}
        subtitle={project ? `ID: ${project.id}` : 'Project details'}
        icon={<FlagIcon />}
        gradientVariant="default"
        subtitleGradientVariant="default"
      />

      <Paper sx={{ p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" className="gradient-subtitle">Features</Typography>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)}>
            Add Feature
          </Button>
        </Box>

        {(loadingProject || loadingFeatures) && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}
        {(projectError || featuresError) && (
          <Typography color="error">Failed to load project or features.</Typography>
        )}

        {!loadingFeatures && features && features.length > 0 ? (
          <Grid container spacing={2}>
            {features.map((f) => (
              <Grid item xs={12} md={6} key={f.id}>
                <Paper
                  onClick={() => openFeatureDetails(f)}
                  sx={{
                    p: 2,
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    cursor: 'pointer',
                    transition: 'box-shadow 0.2s, transform 0.1s',
                    '&:hover': { boxShadow: 4 },
                    '&:active': { transform: 'scale(0.997)' }
                  }}
                  role="button"
                >
                  <Box>
                    <Typography variant="subtitle1">{f.name}</Typography>
                    <Typography variant="body2" color="text.secondary">{f.key}</Typography>
                    <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip size="small" label={`kind: ${f.kind}`} />
                      <Chip size="small" label={`default: ${f.default_variant}`} />
                      <Chip size="small" label={f.enabled ? 'enabled' : 'disabled'} color={f.enabled ? 'success' : 'default'} />
                    </Box>
                  </Box>
                  {canToggleFeature ? (
                    <Tooltip title={f.enabled ? 'Disable feature' : 'Enable feature'}>
                      <Switch
                        checked={f.enabled}
                        onClick={(e) => e.stopPropagation()}
                        onChange={(e) => {
                          e.stopPropagation();
                          const enabled = e.target.checked;
                          toggleMutation.mutate({ featureId: f.id, enabled });
                        }}
                        disabled={toggleMutation.isPending}
                        inputProps={{ 'aria-label': 'toggle feature' }}
                      />
                    </Tooltip>
                  ) : (
                    <Tooltip title="You don't have permission to toggle features in this project">
                      <span onClick={(e) => e.stopPropagation()}>
                        <Switch checked={f.enabled} disabled />
                      </span>
                    </Tooltip>
                  )}
                </Paper>
              </Grid>
            ))}
          </Grid>
        ) : !loadingFeatures ? (
          <Typography variant="body2">No features yet.</Typography>
        ) : null}
      </Paper>

      {/* Feature Details Dialog */}
      <FeatureDetailsDialog open={detailsOpen} onClose={() => setDetailsOpen(false)} feature={selectedFeature} />

      {/* Create Feature Dialog */}
      <CreateFeatureDialog open={open} onClose={() => setOpen(false)} projectId={projectId} />

    </AuthenticatedLayout>
  );
};

export default ProjectPage;
