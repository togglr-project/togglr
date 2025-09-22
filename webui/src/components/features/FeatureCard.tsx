import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  Switch,
  IconButton,
  Tooltip,
  Stack,
} from '@mui/material';
import {
  Edit as EditIcon,
  Schedule as ScheduleIcon,
  Visibility as ViewIcon,
} from '@mui/icons-material';
import type { FeatureExtended } from '../../generated/api/client';
import { getNextStateDescription } from '../../utils/timeUtils';

interface FeatureCardProps {
  feature: FeatureExtended;
  onEdit: (feature: FeatureExtended) => void;
  onView: (feature: FeatureExtended) => void;
  onToggle: (feature: FeatureExtended) => void;
  canToggle?: boolean;
}

const FeatureCard: React.FC<FeatureCardProps> = ({
  feature,
  onEdit,
  onView,
  onToggle,
  canToggle = true,
}) => {
  const getKindColor = (kind: string) => {
    switch (kind) {
      case 'simple': return 'warning';
      case 'multivariant': return 'primary';
      default: return 'default';
    }
  };

  const getStatusColor = (enabled: boolean) => {
    return enabled ? 'success' : 'default';
  };

  return (
    <Card 
      sx={{ 
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
        transition: 'all 0.2s ease-in-out',
        '&:hover': {
          boxShadow: 2,
          transform: 'translateY(-1px)',
        },
        border: '1px solid',
        borderColor: feature.enabled ? 'success.light' : 'divider',
        bgcolor: feature.enabled ? 'success.50' : 'background.paper',
        minHeight: 80,
      }}
    >
      <CardContent sx={{ flexGrow: 1, p: 2, '&:last-child': { pb: 2 } }}>
        {/* Main content - horizontal layout */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
          {/* Left side - Name and key */}
          <Box sx={{ flexGrow: 1, minWidth: 0 }}>
            <Typography 
              variant="h6" 
              sx={{ 
                fontWeight: 600,
                fontSize: '1rem',
                mb: 0.5,
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              }}
              title={feature.name}
            >
              {feature.name}
            </Typography>
            <Typography 
              variant="body2" 
              color="text.secondary"
              sx={{ 
                fontSize: '0.8rem',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              }}
              title={feature.key}
            >
              {feature.key}
            </Typography>
          </Box>

          {/* Middle - Chips with separators */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, overflow: 'hidden' }}>
            <Chip
              size="small"
              label={feature.kind}
              color={getKindColor(feature.kind)}
              sx={{ 
                fontSize: '0.7rem',
                height: 20,
                textTransform: 'capitalize',
              }}
            />
            
            <Box sx={{ width: 1, height: 12, bgcolor: 'divider', opacity: 0.5 }} />
            
            {feature.default_variant && (
              <>
                <Chip
                  size="small"
                  label={`default: ${feature.default_variant}`}
                  variant="outlined"
                  sx={{ 
                    fontSize: '0.7rem',
                    height: 20,
                  }}
                />
                <Box sx={{ width: 1, height: 12, bgcolor: 'divider', opacity: 0.5 }} />
              </>
            )}
            
            <Chip
              size="small"
              label={feature.is_active ? 'active' : 'not active'}
              color={feature.is_active ? 'success' : 'default'}
              sx={{ 
                fontSize: '0.7rem',
                height: 20,
              }}
            />
            
            {feature.next_state !== undefined && feature.next_state_time && (
              <>
                <Box sx={{ width: 1, height: 12, bgcolor: 'divider', opacity: 0.5 }} />
                <Chip 
                  size="small" 
                  icon={<ScheduleIcon />}
                  label={getNextStateDescription(feature.next_state, feature.next_state_time) || 'Scheduled'} 
                  color={feature.next_state ? 'info' : 'warning'}
                  variant="outlined"
                  sx={{ 
                    fontSize: '0.7rem',
                    height: 20,
                  }}
                />
              </>
            )}
          </Box>

          {/* Right side - Switch and actions */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Switch
              size="small"
              checked={feature.enabled}
              onChange={() => onToggle(feature)}
              disabled={!canToggle}
            />
            
            {/* Actions */}
            <Box sx={{ display: 'flex', gap: 0.5, ml: 1 }}>
              <Tooltip title="Edit feature">
                <IconButton size="small" onClick={() => onEdit(feature)}>
                  <EditIcon fontSize="small" />
                </IconButton>
              </Tooltip>
              <Tooltip title="View details">
                <IconButton size="small" onClick={() => onView(feature)}>
                  <ViewIcon fontSize="small" />
                </IconButton>
              </Tooltip>
            </Box>
          </Box>
        </Box>
      </CardContent>
    </Card>
  );
};

export default FeatureCard;