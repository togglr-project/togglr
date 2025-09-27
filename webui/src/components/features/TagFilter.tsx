import React from 'react';
import {
  Autocomplete,
  Box,
  Chip,
  TextField,
  Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { ProjectTag } from '../../generated/api/client';

interface TagFilterProps {
  projectId: string;
  selectedTags: ProjectTag[];
  onChange: (tags: ProjectTag[]) => void;
  disabled?: boolean;
}

const TagFilter: React.FC<TagFilterProps> = ({
  projectId,
  selectedTags,
  onChange,
  disabled = false,
}) => {
  const { data: availableTags, isLoading } = useQuery<ProjectTag[]>({
    queryKey: ['project-tags', projectId],
    queryFn: async () => {
      const res = await apiClient.listProjectTags(projectId);
      return res.data;
    },
    enabled: !!projectId,
  });

  const handleTagChange = (_: unknown, newValue: ProjectTag[]) => {
    onChange(newValue);
  };

  return (
    <Box sx={{ minWidth: 200 }}>
      <Autocomplete
        multiple
        options={availableTags || []}
        value={selectedTags}
        onChange={handleTagChange}
        disabled={disabled || isLoading}
        loading={isLoading}
        getOptionLabel={(option) => option.name}
        isOptionEqualToValue={(option, value) => option.id === value.id}
        renderTags={(value, getTagProps) =>
          value.map((option, index) => (
            <Chip
              {...getTagProps({ index })}
              key={option.id}
              label={option.name}
              size="small"
              sx={{ 
                fontSize: '0.75rem', 
                height: 24,
                backgroundColor: option.color || 'default',
                color: option.color ? 'white' : 'inherit',
                '& .MuiChip-label': {
                  color: option.color ? 'white' : 'inherit'
                }
              }}
            />
          ))
        }
        renderOption={(props, option) => (
          <Box component="li" {...props}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Box
                sx={{
                  width: 12,
                  height: 12,
                  borderRadius: '50%',
                  backgroundColor: option.color || '#ccc',
                }}
              />
              <Box>
                <Typography variant="body2">{option.name}</Typography>
                {option.description && (
                  <Typography variant="caption" color="text.secondary">
                    {option.description}
                  </Typography>
                )}
              </Box>
            </Box>
          </Box>
        )}
        renderInput={(params) => (
          <TextField
            {...params}
            placeholder="Choose tags..."
            size="small"
            variant="outlined"
          />
        )}
        noOptionsText="Tags not found"
        loadingText="Loading tags..."
      />
    </Box>
  );
};

export default TagFilter;
