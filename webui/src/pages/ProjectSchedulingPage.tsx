import React, { useMemo, useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Button,
  Chip,
  CircularProgress,
  IconButton,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  Grid
} from '@mui/material';
import {
  ExpandMore as ExpandMoreIcon,
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Schedule as ScheduleIcon,
  Flag as FlagIcon
} from '@mui/icons-material';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import type { Feature, FeatureSchedule, FeatureScheduleAction, Project } from '../generated/api/client';

interface ProjectResponse { project: Project }

interface ScheduleFormValues {
  action: FeatureScheduleAction;
  timezone: string;
  starts_at?: string;
  ends_at?: string;
  cron_expr?: string;
}

const emptyForm = (): ScheduleFormValues => ({
  action: 'enable',
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
  starts_at: '',
  ends_at: '',
  cron_expr: ''
});

const normalizeDateTimeLocalToISO = (val?: string): string | undefined => {
  if (!val) return undefined;
  // Convert local datetime-local string to ISO UTC
  const d = new Date(val);
  if (isNaN(d.getTime())) return undefined;
  return d.toISOString();
};

const ScheduleDialog: React.FC<{
  open: boolean;
  onClose: () => void;
  onSubmit: (values: ScheduleFormValues) => void;
  initial?: Partial<ScheduleFormValues>;
  title: string;
}> = ({ open, onClose, onSubmit, initial, title }) => {
  const [values, setValues] = useState<ScheduleFormValues>(() => ({ ...emptyForm(), ...initial } as ScheduleFormValues));

  React.useEffect(() => {
    setValues({ ...emptyForm(), ...initial } as ScheduleFormValues);
  }, [initial, open]);

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>{title}</DialogTitle>
      <DialogContent dividers>
        <Box sx={{ mt: 1 }}>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <TextField
                select
                fullWidth
                label="Action"
                value={values.action}
                onChange={(e) => setValues(v => ({ ...v, action: e.target.value as FeatureScheduleAction }))}
              >
                <MenuItem value="enable">Enable</MenuItem>
                <MenuItem value="disable">Disable</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Timezone"
                value={values.timezone}
                onChange={(e) => setValues(v => ({ ...v, timezone: e.target.value }))}
                helperText="IANA timezone, e.g., UTC, Europe/Moscow"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Starts at"
                type="datetime-local"
                InputLabelProps={{ shrink: true }}
                value={values.starts_at || ''}
                onChange={(e) => setValues(v => ({ ...v, starts_at: e.target.value }))}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Ends at"
                type="datetime-local"
                InputLabelProps={{ shrink: true }}
                value={values.ends_at || ''}
                onChange={(e) => setValues(v => ({ ...v, ends_at: e.target.value }))}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Cron expression"
                placeholder="Optional, e.g., 0 8 * * 1-5"
                value={values.cron_expr || ''}
                onChange={(e) => setValues(v => ({ ...v, cron_expr: e.target.value }))}
                helperText="If set, schedule will follow this cron (timezone above)."
              />
            </Grid>
          </Grid>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button variant="contained" onClick={() => {
          const payload: ScheduleFormValues = {
            ...values,
            starts_at: normalizeDateTimeLocalToISO(values.starts_at),
            ends_at: normalizeDateTimeLocalToISO(values.ends_at),
          };
          onSubmit(payload);
        }}>Save</Button>
      </DialogActions>
    </Dialog>
  );
};

const ProjectSchedulingPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const navigate = useNavigate();
  const qc = useQueryClient();

  const { data: projectResp, isLoading: loadingProject } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const res = await apiClient.getProject(projectId);
      return res.data as ProjectResponse;
    },
    enabled: !!projectId,
  });

  const { data: features, isLoading: loadingFeatures } = useQuery<Feature[]>({
    queryKey: ['project-features', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectFeatures(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  const { data: allSchedules, isLoading: loadingSchedules } = useQuery<FeatureSchedule[]>({
    queryKey: ['feature-schedules', projectId],
    queryFn: async () => {
      const res = await apiClient.listAllFeatureSchedules();
      // Filter by project just in case API returns global list
      return (res.data || []).filter((s: FeatureSchedule) => s.project_id === projectId);
    },
    enabled: !!projectId,
  });

  const schedulesByFeature = useMemo(() => {
    const map: Record<string, FeatureSchedule[]> = {};
    (allSchedules || []).forEach(s => {
      if (!map[s.feature_id]) map[s.feature_id] = [];
      map[s.feature_id].push(s);
    });
    // sort by created_at desc
    Object.values(map).forEach(list => list.sort((a, b) => (b.created_at.localeCompare(a.created_at))));
    return map;
  }, [allSchedules]);

  // Dialog state
  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogFeature, setDialogFeature] = useState<Feature | null>(null);
  const [editSchedule, setEditSchedule] = useState<FeatureSchedule | null>(null);

  const openCreate = (feature: Feature) => {
    setDialogFeature(feature);
    setEditSchedule(null);
    setDialogOpen(true);
  };
  const openEdit = (feature: Feature, schedule: FeatureSchedule) => {
    setDialogFeature(feature);
    setEditSchedule(schedule);
    setDialogOpen(true);
  };

  const closeDialog = () => setDialogOpen(false);

  // Mutations
  const createMut = useMutation({
    mutationFn: async ({ featureId, values }: { featureId: string; values: ScheduleFormValues }) => {
      return apiClient.createFeatureSchedule(featureId, values);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
    }
  });

  const updateMut = useMutation({
    mutationFn: async ({ scheduleId, values }: { scheduleId: string; values: ScheduleFormValues }) => {
      return apiClient.updateFeatureSchedule(scheduleId, values);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
    }
  });

  const deleteMut = useMutation({
    mutationFn: async ({ scheduleId }: { scheduleId: string }) => {
      return apiClient.deleteFeatureSchedule(scheduleId);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
    }
  });

  const project = projectResp?.project;

  return (
    <AuthenticatedLayout showBackButton backTo={`/projects/${projectId}`}>
      <PageHeader
        title={project ? `${project.name} — Scheduling` : 'Scheduling'}
        subtitle={project ? `Manage feature schedules in project ${project.name}` : 'Feature schedules'}
        icon={<ScheduleIcon />}
        gradientVariant="default"
        subtitleGradientVariant="default"
      />

      {(loadingProject || loadingFeatures || loadingSchedules) && (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {!loadingFeatures && features && features.length > 0 ? (
        <Box>
          {features.map((f) => (
            <Accordion key={f.id} defaultExpanded sx={{ mb: 2 }}>
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%', justifyContent: 'space-between' }}>
                  <Box>
                    <Typography variant="subtitle1">{f.name}</Typography>
                    <Typography variant="body2" color="text.secondary">{f.key}</Typography>
                    <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip size="small" label={`kind: ${f.kind}`} />
                      <Chip size="small" label={`default: ${f.default_variant}`} />
                      <Chip size="small" label={f.enabled ? 'enabled' : 'disabled'} color={f.enabled ? 'success' : 'default'} />
                    </Box>
                  </Box>
                  <Button variant="contained" startIcon={<AddIcon />} onClick={(e) => { e.stopPropagation(); openCreate(f); }}>
                    Add schedule
                  </Button>
                </Box>
              </AccordionSummary>
              <AccordionDetails>
                <Box>
                  {(schedulesByFeature[f.id] && schedulesByFeature[f.id].length > 0) ? (
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                      {schedulesByFeature[f.id].map((s) => (
                        <Paper key={s.id} sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                          <Box>
                            <Typography variant="body1" sx={{ fontWeight: 600, textTransform: 'capitalize' }}>{s.action}</Typography>
                            <Typography variant="body2" color="text.secondary">
                              {s.cron_expr ? `Cron: ${s.cron_expr}` : `From ${s.starts_at || '—'} to ${s.ends_at || '—'}`}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">Timezone: {s.timezone}</Typography>
                          </Box>
                          <Box>
                            <Tooltip title="Edit schedule">
                              <IconButton onClick={() => openEdit(f, s)}><EditIcon /></IconButton>
                            </Tooltip>
                            <Tooltip title="Delete schedule">
                              <IconButton color="error" onClick={() => deleteMut.mutate({ scheduleId: s.id })}><DeleteIcon /></IconButton>
                            </Tooltip>
                          </Box>
                        </Paper>
                      ))}
                    </Box>
                  ) : (
                    <Typography variant="body2" color="text.secondary">No schedules for this feature yet.</Typography>
                  )}
                </Box>
              </AccordionDetails>
            </Accordion>
          ))}
        </Box>
      ) : !loadingFeatures ? (
        <Typography variant="body2">No features yet.</Typography>
      ) : null}

      {/* Create/Edit Dialog */}
      <ScheduleDialog
        open={dialogOpen}
        onClose={closeDialog}
        title={editSchedule ? 'Edit schedule' : 'Create schedule'}
        initial={editSchedule ? {
          action: editSchedule.action,
          timezone: editSchedule.timezone,
          // Convert ISO to datetime-local format yyyy-MM-ddTHH:mm
          starts_at: editSchedule.starts_at ? new Date(editSchedule.starts_at).toISOString().slice(0,16) : '',
          ends_at: editSchedule.ends_at ? new Date(editSchedule.ends_at).toISOString().slice(0,16) : '',
          cron_expr: editSchedule.cron_expr || ''
        } : undefined}
        onSubmit={(values) => {
          if (!dialogFeature) return;
          if (editSchedule) {
            updateMut.mutate({ scheduleId: editSchedule.id, values });
          } else {
            createMut.mutate({ featureId: dialogFeature.id, values });
          }
          setDialogOpen(false);
        }}
      />
    </AuthenticatedLayout>
  );
};

export default ProjectSchedulingPage;
