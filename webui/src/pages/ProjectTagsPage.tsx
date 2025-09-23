import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  CircularProgress,
  Button,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Alert,
  Card,
  CardContent,
  CardActions,
  Tooltip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Autocomplete,
} from '@mui/material';
import { 
  LocalOffer as TagIcon, 
  Add as AddIcon, 
  Edit as EditIcon, 
  Delete as DeleteIcon,
  Category as CategoryIcon
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams, Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import { userPermissions } from '../hooks/userPermissions';
import type { ProjectTag, Category, CreateProjectTagRequest, UpdateProjectTagRequest } from '../generated/api/client';

const ProjectTagsPage: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const { user } = useAuth();
  const { canManageProject } = userPermissions();
  const queryClient = useQueryClient();
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [selectedTag, setSelectedTag] = useState<ProjectTag | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    slug: '',
    description: '',
    color: '#3B82F6',
    category_id: '',
  });
  const [error, setError] = useState<string | null>(null);

  if (!projectId) {
    return <Navigate to="/projects" replace />;
  }

  const { data: tags, isLoading, error: queryError } = useQuery<ProjectTag[]>({
    queryKey: ['project-tags', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectTags(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  const { data: categories } = useQuery<Category[]>({
    queryKey: ['categories'],
    queryFn: async () => {
      const res = await apiClient.listCategories();
      return res.data;
    },
  });

  const createMutation = useMutation({
    mutationFn: async (data: CreateProjectTagRequest) => {
      const res = await apiClient.createProjectTag(projectId, data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project-tags', projectId] });
      setCreateOpen(false);
      setFormData({ name: '', slug: '', description: '', color: '#3B82F6', category_id: '' });
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to create tag');
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateProjectTagRequest }) => {
      const res = await apiClient.updateProjectTag(projectId, id, data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project-tags', projectId] });
      setEditOpen(false);
      setSelectedTag(null);
      setFormData({ name: '', slug: '', description: '', color: '#3B82F6', category_id: '' });
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to update tag');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.deleteProjectTag(projectId, id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['project-tags', projectId] });
      setDeleteOpen(false);
      setSelectedTag(null);
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to delete tag');
    },
  });

  const canManage = user?.is_superuser || (projectId && canManageProject(parseInt(projectId)));

  const handleCreate = () => {
    if (!formData.name.trim() || !formData.slug.trim()) {
      setError('Name and slug are required');
      return;
    }
    createMutation.mutate({
      name: formData.name,
      slug: formData.slug,
      description: formData.description || undefined,
      color: formData.color || undefined,
      category_id: formData.category_id || undefined,
    });
  };

  const handleEdit = (tag: ProjectTag) => {
    setSelectedTag(tag);
    setFormData({
      name: tag.name,
      slug: tag.slug,
      description: tag.description || '',
      color: tag.color || '#3B82F6',
      category_id: tag.category?.id || tag.category_id || '',
    });
    setEditOpen(true);
  };

  const handleUpdate = () => {
    if (!selectedTag || !formData.name.trim() || !formData.slug.trim()) {
      setError('Name and slug are required');
      return;
    }
    updateMutation.mutate({ 
      id: selectedTag.id, 
      data: {
        name: formData.name,
        slug: formData.slug,
        description: formData.description || undefined,
        color: formData.color || undefined,
        category_id: formData.category_id || undefined,
      }
    });
  };

  const handleDelete = (tag: ProjectTag) => {
    setSelectedTag(tag);
    setDeleteOpen(true);
  };

  const confirmDelete = () => {
    if (selectedTag) {
      deleteMutation.mutate(selectedTag.id);
    }
  };

  const getCategoryName = (tag: ProjectTag) => {
    // First try to use the category object if available
    if (tag.category) {
      return tag.category.name;
    }
    
    // Fallback to category_id lookup
    if (!tag.category_id || !categories) return 'No category';
    const category = categories.find(c => c.id === tag.category_id);
    return category ? category.name : 'Unknown category';
  };

  const getCategoryColor = (tag: ProjectTag) => {
    // First try to use the category object if available
    if (tag.category) {
      return tag.category.color || 'default';
    }
    
    // Fallback to category_id lookup
    if (!tag.category_id || !categories) return 'default';
    const category = categories.find(c => c.id === tag.category_id);
    return category?.color || 'default';
  };

  const getCategoryChipProps = (tag: ProjectTag) => {
    const color = getCategoryColor(tag);
    
    // If it's a hex color, return custom styling
    if (color && color.startsWith('#')) {
      return {
        sx: {
          backgroundColor: color,
          color: 'white',
          '& .MuiChip-label': {
            color: 'white'
          }
        }
      };
    }
    
    // Otherwise use Material-UI color prop
    return {
      color: color as any
    };
  };

  return (
    <AuthenticatedLayout>
      <Box sx={{ p: 3 }}>
        <PageHeader
          title="Project Tags"
          subtitle="Manage tags for this project"
          icon={<TagIcon />}
        >
          {canManage && (
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => setCreateOpen(true)}
            >
              Create Tag
            </Button>
          )}
        </PageHeader>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : queryError ? (
          <Alert severity="error">
            Error loading tags. Please try again.
          </Alert>
        ) : tags && tags.length > 0 ? (
          <Grid container spacing={2}>
            {tags.map((tag) => (
              <Grid item xs={12} sm={6} md={4} key={tag.id}>
                <Card>
                  <CardContent>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                      <Box
                        sx={{
                          width: 20,
                          height: 20,
                          borderRadius: '50%',
                          backgroundColor: tag.color || '#3B82F6',
                          mr: 1,
                        }}
                      />
                      <Typography variant="h6" sx={{ flexGrow: 1 }}>
                        {tag.name}
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                      {tag.slug}
                    </Typography>
                    {tag.description && (
                      <Typography variant="body2" sx={{ mb: 1 }}>
                        {tag.description}
                      </Typography>
                    )}
                    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip
                        icon={<CategoryIcon />}
                        label={getCategoryName(tag)}
                        size="small"
                        {...getCategoryChipProps(tag)}
                      />
                    </Box>
                  </CardContent>
                  {canManage && (
                    <CardActions>
                      <Tooltip title="Edit tag">
                        <IconButton
                          size="small"
                          onClick={() => handleEdit(tag)}
                        >
                          <EditIcon />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Delete tag">
                        <IconButton
                          size="small"
                          onClick={() => handleDelete(tag)}
                          color="error"
                        >
                          <DeleteIcon />
                        </IconButton>
                      </Tooltip>
                    </CardActions>
                  )}
                </Card>
              </Grid>
            ))}
          </Grid>
        ) : (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <TagIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
              No tags found
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              {canManage 
                ? 'Create your first tag to organize features.'
                : 'No tags are available for this project.'
              }
            </Typography>
            {canManage && (
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={() => setCreateOpen(true)}
              >
                Create Tag
              </Button>
            )}
          </Paper>
        )}

        {/* Create Dialog */}
        <Dialog open={createOpen} onClose={() => setCreateOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Create Tag</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="Name"
              fullWidth
              variant="outlined"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Slug"
              fullWidth
              variant="outlined"
              value={formData.slug}
              onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Description"
              fullWidth
              variant="outlined"
              multiline
              rows={3}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Color"
              fullWidth
              variant="outlined"
              type="color"
              value={formData.color}
              onChange={(e) => setFormData({ ...formData, color: e.target.value })}
              sx={{ mb: 2 }}
            />
            <FormControl fullWidth margin="dense">
              <InputLabel>Category</InputLabel>
              <Select
                value={formData.category_id}
                onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                label="Category"
              >
                <MenuItem value="">
                  <em>No category</em>
                </MenuItem>
                {categories?.map((category) => (
                  <MenuItem key={category.id} value={category.id}>
                    <Box sx={{ display: 'flex', alignItems: 'center', width: '100%' }}>
                      <Box
                        sx={{
                          width: 12,
                          height: 12,
                          borderRadius: '50%',
                          backgroundColor: category.color || '#3B82F6',
                          mr: 1,
                        }}
                      />
                      {category.name}
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setCreateOpen(false)}>Cancel</Button>
            <Button 
              onClick={handleCreate} 
              variant="contained"
              disabled={createMutation.isPending}
            >
              {createMutation.isPending ? 'Creating...' : 'Create'}
            </Button>
          </DialogActions>
        </Dialog>

        {/* Edit Dialog */}
        <Dialog open={editOpen} onClose={() => setEditOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Edit Tag</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="Name"
              fullWidth
              variant="outlined"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Slug"
              fullWidth
              variant="outlined"
              value={formData.slug}
              onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Description"
              fullWidth
              variant="outlined"
              multiline
              rows={3}
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Color"
              fullWidth
              variant="outlined"
              type="color"
              value={formData.color}
              onChange={(e) => setFormData({ ...formData, color: e.target.value })}
              sx={{ mb: 2 }}
            />
            <FormControl fullWidth margin="dense">
              <InputLabel>Category</InputLabel>
              <Select
                value={formData.category_id}
                onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                label="Category"
              >
                <MenuItem value="">
                  <em>No category</em>
                </MenuItem>
                {categories?.map((category) => (
                  <MenuItem key={category.id} value={category.id}>
                    <Box sx={{ display: 'flex', alignItems: 'center', width: '100%' }}>
                      <Box
                        sx={{
                          width: 12,
                          height: 12,
                          borderRadius: '50%',
                          backgroundColor: category.color || '#3B82F6',
                          mr: 1,
                        }}
                      />
                      {category.name}
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setEditOpen(false)}>Cancel</Button>
            <Button 
              onClick={handleUpdate} 
              variant="contained"
              disabled={updateMutation.isPending}
            >
              {updateMutation.isPending ? 'Updating...' : 'Update'}
            </Button>
          </DialogActions>
        </Dialog>

        {/* Delete Dialog */}
        <Dialog open={deleteOpen} onClose={() => setDeleteOpen(false)}>
          <DialogTitle>Delete Tag</DialogTitle>
          <DialogContent>
            <Typography>
              Are you sure you want to delete the tag "{selectedTag?.name}"? 
              This action cannot be undone.
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setDeleteOpen(false)}>Cancel</Button>
            <Button 
              onClick={confirmDelete} 
              color="error"
              variant="contained"
              disabled={deleteMutation.isPending}
            >
              {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </AuthenticatedLayout>
  );
};

export default ProjectTagsPage;
