import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  CircularProgress,
  Button,
} from '@mui/material';
import { Folder as ProjectsIcon, Add as AddIcon } from '@mui/icons-material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate, Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import CreateProjectDialog from '../components/admin/CreateProjectDialog';
import { userPermissions } from '../hooks/userPermissions';
import ProjectCard from '../components/projects/ProjectCard';

interface Project {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

const ProjectsPage: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const queryClient = useQueryClient();
  const { canCreateProjects } = userPermissions();
  const [createOpen, setCreateOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);

  const { data: projects, isLoading, error } = useQuery<Project[]>({
    queryKey: ['projects'],
    queryFn: async () => {
      const res = await apiClient.listProjects();
      return res.data;
    },
  });

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  const handleCreateProject = async (name: string, description: string) => {
    try {
      setCreating(true);
      setCreateError(null);
      await apiClient.addProject({ name, description });
      setCreateOpen(false);
      await queryClient.invalidateQueries({ queryKey: ['projects'] });
    } catch (e: any) {
      const message = e?.response?.data?.error?.message || 'Failed to create project';
      setCreateError(message);
    } finally {
      setCreating(false);
    }
  };

  return (
    <AuthenticatedLayout>
      <PageHeader
        title="Projects"
        subtitle="Browse and manage your projects."
        icon={<ProjectsIcon />}
        gradientVariant="default"
        subtitleGradientVariant="default"
      />

      <Paper
        sx={{
          p: 3,
          background: (theme) =>
            theme.palette.mode === 'dark'
              ? 'linear-gradient(to bottom, rgba(65, 68, 74, 0.5), rgba(55, 58, 64, 0.5))'
              : 'linear-gradient(to bottom, rgba(255, 255, 255, 0.9), rgba(245, 245, 245, 0.9))',
          backdropFilter: 'blur(10px)',
          boxShadow: '0 4px 20px 0 rgba(0, 0, 0, 0.05)'
        }}
      >
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
          <Typography variant="h6" sx={{ color: 'primary.light' }}>
            Projects
          </Typography>
          {canCreateProjects() && (
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => setCreateOpen(true)}
              disabled={creating}
            >
              Add Project
            </Button>
          )}
        </Box>

        {isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Typography color="error">Failed to load projects.</Typography>
        ) : projects && projects.length > 0 ? (
          <Grid container spacing={2}>
            {projects.map((project: Project) => (
              <Grid item xs={12} sm={6} md={4} key={project.id}>
                <ProjectCard
                  id={project.id}
                  name={project.name}
                  description={project.description}
                  onClick={() => navigate(`/projects/${project.id}`)}
                />
              </Grid>
            ))}
          </Grid>
        ) : (
          <Typography variant="body2">No projects to display.</Typography>
        )}
      </Paper>

      {/* Create Project Dialog */}
      <CreateProjectDialog
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreateProject={handleCreateProject}
        isLoadingTeams={false}
      />
      {createError && (
        <Box sx={{ mt: 2 }}>
          <Typography color="error">{createError}</Typography>
        </Box>
      )}
    </AuthenticatedLayout>
  );
};

export default ProjectsPage;
