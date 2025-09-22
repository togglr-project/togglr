import React, { useEffect, useMemo, useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  CircularProgress,
  IconButton,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Grid,
  Chip,
  MenuItem,
  Checkbox,
  FormControl,
  InputLabel,
  Select,
  Stack,
  Pagination,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PeopleOutline as PeopleIcon
} from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import SearchPanel from '../components/SearchPanel';
import ConditionExpressionBuilder from '../components/conditions/ConditionExpressionBuilder';
import apiClient from '../api/apiClient';
import type { Project, Segment, RuleConditionExpression, ListProjectSegmentsSortByEnum, SortOrder, ListSegmentsResponse } from '../generated/api/client';

interface ProjectResponse { project: Project }

const countLeaves = (e?: RuleConditionExpression): number => {
  if (!e) return 0;
  if ((e as any).condition) {
    const c = (e as any).condition as { attribute?: string };
    return c.attribute && c.attribute.trim() ? 1 : 0;
  }
  if ((e as any).group) {
    const g = (e as any).group as { children?: RuleConditionExpression[] };
    return (g.children || []).reduce((sum, ch) => sum + countLeaves(ch), 0);
  }
  return 0;
}

// Condition type matches RuleCondition from API, but we only need shape here
// attribute: string; operator: string; value: any
interface ConditionFormItem { attribute: string; operator: OperatorOption; value: string }

type OperatorOption = 'eq' | 'neq' | 'in' | 'not_in' | 'gt' | 'gte' | 'lt' | 'lte' | 'regex' | 'percentage';
const operatorOptions: OperatorOption[] = ['eq','neq','in','not_in','gt','gte','lt','lte','regex','percentage'];

const emptyCondition = (): ConditionFormItem => ({ attribute: '', operator: 'eq', value: '' });

const SegmentDesyncCount: React.FC<{ segmentId: string }> = ({ segmentId }) => {
  const { data, isLoading, isError } = useQuery<string[]>({
    queryKey: ['segment-desync-ids', segmentId],
    queryFn: async () => {
      const res = await apiClient.listSegmentDesyncFeatureIDs(segmentId);
      return res.data as unknown as string[];
    },
    enabled: !!segmentId,
  });

  if (isLoading) {
    return <Chip size="small" label="customized: …" />;
  }
  if (isError) {
    return <Chip size="small" label="customized: error" color="error" />;
  }
  const count = Array.isArray(data) ? data.length : 0;
  return <Chip size="small" label={`customized: ${count}`} color={count > 0 ? 'warning' as any : undefined} />;
};

// Dialog to batch sync customized features for a segment
const SyncCustomizedFeaturesDialog: React.FC<{ open: boolean; onClose: () => void; segmentId: string; }> = ({ open, onClose, segmentId }) => {
  const queryClient = useQueryClient();
  const { data: idsData, isLoading: loadingIds, isError: idsError } = useQuery<string[]>({
    queryKey: ['segment-desync-ids', segmentId],
    queryFn: async () => {
      const res = await apiClient.listSegmentDesyncFeatureIDs(segmentId);
      return res.data as unknown as string[];
    },
    enabled: Boolean(open && segmentId),
  });
  const ids = Array.isArray(idsData) ? idsData : [];

  const { data: featuresData, isLoading: loadingFeatures } = useQuery<any[]>({
    queryKey: ['segment-desync-features', segmentId, ids],
    queryFn: async () => {
      const arr = await Promise.all(ids.map(async (fid) => {
        const r = await apiClient.getFeature(fid);
        return r.data as any;
      }));
      return arr;
    },
    enabled: Boolean(open && segmentId && ids.length > 0),
  });

  const [selected, setSelected] = useState<Record<string, boolean>>({});
  useEffect(() => {
    if (!open) return;
    const next: Record<string, boolean> = {};
    (featuresData || []).forEach((fd: any) => { next[fd.feature.id] = true; });
    setSelected(next);
  }, [open, featuresData]);

  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const toggle = (fid: string) => setSelected(prev => ({ ...prev, [fid]: !prev[fid] }));

  const findRuleId = (fd: any): string | null => {
    const rules = fd?.rules || [];
    const preferred = rules.find((r: any) => r.segment_id === segmentId && r.is_customized);
    const anyMatch = preferred || rules.find((r: any) => r.segment_id === segmentId);
    return anyMatch ? anyMatch.id : null;
  };

  const handleSync = async () => {
    setSubmitting(true);
    setSubmitError(null);
    try {
      const items = (featuresData || []) as any[];
      await Promise.allSettled(items.map(async (fd: any) => {
        const fid = fd.feature.id as string;
        if (!selected[fid]) return;
        const ruleId = findRuleId(fd);
        if (!ruleId) return;
        await apiClient.syncCustomizedFeatureRule(fid, ruleId);
        // refresh specific feature cache if used elsewhere
        await queryClient.invalidateQueries({ queryKey: ['feature-details', fid] });
      }));
      await queryClient.invalidateQueries({ queryKey: ['segment-desync-ids', segmentId] });
      onClose();
    } catch (e: any) {
      setSubmitError(e?.message || 'Failed to synchronize');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>Sync customized features</DialogTitle>
      <DialogContent dividers>
        {loadingIds || loadingFeatures ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
            <CircularProgress />
          </Box>
        ) : idsError ? (
          <Typography color="error">Failed to load customized features.</Typography>
        ) : ids.length === 0 ? (
          <Typography variant="body2">No customized features for this segment.</Typography>
        ) : (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
            {(featuresData || []).map((fd: any) => {
              const fid = fd.feature.id as string;
              const name = fd.feature.name || fd.feature.key || fid;
              const rid = findRuleId(fd);
              const disabled = !rid;
              return (
                <Box key={fid} sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Checkbox checked={!!selected[fid]} onChange={() => toggle(fid)} disabled={disabled} />
                  <Typography variant="body2" sx={{ flex: 1 }}>{name}</Typography>
                  {!rid && <Chip size="small" color="default" label="No rule for this segment" />}
                </Box>
              );
            })}
          </Box>
        )}
        {submitError && <Typography color="error" sx={{ mt: 1 }}>{submitError}</Typography>}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={submitting} size="small">Cancel</Button>
        <Button onClick={handleSync} variant="contained" disabled={submitting || (ids?.length || 0) === 0} size="small">
          {submitting ? 'Synchronizing…' : 'Synchronize'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

const CreateEditSegmentDialog: React.FC<{
  open: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; description?: string; conditions: RuleConditionExpression }) => void;
  initial?: Partial<{ name: string; description?: string; conditions: RuleConditionExpression }>;
  title: string;
  submitting?: boolean;
  isEdit?: boolean;
  segmentId?: string;
}> = ({ open, onClose, onSubmit, initial, title, submitting, isEdit, segmentId }) => {
  const [name, setName] = useState<string>(initial?.name || '');
  const [description, setDescription] = useState<string>(initial?.description || '');
  const [expr, setExpr] = useState<RuleConditionExpression>(initial?.conditions || { group: { operator: 'and', children: [{ condition: { attribute: '', operator: 'eq', value: '' } }] } as any });
  const [error, setError] = useState<string>('');

  // Desync features for this segment (to show sync button conditionally)
  const { data: desyncIds } = useQuery<string[]>({
    queryKey: ['segment-desync-ids', segmentId],
    queryFn: async () => {
      if (!segmentId) return [] as string[];
      const res = await apiClient.listSegmentDesyncFeatureIDs(segmentId);
      return res.data as unknown as string[];
    },
    enabled: Boolean(open && isEdit && segmentId),
  });
  const desyncCount = Array.isArray(desyncIds) ? desyncIds.length : 0;
  const [syncOpen, setSyncOpen] = useState(false);

  React.useEffect(() => {
    setName(initial?.name || '');
    setDescription(initial?.description || '');
    setExpr(initial?.conditions || { group: { operator: 'and', children: [{ condition: { attribute: '', operator: 'eq', value: '' } }] } as any });
    setError('');
  }, [initial, open]);

  const hasValidLeaf = (e?: RuleConditionExpression): boolean => {
    if (!e) return false;
    if ((e as any).condition) {
      const c = (e as any).condition as { attribute?: string };
      return Boolean(c.attribute && c.attribute.trim().length > 0);
    }
    if ((e as any).group) {
      const g = (e as any).group as { children?: RuleConditionExpression[] };
      return Array.isArray(g.children) && g.children.some(ch => hasValidLeaf(ch));
    }
    return false;
  };

  const canSubmit = useMemo(() => {
    if (!name.trim()) return false;
    if (!hasValidLeaf(expr)) return false;
    return !submitting;
  }, [name, expr, submitting]);

  const handleSubmit = () => {
    setError('');
    try {
      onSubmit({ name: name.trim(), description: description.trim() || undefined, conditions: expr });
    } catch (e) {
      setError('Failed to prepare data');
    }
  };

  return (
    <>
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle sx={{ color: 'primary.main' }}>{title}</DialogTitle>
      <DialogContent dividers>
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 1.5 }}>
          <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} required fullWidth size="small" />
          <TextField label="Description" value={description} onChange={(e) => setDescription(e.target.value)} fullWidth multiline minRows={2} />
        </Box>

        <Box sx={{ mt: 3 }}>
          <Typography variant="subtitle1" sx={{ mb: 1 }}>Conditions</Typography>
          <ConditionExpressionBuilder value={expr} onChange={setExpr} />
        </Box>

        {error && (
          <Typography color="error" sx={{ mt: 2 }}>{error}</Typography>
        )}
      </DialogContent>
      <DialogActions>
        {isEdit && desyncCount > 0 && (
          <Button variant="outlined" color="secondary" onClick={() => setSyncOpen(true)} size="small">
            Sync customized features
          </Button>
        )}
        <Button onClick={onClose} size="small">Cancel</Button>
        <Button onClick={handleSubmit} disabled={!canSubmit} variant="contained" size="small">Save</Button>
      </DialogActions>
    </Dialog>
    <SyncCustomizedFeaturesDialog open={syncOpen} onClose={() => setSyncOpen(false)} segmentId={segmentId || ''} />
  </>
  );
};

const ProjectSegmentsPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const queryClient = useQueryClient();

  const { data: projectResp, isLoading: loadingProject } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const res = await apiClient.getProject(projectId);
      return res.data as ProjectResponse;
    },
    enabled: !!projectId,
  });

  // Filters, sorting and pagination state for segments
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<ListProjectSegmentsSortByEnum>('name');
  const [sortOrder, setSortOrder] = useState<SortOrder>('asc');
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);

  const effectiveSearch = search.trim();
  const minSearch = effectiveSearch.length >= 3 ? effectiveSearch : undefined;

  const { data: segmentsResp, isLoading, error } = useQuery<ListSegmentsResponse>({
    queryKey: ['project-segments', projectId, { search: minSearch, sortBy, sortOrder, page, perPage }],
    queryFn: async () => {
      const res = await apiClient.listProjectSegments(
        projectId,
        minSearch,
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

  // Dialog state
  const [openCreate, setOpenCreate] = useState(false);
  const [editData, setEditData] = useState<Segment | null>(null);
  const [confirmDelete, setConfirmDelete] = useState<Segment | null>(null);

  const createMutation = useMutation({
    mutationFn: async (payload: { name: string; description?: string; conditions: RuleConditionExpression }) => {
      await apiClient.createProjectSegment(projectId, payload);
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['project-segments', projectId] });
      setOpenCreate(false);
    }
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, payload }: { id: string; payload: { name: string; description?: string; conditions: RuleConditionExpression } }) => {
      await apiClient.updateSegment(id, payload);
    },
    onSuccess: async (_d, vars) => {
      await queryClient.invalidateQueries({ queryKey: ['project-segments', projectId] });
      await queryClient.invalidateQueries({ queryKey: ['segment', vars.id] });
      setEditData(null);
    }
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.deleteSegment(id);
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['project-segments', projectId] });
      setConfirmDelete(null);
    }
  });

  const project = projectResp?.project;

  return (
    <AuthenticatedLayout showBackButton backTo="/dashboard">
      <PageHeader
        title={project ? `${project.name} - Segments` : 'Project Segments'}
        subtitle={project ? `Manage user segments for ${project.name}` : 'Segments'}
        icon={<PeopleIcon />}
      />

      <Paper sx={{ p: 2 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1.5 }}>
          <Typography variant="h6" sx={{ color: 'primary.light' }}>Segments</Typography>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpenCreate(true)} size="small">
            Add Segment
          </Button>
        </Box>

        {/* Search and filters */}
        <SearchPanel
          searchValue={search}
          onSearchChange={(value) => { setSearch(value); setPage(1); }}
          placeholder="Search by name or description"
          filters={[
            {
              key: 'sortBy',
              label: 'Sort by',
              value: sortBy,
              options: [
                { value: 'name', label: 'Name' },
                { value: 'created_at', label: 'Created' },
                { value: 'updated_at', label: 'Updated' },
              ],
              onChange: (value) => { setSortBy(value as any); setPage(1); },
            },
            {
              key: 'sortOrder',
              label: 'Order',
              value: sortOrder,
              options: [
                { value: 'asc', label: 'Ascending' },
                { value: 'desc', label: 'Descending' },
              ],
              onChange: (value) => { setSortOrder(value as any); setPage(1); },
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
              onChange: (value) => { setPerPage(Number(value)); setPage(1); },
            },
          ]}
        />

        {(loadingProject || isLoading) && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}
        {error && (
          <Typography color="error">Failed to load segments.</Typography>
        )}

        {!isLoading && (segmentsResp?.items?.length || 0) > 0 ? (
          <>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              {segmentsResp!.items.map((s) => (
                <Paper key={s.id} sx={{ p: 1.5, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Box>
                    <Typography variant="subtitle1">{s.name}</Typography>
                    {s.description && (
                      <Typography variant="body2" color="text.secondary">{s.description}</Typography>
                    )}
                    <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip size="small" label={`conditions: ${countLeaves(s.conditions)}`} />
                      <SegmentDesyncCount segmentId={s.id} />
                    </Box>
                  </Box>
                  <Box sx={{ display: 'flex', gap: 1 }}>
                    <Tooltip title="Edit">
                      <span>
                        <IconButton onClick={() => setEditData(s)} aria-label="edit-segment">
                          <EditIcon />
                        </IconButton>
                      </span>
                    </Tooltip>
                    <Tooltip title="Delete">
                      <span>
                        <IconButton 
                          onClick={() => setConfirmDelete(s)} 
                          aria-label="delete-segment" 
                          disabled={deleteMutation.isPending}
                          sx={{ color: 'error.main' }}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </span>
                    </Tooltip>
                  </Box>
                </Paper>
              ))}
            </Box>

            {/* Pagination */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mt: 1.5 }}>
              <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.8rem' }}>
                {segmentsResp?.pagination ? `Total: ${segmentsResp.pagination.total}` : ''}
              </Typography>
              <Pagination
                size="small"
                page={page}
                count={segmentsResp?.pagination ? Math.max(1, Math.ceil(segmentsResp.pagination.total / (segmentsResp.pagination.per_page || perPage))) : 1}
                onChange={(_e, p) => setPage(p)}
                shape="rounded"
                color="primary"
              />
            </Box>
          </>
        ) : !isLoading ? (
          <Typography variant="body2">No segments yet.</Typography>
        ) : null}
      </Paper>

      {/* Create Dialog */}
      <CreateEditSegmentDialog
        open={openCreate}
        onClose={() => setOpenCreate(false)}
        title="Create Segment"
        submitting={createMutation.isPending}
        isEdit={false}
        onSubmit={(values) => createMutation.mutate(values)}
      />

      {/* Edit Dialog */}
      <CreateEditSegmentDialog
        open={!!editData}
        onClose={() => setEditData(null)}
        title="Edit Segment"
        submitting={updateMutation.isPending}
        isEdit={true}
        segmentId={editData?.id}
        initial={editData ? { name: editData.name, description: editData.description, conditions: editData.conditions as any } : undefined}
        onSubmit={(values) => editData && updateMutation.mutate({ id: editData.id, payload: values })}
      />

      {/* Delete confirm */}
      <Dialog open={!!confirmDelete} onClose={() => setConfirmDelete(null)}>
        <DialogTitle>Delete Segment</DialogTitle>
        <DialogContent>
          <Typography>Are you sure you want to delete segment "{confirmDelete?.name}"?</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmDelete(null)} size="small">Cancel</Button>
          <Button color="error" variant="contained" onClick={() => confirmDelete && deleteMutation.mutate(confirmDelete.id)} disabled={deleteMutation.isPending} size="small">Delete</Button>
        </DialogActions>
      </Dialog>

    </AuthenticatedLayout>
  );
};

export default ProjectSegmentsPage;
