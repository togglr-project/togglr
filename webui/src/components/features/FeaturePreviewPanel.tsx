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
  Tooltip,
} from '@mui/material';
import {
  AccessTime as TimeIcon,
  Person as PersonIcon,
  Edit as EditIcon,
  Add as AddIcon,
  Schedule as ScheduleIcon,
} from '@mui/icons-material';
import type { FeatureExtended, FeatureDetailsResponse, FlagVariant } from '../../generated/api/client';
import SimpleTimelinePreview from './SimpleTimelinePreview';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';

interface FeaturePreviewPanelProps {
  selectedFeature: FeatureExtended | null;
  projectId: string;
  onClose: () => void;
}

// Mock data for tags
const getMockTags = (featureId: string) => {
  const tagSets = [
    [
      { label: 'frontend', color: 'primary' as const, slug: 'frontend' },
      { label: 'experiment', color: 'secondary' as const, slug: 'experiment' },
      { label: 'beta', color: 'warning' as const, slug: 'beta' },
    ],
    [
      { label: 'backend', color: 'default' as const, slug: 'backend' },
      { label: 'critical', color: 'error' as const, slug: 'critical' },
      { label: 'stable', color: 'success' as const, slug: 'stable' },
    ],
    [
      { label: 'mobile', color: 'info' as const, slug: 'mobile' },
      { label: 'analytics', color: 'secondary' as const, slug: 'analytics' },
      { label: 'deprecated', color: 'default' as const, slug: 'deprecated' },
    ],
  ];
  
  const index = parseInt(featureId) % tagSets.length;
  return tagSets[index] || tagSets[0];
};

// Mock data for variants
const getMockVariants = (featureId: string, kind: string) => {
  if (kind !== 'multivariant') return [];
  
  const variantSets = [
    ['control', 'treatment', 'variant_a'],
    ['default', 'premium', 'beta'],
    ['v1', 'v2', 'experimental'],
  ];
  
  const index = parseInt(featureId) % variantSets.length;
  return variantSets[index] || variantSets[0];
};

// Mock data for history
const getMockHistory = (featureId: string) => {
  const histories = [
    [
      { action: 'Created', user: 'john.doe', timestamp: '2 hours ago', icon: AddIcon },
      { action: 'Updated', user: 'jane.smith', timestamp: '1 hour ago', icon: EditIcon },
      { action: 'Enabled', user: 'admin', timestamp: '30 minutes ago', icon: ScheduleIcon },
    ],
    [
      { action: 'Created', user: 'alice.wilson', timestamp: '1 day ago', icon: AddIcon },
      { action: 'Updated', user: 'bob.johnson', timestamp: '6 hours ago', icon: EditIcon },
      { action: 'Disabled', user: 'admin', timestamp: '2 hours ago', icon: ScheduleIcon },
      { action: 'Updated', user: 'charlie.brown', timestamp: '1 hour ago', icon: EditIcon },
    ],
    [
      { action: 'Created', user: 'diana.prince', timestamp: '3 days ago', icon: AddIcon },
      { action: 'Updated', user: 'bruce.wayne', timestamp: '1 day ago', icon: EditIcon },
    ],
  ];
  
  const index = parseInt(featureId) % histories.length;
  return histories[index] || histories[0];
};

const FeaturePreviewPanel: React.FC<FeaturePreviewPanelProps> = ({
  selectedFeature,
  projectId,
  onClose,
}) => {
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

  // Load feature details to get variants
  const { data: featureDetails } = useQuery<FeatureDetailsResponse>({
    queryKey: ['feature-details', selectedFeature.id],
    queryFn: async () => {
      const response = await apiClient.getFeature(selectedFeature.id);
      return response.data;
    },
    enabled: !!selectedFeature && selectedFeature.kind === 'multivariant',
  });

  const tags = getMockTags(selectedFeature.id);
  const variants = featureDetails?.variants?.map(v => v.name) || getMockVariants(selectedFeature.id, selectedFeature.kind);
  const history = getMockHistory(selectedFeature.id);

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
      <Box sx={{ mb: 2 }}>
        <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
          Tags
        </Typography>
        <Stack direction="row" spacing={0.5} flexWrap="wrap" gap={0.5}>
          {tags.map((tag) => (
            <Chip
              key={tag.slug}
              label={tag.label}
              color={tag.color}
              size="small"
              sx={{ fontSize: '0.7rem', height: 20 }}
            />
          ))}
        </Stack>
      </Box>

      {/* Default Variant/Value and Available Variants */}
      {selectedFeature.default_variant && (
        <Box sx={{ mb: 2 }}>
          <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
            {selectedFeature.kind === 'multivariant' ? 'Variants' : 'Default Value'}
          </Typography>
          <Stack direction="row" spacing={0.5} flexWrap="wrap" gap={0.5}>
            <Chip
              label={`default: ${selectedFeature.default_variant}`}
              variant="outlined"
              size="small"
              sx={{ fontSize: '0.7rem', height: 20 }}
            />
            {selectedFeature.kind === 'multivariant' && variants.length > 0 && (
              variants
                .filter(v => v !== selectedFeature.default_variant)
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
        />
      </Box>

      <Divider sx={{ my: 2 }} />

      {/* History */}
      <Box>
        <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
          History
        </Typography>
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
                  <Typography variant="body2" sx={{ fontSize: '0.8rem' }}>
                    {item.action} by {item.user}
                  </Typography>
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
      </Box>
    </Paper>
  );
};

export default FeaturePreviewPanel;
