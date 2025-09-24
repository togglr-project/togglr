import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Alert,
  Box,
} from '@mui/material';
import type { ProjectTag, Category, CreateProjectTagRequest, UpdateProjectTagRequest } from '../../generated/api/client';

interface TagFormData {
  name: string;
  slug: string;
  description: string;
  color: string;
  category_id: string;
}

interface TagFormDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateProjectTagRequest | UpdateProjectTagRequest) => void;
  categories: Category[];
  initialData?: ProjectTag | null;
  mode: 'create' | 'edit';
  error?: string | null;
}

const TagFormDialog: React.FC<TagFormDialogProps> = ({
  open,
  onClose,
  onSubmit,
  categories,
  initialData,
  mode,
  error,
}) => {
  const [formData, setFormData] = useState<TagFormData>({
    name: '',
    slug: '',
    description: '',
    color: '#3B82F6',
    category_id: '',
  });

  useEffect(() => {
    if (mode === 'edit' && initialData) {
      setFormData({
        name: initialData.name || '',
        slug: initialData.slug || '',
        description: initialData.description || '',
        color: initialData.color || '#3B82F6',
        category_id: initialData.category?.id?.toString() || '',
      });
    } else {
      setFormData({
        name: '',
        slug: '',
        description: '',
        color: '#3B82F6',
        category_id: '',
      });
    }
  }, [mode, initialData, open]);

  const handleSubmit = () => {
    const submitData = {
      name: formData.name.trim(),
      slug: formData.slug.trim(),
      description: formData.description.trim(),
      color: formData.color,
      category_id: parseInt(formData.category_id),
    };

    onSubmit(submitData);
  };

  const isFormValid = () => {
    return (
      formData.name.trim() &&
      formData.slug.trim() &&
      /^[a-z0-9_-]+$/.test(formData.slug.trim()) &&
      formData.category_id
    );
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {mode === 'create' ? 'Create Tag' : 'Edit Tag'}
      </DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
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
          <InputLabel>Category</InputLabel>
          <Select
            value={formData.category_id}
            onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
            label="Category"
            error={!formData.category_id}
          >
            {categories
              .filter(category => 
                category.category_type === 'user' || category.category_type === 'domain'
              )
              .map((category) => (
                <MenuItem key={category.id} value={category.id?.toString()}>
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
        <Button onClick={onClose}>Cancel</Button>
        <Button 
          onClick={handleSubmit} 
          variant="contained"
          disabled={!isFormValid()}
        >
          {mode === 'create' ? 'Create' : 'Update'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default TagFormDialog;
