import React, { useMemo, useState, useEffect } from 'react';
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
  Grid,
  Tabs,
  Tab,
  Stack,
  Pagination,
  FormControl,
  InputLabel,
  Select
} from '@mui/material';
import {
  ExpandMore as ExpandMoreIcon,
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Schedule as ScheduleIcon,
  Refresh as RefreshIcon,
  Help as HelpIcon
} from '@mui/icons-material';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import SearchPanel from '../components/SearchPanel';
import TimelineChart from '../components/TimelineChart';
import ScheduleBuilder from '../components/ScheduleBuilder';
import OneShotScheduleDialog from '../components/OneShotScheduleDialog';
import EditRecurringScheduleBuilder from '../components/EditRecurringScheduleBuilder';
import EditOneShotScheduleDialog from '../components/EditOneShotScheduleDialog';
import ScheduleHelpDialog from '../components/ScheduleHelpDialog';
import ScheduleDialog from '../components/schedules/ScheduleDialog';
import apiClient from '../api/apiClient';
import { canAddRecurringSchedule, canAddOneShotSchedule, getScheduleType } from '../utils/scheduleHelpers';
import type { ScheduleBuilderData } from '../utils/cronGenerator';
import type { FeatureExtended, FeatureSchedule, FeatureScheduleAction, Project, ListProjectFeaturesKindEnum, ListProjectFeaturesSortByEnum, SortOrder, ListFeaturesResponse, FeatureTimelineResponse, FeatureTimelineEvent, ProjectTag } from '../generated/api/client';
import { isValidCron } from 'cron-validator';
import cronstrue from 'cronstrue';
import { listTimeZones, findTimeZone, getZonedTime, getUTCOffset } from 'timezone-support';

interface ProjectResponse { project: Project }

interface ScheduleFormValues {
  action: FeatureScheduleAction;
  timezone: string;
  starts_at?: string;
  ends_at?: string;
  cron_expr?: string;
  cron_duration?: string;
}

const emptyForm = (): ScheduleFormValues => ({
  action: 'enable',
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
  starts_at: '',
  ends_at: '',
  cron_expr: '',
  cron_duration: ''
});

const pad2 = (n: number) => (n < 10 ? '0' + n : '' + n);

const toDatetimeLocalInZone = (iso?: string, tz?: string): string => {
  if (!iso) return '';
  try {
    const date = new Date(iso);
    if (isNaN(date.getTime())) return '';
    const tzObj = findTimeZone(tz || 'UTC');
    const z = getZonedTime(date, tzObj);
    const yyyy = z.year;
    const MM = pad2(z.month);
    const dd = pad2(z.day);
    const HH = pad2(z.hours);
    const mm = pad2(z.minutes);
    return `${yyyy}-${MM}-${dd}T${HH}:${mm}`;
  } catch (_) {
    return '';
  }
};

const fromDatetimeLocalInZoneToISO = (val?: string, tz?: string): string | undefined => {
  if (!val) return undefined;
  try {
    const [datePart, timePart] = val.split('T');
    if (!datePart || !timePart) return undefined;
    const [y, m, d] = datePart.split('-').map((s) => parseInt(s, 10));
    const [hh, mm] = timePart.split(':').map((s) => parseInt(s, 10));
    if (!y || !m || !d || isNaN(hh) || isNaN(mm)) return undefined;
    const tzObj = findTimeZone(tz || 'UTC');
    // Build a Date from the provided wall-time components as if they were in UTC,
    // then subtract the timezone offset to get the real UTC instant.
    const wallAsUTC = new Date(Date.UTC(y, m - 1, d, hh, mm, 0, 0));
    const offsetMinutes = getUTCOffset(tzObj, wallAsUTC);
    const utcDate = new Date(wallAsUTC.getTime() - offsetMinutes * 60 * 1000);
    return utcDate.toISOString();
  } catch (_) {
    return undefined;
  }
};


const allTimezones = listTimeZones();

// Helper function to format duration for cron_duration field
const formatDuration = (duration: { value: number; unit: string }): string => {
  const { value, unit } = duration;
  
  switch (unit) {
    case 'minutes':
      return `${value}m`;
    case 'hours':
      return `${value}h`;
    case 'days':
      return `${value}d`;
    default:
      return `${value}${unit}`;
  }
};

// Helper function to format date for datetime-local input
const formatLocalDateTime = (date: Date): string => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day}T${hours}:${minutes}`;
};


const ProjectSchedulingPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const navigate = useNavigate();
  const qc = useQueryClient();
  const [environmentKey, setEnvironmentKey] = useState<string>('prod'); // Default to prod environment

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

  // Filters, sorting and pagination state for features
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
  const { data: featuresResp, isLoading: loadingFeatures } = useQuery<ListFeaturesResponse>({
    queryKey: ['project-features', projectId, { search: minSearch, enabledFilter, kindFilter, sortBy, sortOrder, page, perPage, selectedTags }],
    queryFn: async () => {
      const tagIds = selectedTags.length > 0 ? selectedTags.map(tag => tag.id).join(',') : undefined;
      const res = await apiClient.listProjectFeatures(
        projectId,
        environmentKey,
        kindFilter === 'all' ? undefined : kindFilter,
        enabledFilter === 'all' ? undefined : (enabledFilter === 'enabled'),
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

  // Timeline date range state
  const [timelineFrom, setTimelineFrom] = useState<string>(() => {
    const now = new Date();
    return formatLocalDateTime(now);
  });
  const [timelineTo, setTimelineTo] = useState<string>(() => {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    return formatLocalDateTime(tomorrow);
  });
  const [timelineError, setTimelineError] = useState<string>('');

  // Initialize validation on component mount
  React.useEffect(() => {
    setTimelineError(validateTimelineRange(timelineFrom, timelineTo));
  }, []);

  // Validate timeline date range (max 1 week)
  const validateTimelineRange = (from: string, to: string): string => {
    if (!from || !to) return '';
    
    const fromDate = new Date(from);
    const toDate = new Date(to);
    
    if (isNaN(fromDate.getTime()) || isNaN(toDate.getTime())) {
      return 'Invalid date format';
    }
    
    if (fromDate >= toDate) {
      return 'From date must be before To date';
    }
    
    const diffMs = toDate.getTime() - fromDate.getTime();
    const diffDays = diffMs / (1000 * 60 * 60 * 24);
    
    if (diffDays > 7) {
      return 'Maximum time range is 1 week (7 days)';
    }
    
    return '';
  };

  const { data: allSchedules, isLoading: loadingSchedules } = useQuery<FeatureSchedule[]>({
    queryKey: ['feature-schedules', projectId],
    queryFn: async () => {
      const res = await apiClient.listAllFeatureSchedules();
      // Filter by project just in case API returns global list
      return (res.data || []).filter((s: FeatureSchedule) => s.project_id === projectId);
    },
    enabled: !!projectId,
  });

  // Timeline data for selected features
  const selectedFeatures = useMemo(() => {
    return features || [];
  }, [features]);

  const { data: timelinesData, isLoading: loadingTimelines, error: timelinesError, refetch: refetchTimelines } = useQuery<Record<string, FeatureTimelineEvent[]>>({
    queryKey: ['feature-timelines', projectId, environmentKey, selectedFeatures.map(f => f.id), timelineFrom, timelineTo, timelineError],
    queryFn: async () => {
      if (selectedFeatures.length === 0) return {};
      
      // Don't fetch if there's a validation error
      if (timelineError) return {};

      // Convert local datetime to ISO string for API
      const from = new Date(timelineFrom).toISOString();
      const to = new Date(timelineTo).toISOString();

      const timelinePromises = selectedFeatures.map(async (feature) => {
        try {
          // Get user's timezone
          const location = Intl.DateTimeFormat().resolvedOptions().timeZone;
          const res = await apiClient.getFeatureTimeline(feature.id, environmentKey, from, to, location);
          return { featureId: feature.id, events: res.data.events };
        } catch (error) {
          console.error(`Failed to load timeline for feature ${feature.id}:`, error);
          return { featureId: feature.id, events: [] };
        }
      });

      const results = await Promise.all(timelinePromises);
      
      const timelines: Record<string, FeatureTimelineEvent[]> = {};
      results.forEach(({ featureId, events }) => {
        timelines[featureId] = events;
      });

      return timelines;
    },
    enabled: !!projectId && selectedFeatures.length > 0
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

  // Function to test timeline with mock schedules
  const testTimelineWithSchedules = async (featureId: string, schedules: any[], from: Date, to: Date) => {
    try {
      const response = await apiClient.testFeatureTimeline(
        featureId,
        environmentKey,
        from.toISOString(),
        to.toISOString(),
        Intl.DateTimeFormat().resolvedOptions().timeZone,
        {
          schedules: schedules.map(schedule => ({
            starts_at: schedule.startsAt ? new Date(schedule.startsAt).toISOString() : undefined,
            ends_at: schedule.endsAt ? new Date(schedule.endsAt).toISOString() : undefined,
            cron_expr: schedule.cronExpr,
            timezone: schedule.timezone,
            action: schedule.action,
            cron_duration: schedule.cronDuration
          }))
        }
      );
      return response.data;
    } catch (error) {
      console.error('Failed to test timeline:', error);
      throw error;
    }
  };

  // Separate features into two groups
  const { featuresWithSchedules, featuresWithoutSchedules } = useMemo(() => {
    const list = features || [];

    const withSchedules: FeatureExtended[] = [];
    const withoutSchedules: FeatureExtended[] = [];

    list.forEach(feature => {
      if (schedulesByFeature[feature.id] && schedulesByFeature[feature.id].length > 0) {
        withSchedules.push(feature);
      } else {
        withoutSchedules.push(feature);
      }
    });

    return { featuresWithSchedules: withSchedules, featuresWithoutSchedules: withoutSchedules };
  }, [features, schedulesByFeature]);

  // Dialog state
  const [dialogOpen, setDialogOpen] = useState(false);
  const [dialogFeature, setDialogFeature] = useState<FeatureExtended | null>(null);
  const [editSchedule, setEditSchedule] = useState<FeatureSchedule | null>(null);
  
  // New schedule builder dialogs
  const [scheduleBuilderOpen, setScheduleBuilderOpen] = useState(false);
  const [oneShotDialogOpen, setOneShotDialogOpen] = useState(false);
  const [scheduleType, setScheduleType] = useState<'cron' | 'one-shot' | null>(null);
  
  // Edit dialogs
  const [editRecurringBuilderOpen, setEditRecurringBuilderOpen] = useState(false);
  const [editOneShotDialogOpen, setEditOneShotDialogOpen] = useState(false);
  const [editingSchedule, setEditingSchedule] = useState<FeatureSchedule | null>(null);
  const [helpDialogOpen, setHelpDialogOpen] = useState(false);
  
  // Tab state
  const [activeTab, setActiveTab] = useState(0);

  // Auto-switch tabs based on data changes
  useEffect(() => {
    if (featuresWithSchedules.length > 0 && activeTab === 1) {
      // If we're on "without schedules" tab but there are features with schedules,
      // and the current feature list is empty, switch to "with schedules" tab
      if (featuresWithoutSchedules.length === 0) {
        setActiveTab(0);
      }
    }
  }, [featuresWithSchedules.length, featuresWithoutSchedules.length, activeTab]);

  const openCreate = (feature: FeatureExtended) => {
    setDialogFeature(feature);
    setEditSchedule(null);
    setDialogOpen(true);
  };
  
  const openEdit = (feature: FeatureExtended, schedule: FeatureSchedule) => {
    setDialogFeature(feature);
    setEditingSchedule(schedule);
    
    // Determine which edit dialog to open based on schedule type
    if (schedule.cron_expr) {
      setEditRecurringBuilderOpen(true);
    } else {
      setEditOneShotDialogOpen(true);
    }
  };

  const openScheduleBuilder = (feature: FeatureExtended) => {
    setDialogFeature(feature);
    setScheduleType('cron');
    setScheduleBuilderOpen(true);
  };

  const openOneShotDialog = (feature: FeatureExtended) => {
    setDialogFeature(feature);
    setScheduleType('one-shot');
    setOneShotDialogOpen(true);
  };

  const closeDialog = () => setDialogOpen(false);
  const closeScheduleBuilder = () => setScheduleBuilderOpen(false);
  const closeOneShotDialog = () => setOneShotDialogOpen(false);
  const closeEditRecurringBuilder = () => setEditRecurringBuilderOpen(false);
  const closeEditOneShotDialog = () => setEditOneShotDialogOpen(false);

  // Mutations
  const createMut = useMutation({
    mutationFn: async ({ featureId, values }: { featureId: string; values: ScheduleFormValues }) => {
      return apiClient.createFeatureSchedule(featureId, environmentKey, values);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
      // Switch to "with schedules" tab after creating a schedule
      setActiveTab(0);
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
      // After deleting, check if we need to switch tabs
      // This will be handled by the useMemo that recalculates the feature lists
    }
  });

  // New mutations for schedule builder
  const createCronScheduleMut = useMutation({
    mutationFn: async ({ featureId, data }: { featureId: string; data: ScheduleBuilderData & { cronExpression: string } }) => {
      const payload = {
        timezone: data.timezone,
        starts_at: data.startsAt,
        ends_at: data.endsAt,
        action: data.action,
        cron_expr: data.cronExpression,
        cron_duration: formatDuration(data.duration)
      };
      return apiClient.createFeatureSchedule(featureId, environmentKey, payload);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
      setScheduleBuilderOpen(false);
      setActiveTab(0);
    }
  });

  const createOneShotScheduleMut = useMutation({
    mutationFn: async ({ featureId, data }: { featureId: string; data: { timezone: string; starts_at: string; ends_at: string; action: FeatureScheduleAction } }) => {
      return apiClient.createFeatureSchedule(featureId, environmentKey, data);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
      setOneShotDialogOpen(false);
      setActiveTab(0);
    }
  });

  // Edit mutations
  const editRecurringScheduleMut = useMutation({
    mutationFn: async ({ scheduleId, data }: { scheduleId: string; data: ScheduleBuilderData & { cronExpression: string } }) => {
      const payload = {
        timezone: data.timezone,
        starts_at: data.startsAt,
        ends_at: data.endsAt,
        action: data.action,
        cron_expr: data.cronExpression,
        cron_duration: formatDuration(data.duration)
      };
      return apiClient.updateFeatureSchedule(scheduleId, payload);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
      setEditRecurringBuilderOpen(false);
    }
  });

  const editOneShotScheduleMut = useMutation({
    mutationFn: async ({ scheduleId, data }: { scheduleId: string; data: { timezone: string; starts_at: string; ends_at: string; action: FeatureScheduleAction } }) => {
      return apiClient.updateFeatureSchedule(scheduleId, data);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['feature-schedules', projectId] });
      setEditOneShotDialogOpen(false);
    }
  });

  const project = projectResp?.project;

  // Component to render feature list
  const renderFeatureList = (featureList: FeatureExtended[]) => (
    <Box>
      {featureList.map((f) => (
        <Accordion key={f.id} defaultExpanded sx={{ mb: 2 }}>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%', justifyContent: 'space-between' }}>
              <Box>
                <Typography variant="subtitle1">{f.name}</Typography>
                <Typography variant="body2" color="text.secondary">{f.key}</Typography>
                <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                  <Chip size="small" label={`kind: ${f.kind}`} />
                  <Chip size="small" label={`default: ${f.default_value}`} />
                  <Chip size="small" label={f.is_active ? 'active' : 'not active'} color={f.is_active ? 'success' : 'default'} />
                </Box>
              </Box>
              {canAddRecurringSchedule(schedulesByFeature[f.id] || []) ? (
                <Box sx={{ display: 'flex', gap: 1 }}>
                  <Button 
                    variant="contained" 
                    startIcon={<AddIcon />} 
                    onClick={(e) => { e.stopPropagation(); openScheduleBuilder(f); }}
                    size="small"
                  >
                    Recurring
                  </Button>
                  <Button 
                    variant="outlined" 
                    startIcon={<AddIcon />} 
                    onClick={(e) => { e.stopPropagation(); openOneShotDialog(f); }}
                    size="small"
                  >
                    One-shot
                  </Button>
                </Box>
              ) : canAddOneShotSchedule(schedulesByFeature[f.id] || []) ? (
                <Button 
                  variant="outlined" 
                  startIcon={<AddIcon />} 
                  onClick={(e) => { e.stopPropagation(); openOneShotDialog(f); }}
                  size="small"
                >
                  Add One-shot
                </Button>
              ) : (
                <Chip 
                  size="small" 
                  label="Recurring schedule exists" 
                  color="warning" 
                  variant="outlined"
                  icon={<ScheduleIcon />}
                />
              )}
            </Box>
          </AccordionSummary>
          <AccordionDetails>
            <Box>
              {(schedulesByFeature[f.id] && schedulesByFeature[f.id].length > 0) ? (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1.5 }}>
                  {schedulesByFeature[f.id].map((s) => (
                    <Paper key={s.id} sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                      <Box>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                          <Typography variant="body1" sx={{ fontWeight: 600, textTransform: 'capitalize' }}>
                            {s.action === 'enable' ? 'Activate' : s.action === 'disable' ? 'Deactivate' : s.action}
                          </Typography>
                          <Chip 
                            size="small" 
                            label={getScheduleType(s) === 'cron' ? 'Recurring' : 'One-shot'} 
                            color={getScheduleType(s) === 'cron' ? 'primary' : 'secondary'}
                            variant="outlined"
                          />
                        </Box>
                        <Typography variant="body2" color="text.secondary">
                          {s.cron_expr ? (() => {
                            try {
                              return cronstrue.toString(s.cron_expr);
                            } catch (error) {
                              return `Cron: ${s.cron_expr}`;
                            }
                          })() : `From ${s.starts_at || '—'} to ${s.ends_at || '—'}`}
                        </Typography>
                        {s.cron_expr && s.cron_duration && (
                          <Typography variant="body2" color="text.secondary">
                            Duration: {s.cron_duration}
                          </Typography>
                        )}
                        <Typography variant="caption" color="text.secondary">Timezone: {s.timezone}</Typography>
                      </Box>
                      <Box>
                        <Tooltip title="Edit schedule">
                          <IconButton onClick={() => openEdit(f, s)} size="small"><EditIcon /></IconButton>
                        </Tooltip>
                        <Tooltip title="Delete schedule">
                          <IconButton
                            color="error"
                            onClick={() => {
                              const confirmed = window.confirm('Delete the schedule? This action is irreversible.');
                              if (confirmed) {
                                deleteMut.mutate({ scheduleId: s.id });
                              }
                            }}
                          >
                            <DeleteIcon />
                          </IconButton>
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
  );

  return (
    <AuthenticatedLayout showBackButton backTo={`/projects/${projectId}`}>
      {/* Help Link - positioned as part of subtitle */}
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'flex-start',
        mb: 2, 
        mt: 2,
        px: 2
      }}>
        <Box
          component="span"
          onClick={() => setHelpDialogOpen(true)}
          sx={{
            display: 'inline-flex',
            alignItems: 'center',
            gap: 0.5,
            color: 'primary.main',
            cursor: 'pointer',
            textDecoration: 'underline',
            fontSize: '0.875rem',
            '&:hover': {
              color: 'primary.dark',
              textDecoration: 'underline'
            }
          }}
        >
          <HelpIcon fontSize="small" />
          Understanding Feature Enablement and Schedules
        </Box>
      </Box>

      {(loadingProject || loadingFeatures || loadingSchedules) && (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      )}

      {!loadingFeatures && features && features.length > 0 ? (
        <Box>
          {/* Environment selector */}
          <Box sx={{ display: 'flex', gap: 2, alignItems: 'center', mb: 2 }}>
            <FormControl size="small" sx={{ minWidth: 200 }}>
              <InputLabel>Environment</InputLabel>
              <Select
                label="Environment"
                value={environmentKey}
                onChange={(e) => setEnvironmentKey(e.target.value)}
                disabled={loadingEnvironments}
              >
                {environments.map((env: any) => (
                  <MenuItem key={env.id} value={env.key}>
                    {env.name} ({env.key})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>

          {/* Search and filters */}
          <SearchPanel
            searchValue={search}
            onSearchChange={(value) => { setSearch(value); setPage(1); }}
            placeholder="Search by name or key"
            projectId={projectId}
            selectedTags={selectedTags}
            onTagsChange={(tags) => { setSelectedTags(tags); setPage(1); }}
            showTagFilter={true}
            filters={[
              {
                key: 'enabledFilter',
                label: 'Enabled',
                value: enabledFilter,
                options: [
                  { value: 'all', label: 'All' },
                  { value: 'enabled', label: 'Enabled' },
                  { value: 'disabled', label: 'Disabled' },
                ],
                onChange: (value) => { setEnabledFilter(value as any); setPage(1); },
              },
              {
                key: 'kindFilter',
                label: 'Kind',
                value: kindFilter,
                options: [
                  { value: 'all', label: 'All' },
                  { value: 'simple', label: 'Simple' },
                  { value: 'multivariant', label: 'Multivariant' },
                ],
                onChange: (value) => { setKindFilter(value as any); setPage(1); },
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

          <Tabs value={activeTab} onChange={(_, newValue) => setActiveTab(newValue)} sx={{ mb: 2 }}>
            <Tab 
              label={`Features with schedules (${featuresWithSchedules.length})`} 
              sx={{ textTransform: 'none' }}
            />
            <Tab 
              label={`Features without schedules (${featuresWithoutSchedules.length})`} 
              sx={{ textTransform: 'none' }}
            />
            <Tab 
              label={`Timeline (${selectedFeatures.length})`} 
              sx={{ textTransform: 'none' }}
            />
          </Tabs>
          
          {activeTab === 0 && (
            <Box>
              {featuresWithSchedules.length > 0 ? (
                renderFeatureList(featuresWithSchedules)
              ) : (
                <Typography variant="body2" color="text.secondary">
                  No features with schedules yet. Add schedules to features in the second tab.
                </Typography>
              )}
            </Box>
          )}
          
          {activeTab === 1 && (
            <Box>
              {featuresWithoutSchedules.length > 0 ? (
                renderFeatureList(featuresWithoutSchedules)
              ) : (
                <Typography variant="body2" color="text.secondary">
                  All features have schedules! Great job organizing your feature schedules.
                </Typography>
              )}
            </Box>
          )}

          {activeTab === 2 && (
            <Box>
              {/* Date range selector for timeline */}
              <Paper sx={{ p: 2, mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="h6">
                    Timeline Settings
                  </Typography>
                  <Tooltip title="Refresh timeline data">
                    <IconButton
                      onClick={() => {
                        // Refetch timeline data
                        refetchTimelines();
                      }}
                      disabled={loadingTimelines}
                      size="small"
                    >
                      <RefreshIcon />
                    </IconButton>
                  </Tooltip>
                </Box>
                <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} sx={{ mb: 2 }}>
                  <TextField
                    label="From"
                    type="datetime-local"
                    size="small"
                    value={timelineFrom}
                    onChange={(e) => {
                      const newFrom = e.target.value;
                      setTimelineFrom(newFrom);
                      setTimelineError(validateTimelineRange(newFrom, timelineTo));
                    }}
                    InputLabelProps={{ shrink: true }}
                    sx={{ minWidth: 200 }}
                  />
                  <TextField
                    label="To"
                    type="datetime-local"
                    size="small"
                    value={timelineTo}
                    onChange={(e) => {
                      const newTo = e.target.value;
                      setTimelineTo(newTo);
                      setTimelineError(validateTimelineRange(timelineFrom, newTo));
                    }}
                    InputLabelProps={{ shrink: true }}
                    sx={{ minWidth: 200 }}
                  />
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={() => {
                      const now = new Date();
                      const tomorrow = new Date();
                      tomorrow.setDate(tomorrow.getDate() + 1);
                      
                      const newFrom = formatLocalDateTime(now);
                      const newTo = formatLocalDateTime(tomorrow);
                      
                      setTimelineFrom(newFrom);
                      setTimelineTo(newTo);
                      setTimelineError(validateTimelineRange(newFrom, newTo));
                    }}
                    sx={{ alignSelf: 'flex-start' }}
                  >
                    Reset to Now + 1 Day
                  </Button>
                </Stack>
                {timelineError && (
                  <Typography variant="body2" color="error" sx={{ mt: 1 }}>
                    {timelineError}
                  </Typography>
                )}
                <Typography variant="body2" color="text.secondary">
                  Select the time range to view feature timelines. Data will be loaded for all selected features.
                  <br />
                  <strong>Maximum range:</strong> 1 week (7 days)
                  <br />
                  <strong>Timezone:</strong> {Intl.DateTimeFormat().resolvedOptions().timeZone}
                </Typography>
              </Paper>

              <TimelineChart
                features={selectedFeatures}
                timelines={timelinesData || {}}
                isLoading={loadingTimelines}
                error={timelinesError?.message}
                from={timelineFrom}
                to={timelineTo}
              />
            </Box>
          )}

          {/* Pagination */}
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mt: 1.5 }}>
            <Typography variant="body2" color="text.secondary" sx={{ fontSize: '0.8rem' }}>
              {pagination ? `Total: ${pagination.total}` : ''}
            </Typography>
            <Pagination
              size="small"
              page={page}
              count={pagination ? Math.max(1, Math.ceil(pagination.total / (pagination.per_page || perPage))) : 1}
              onChange={(_e, p) => setPage(p)}
              shape="rounded"
              color="primary"
            />
          </Box>
        </Box>
      ) : !loadingFeatures ? (
        <Box>
          {/* Search and filters */}
          <SearchPanel
            searchValue={search}
            onSearchChange={(value) => { setSearch(value); setPage(1); }}
            placeholder="Search by name or key"
            projectId={projectId}
            selectedTags={selectedTags}
            onTagsChange={(tags) => { setSelectedTags(tags); setPage(1); }}
            showTagFilter={true}
            filters={[
              {
                key: 'enabledFilter',
                label: 'Enabled',
                value: enabledFilter,
                options: [
                  { value: 'all', label: 'All' },
                  { value: 'enabled', label: 'Enabled' },
                  { value: 'disabled', label: 'Disabled' },
                ],
                onChange: (value) => { setEnabledFilter(value as any); setPage(1); },
              },
              {
                key: 'kindFilter',
                label: 'Kind',
                value: kindFilter,
                options: [
                  { value: 'all', label: 'All' },
                  { value: 'simple', label: 'Simple' },
                  { value: 'multivariant', label: 'Multivariant' },
                ],
                onChange: (value) => { setKindFilter(value as any); setPage(1); },
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
          <Typography variant="body2">No features yet.</Typography>
        </Box>
      ) : null}

      {/* Create/Edit Dialog */}
      <ScheduleDialog
        open={dialogOpen}
        onClose={closeDialog}
        title={editSchedule ? 'Edit schedule' : 'Create schedule'}
        initial={editSchedule ? {
          action: editSchedule.action,
          timezone: editSchedule.timezone,
          // Pass raw ISO; dialog will convert according to timezone
          starts_at: editSchedule.starts_at || '',
          ends_at: editSchedule.ends_at || '',
          cron_expr: editSchedule.cron_expr || '',
          cron_duration: editSchedule.cron_duration || ''
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

      {/* Schedule Builder Dialog */}
      <Dialog 
        open={scheduleBuilderOpen} 
        onClose={closeScheduleBuilder}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle>Create recurring schedule</DialogTitle>
        <DialogContent>
          <ScheduleBuilder
            open={scheduleBuilderOpen}
            featureId={dialogFeature?.id || ''}
            environmentKey={environmentKey}
            onSubmit={(data) => {
              if (dialogFeature) {
                createCronScheduleMut.mutate({ featureId: dialogFeature.id, data });
              }
            }}
            featureCreatedAt={dialogFeature?.created_at}
          />
        </DialogContent>
      </Dialog>

      {/* One-Shot Schedule Dialog */}
      <OneShotScheduleDialog
        open={oneShotDialogOpen}
        onClose={closeOneShotDialog}
        onSubmit={(data) => {
          if (dialogFeature) {
            createOneShotScheduleMut.mutate({ featureId: dialogFeature.id, data });
          }
        }}
        existingSchedules={allSchedules?.filter(s => s.feature_id === dialogFeature?.id) || []}
      />

      {/* Edit Recurring Schedule Builder */}
      <Dialog 
        open={editRecurringBuilderOpen} 
        onClose={closeEditRecurringBuilder}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle>Edit recurring schedule</DialogTitle>
        <DialogContent>
          <EditRecurringScheduleBuilder
            open={editRecurringBuilderOpen}
            featureId={editingSchedule?.feature_id || ''}
            environmentKey={environmentKey}
            onSubmit={(data) => {
              if (editingSchedule) {
                editRecurringScheduleMut.mutate({ scheduleId: editingSchedule.id, data });
              }
            }}
            initialData={editingSchedule || undefined}
          />
        </DialogContent>
      </Dialog>

      {/* Edit One-Shot Schedule Dialog */}
      <EditOneShotScheduleDialog
        open={editOneShotDialogOpen}
        onClose={closeEditOneShotDialog}
        onSubmit={(data) => {
          if (editingSchedule) {
            editOneShotScheduleMut.mutate({ scheduleId: editingSchedule.id, data });
          }
        }}
        initialData={editingSchedule || undefined}
        existingSchedules={allSchedules?.filter(s => s.feature_id === dialogFeature?.id) || []}
      />

      {/* Help Dialog */}
      <ScheduleHelpDialog
        open={helpDialogOpen}
        onClose={() => setHelpDialogOpen(false)}
      />
    </AuthenticatedLayout>
  );
};

export default ProjectSchedulingPage;
