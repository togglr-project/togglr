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
  Tabs,
  Tab,
  Menu,
  ListItemIcon,
  ListItemText,
} from '@mui/material';
import { 
  LocalOffer as TagIcon, 
  Add as AddIcon, 
  Edit as EditIcon, 
  Delete as DeleteIcon,
  Category as CategoryIcon,
  MoreVert as MoreVertIcon
} from '@mui/icons-material';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useParams, Navigate } from 'react-router-dom';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import PageHeader from '../components/PageHeader';
import TagFormDialog from '../components/tags/TagFormDialog';
import apiClient from '../api/apiClient';
import { useAuth } from '../auth/AuthContext';
import { useRBAC } from '../auth/permissions';
import type { ProjectTag, Category, CreateProjectTagRequest, UpdateProjectTagRequest } from '../generated/api/client';

const ProjectTagsPage: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const { user } = useAuth();
  const rbac = useRBAC(projectId);
  const queryClient = useQueryClient();
  const [createOpen, setCreateOpen] = useState(false);
  const [editOpen, setEditOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [selectedTag, setSelectedTag] = useState<ProjectTag | null>(null);
  const [activeTab, setActiveTab] = useState(0);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const [menuTag, setMenuTag] = useState<ProjectTag | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    slug: '',
    description: '',
    color: '#3B82F6',
    category_id: '',
  });
  const [error, setError] = useState<string | null>(null);

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

  if (!projectId) {
    return <Navigate to="/projects" replace />;
  }

  // Check project access and tag management permissions
  if (!rbac.canViewProject()) {
    return (
      <AuthenticatedLayout showBackButton backTo="/dashboard">
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Typography variant="h6" color="error" gutterBottom>
            Access Denied
          </Typography>
          <Typography variant="body2" color="text.secondary">
            You don't have permission to view this project.
          </Typography>
        </Box>
      </AuthenticatedLayout>
    );
  }

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
    if (!formData.category_id) {
      setError('Category is required');
      return false;
    }
    return true;
  };

  const handleCreate = () => {
    if (!validateForm()) {
      return;
    }
    
    // Проверяем, что выбранная категория имеет правильный kind
    if (formData.category_id) {
      const selectedCategory = categories?.find(cat => cat.id === formData.category_id);
      if (selectedCategory && selectedCategory.kind !== 'domain' && selectedCategory.kind !== 'user') {
        setError('Can only create tags for domain or user categories');
        return;
      }
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
    // Проверяем, можно ли редактировать тег (категория не должна быть system)
    if (tag.category && (tag.category.kind as string) === 'system') {
      setError('Cannot edit tags with system categories');
      return;
    }
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
    if (!selectedTag) {
      setError('No tag selected');
      return;
    }
    if (!validateForm()) {
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
    // Проверяем, можно ли удалить тег (категория не должна быть system)
    if (tag.category && (tag.category.kind as string) === 'system') {
      setError('Cannot delete tags with system categories');
      return;
    }
    setSelectedTag(tag);
    setDeleteOpen(true);
  };

  const confirmDelete = () => {
    if (selectedTag) {
      deleteMutation.mutate(selectedTag.id);
    }
  };

  const getFilteredTags = () => {
    if (!tags) return [];
    
    const tabKinds = ['domain', 'user', 'system'];
    const selectedKind = tabKinds[activeTab];
    
    return tags.filter(tag => tag.category?.kind === selectedKind);
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
  };

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, tag: ProjectTag) => {
    setMenuAnchor(event.currentTarget);
    setMenuTag(tag);
  };

  const handleMenuClose = () => {
    setMenuAnchor(null);
    setMenuTag(null);
  };

  const handleEditFromMenu = () => {
    if (menuTag) {
      handleEdit(menuTag);
    }
    handleMenuClose();
  };

  const handleDeleteFromMenu = () => {
    if (menuTag) {
      handleDelete(menuTag);
    }
    handleMenuClose();
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
          {rbac.canManageTags() && (
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

        {/* Tabs */}
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={activeTab} onChange={handleTabChange} aria-label="tag category type tabs">
            <Tab label="Domain" />
            <Tab label="User" />
            <Tab label="System" />
          </Tabs>
        </Box>

        {isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : queryError ? (
          <Alert severity="error">
            Error loading tags. Please try again.
          </Alert>
        ) : tags && tags.length > 0 ? (
          getFilteredTags().length > 0 ? (
            <Grid container spacing={2}>
              {getFilteredTags().map((tag) => (
              <Grid item xs={12} sm={6} md={4} key={tag.id}>
                <Card>
                  <CardContent sx={{ position: 'relative' }}>
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
                      {rbac.canManageTags() && (
                        <IconButton
                          size="small"
                          onClick={(e) => handleMenuOpen(e, tag)}
                          sx={{ position: 'absolute', top: 8, right: 8 }}
                        >
                          <MoreVertIcon />
                        </IconButton>
                      )}
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
                </Card>
              </Grid>
            ))}
          </Grid>
          ) : (
            <Paper sx={{ p: 4, textAlign: 'center' }}>
              <TagIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
              <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
                No {['Domain', 'User', 'System'][activeTab]} tags found
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                {rbac.canManageTags() 
                  ? 'Create your first tag to organize features.'
                  : 'No tags are available for this project.'
                }
              </Typography>
              {rbac.canManageTags() && (
                <Button
                  variant="contained"
                  startIcon={<AddIcon />}
                  onClick={() => setCreateOpen(true)}
                >
                  Create Tag
                </Button>
              )}
            </Paper>
          )
        ) : (
          <Paper sx={{ p: 4, textAlign: 'center' }}>
            <TagIcon sx={{ fontSize: 64, color: 'text.secondary', mb: 2 }} />
            <Typography variant="h6" color="text.secondary" sx={{ mb: 1 }}>
              No tags found
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              {rbac.canManageTags() 
                ? 'Create your first tag to organize features.'
                : 'No tags are available for this project.'
              }
            </Typography>
            {rbac.canManageTags() && (
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
        <TagFormDialog
          open={createOpen}
          onClose={() => setCreateOpen(false)}
          onSubmit={(data) => createMutation.mutate(data)}
          categories={categories || []}
          mode="create"
          error={error}
        />
        {/* Edit Dialog */}
        <TagFormDialog
          open={editOpen}
          onClose={() => setEditOpen(false)}
          onSubmit={(data) => updateMutation.mutate({ id: selectedTag?.id || '', data })}
          categories={categories || []}
          mode="edit"
          initialData={selectedTag}
          error={error}
        />
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
          {/* Edit option - only for non-system categories */}
          {menuTag && menuTag.category && (menuTag.category.kind as string) !== 'system' && (
            <MenuItem onClick={handleEditFromMenu}>
              <ListItemIcon>
                <EditIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Edit tag</ListItemText>
            </MenuItem>
          )}
          
          {/* Delete option */}
          <MenuItem 
            onClick={handleDeleteFromMenu}
            disabled={!!(menuTag && menuTag.category && (menuTag.category.kind as string) === 'system')}
          >
            <ListItemIcon>
              <DeleteIcon fontSize="small" />
            </ListItemIcon>
            <ListItemText>
              {menuTag && menuTag.category && (menuTag.category.kind as string) === 'system'
                ? 'Cannot delete tags with system categories'
                : 'Delete tag'
              }
            </ListItemText>
          </MenuItem>
        </Menu>
      </Box>
    </AuthenticatedLayout>
  );
};

export default ProjectTagsPage;
