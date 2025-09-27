import React from 'react';
import {
  Box,
  Chip,
  Autocomplete,
  TextField,
  Typography,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { ProjectTag } from '../../generated/api/client';

interface TagSelectorProps {
  projectId: string;
  selectedTags: ProjectTag[];
  onChange: (tags: ProjectTag[]) => void;
  disabled?: boolean;
}

const TagSelector: React.FC<TagSelectorProps> = ({
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

  const handleTagChange = (_: any, newValue: ProjectTag[]) => {
    onChange(newValue);
  };

  return (
    <Box>
      <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
        Tags
      </Typography>
      <Autocomplete
        multiple
        options={availableTags || []}
        value={selectedTags}
        onChange={handleTagChange}
        disabled={disabled || isLoading}
        loading={isLoading}
        getOptionLabel={(option) => option.slug}
        isOptionEqualToValue={(option, value) => option.id === value.id}
        renderTags={(value, getTagProps) =>
          value.map((option, index) => {
            const tagProps = getTagProps({ index });
            const { key, ...restTagProps } = tagProps as any;
            return (
              <Chip
                {...restTagProps}
                key={option.id}
                label={option.slug}
                size="small"
                sx={{
                  fontSize: '0.7rem',
                  height: 20,
                  backgroundColor: option.color || 'default',
                  color: option.color ? 'white' : 'inherit',
                  '& .MuiChip-label': {
                    color: option.color ? 'white' : 'inherit',
                  },
                }}
              />
            );
          })
        }
        renderOption={(props, option) => {
          const { key, ...rest } = props as any;
          return (
            <Box component="li" key={key} {...rest}>
              <Box sx={{ display: 'flex', alignItems: 'center', width: '100%' }}>
                <Box
                  sx={{
                    width: 12,
                    height: 12,
                    borderRadius: '50%',
                    backgroundColor: option.color || '#3B82F6',
                    mr: 1,
                  }}
                />
                <Box sx={{ flexGrow: 1 }}>
                  <Typography variant="body2">{option.slug}</Typography>
                  {option.name !== option.slug && (
                    <Typography variant="caption" color="text.secondary">
                      {option.name}
                    </Typography>
                  )}
                </Box>
              </Box>
            </Box>
          );
        }}
        renderInput={(params) => (
          <TextField
            {...params}
            placeholder="Select tags..."
            variant="outlined"
            size="small"
          />
        )}
      />
    </Box>
  );
};

export default TagSelector;
