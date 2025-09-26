import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import type { Category, CreateCategoryRequest, UpdateCategoryRequest, CreateCategoryRequestKindEnum } from '../../generated/api/client';

interface CategoryFormData {
  name: string;
  slug: string;
  description: string;
  color: string;
  kind: CreateCategoryRequestKindEnum;
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
    kind: 'user' as CreateCategoryRequestKindEnum,
  });

  useEffect(() => {
    if (mode === 'edit' && initialData) {
      setFormData({
        name: initialData.name || '',
        slug: initialData.slug || '',
        description: initialData.description || '',
        color: initialData.color || '#3B82F6',
        kind: (initialData.kind as CreateCategoryRequestKindEnum) || 'user' as CreateCategoryRequestKindEnum,
      });
    } else {
      setFormData({
        name: '',
        slug: '',
        description: '',
        color: '#3B82F6',
        kind: 'user' as CreateCategoryRequestKindEnum,
      });
    }
  }, [mode, initialData, open]);

  const handleSubmit = () => {
    if (mode === 'create') {
      const submitData: CreateCategoryRequest = {
        name: formData.name.trim(),
        slug: formData.slug.trim(),
        description: formData.description.trim(),
        color: formData.color,
        kind: formData.kind,
      };
      onSubmit(submitData);
    } else {
      const submitData: UpdateCategoryRequest = {
        name: formData.name.trim(),
        slug: formData.slug.trim(),
        description: formData.description.trim(),
        color: formData.color,
      };
      onSubmit(submitData);
    }
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
        
        <FormControl fullWidth margin="dense">
          <InputLabel>Type</InputLabel>
          <Select
            value={formData.kind}
            onChange={(e) => setFormData({ ...formData, kind: e.target.value as CreateCategoryRequestKindEnum })}
            label="Type"
            disabled={mode === 'edit'} // Don't allow changing type when editing
          >
            <MenuItem value="user">User</MenuItem>
            <MenuItem value="domain">Domain</MenuItem>
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

export default CategoryFormDialog;
