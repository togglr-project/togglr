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
  LinearProgress,
  Fade,
} from '@mui/material';
import {
  Edit as EditIcon,
  Schedule as ScheduleIcon,
  Visibility as ViewIcon,
  Pending as PendingIcon,
} from '@mui/icons-material';
import type { FeatureExtended } from '../../generated/api/client';
import { getNextStateDescription } from '../../utils/timeUtils';
import { useFeatureHasPendingChanges } from '../../hooks/useProjectPendingChanges';

interface FeatureCardProps {
  feature: FeatureExtended;
  onEdit: (feature: FeatureExtended) => void;
  onView: (feature: FeatureExtended) => void;
  onToggle: (feature: FeatureExtended) => void;
  onSelect?: (feature: FeatureExtended) => void;
  canToggle?: boolean;
  isToggling?: boolean;
  isSelected?: boolean;
  projectId?: string;
}

const FeatureCard: React.FC<FeatureCardProps> = ({
  feature,
  onEdit,
  onView,
  onToggle,
  onSelect,
  canToggle = true,
  isToggling = false,
  isSelected = false,
  projectId,
}) => {
  // Check if feature has pending changes
  const hasPendingChanges = useFeatureHasPendingChanges(feature.id, projectId);
  const getKindColor = (kind: string) => {
    switch (kind) {
      case 'simple': return 'warning';
      case 'multivariant': return 'primary';
      default: return 'default';
    }
  };



  const handleCardClick = (e: React.MouseEvent) => {
    // Don't trigger selection if clicking on buttons or switch
    if (
      (e.target as HTMLElement).closest('button') ||
      (e.target as HTMLElement).closest('[role="switch"]')
    ) {
      return;
    }
    onSelect?.(feature);
  };

  return (
    <Card 
      onClick={handleCardClick}
      sx={{ 
        display: 'flex',
        flexDirection: 'row',
        alignItems: 'center',
        transition: 'all 0.2s ease-in-out',
        cursor: hasPendingChanges ? 'not-allowed' : 'pointer',
        '&:hover': {
          boxShadow: hasPendingChanges ? 1 : 2,
          transform: hasPendingChanges ? 'none' : 'translateY(-1px)',
        },
        border: '2px solid',
        borderColor: isSelected 
          ? 'primary.main' 
          : hasPendingChanges
            ? 'warning.main'
            : feature.enabled 
              ? 'success.light' 
              : 'divider',
        bgcolor: isSelected 
          ? 'primary.50' 
          : hasPendingChanges
            ? 'warning.50'
            : feature.enabled 
              ? 'success.50' 
              : 'background.paper',
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

            {hasPendingChanges && (
              <>
                <Box sx={{ width: 1, height: 12, bgcolor: 'divider', opacity: 0.5 }} />
                <Tooltip title="This feature has pending changes awaiting approval">
                  <Chip 
                    size="small" 
                    icon={<PendingIcon />}
                    label="Pending" 
                    color="warning"
                    variant="filled"
                    sx={{ 
                      fontSize: '0.7rem',
                      height: 20,
                    }}
                  />
                </Tooltip>
              </>
            )}
          </Box>

          {/* Right side - Switch and actions */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Box sx={{ position: 'relative' }}>
              <Switch
                size="small"
                checked={feature.enabled}
                onChange={() => onToggle(feature)}
                disabled={!canToggle || isToggling || hasPendingChanges}
              />
              <Fade in={isToggling}>
                <Box
                  sx={{
                    position: 'absolute',
                    top: '50%',
                    left: '50%',
                    transform: 'translate(-50%, -50%)',
                    width: 20,
                    height: 20,
                    borderRadius: '50%',
                    bgcolor: 'background.paper',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                >
                  <LinearProgress
                    sx={{
                      width: 16,
                      height: 16,
                      borderRadius: '50%',
                      color: 'primary.main',
                      '& .MuiLinearProgress-bar': {
                        borderRadius: '50%',
                      },
                    }}
                  />
                </Box>
              </Fade>
            </Box>
            
            {/* Actions */}
            <Box sx={{ display: 'flex', gap: 0.5, ml: 1 }}>
              <Tooltip title={hasPendingChanges ? "Cannot edit: feature has pending changes" : "Edit feature"}>
                <span>
                  <IconButton 
                    size="small" 
                    onClick={() => onEdit(feature)}
                    disabled={hasPendingChanges}
                    sx={{ 
                      opacity: hasPendingChanges ? 0.5 : 1,
                      cursor: hasPendingChanges ? 'not-allowed' : 'pointer'
                    }}
                  >
                    <EditIcon fontSize="small" />
                  </IconButton>
                </span>
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