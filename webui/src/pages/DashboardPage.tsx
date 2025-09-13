import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  Card,
  CardContent,
  CardActionArea,
  CircularProgress,
  Chip,
  Button,
} from '@mui/material';
import { Dashboard as DashboardIcon, Add as AddIcon } from '@mui/icons-material';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate, Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import CreateProjectDialog from '../components/admin/CreateProjectDialog';
import { userPermissions } from '../hooks/userPermissions';

interface Project {
  'id': string;
  'name': string;
  'description': string;
  'created_at': string;
}

const DashboardPage: React.FC = () => {
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
      // Keep dialog open to let user adjust input
    } finally {
      setCreating(false);
    }
  };

  return (
    <AuthenticatedLayout>
      <PageHeader
        title="Dashboard"
        subtitle="Overview of your projects."
        icon={<DashboardIcon />}
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
          <Typography variant="h6" className="gradient-subtitle">
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
                <Card
                  sx={{
                    background: (theme) =>
                      theme.palette.mode === 'dark'
                        ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
                        : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
                    backdropFilter: 'blur(8px)',
                    boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)',
                    transition: 'all 0.2s ease-in-out',
                    '&:hover': {
                      background: (theme) =>
                        theme.palette.mode === 'dark'
                          ? 'linear-gradient(135deg, rgba(65, 68, 75, 0.7) 0%, rgba(60, 63, 70, 0.7) 100%)'
                          : 'linear-gradient(135deg, rgba(255, 255, 255, 1) 0%, rgba(250, 250, 250, 1) 100%)',
                      boxShadow: '0 5px 15px 0 rgba(0, 0, 0, 0.1)',
                      transform: 'translateY(-3px)'
                    }
                  }}
                >
                  <CardActionArea onClick={() => navigate(`/projects/${project.id}`)}>
                    <CardContent>
                      <Typography variant="h6" component="div">
                        {project.name}
                      </Typography>
                      <Box sx={{ mt: 1, display: 'flex', gap: 1 }}>
                        <Chip label={`ID: ${project.id}`} size="small" variant="outlined" />
                      </Box>
                    </CardContent>
                  </CardActionArea>
                </Card>
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

export default DashboardPage;
