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
  Menu,
  ListItemIcon,
  ListItemText,
  Tabs,
  Tab,
} from '@mui/material';
import { 
  Category as CategoryIcon, 
  Add as AddIcon, 
  Edit as EditIcon,
  Delete as DeleteIcon,
  MoreVert as MoreVertIcon
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate, Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import { userPermissions } from '../hooks/userPermissions';
import type { Category, CreateCategoryRequest, UpdateCategoryRequest, CreateCategoryRequestCategoryTypeEnum } from '../generated/api/client';

const CategoriesPage: React.FC = () => {
  const navigate = useNavigate();
  const { isAuthenticated, user } = useAuth();
  const queryClient = useQueryClient();
  const { canCreateProjects } = userPermissions();
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const [menuCategory, setMenuCategory] = useState<Category | null>(null);
  const [activeTab, setActiveTab] = useState(0);
  const [formData, setFormData] = useState({
    name: '',
    slug: '',
    description: '',
    color: '#3B82F6',
    category_type: 'user',
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
      setFormData({ name: '', slug: '', description: '', color: '#3B82F6', category_type: 'user' });
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
      setFormData({ name: '', slug: '', description: '', color: '#3B82F6', category_type: 'user' });
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

  const validateForm = () => {
    if (!formData.name.trim()) {
      setError('Name is required');
      return false;
    }
    if (!formData.slug.trim()) {
      setError('Slug is required');
      return false;
    }
    // Slug validation: only a-z, 0-9, dash and underscore
    const slugRegex = /^[a-z0-9_-]+$/;
    if (!slugRegex.test(formData.slug)) {
      setError('Slug can only contain lowercase letters, numbers, dashes and underscores');
      return false;
    }
    return true;
  };

  const handleCreate = () => {
    if (!validateForm()) {
      return;
    }
    createMutation.mutate({
      name: formData.name,
      slug: formData.slug,
      description: formData.description || undefined,
      color: formData.color || undefined,
      category_type: formData.category_type as CreateCategoryRequestCategoryTypeEnum,
    });
  };

  const handleEdit = (category: Category) => {
    setSelectedCategory(category);
    setFormData({
      name: category.name,
      slug: category.slug,
      description: category.description || '',
      color: category.color || '#3B82F6',
      category_type: category.category_type,
    });
    setEditOpen(true);
  };

  const handleUpdate = () => {
    if (!selectedCategory) {
      setError('No category selected');
      return;
    }
    if (!validateForm()) {
      return;
    }
    updateMutation.mutate({ 
      id: selectedCategory.id, 
      data: {
        name: formData.name,
        slug: formData.slug,
        description: formData.description || undefined,
        color: formData.color || undefined,
        // category_type и kind нельзя обновлять
      }
    });
  };

  const handleDelete = (category: Category) => {
    // Проверяем, можно ли удалить категорию
    if ((category.kind as string) === 'system' || (category.kind as string) === 'nocopy') {
      setError('Cannot delete system or nocopy categories');
      return;
    }
    setSelectedCategory(category);
    setDeleteOpen(true);
  };

  const confirmDelete = () => {
    if (selectedCategory) {
      deleteMutation.mutate(selectedCategory.id);
    }
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, category: Category) => {
    setMenuAnchor(event.currentTarget);
    setMenuCategory(category);
  };

  const handleMenuClose = () => {
    setMenuAnchor(null);
    setMenuCategory(null);
  };

  const handleEditFromMenu = () => {
    if (menuCategory) {
      handleEdit(menuCategory);
    }
    handleMenuClose();
  };

  const handleDeleteFromMenu = () => {
    if (menuCategory) {
      handleDelete(menuCategory);
    }
    handleMenuClose();
  };

  const getCategoryTypeColor = (categoryType: string) => {
    switch (categoryType) {
      case 'domain': return 'primary';
      case 'safety': return 'error';
      case 'user': return 'secondary';
      default: return 'default';
    }
  };

  const getCategoryTypeLabel = (categoryType: string) => {
    switch (categoryType) {
      case 'domain': return 'Domain';
      case 'safety': return 'Safety';
      case 'user': return 'User';
      default: return categoryType;
    }
  };

  const getFilteredCategories = () => {
    if (!categories) return [];
    
    const tabTypes = ['domain', 'user', 'safety'];
    const selectedType = tabTypes[activeTab];
    
    return categories.filter(category => category.category_type === selectedType);
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  return (
    <AuthenticatedLayout>
      <Box sx={{ p: 3 }}>
        <PageHeader
          title="Categories"
          subtitle="Manage feature tag categories"
          icon={<CategoryIcon />}
        >
          {isSuperuser && (
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={() => setCreateOpen(true)}
            >
              Create Category
            </Button>
          )}
        </PageHeader>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {/* Tabs */}
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={activeTab} onChange={handleTabChange} aria-label="category type tabs">
            <Tab label="Domain" />
            <Tab label="User" />
            <Tab label="Safety" />
          </Tabs>
        </Box>

        {isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : queryError ? (
          <Alert severity="error">
            Error loading categories. Please try again.
          </Alert>
        ) : categories && categories.length > 0 ? (
          getFilteredCategories().length > 0 ? (
            <Grid container spacing={2}>
              {getFilteredCategories().map((category) => (
                <Grid item xs={12} sm={6} md={4} key={category.id}>
                  <Card>
                    <CardContent sx={{ position: 'relative' }}>
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
                        {isSuperuser && (
                          <IconButton
                            size="small"
                            onClick={(e) => handleMenuOpen(e, category)}
                            sx={{ position: 'absolute', top: 8, right: 8 }}
                          >
                            <MoreVertIcon />
                          </IconButton>
                        )}
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
                          label={getCategoryTypeLabel(category.category_type)}
                          color={getCategoryTypeColor(category.category_type) as any}
                          size="small"
                        />
                      </Box>
                    </CardContent>
                  </Card>
                </Grid>
              ))}
            </Grid>
          ) : (
            <Paper sx={{ p: 4, textAlign: 'center' }}>
              <CategoryIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
                No {['Domain', 'User', 'Safety'][activeTab]} categories found
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                {isSuperuser 
                  ? 'Create your first category to get started.'
                  : 'No categories have been created yet.'
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
          )
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
              error={!formData.name.trim()}
              helperText={!formData.name.trim() ? 'Name is required' : undefined}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Slug"
              fullWidth
              variant="outlined"
              value={formData.slug}
              onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
              error={Boolean(!formData.slug.trim() || (formData.slug.trim() && !/^[a-z0-9_-]+$/.test(formData.slug)))}
              helperText={
                !formData.slug.trim() 
                  ? 'Slug is required' 
                  : formData.slug.trim() && !/^[a-z0-9_-]+$/.test(formData.slug)
                    ? 'Slug can only contain lowercase letters, numbers, dashes and underscores'
                    : undefined
              }
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
              <InputLabel>Category Type</InputLabel>
              <Select
                value={formData.category_type}
                onChange={(e) => setFormData({ ...formData, category_type: e.target.value })}
                label="Category Type"
              >
                <MenuItem value="user">User</MenuItem>
                <MenuItem value="domain">Domain</MenuItem>
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
          <DialogTitle>Edit Category</DialogTitle>
          <DialogContent>
            <Alert severity="info" sx={{ mb: 2 }}>
              You can only edit name, slug, description, and color. Kind and category type cannot be changed.
            </Alert>
            <TextField
              autoFocus
              margin="dense"
              label="Name"
              fullWidth
              variant="outlined"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              error={!formData.name.trim()}
              helperText={!formData.name.trim() ? 'Name is required' : undefined}
              sx={{ mb: 2 }}
            />
            <TextField
              margin="dense"
              label="Slug"
              fullWidth
              variant="outlined"
              value={formData.slug}
              onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
              error={Boolean(!formData.slug.trim() || (formData.slug.trim() && !/^[a-z0-9_-]+$/.test(formData.slug)))}
              helperText={
                !formData.slug.trim() 
                  ? 'Slug is required' 
                  : formData.slug.trim() && !/^[a-z0-9_-]+$/.test(formData.slug)
                    ? 'Slug can only contain lowercase letters, numbers, dashes and underscores'
                    : undefined
              }
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

        {/* Context Menu */}
        <Menu
          anchorEl={menuAnchor}
          open={Boolean(menuAnchor)}
          onClose={handleMenuClose}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'right',
          }}
        >
          {/* Edit option - only for user categories */}
          {menuCategory && (menuCategory.kind as string) === 'user' && (
            <MenuItem onClick={handleEditFromMenu}>
              <ListItemIcon>
                <EditIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Edit category</ListItemText>
            </MenuItem>
          )}
          
          {/* Delete option */}
          <MenuItem 
            onClick={handleDeleteFromMenu}
            disabled={!!(menuCategory && ((menuCategory.kind as string) === 'system' || (menuCategory.kind as string) === 'nocopy'))}
          >
            <ListItemIcon>
              <DeleteIcon fontSize="small" />
            </ListItemIcon>
            <ListItemText>
              {menuCategory && ((menuCategory.kind as string) === 'system' || (menuCategory.kind as string) === 'nocopy')
                ? 'Cannot delete system category'
                : 'Delete category'
              }
            </ListItemText>
          </MenuItem>
        </Menu>
      </Box>
    </AuthenticatedLayout>
  );
};

export default CategoriesPage;
