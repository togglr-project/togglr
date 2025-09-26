import React, { useState } from 'react';
import { Box, Paper, Typography, Button, CircularProgress, Pagination } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import SearchPanel from '../components/SearchPanel';
import FeaturePreviewPanel from '../components/features/FeaturePreviewPanel';
import apiClient from '../api/apiClient';
import type { FeatureExtended, Project, ListProjectFeaturesKindEnum, ListProjectFeaturesSortByEnum, SortOrder, ListFeaturesResponse, ProjectTag } from '../generated/api/client';
import CreateFeatureDialog from '../components/features/CreateFeatureDialog';
import FeatureDetailsDialog from '../components/features/FeatureDetailsDialog';
import EditFeatureDialog from '../components/features/EditFeatureDialog';
import FeatureCard from '../components/features/FeatureCard';
import { useAuth } from '../auth/AuthContext';
import GuardResponseHandler from '../components/pending-changes/GuardResponseHandler';
import { useApprovePendingChange } from '../hooks/usePendingChanges';
import { useProjectPendingChanges } from '../hooks/useProjectPendingChanges';
import type { AuthCredentialsMethodEnum } from '../generated/api/client';

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
  const [selectedTags, setSelectedTags] = useState<ProjectTag[]>([]);

  const effectiveSearch = search.trim();
  const minSearch = effectiveSearch.length >= 3 ? effectiveSearch : undefined;
  const { data: featuresResp, isLoading: loadingFeatures, error: featuresError } = useQuery<ListFeaturesResponse>({
    queryKey: ['project-features', projectId, { search: minSearch, enabledFilter, kindFilter, sortBy, sortOrder, page, perPage, selectedTags }],
    queryFn: async () => {
      const tagIds = selectedTags.length > 0 ? selectedTags.map(tag => tag.id).join(',') : undefined;
      const res = await apiClient.listProjectFeatures(
        projectId,
        kindFilter === 'all' ? undefined : kindFilter,
        enabledFilter === 'all' ? undefined : enabledFilter === 'enabled',
        minSearch,
        tagIds,
        sortBy,
        sortOrder,
        page,
        perPage
      );
      return res.data;
    },
    enabled: !!projectId,
    placeholderData: keepPreviousData,
  });

  const features = featuresResp?.items ?? [];
  const pagination = featuresResp?.pagination;


  // Create Feature Dialog state
  const [open, setOpen] = useState(false);

  // Feature details dialog state & data
  const [detailsOpen, setDetailsOpen] = useState(false);
  const [selectedFeature, setSelectedFeature] = useState<FeatureExtended | null>(null);
  
  // Feature edit dialog state
  const [editFeature, setEditFeature] = useState<FeatureExtended | null>(null);
  const [editOpen, setEditOpen] = useState(false);
  
  // Feature preview panel state
  const [previewFeature, setPreviewFeature] = useState<FeatureExtended | null>(null);

  // Guard workflow state
  const [guardResponse, setGuardResponse] = useState<{
    pendingChange?: any;
    conflictError?: string;
    forbiddenError?: string;
  }>({});

  // Permission to toggle features in this project (superuser can always toggle)
  const canToggleFeature = Boolean(user?.is_superuser || user?.project_permissions?.[projectId]?.includes('feature.toggle'));

  // Get pending changes for the project
  const { data: pendingChanges } = useProjectPendingChanges(projectId);

  const approveMutation = useApprovePendingChange();

  // Toggle mutation
  const toggleMutation = useMutation({
    mutationFn: async ({ featureId, enabled }: { featureId: string; enabled: boolean }) => {
      try {
        const res = await apiClient.toggleFeature(featureId, { enabled });
        return { data: res.data, status: res.status };
      } catch (error: any) {
        // Handle guard workflow responses
        if (error.response?.status === 202) {
          // Pending change created
          setGuardResponse({ pendingChange: error.response.data });
          return null;
        }
        if (error.response?.status === 409) {
          // Conflict
          setGuardResponse({ conflictError: error.response.data.message });
          return null;
        }
        throw error;
      }
    },
    onSuccess: (result, variables) => {
      if (result) {
        if (result.status === 202) {
          // Pending change created - handle as guard workflow
          setGuardResponse({ pendingChange: result.data });
        } else if (result.status === 409) {
          // Conflict - feature locked by another pending change
          setGuardResponse({ conflictError: 'Feature is already locked by another pending change' });
        } else if (result.status === 403) {
          // Forbidden - user doesn't have permission to modify guarded feature
          setGuardResponse({ forbiddenError: 'You don\'t have permission to modify this guarded feature' });
        } else {
          // Normal success - toggle applied immediately
          queryClient.invalidateQueries({ queryKey: ['feature-details'] });
          queryClient.invalidateQueries({ queryKey: ['project-features'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
          queryClient.invalidateQueries({ queryKey: ['project-features', projectId] });
          queryClient.invalidateQueries({ queryKey: ['feature-details', variables.featureId] });
        }
      }
      // If result is null, guard workflow is handling the response
    },
  });


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
          queryClient.invalidateQueries({ queryKey: ['feature-details'] });
          queryClient.invalidateQueries({ queryKey: ['project-features'] });
          queryClient.invalidateQueries({ queryKey: ['pending-changes'] });
          queryClient.invalidateQueries({ queryKey: ['project-features', projectId] });
        },
      }
    );
  };

  const openFeatureDetails = (f: FeatureExtended) => {
    setSelectedFeature(f);
    setDetailsOpen(true);
  };

  const openFeatureEdit = (f: FeatureExtended) => {
    setEditFeature(f);
    setEditOpen(true);
  };

  const handleFeatureSelect = (f: FeatureExtended) => {
    // If clicking the same feature, deselect it
    if (previewFeature?.id === f.id) {
      setPreviewFeature(null);
    } else {
      setPreviewFeature(f);
    }
  };

  return (
    <AuthenticatedLayout showBackButton backTo="/dashboard">
      <Paper id="features" sx={{ p: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5 }}>
          <Typography variant="h6" sx={{ color: 'primary.light' }}>Features</Typography>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)} size="small">
            Add Feature
          </Button>
        </Box>

        {/* Search and filters */}
        <SearchPanel
          searchValue={search}
          onSearchChange={(value) => { setSearch(value); setPage(1); }}
          placeholder="Search features by name or key..."
          projectId={projectId}
          selectedTags={selectedTags}
          onTagsChange={(tags) => { setSelectedTags(tags); setPage(1); }}
          showTagFilter={true}
          quickFilters={[
            {
              label: 'Enabled',
              value: 'enabled',
              active: enabledFilter === 'enabled',
              onClick: () => { setEnabledFilter(enabledFilter === 'enabled' ? 'all' : 'enabled'); setPage(1); }
            },
            {
              label: 'Disabled',
              value: 'disabled',
              active: enabledFilter === 'disabled',
              onClick: () => { setEnabledFilter(enabledFilter === 'disabled' ? 'all' : 'disabled'); setPage(1); }
            },
            {
              label: 'Simple',
              value: 'simple',
              active: kindFilter === 'simple',
              onClick: () => { setKindFilter(kindFilter === 'simple' ? 'all' : 'simple'); setPage(1); }
            },
            {
              label: 'Multivariant',
              value: 'multivariant',
              active: kindFilter === 'multivariant',
              onClick: () => { setKindFilter(kindFilter === 'multivariant' ? 'all' : 'multivariant'); setPage(1); }
            },
          ]}
          filters={[
            {
              key: 'enabled',
              label: 'Status',
              value: enabledFilter,
              options: [
                { value: 'all', label: 'All' },
                { value: 'enabled', label: 'Enabled' },
                { value: 'disabled', label: 'Disabled' },
              ],
              onChange: (value) => { setEnabledFilter(value); setPage(1); }
            },
            {
              key: 'kind',
              label: 'Kind',
              value: kindFilter,
              options: [
                { value: 'all', label: 'All' },
                { value: 'simple', label: 'Simple' },
                { value: 'multivariant', label: 'Multivariant' },
              ],
              onChange: (value) => { setKindFilter(value); setPage(1); }
            },
            {
              key: 'sortBy',
              label: 'Sort by',
              value: sortBy,
              options: [
                { value: 'name', label: 'Name' },
                { value: 'key', label: 'Key' },
                { value: 'enabled', label: 'Enabled' },
                { value: 'kind', label: 'Kind' },
                { value: 'created_at', label: 'Created' },
                { value: 'updated_at', label: 'Updated' },
              ],
              onChange: (value) => { setSortBy(value); setPage(1); }
            },
            {
              key: 'sortOrder',
              label: 'Order',
              value: sortOrder,
              options: [
                { value: 'asc', label: 'Ascending' },
                { value: 'desc', label: 'Descending' },
              ],
              onChange: (value) => { setSortOrder(value); setPage(1); }
            },
            {
              key: 'perPage',
              label: 'Per page',
              value: perPage,
              options: [
                { value: 10, label: '10' },
                { value: 20, label: '20' },
                { value: 50, label: '50' },
                { value: 100, label: '100' },
              ],
              onChange: (value) => { setPerPage(Number(value)); setPage(1); }
            },
          ]}
        />

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
            {/* Features list and preview panel */}
            <Box sx={{ display: 'flex', gap: 2, alignItems: 'flex-start' }}>
              {/* Features list - 2/3 width */}
              <Box sx={{ flex: '2', display: 'flex', flexDirection: 'column', gap: 1 }}>
                {features.map((f) => (
                    <FeatureCard
                      key={f.id}
                      feature={f}
                      onEdit={openFeatureEdit}
                      onView={openFeatureDetails}
                      onSelect={handleFeatureSelect}
                      onToggle={(feature) => {
                        if (canToggleFeature) {
                          toggleMutation.mutate({ 
                            featureId: feature.id, 
                            enabled: !feature.enabled 
                          });
                        }
                      }}
                      canToggle={canToggleFeature}
                      isToggling={toggleMutation.isPending}
                      isSelected={previewFeature?.id === f.id}
                      projectId={projectId}
                    />
                ))}
              </Box>
              
              {/* Preview panel - 1/3 width */}
              <Box sx={{ flex: '1', minWidth: 300 }}>
                <FeaturePreviewPanel
                  selectedFeature={previewFeature}
                  projectId={projectId!}
                  onClose={() => setPreviewFeature(null)}
                />
              </Box>
            </Box>

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

      {/* Feature Edit Dialog */}
      <EditFeatureDialog 
        open={editOpen} 
        onClose={() => setEditOpen(false)} 
        featureDetails={null}
        feature={editFeature}
      />

      {/* Create Feature Dialog */}
      <CreateFeatureDialog open={open} onClose={() => setOpen(false)} projectId={projectId} />

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

export default ProjectPage;
