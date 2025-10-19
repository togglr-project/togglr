import React, { useMemo, useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  TextField,
  Button,
  CircularProgress,
  IconButton,
  Tooltip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Tabs,
  Tab
} from '@mui/material';
import type { AlertColor } from '@mui/material';
import {
  Settings as SettingsIcon,
  Visibility as VisibilityIcon,
  VisibilityOff as VisibilityOffIcon,
  ContentCopy as ContentCopyIcon,
  DeleteOutline as DeleteIcon,
  Save as SaveIcon,
  Notifications as NotificationsIcon,
  Tune as TuneIcon
} from '@mui/icons-material';
import { useParams, useNavigate } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import NotificationSettings from '../components/NotificationSettings';
import ProjectSettingsList from '../components/projects/ProjectSettingsList';
import apiClient from '../api/apiClient';
import type { Project } from '../generated/api/client';
import { useNotification } from '../App';
import { useRBAC } from '../auth/permissions';

interface ProjectResponse { project: Project }

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

interface ApiKeysSectionProps {
  projectId?: string;
  showNotification: (message: string, type: AlertColor, duration?: number) => void;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`settings-tabpanel-${index}`}
      aria-labelledby={`settings-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 3 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const ApiKeysSection: React.FC<ApiKeysSectionProps> = ({ projectId = '', showNotification }) => {
  const { data: envResp, isLoading: loadingEnvs, error: envError } = useQuery({
    queryKey: ['project-environments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectEnvironments(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  const environments = envResp?.items ?? [];
  const [visible, setVisible] = useState<Record<number, boolean>>({});

  const toggleVisible = (id: number) => setVisible(v => ({ ...v, [id]: !v[id] }));

  const copy = async (value: string) => {
    try {
      await navigator.clipboard.writeText(value);
      showNotification('API key copied to clipboard', 'success', 2000);
    } catch {
      showNotification('Failed to copy API key', 'error', 3000);
    }
  };

  return (
    <Box>
      <Typography variant="subtitle2" color="text.secondary">API keys</Typography>
      {loadingEnvs && (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 1 }}>
          <CircularProgress size={18} />
          <Typography variant="body2" color="text.secondary">Loading environments…</Typography>
        </Box>
      )}
      {envError && (
        <Typography variant="body2" color="error" sx={{ mt: 1 }}>Failed to load environments.</Typography>
      )}
      {!loadingEnvs && environments.length === 0 && (
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>No environments found.</Typography>
      )}

      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mt: 1 }}>
        {environments.map(env => (
          <Paper key={env.id} variant="outlined" sx={{ p: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box sx={{ minWidth: 160 }}>
              <Typography variant="body2" sx={{ fontWeight: 500 }}>{env.name}</Typography>
              <Typography variant="caption" color="text.secondary">{env.key}</Typography>
            </Box>
            <Box sx={{ flexGrow: 1, fontFamily: 'monospace', bgcolor: 'action.hover', px: 1.5, py: 0.75, borderRadius: 1, overflowX: 'auto' }}>
              {visible[env.id] ? env.api_key : '•'.repeat(Math.max(12, env.api_key?.length || 16))}
            </Box>
            <Tooltip title={visible[env.id] ? 'Hide' : 'Show'}>
              <IconButton size="small" onClick={() => toggleVisible(env.id)}>
                {visible[env.id] ? <VisibilityOffIcon fontSize="small" /> : <VisibilityIcon fontSize="small" />}
              </IconButton>
            </Tooltip>
            <Tooltip title="Copy">
              <IconButton size="small" onClick={() => copy(env.api_key)}>
                <ContentCopyIcon fontSize="small" />
              </IconButton>
            </Tooltip>
          </Paper>
        ))}
      </Box>
    </Box>
  );
};

const ProjectSettingsPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const navigate = useNavigate();
  const qc = useQueryClient();
  const { showNotification } = useNotification();
  const rbac = useRBAC(projectId);

  const { data, isLoading, error } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const res = await apiClient.getProject(projectId);
      return res.data as ProjectResponse;
    },
    enabled: !!projectId,
  });

  // Get environments for the project
  const { data: envResp, isLoading: loadingEnvs } = useQuery({
    queryKey: ['project-environments', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectEnvironments(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  const environments = envResp?.items ?? [];
  const [environmentKey, setEnvironmentKey] = useState<string>(() => {
    return localStorage.getItem('currentEnvironmentKey') || 'prod';
  });

  const project = data?.project;

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [tabValue, setTabValue] = useState(0);

  React.useEffect(() => {
    if (project) {
      setName(project.name);
      setDescription(project.description || '');
    }
  }, [project?.id]);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const changed = useMemo(() => {
    return project ? (name !== project.name || (description || '') !== (project.description || '')) : false;
  }, [name, description, project]);

  const saveMut = useMutation({
    mutationFn: async () => {
      if (!projectId) return;
      await apiClient.updateProject(projectId, { name, description });
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['project', projectId] });
      showNotification('Project updated', 'success', 3000);
    },
    onError: () => {
      showNotification('Failed to update project', 'error', 5000);
    }
  });

  const deleteMut = useMutation({
    mutationFn: async () => {
      if (!projectId) return;
      await apiClient.archiveProject(projectId);
    },
    onSuccess: () => {
      showNotification('Project deleted', 'success', 3000);
      navigate('/projects');
    },
    onError: () => {
      showNotification('Failed to delete project', 'error', 5000);
    }
  });

  return (
    <AuthenticatedLayout showBackButton backTo={`/projects/${projectId}`}>
      <PageHeader
        title={project ? `${project.name} — Settings` : 'Project Settings'}
        subtitle={project ? `Project ID: ${project.id}` : 'Edit project'}
        icon={<SettingsIcon />}
      />

      {isLoading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      )}
      {error && (
        <Typography color="error">Failed to load project.</Typography>
      )}

      {project && (
        <Paper sx={{ width: '100%' }}>
          <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Tabs value={tabValue} onChange={handleTabChange} aria-label="project settings tabs">
              <Tab 
                icon={<SettingsIcon />} 
                label="General" 
                id="settings-tab-0"
                aria-controls="settings-tabpanel-0"
              />
              <Tab 
                icon={<TuneIcon />} 
                label="Settings" 
                id="settings-tab-1"
                aria-controls="settings-tabpanel-1"
              />
              <Tab 
                icon={<NotificationsIcon />} 
                label="Feature Notifications"
                id="settings-tab-2"
                aria-controls="settings-tabpanel-2"
              />
            </Tabs>
          </Box>

          <TabPanel value={tabValue} index={0}>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
              <Box
                sx={{
                  display: 'grid',
                  gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' },
                  gap: 3,
                  alignItems: 'start'
                }}
              >
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <TextField
                    label="Project name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    fullWidth
                  />
                  <TextField
                    label="Description"
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    fullWidth
                    multiline
                    minRows={3}
                  />
                </Box>

                <Box>
                  {/* API Keys per environment */}
                  <ApiKeysSection projectId={projectId} showNotification={showNotification} />
                </Box>
              </Box>

              <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
                {rbac.canManageProject() && (
                  <Button
                    color="error"
                    startIcon={<DeleteIcon />}
                    onClick={() => setConfirmOpen(true)}
                  >
                    Delete project
                  </Button>
                )}

                {rbac.canManageProject() && (
                  <Button
                    variant="contained"
                    startIcon={<SaveIcon />}
                    disabled={!changed || saveMut.isPending}
                    onClick={() => saveMut.mutate()}
                  >
                    Save changes
                  </Button>
                )}
              </Box>
            </Box>
          </TabPanel>

          <TabPanel value={tabValue} index={1}>
            <ProjectSettingsList projectId={projectId} />
          </TabPanel>

          <TabPanel value={tabValue} index={2}>
            <NotificationSettings 
              projectId={projectId} 
              environmentKey={environmentKey}
              environments={environments}
              loadingEnvironments={loadingEnvs}
            />
          </TabPanel>
        </Paper>
      )}

      <Dialog open={confirmOpen} onClose={() => setConfirmOpen(false)}>
        <DialogTitle>Delete project</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Are you sure you want to delete this project? This action cannot be undone.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmOpen(false)} size="small">Cancel</Button>
          <Button color="error" onClick={() => { setConfirmOpen(false); deleteMut.mutate(); }} size="small">Delete</Button>
        </DialogActions>
      </Dialog>
    </AuthenticatedLayout>
  );
};

export default ProjectSettingsPage;
