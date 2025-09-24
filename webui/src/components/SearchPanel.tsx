import React, { useState } from 'react';
import {
  Box,
  TextField,
  IconButton,
  Collapse,
  Stack,
  Chip,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Typography,
  Tooltip,
} from '@mui/material';
import {
  Search as SearchIcon,
  FilterList as FilterIcon,
  Clear as ClearIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
} from '@mui/icons-material';
import type { ProjectTag } from '../generated/api/client';
import TagFilter from './features/TagFilter';

export interface FilterOption {
  key: string;
  label: string;
  value: any;
  options: { value: any; label: string }[];
  onChange: (value: any) => void;
}

export interface SearchPanelProps {
  searchValue: string;
  onSearchChange: (value: string) => void;
  filters?: FilterOption[];
  quickFilters?: { label: string; value: any; active: boolean; onClick: () => void }[];
  placeholder?: string;
  showAdvancedFilters?: boolean;
  onToggleAdvanced?: (expanded: boolean) => void;
  // Tag filter props
  projectId?: string;
  selectedTags?: ProjectTag[];
  onTagsChange?: (tags: ProjectTag[]) => void;
  showTagFilter?: boolean;
}

const SearchPanel: React.FC<SearchPanelProps> = ({
  searchValue,
  onSearchChange,
  filters = [],
  quickFilters = [],
  placeholder = "Search...",
  showAdvancedFilters = true,
  onToggleAdvanced,
  projectId,
  selectedTags = [],
  onTagsChange,
  showTagFilter = false,
}) => {
  const [expanded, setExpanded] = useState(false);

  const handleToggle = () => {
    const newExpanded = !expanded;
    setExpanded(newExpanded);
    onToggleAdvanced?.(newExpanded);
  };

  const handleClearSearch = () => {
    onSearchChange('');
  };

  const activeQuickFilters = quickFilters.filter(f => f.active);
  const hasActiveFilters = activeQuickFilters.length > 0 || 
    filters.some(f => f.value !== '' && f.value !== null && f.value !== undefined) ||
    (showTagFilter && selectedTags.length > 0);

  return (
    <Box sx={{ mb: 2 }}>
      {/* Main search bar */}
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <TextField
          fullWidth
          size="small"
          placeholder={placeholder}
          value={searchValue}
          onChange={(e) => onSearchChange(e.target.value)}
          InputProps={{
            startAdornment: <SearchIcon sx={{ color: 'text.secondary', mr: 1 }} />,
            endAdornment: searchValue && (
              <IconButton size="small" onClick={handleClearSearch}>
                <ClearIcon fontSize="small" />
              </IconButton>
            ),
          }}
        />
        
        {showAdvancedFilters && (
          <Tooltip title={expanded ? "Hide filters" : "Show filters"}>
            <IconButton
              size="small"
              onClick={handleToggle}
              color={hasActiveFilters ? "primary" : "default"}
              sx={{ 
                border: hasActiveFilters ? '1px solid' : '1px solid transparent',
                borderColor: hasActiveFilters ? 'primary.main' : 'divider',
                minWidth: 40,
                height: 40,
              }}
            >
              <FilterIcon fontSize="small" />
              {expanded ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
            </IconButton>
          </Tooltip>
        )}
      </Box>

      {/* Quick filters and selected tags */}
      {(quickFilters.length > 0 || (showTagFilter && selectedTags.length > 0)) && (
        <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
          {quickFilters.map((filter, index) => (
            <Chip
              key={index}
              label={filter.label}
              size="small"
              variant={filter.active ? "filled" : "outlined"}
              color={filter.active ? "primary" : "default"}
              onClick={filter.onClick}
              sx={{ 
                fontSize: '0.75rem',
                height: 24,
                '& .MuiChip-label': { px: 1 },
              }}
            />
          ))}
          {showTagFilter && selectedTags.map((tag) => (
            <Chip
              key={tag.id}
              label={tag.name}
              size="small"
              variant="filled"
              color="secondary"
              onDelete={() => onTagsChange?.(selectedTags.filter(t => t.id !== tag.id))}
              sx={{ 
                fontSize: '0.75rem',
                height: 24,
                backgroundColor: tag.color || 'secondary.main',
                color: tag.color ? 'white' : 'inherit',
                '& .MuiChip-label': { px: 1 },
                '& .MuiChip-deleteIcon': {
                  color: tag.color ? 'white' : 'inherit',
                },
              }}
            />
          ))}
        </Box>
      )}

      {/* Advanced filters */}
      <Collapse in={expanded}>
        <Box sx={{ mt: 2, p: 2, bgcolor: 'background.paper', borderRadius: 1, border: '1px solid', borderColor: 'divider' }}>
          <Typography variant="subtitle2" sx={{ mb: 1.5, color: 'text.secondary' }}>
            Advanced Filters
          </Typography>
          <Stack direction={{ xs: 'column', md: 'row' }} spacing={1.5}>
            {filters.map((filter) => (
              <FormControl key={filter.key} size="small" sx={{ minWidth: 160 }}>
                <InputLabel id={`${filter.key}-label`}>{filter.label}</InputLabel>
                <Select
                  labelId={`${filter.key}-label`}
                  label={filter.label}
                  value={filter.value}
                  onChange={(e) => filter.onChange(e.target.value)}
                  size="small"
                >
                  {filter.options.map((option) => (
                    <MenuItem key={option.value} value={option.value}>
                      {option.label}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            ))}
            {showTagFilter && projectId && onTagsChange && (
              <Box>
                <TagFilter
                  projectId={projectId}
                  selectedTags={selectedTags}
                  onChange={onTagsChange}
                />
              </Box>
            )}
          </Stack>
        </Box>
      </Collapse>
    </Box>
  );
};

export default SearchPanel;
