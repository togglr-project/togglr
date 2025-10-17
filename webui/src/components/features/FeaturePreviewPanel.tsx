import React from 'react';
import {
  Box,
  Paper,
  Typography,
  Chip,
  Stack,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Avatar,
} from '@mui/material';
import {
  Edit as EditIcon,
  Add as AddIcon,
  Schedule as ScheduleIcon,
} from '@mui/icons-material';
import type { FeatureExtended, FeatureDetailsResponse, ListChangesResponse, ChangeGroup } from '../../generated/api/client';
import SimpleTimelinePreview from './SimpleTimelinePreview';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import { getFirstEnabledAlgorithmSlug } from '../../utils/algorithmUtils';

interface FeaturePreviewPanelProps {
  selectedFeature: FeatureExtended | null;
  projectId: string;
  environmentKey: string;
  onClose: () => void;
}

// Helper function to format timestamp
const formatTimestamp = (timestamp: string) => {
  const date = new Date(timestamp);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMinutes = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMinutes < 60) {
    return `${diffMinutes} minutes ago`;
  } else if (diffHours < 24) {
    return `${diffHours} hours ago`;
  } else {
    return `${diffDays} days ago`;
  }
};

// Helper function to get action icon
const getActionIcon = (action: string) => {
  switch (action) {
    case 'create':
      return AddIcon;
    case 'update':
      return EditIcon;
    case 'delete':
      return ScheduleIcon;
    default:
      return EditIcon;
  }
};

// Helper function to format action text
const formatActionText = (action: string, entity: string) => {
  const actionMap: { [key: string]: string } = {
    create: 'Created',
    update: 'Updated',
    delete: 'Deleted',
  };
  
  const entityMap: { [key: string]: string } = {
    feature: 'feature',
    rule: 'rule',
    flag_variant: 'variant',
    feature_schedule: 'schedule',
  };

  const actionText = actionMap[action] || action;
  const entityText = entityMap[entity] || entity;
  
  return `${actionText} ${entityText}`;
};

const FeaturePreviewPanel: React.FC<FeaturePreviewPanelProps> = ({
  selectedFeature,
  projectId,
  environmentKey,
}) => {
  // Load feature details to get variants and tags
  const { data: featureDetails } = useQuery<FeatureDetailsResponse>({
    queryKey: ['feature-details', selectedFeature?.id],
    queryFn: async () => {
      if (!selectedFeature) return null;
      const response = await apiClient.getFeature(selectedFeature.id, environmentKey);
      return response.data;
    },
    enabled: !!selectedFeature,
    staleTime: 0, // No caching - always fetch fresh data
    refetchOnWindowFocus: true, // Refetch when window gains focus
  });

  // Load feature changes history
  const { data: changesData } = useQuery<ListChangesResponse>({
    queryKey: ['feature-changes', selectedFeature?.id, projectId],
    queryFn: async () => {
      if (!selectedFeature) return null;
      const response = await apiClient.listProjectChanges(
        projectId,
        1, // page
        3, // perPage - limit to 3 events as requested
        undefined, // sortBy
        'desc', // sortOrder - newest first
        undefined, // actor
        undefined, // entity
        undefined, // action
        selectedFeature.id, // featureId - filter by specific feature
        undefined, // from
        undefined  // to
      );
      return response.data;
    },
    enabled: !!selectedFeature,
  });

  if (!selectedFeature) {
    return (
      <Paper
        sx={{
          p: 3,
          height: 'fit-content',
          minHeight: 400,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          border: '1px dashed',
          borderColor: 'divider',
        }}
      >
        <Typography variant="body1" color="text.secondary">
          Select any feature to view details
        </Typography>
      </Paper>
    );
  }

  // Use real tags from API response
  const tags = featureDetails?.tags || [];
  const variants = featureDetails?.variants?.map(v => v.name) || [];
  
  // Process changes data into history format
  const history = changesData?.items?.map((changeGroup: ChangeGroup) => {
    // Get the first change from the group to determine the main action
    const firstChange = changeGroup.changes[0];
    const actionText = firstChange ? formatActionText(firstChange.action, firstChange.entity) : 'Changed';
    const actionIcon = firstChange ? getActionIcon(firstChange.action) : EditIcon;
    
    return {
      action: actionText,
      user: changeGroup.username || 'Unknown',
      timestamp: formatTimestamp(changeGroup.created_at),
      icon: actionIcon,
      changesCount: changeGroup.changes.length,
    };
  }) || [];

  return (
    <Paper sx={{ p: 2, height: 'fit-content', minHeight: 400 }}>
      {/* Header */}
      <Box sx={{ mb: 2 }}>
        <Typography variant="h6" sx={{ fontWeight: 600, mb: 0.5 }}>
          {selectedFeature.name}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ fontFamily: 'monospace' }}>
          {selectedFeature.key}
        </Typography>
      </Box>

      {/* Tags */}
      {tags.length > 0 && (
        <Box sx={{ mb: 2 }}>
          <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
            Tags
          </Typography>
          <Stack direction="row" spacing={0.5} flexWrap="wrap" gap={0.5}>
            {tags.map((tag) => (
              <Chip
                key={tag.id}
                label={tag.slug}
                size="small"
                sx={{ 
                  fontSize: '0.7rem', 
                  height: 20,
                  backgroundColor: tag.color || 'default',
                  color: tag.color ? 'white' : 'inherit',
                  '& .MuiChip-label': {
                    color: tag.color ? 'white' : 'inherit'
                  }
                }}
              />
            ))}
          </Stack>
        </Box>
      )}

      {/* Algorithm */}
      {getFirstEnabledAlgorithmSlug(selectedFeature.algorithms) && (
        <Box sx={{ mb: 2 }}>
          <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
            Algorithm
          </Typography>
          <Chip
            label={`Algorithm: ${getFirstEnabledAlgorithmSlug(selectedFeature.algorithms)}`}
            color="info"
            variant="outlined"
            size="small"
            sx={{ fontSize: '0.7rem', height: 20 }}
          />
        </Box>
      )}

      {/* Default Variant/Value and Available Variants */}
      {selectedFeature.default_value && (
        <Box sx={{ mb: 2 }}>
          <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
            {selectedFeature.kind === 'multivariant' ? 'Variants' : 'Default Value'}
          </Typography>
          <Stack direction="row" spacing={0.5} flexWrap="wrap" gap={0.5}>
            <Chip
              label={`default: ${selectedFeature.default_value}`}
              variant="outlined"
              size="small"
              sx={{ fontSize: '0.7rem', height: 20 }}
            />
            {selectedFeature.kind === 'multivariant' && variants.length > 0 && (
              variants
                .filter(v => v !== selectedFeature.default_value)
                .map((variant) => (
                  <Chip
                    key={variant}
                    label={variant}
                    color="default"
                    size="small"
                    sx={{ fontSize: '0.7rem', height: 20 }}
                  />
                ))
            )}
          </Stack>
        </Box>
      )}

      {/* Timeline */}
      <Box sx={{ mb: 2 }}>
        <SimpleTimelinePreview
          featureId={selectedFeature.id}
          projectId={projectId}
          featureEnabled={selectedFeature.enabled}
          environmentKey={environmentKey}
        />
      </Box>

      <Divider sx={{ my: 2 }} />

      {/* History */}
      <Box>
        <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
          History
        </Typography>
        {history.length > 0 ? (
          <List dense sx={{ p: 0 }}>
            {history.map((item, index) => (
              <ListItem key={index} sx={{ px: 0, py: 0.5 }}>
                <ListItemIcon sx={{ minWidth: 32 }}>
                  <Avatar sx={{ width: 20, height: 20, bgcolor: 'action.hover' }}>
                    <item.icon sx={{ fontSize: 12 }} />
                  </Avatar>
                </ListItemIcon>
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" sx={{ fontSize: '0.8rem' }}>
                        {item.action} by {item.user}
                      </Typography>
                      {item.changesCount > 1 && (
                        <Chip
                          label={`+${item.changesCount - 1}`}
                          size="small"
                          sx={{ 
                            height: 16, 
                            fontSize: '0.6rem',
                            bgcolor: 'action.hover',
                            color: 'text.secondary'
                          }}
                        />
                      )}
                    </Box>
                  }
                  secondary={
                    <Typography variant="caption" color="text.secondary">
                      {item.timestamp}
                    </Typography>
                  }
                />
              </ListItem>
            ))}
          </List>
        ) : (
          <Typography variant="body2" color="text.secondary" sx={{ fontStyle: 'italic' }}>
            No changes recorded
          </Typography>
        )}
      </Box>
    </Paper>
  );
};

export default FeaturePreviewPanel;
