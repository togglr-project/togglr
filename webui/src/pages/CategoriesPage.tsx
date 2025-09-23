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
} from '@mui/material';
import { 
  Category as CategoryIcon, 
  Add as AddIcon, 
  Edit as EditIcon, 
  Delete as DeleteIcon
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate, Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import { userPermissions } from '../hooks/userPermissions';
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '../generated/api/client';

const CategoriesPage: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated, user } = useAuth();
  const queryClient = useQueryClient();
  const { canCreateProjects } = userPermissions();
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    slug: '',
    description: '',
    color: '#3B82F6',
  });
  const [error, setError] = useState<string | null>(null);

  const { data: categories, isLoading, error: queryError } = useQuery<Category[]>({
    queryKey: ['categories'],
    queryFn: async () => {
      const res = await apiClient.listCategories();
      return res.data;
    },
  });

  const createMutation = useMutation({
    mutationFn: async (data: CreateCategoryRequest) => {
      const res = await apiClient.createCategory(data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      setCreateOpen(false);
      setFormData({ name: '', slug: '', description: '', color: '#3B82F6' });
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to create category');
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateCategoryRequest }) => {
      const res = await apiClient.updateCategory(id, data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      setEditOpen(false);
      setSelectedCategory(null);
      setFormData({ name: '', slug: '', description: '', color: '#3B82F6' });
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to update category');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.deleteCategory(id);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      setDeleteOpen(false);
      setSelectedCategory(null);
      setError(null);
    },
    onError: (err: any) => {
      setError(err.response?.data?.error?.message || 'Failed to delete category');
    },
  });

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  const isSuperuser = user?.is_superuser;

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
    });
  };

  const handleEdit = (category: Category) => {
    setSelectedCategory(category);
    setFormData({
      name: category.name,
      slug: category.slug,
      description: category.description || '',
      color: category.color || '#3B82F6',
    });
    setEditOpen(true);
  };

  const handleUpdate = () => {
    if (!selectedCategory || !formData.name.trim() || !formData.slug.trim()) {
      setError('Name and slug are required');
      return;
    }
    updateMutation.mutate({ 
      id: selectedCategory.id, 
      data: {
        name: formData.name,
        slug: formData.slug,
        description: formData.description || undefined,
        color: formData.color || undefined,
      }
    });
  };

  const handleDelete = (category: Category) => {
    setSelectedCategory(category);
    setDeleteOpen(true);
  };

  const confirmDelete = () => {
    if (selectedCategory) {
      deleteMutation.mutate(selectedCategory.id);
    }
  };

  const getKindColor = (kind: string) => {
    switch (kind) {
      case 'system': return 'primary';
      case 'user': return 'secondary';
      default: return 'default';
    }
  };

  const getKindLabel = (kind: string) => {
    switch (kind) {
      case 'system': return 'System';
      case 'user': return 'User';
      default: return kind;
    }
  };

  return (
    <AuthenticatedLayout>
      <Box sx={{ p: 3 }}>
        <PageHeader
          title="Categories"
          subtitle="Manage feature tag categories"
          icon={<CategoryIcon />}
          action={
            isSuperuser ? (
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={() => setCreateOpen(true)}
              >
                Create Category
              </Button>
            ) : null
          }
        />

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
            Error loading categories. Please try again.
          </Alert>
        ) : categories && categories.length > 0 ? (
          <Grid container spacing={2}>
            {categories.map((category) => (
              <Grid item xs={12} sm={6} md={4} key={category.id}>
                <Card>
                  <CardContent>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                      <Box
                        sx={{
                          width: 20,
                          height: 20,
                          borderRadius: '50%',
                          backgroundColor: category.color || '#3B82F6',
                          mr: 1,
                        }}
                      />
                      <Typography variant="h6" sx={{ flexGrow: 1 }}>
                        {category.name}
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                      {category.slug}
                    </Typography>
                    {category.description && (
                      <Typography variant="body2" sx={{ mb: 1 }}>
                        {category.description}
                      </Typography>
                    )}
                    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
                      <Chip
                        label={getKindLabel(category.kind)}
                        color={getKindColor(category.kind) as any}
                        size="small"
                      />
                    </Box>
                  </CardContent>
                  {isSuperuser && (
                    <CardActions>
                      <Tooltip title="Edit category">
                        <IconButton
                          size="small"
                          onClick={() => handleEdit(category)}
                        >
                          <EditIcon />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Delete category">
                        <IconButton
                          size="small"
                          onClick={() => handleDelete(category)}
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
            <CategoryIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
              No categories found
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              {isSuperuser 
                ? 'Create your first category to organize feature tags.'
                : 'No categories are available at the moment.'
              }
            </Typography>
            {isSuperuser && (
              <Button
                variant="contained"
                startIcon={<AddIcon />}
                onClick={() => setCreateOpen(true)}
              >
                Create Category
              </Button>
            )}
          </Paper>
        )}

        {/* Create Dialog */}
        <Dialog open={createOpen} onClose={() => setCreateOpen(false)} maxWidth="sm" fullWidth>
          <DialogTitle>Create Category</DialogTitle>
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
            />
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
          <DialogTitle>Edit Category</DialogTitle>
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
            />
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
          <DialogTitle>Delete Category</DialogTitle>
          <DialogContent>
            <Typography>
              Are you sure you want to delete the category "{selectedCategory?.name}"? 
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

export default CategoriesPage;
