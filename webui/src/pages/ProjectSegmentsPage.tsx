import React, { useMemo, useState } from 'react';
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
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PeopleOutline as PeopleIcon
} from '@mui/icons-material';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import ConditionExpressionBuilder from '../components/conditions/ConditionExpressionBuilder';
import apiClient from '../api/apiClient';
import type { Project, Segment, RuleConditionExpression } from '../generated/api/client';

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

const CreateEditSegmentDialog: React.FC<{
  open: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; description?: string; conditions: RuleConditionExpression }) => void;
  initial?: Partial<{ name: string; description?: string; conditions: RuleConditionExpression }>;
  title: string;
  submitting?: boolean;
  isEdit?: boolean;
}> = ({ open, onClose, onSubmit, initial, title, submitting, isEdit }) => {
  const [name, setName] = useState<string>(initial?.name || '');
  const [description, setDescription] = useState<string>(initial?.description || '');
  const [expr, setExpr] = useState<RuleConditionExpression>(initial?.conditions || { group: { operator: 'and', children: [{ condition: { attribute: '', operator: 'eq', value: '' } }] } as any });
  const [error, setError] = useState<string>('');

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
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="md">
      <DialogTitle className="gradient-text-purple">{title}</DialogTitle>
      <DialogContent dividers>
        <Box sx={{ display: 'grid', gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' }, gap: 2 }}>
          <TextField label="Name" value={name} onChange={(e) => setName(e.target.value)} required fullWidth />
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
        {isEdit && (
          <Button variant="outlined" color="secondary" onClick={() => alert('Sync of customized features is not implemented yet')}>
            Sync customized features
          </Button>
        )}
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSubmit} disabled={!canSubmit} variant="contained">Save</Button>
      </DialogActions>
    </Dialog>
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

  const { data: segments, isLoading, error } = useQuery<Segment[]>({
    queryKey: ['project-segments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectSegments(projectId);
      return res.data;
    },
    enabled: !!projectId,
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
        gradientVariant="default"
        subtitleGradientVariant="default"
      />

      <Paper sx={{ p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" className="gradient-subtitle">Segments</Typography>
          <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpenCreate(true)}>
            Add Segment
          </Button>
        </Box>

        {(loadingProject || isLoading) && (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        )}
        {error && (
          <Typography color="error">Failed to load segments.</Typography>
        )}

        {!isLoading && segments && segments.length > 0 ? (
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
            {segments.map((s) => (
              <Paper key={s.id} sx={{ p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
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
                      <IconButton onClick={() => setConfirmDelete(s)} aria-label="delete-segment" disabled={deleteMutation.isPending}>
                        <DeleteIcon />
                      </IconButton>
                    </span>
                  </Tooltip>
                </Box>
              </Paper>
            ))}
          </Box>
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
          <Button onClick={() => setConfirmDelete(null)}>Cancel</Button>
          <Button color="error" variant="contained" onClick={() => confirmDelete && deleteMutation.mutate(confirmDelete.id)} disabled={deleteMutation.isPending}>Delete</Button>
        </DialogActions>
      </Dialog>

    </AuthenticatedLayout>
  );
};

export default ProjectSegmentsPage;
