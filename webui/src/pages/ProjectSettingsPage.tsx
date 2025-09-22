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
  DialogActions
} from '@mui/material';
import {
  Settings as SettingsIcon,
  Visibility as VisibilityIcon,
  VisibilityOff as VisibilityOffIcon,
  ContentCopy as ContentCopyIcon,
  DeleteOutline as DeleteIcon,
  Save as SaveIcon
} from '@mui/icons-material';
import { useParams, useNavigate } from 'react-router-dom';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import type { Project } from '../generated/api/client';
import { useNotification } from '../App';

interface ProjectResponse { project: Project }

const ProjectSettingsPage: React.FC = () => {
  const { projectId = '' } = useParams();
  const navigate = useNavigate();
  const qc = useQueryClient();
  const { showNotification } = useNotification();

  const { data, isLoading, error } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const res = await apiClient.getProject(projectId);
      return res.data as ProjectResponse;
    },
    enabled: !!projectId,
  });

  const project = data?.project;

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [showApiKey, setShowApiKey] = useState(false);
  const [copied, setCopied] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  React.useEffect(() => {
    if (project) {
      setName(project.name);
      setDescription(project.description || '');
    }
  }, [project?.id]);

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
        <Paper sx={{ p: 3, display: 'flex', flexDirection: 'column', gap: 3 }}>
          <Box>
            <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1 }}>API Key</Typography>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, flexWrap: 'wrap' }}>
              <Typography
                variant="body2"
                sx={{ fontFamily: 'monospace', userSelect: 'text', maxWidth: '100%', overflow: 'hidden', textOverflow: 'ellipsis' }}
              >
                {showApiKey ? project.api_key : '*'.repeat(project.api_key?.length || 8)}
              </Typography>
              <IconButton aria-label={showApiKey ? 'hide api key' : 'show api key'} size="small" onClick={() => setShowApiKey(v => !v)}>
                {showApiKey ? <VisibilityOffIcon fontSize="small" /> : <VisibilityIcon fontSize="small" />}
              </IconButton>
              <Tooltip title={copied ? 'Copied!' : 'Copy API Key'} placement="top" onClose={() => setCopied(false)}>
                <IconButton
                  aria-label="copy api key"
                  size="small"
                  onClick={async () => {
                    try {
                      if (project.api_key) {
                        await navigator.clipboard.writeText(project.api_key);
                        setCopied(true);
                        setTimeout(() => setCopied(false), 1200);
                      }
                    } catch (e) {
                      // ignore
                    }
                  }}
                >
                  <ContentCopyIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Box>
          </Box>

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

          <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
            <Button
              color="error"
              startIcon={<DeleteIcon />}
              onClick={() => setConfirmOpen(true)}
            >
              Delete project
            </Button>

            <Button
              variant="contained"
              startIcon={<SaveIcon />}
              disabled={!changed || saveMut.isPending}
              onClick={() => saveMut.mutate()}
            >
              Save changes
            </Button>
          </Box>
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
