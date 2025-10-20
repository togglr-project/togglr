import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Alert,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormHelperText,
  Chip,
} from '@mui/material';
import {
  Code as CodeIcon,
  FormatAlignLeft as TextIcon,
  Check as CheckIcon,
  Close as CloseIcon,
  ToggleOn as BooleanIcon,
  Numbers as NumberIcon,
  Functions as DoubleIcon,
} from '@mui/icons-material';
import { useMutation } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { ProjectSetting, CreateProjectSettingRequest, UpdateProjectSettingRequest } from '../../generated/api/client';
import { useNotification } from '../../App';
import { getSettingDefinition, isPredefinedSetting } from '../../constants/projectSettings';

interface ProjectSettingFormDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: () => void;
  projectId: string;
  mode: 'create' | 'edit';
  initialData?: ProjectSetting | null;
}

interface FormData {
  name: string;
  value: string;
  valueType: 'text' | 'json' | 'boolean' | 'integer' | 'double';
}

const ProjectSettingFormDialog: React.FC<ProjectSettingFormDialogProps> = ({
  open,
  onClose,
  onSubmit,
  projectId,
  mode,
  initialData,
}) => {
  const { showNotification } = useNotification();
  const [formData, setFormData] = useState<FormData>({
    name: '',
    value: '',
    valueType: 'text',
  });
  const [errors, setErrors] = useState<Partial<FormData>>({});
  const [jsonError, setJsonError] = useState<string>('');

  useEffect(() => {
    if (open) {
      if (mode === 'edit' && initialData) {
        const settingDef = getSettingDefinition(initialData.name);
        const isJson = isValidJson(initialData.value);
        setFormData({
          name: initialData.name,
          value: initialData.value,
          valueType: settingDef?.type || (isJson ? 'json' : 'text'),
        });
      } else {
        setFormData({
          name: '',
          value: '',
          valueType: 'text',
        });
      }
      setErrors({});
      setJsonError('');
    }
  }, [open, mode, initialData]);

  const isValidJson = (value: string): boolean => {
    try {
      JSON.parse(value);
      return true;
    } catch {
      return false;
    }
  };

  const isValidInteger = (value: string): boolean => {
    const num = parseInt(value, 10);
    return !isNaN(num) && num.toString() === value;
  };

  const isValidDouble = (value: string): boolean => {
    const num = parseFloat(value);
    return !isNaN(num) && isFinite(num);
  };

  const isValidBoolean = (value: string): boolean => {
    return value === 'true' || value === 'false';
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'json':
        return 'JSON';
      case 'boolean':
        return 'Boolean';
      case 'integer':
        return 'Integer';
      case 'double':
        return 'Double';
      default:
        return 'Text';
    }
  };

  const formatJson = (value: string): string => {
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch {
      return value;
    }
  };

  const validateForm = (): boolean => {
    const newErrors: Partial<FormData> = {};
    const settingDef = getSettingDefinition(formData.name);

    if (!formData.name.trim()) {
      newErrors.name = 'Name is required';
    }

    if (!formData.value.trim()) {
      newErrors.value = 'Value is required';
    }

    if (formData.valueType === 'json' && !isValidJson(formData.value)) {
      setJsonError('Invalid JSON format');
      return false;
    }

    if (formData.valueType === 'integer' && !isValidInteger(formData.value)) {
      newErrors.value = 'Invalid integer value';
    }

    if (formData.valueType === 'double' && !isValidDouble(formData.value)) {
      newErrors.value = 'Invalid number value';
    }

    if (formData.valueType === 'boolean' && !isValidBoolean(formData.value)) {
      newErrors.value = 'Invalid boolean value (use true or false)';
    }

    if (settingDef?.validation) {
      const numValue = parseFloat(formData.value);
      if (!isNaN(numValue)) {
        if (settingDef.validation.min !== undefined && numValue < settingDef.validation.min) {
          newErrors.value = `Value must be at least ${settingDef.validation.min}`;
        }
        if (settingDef.validation.max !== undefined && numValue > settingDef.validation.max) {
          newErrors.value = `Value must be at most ${settingDef.validation.max}`;
        }
      }
    }

    setErrors(newErrors);
    setJsonError('');
    return Object.keys(newErrors).length === 0;
  };

  const createMutation = useMutation({
    // mutationFn: async (data: CreateProjectSettingRequest) => {
    //   const response = await apiClient.createProjectSetting(projectId, data);
    //   return response.data;
    // },
    // onSuccess: () => {
    //   showNotification('Setting created successfully', 'success');
    //   onSubmit();
    // },
    // onError: (error: any) => {
    //   showNotification(`Error creating setting: ${error.message}`, 'error');
    // },
  });

  const updateMutation = useMutation({
    mutationFn: async (data: UpdateProjectSettingRequest) => {
      if (!initialData) throw new Error('No data to update');
      const response = await apiClient.updateProjectSetting(projectId, initialData.name, data);
      return response.data;
    },
    onSuccess: () => {
      showNotification('Setting updated successfully', 'success');
      onSubmit();
    },
    onError: (error: any) => {
      showNotification(`Error updating setting: ${error.message}`, 'error');
    },
  });

  const handleSubmit = () => {
    if (!validateForm()) {
      return;
    }

    const submitData = {
      value: formData.value,
    };

    if (mode === 'create') {
      // createMutation.mutate({
      //   name: formData.name,
      //   value: formData.value,
      // });
    } else {
      updateMutation.mutate(submitData);
    }
  };

  const handleValueTypeChange = (newType: 'text' | 'json' | 'boolean' | 'integer' | 'double') => {
    const settingDef = getSettingDefinition(formData.name);
    let newValue = formData.value;

    if (newType === 'json' && formData.value) {
      try {
        newValue = formatJson(formData.value);
      } catch {
      }
    } else if (newType === 'boolean' && !isValidBoolean(formData.value)) {
      newValue = 'false';
    } else if (newType === 'integer' && !isValidInteger(formData.value)) {
      newValue = settingDef?.defaultValue?.toString() || '0';
    } else if (newType === 'double' && !isValidDouble(formData.value)) {
      newValue = settingDef?.defaultValue?.toString() || '0.0';
    }

    setFormData(prev => ({
      ...prev,
      valueType: newType,
      value: newValue,
    }));
  };

  const handleValueChange = (value: string) => {
    setFormData(prev => ({
      ...prev,
      value,
    }));

    if (formData.valueType === 'json') {
      setJsonError(isValidJson(value) ? '' : 'Invalid JSON format');
    }
  };

  const isSubmitting = createMutation.isPending || updateMutation.isPending;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>
        {mode === 'create' ? 'Create Setting' : 'Edit Setting'}
      </DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3, pt: 1 }}>
          <TextField
            label="Setting Name"
            value={formData.name}
            onChange={(e) => {
              const newName = e.target.value;
              const settingDef = getSettingDefinition(newName);
              setFormData(prev => ({
                ...prev,
                name: newName,
                valueType: settingDef?.type || prev.valueType,
                value: settingDef?.defaultValue?.toString() || prev.value,
              }));
            }}
            fullWidth
            disabled={mode === 'edit'}
            error={!!errors.name}
            helperText={errors.name || (isPredefinedSetting(formData.name) ? 'Predefined setting' : '')}
            placeholder="e.g., auto_disable_enabled"
          />

          {mode === 'create' && (
            <Box>
              <FormControl fullWidth>
                <InputLabel>Value Type</InputLabel>
                <Select
                  value={formData.valueType}
                  onChange={(e) => handleValueTypeChange(e.target.value as 'text' | 'json' | 'boolean' | 'integer' | 'double')}
                  label="Value Type"
                >
                  <MenuItem value="text">
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <TextIcon fontSize="small" />
                      Text
                    </Box>
                  </MenuItem>
                  <MenuItem value="json">
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <CodeIcon fontSize="small" />
                      JSON
                    </Box>
                  </MenuItem>
                  <MenuItem value="boolean">
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <BooleanIcon fontSize="small" />
                      Boolean
                    </Box>
                  </MenuItem>
                  <MenuItem value="integer">
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <NumberIcon fontSize="small" />
                      Integer
                    </Box>
                  </MenuItem>
                  <MenuItem value="double">
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <DoubleIcon fontSize="small" />
                      Double
                    </Box>
                  </MenuItem>
                </Select>
              </FormControl>
            </Box>
          )}

          {mode === 'edit' && (
            <Box>
              <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1 }}>
                Value Type: {getTypeLabel(formData.valueType)}
              </Typography>
            </Box>
          )}

          <Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
              <Typography variant="subtitle2">
                Value
              </Typography>
              {formData.valueType === 'json' && (
                <Chip
                  icon={<CodeIcon />}
                  label="JSON"
                  size="small"
                  color="primary"
                  variant="outlined"
                />
              )}
            </Box>
            {formData.valueType === 'boolean' ? (
              <FormControl fullWidth>
                <Select
                  value={formData.value}
                  onChange={(e) => handleValueChange(e.target.value)}
                  error={!!errors.value}
                >
                  <MenuItem value="true">True</MenuItem>
                  <MenuItem value="false">False</MenuItem>
                </Select>
                {errors.value && <FormHelperText error>{errors.value}</FormHelperText>}
              </FormControl>
            ) : (
              <TextField
                multiline={formData.valueType === 'json' || formData.valueType === 'text'}
                rows={formData.valueType === 'json' ? 8 : 1}
                type={formData.valueType === 'integer' || formData.valueType === 'double' ? 'number' : 'text'}
                value={formData.value}
                onChange={(e) => handleValueChange(e.target.value)}
                fullWidth
                error={!!errors.value || !!jsonError}
                helperText={errors.value || jsonError}
                placeholder={
                  formData.valueType === 'json'
                    ? '{\n  "key": "value",\n  "enabled": true\n}'
                    : formData.valueType === 'integer'
                    ? 'Enter integer value'
                    : formData.valueType === 'double'
                    ? 'Enter decimal value'
                    : 'Enter setting value'
                }
                sx={{
                  '& .MuiInputBase-input': {
                    fontFamily: formData.valueType === 'json' ? 'monospace' : 'inherit',
                  },
                }}
              />
            )}
            {formData.valueType === 'json' && formData.value && (
              <Box sx={{ mt: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                {isValidJson(formData.value) ? (
                  <Chip
                    icon={<CheckIcon />}
                    label="Valid JSON"
                    size="small"
                    color="success"
                    variant="outlined"
                  />
                ) : (
                  <Chip
                    icon={<CloseIcon />}
                    label="Invalid JSON"
                    size="small"
                    color="error"
                    variant="outlined"
                  />
                )}
              </Box>
            )}
          </Box>

          {formData.valueType === 'json' && (
            <Alert severity="info" sx={{ mt: 1 }}>
              <Typography variant="body2">
                For JSON values, use valid syntax. The value will be automatically formatted.
              </Typography>
            </Alert>
          )}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={isSubmitting || !!jsonError}
        >
          {isSubmitting ? 'Saving...' : (mode === 'create' ? 'Create' : 'Save')}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ProjectSettingFormDialog;
