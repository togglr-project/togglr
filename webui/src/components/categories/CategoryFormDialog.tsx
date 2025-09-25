import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Alert,
} from '@mui/material';
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from '../../generated/api/client';

interface CategoryFormData {
  name: string;
  slug: string;
  description: string;
  color: string;
}

interface CategoryFormDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: CreateCategoryRequest | UpdateCategoryRequest) => void;
  initialData?: Category | null;
  mode: 'create' | 'edit';
  error?: string | null;
}

const CategoryFormDialog: React.FC<CategoryFormDialogProps> = ({
  open,
  onClose,
  onSubmit,
  initialData,
  mode,
  error,
}) => {
  const [formData, setFormData] = useState<CategoryFormData>({
    name: '',
    slug: '',
    description: '',
    color: '#3B82F6',
  });

  useEffect(() => {
    if (mode === 'edit' && initialData) {
      setFormData({
        name: initialData.name || '',
        slug: initialData.slug || '',
        description: initialData.description || '',
        color: initialData.color || '#3B82F6',
      });
    } else {
      setFormData({
        name: '',
        slug: '',
        description: '',
        color: '#3B82F6',
      });
    }
  }, [mode, initialData, open]);

  const handleSubmit = () => {
    const submitData = {
      name: formData.name.trim(),
      slug: formData.slug.trim(),
      description: formData.description.trim(),
      color: formData.color,
    };

    onSubmit(submitData);
  };

  const isFormValid = () => {
    return (
      formData.name.trim() &&
      formData.slug.trim() &&
      /^[a-z0-9_-]+$/.test(formData.slug.trim())
    );
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {mode === 'create' ? 'Create Category' : 'Edit Category'}
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

export default CategoryFormDialog;
