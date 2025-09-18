import React, { useState } from 'react';
import { Box, Paper, Typography, Button, CircularProgress, Grid, Chip, Switch, Tooltip, TextField, FormControl, InputLabel, Select, MenuItem, Stack, Pagination } from '@mui/material';
import { Add as AddIcon, Flag as FlagIcon } from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import type { Feature, Project, ListProjectFeaturesKindEnum, ListProjectFeaturesSortByEnum, SortOrder, ListFeaturesResponse } from '../generated/api/client';
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

  // Filters, sorting and pagination state
  const [search, setSearch] = useState('');
  const [enabledFilter, setEnabledFilter] = useState<'all' | 'enabled' | 'disabled'>('all');
  const [kindFilter, setKindFilter] = useState<ListProjectFeaturesKindEnum | 'all'>('all');
  const [sortBy, setSortBy] = useState<ListProjectFeaturesSortByEnum>('name');
  const [sortOrder, setSortOrder] = useState<SortOrder>('asc');
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);

  const effectiveSearch = search.trim();
  const minSearch = effectiveSearch.length >= 3 ? effectiveSearch : undefined;
  const { data: featuresResp, isLoading: loadingFeatures, error: featuresError } = useQuery<ListFeaturesResponse>({
    queryKey: ['project-features', projectId, { search: minSearch, enabledFilter, kindFilter, sortBy, sortOrder, page, perPage }],
    queryFn: async () => {
      const res = await apiClient.listProjectFeatures(
        projectId,
        kindFilter === 'all' ? undefined : kindFilter,
        enabledFilter === 'all' ? undefined : enabledFilter === 'enabled',
        minSearch,
        sortBy,
        sortOrder,
        page,
        perPage
      );
      return res.data;
    },
    enabled: !!projectId,
    keepPreviousData: true,
  });

  const features = featuresResp?.items ?? [];
  const pagination = featuresResp?.pagination;


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
        title={project ? `${project.name} - Features` : 'Project'}
        subtitle={project ? `Manage features in project ${project.name}` : 'Features'}
        icon={<FlagIcon />}
        gradientVariant="default"
        subtitleGradientVariant="default"
      />


      <Paper id="features" sx={{ p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" className="gradient-subtitle">Features</Typography>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)}>
            Add Feature
          </Button>
        </Box>

        {/* Filters and controls */}
        <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} sx={{ mb: 2 }}>
          <TextField
            label="Search by name or key"
            size="small"
            value={search}
            onChange={(e) => { setSearch(e.target.value); setPage(1); }}
            sx={{ minWidth: 240 }}
          />

          <FormControl size="small" sx={{ minWidth: 160 }}>
            <InputLabel id="enabled-filter-label">Enabled</InputLabel>
            <Select
              labelId="enabled-filter-label"
              label="Enabled"
              value={enabledFilter}
              onChange={(e) => { setEnabledFilter(e.target.value as any); setPage(1); }}
            >
              <MenuItem value="all">All</MenuItem>
              <MenuItem value="enabled">Enabled</MenuItem>
              <MenuItem value="disabled">Disabled</MenuItem>
            </Select>
          </FormControl>

          <FormControl size="small" sx={{ minWidth: 180 }}>
            <InputLabel id="kind-filter-label">Kind</InputLabel>
            <Select
              labelId="kind-filter-label"
              label="Kind"
              value={kindFilter}
              onChange={(e) => { setKindFilter(e.target.value as any); setPage(1); }}
            >
              <MenuItem value="all">All</MenuItem>
              <MenuItem value="simple">simple</MenuItem>
              <MenuItem value="multivariant">multivariant</MenuItem>
            </Select>
          </FormControl>

          <FormControl size="small" sx={{ minWidth: 180 }}>
            <InputLabel id="sort-by-label">Sort by</InputLabel>
            <Select
              labelId="sort-by-label"
              label="Sort by"
              value={sortBy}
              onChange={(e) => { setSortBy(e.target.value as any); setPage(1); }}
            >
              <MenuItem value="name">name</MenuItem>
              <MenuItem value="key">key</MenuItem>
              <MenuItem value="enabled">enabled</MenuItem>
              <MenuItem value="kind">kind</MenuItem>
              <MenuItem value="created_at">created_at</MenuItem>
              <MenuItem value="updated_at">updated_at</MenuItem>
            </Select>
          </FormControl>

          <FormControl size="small" sx={{ minWidth: 140 }}>
            <InputLabel id="sort-order-label">Order</InputLabel>
            <Select
              labelId="sort-order-label"
              label="Order"
              value={sortOrder}
              onChange={(e) => { setSortOrder(e.target.value as any); setPage(1); }}
            >
              <MenuItem value="asc">asc</MenuItem>
              <MenuItem value="desc">desc</MenuItem>
            </Select>
          </FormControl>

          <FormControl size="small" sx={{ minWidth: 120, ml: { xs: 0, md: 'auto' } }}>
            <InputLabel id="per-page-label">Per page</InputLabel>
            <Select
              labelId="per-page-label"
              label="Per page"
              value={perPage}
              onChange={(e) => { setPerPage(Number(e.target.value)); setPage(1); }}
            >
              <MenuItem value={10}>10</MenuItem>
              <MenuItem value={20}>20</MenuItem>
              <MenuItem value={50}>50</MenuItem>
              <MenuItem value={100}>100</MenuItem>
            </Select>
          </FormControl>
        </Stack>

        {(loadingProject || loadingFeatures) && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}
        {(projectError || featuresError) && (
          <Typography color="error">Failed to load project or features.</Typography>
        )}

        {!loadingFeatures && features && features.length > 0 ? (
          <>
            <Grid container spacing={2}>
              {features.map((f) => (
                <Grid item xs={12} key={f.id}>
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

            {/* Pagination */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mt: 2 }}>
              <Typography variant="body2" color="text.secondary">
                {pagination ? `Total: ${pagination.total}` : ''}
              </Typography>
              <Pagination
                page={page}
                count={pagination ? Math.max(1, Math.ceil(pagination.total / (pagination.per_page || perPage))) : 1}
                onChange={(_e, p) => setPage(p)}
                shape="rounded"
                color="primary"
              />
            </Box>
          </>
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
