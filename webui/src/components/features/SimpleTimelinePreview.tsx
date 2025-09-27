import React from 'react';
import {
  Box,
  Typography,
  CircularProgress,
  Alert,
} from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import type { FeatureTimelineResponse } from '../../generated/api/client';

interface SimpleTimelinePreviewProps {
  featureId: string;
  projectId: string;
  featureEnabled: boolean;
  environmentKey: string;
}

const SimpleTimelinePreview: React.FC<SimpleTimelinePreviewProps> = ({
  featureId,
  projectId,
  featureEnabled,
  environmentKey,
}) => {
  const { data: timeline, isLoading, error } = useQuery<FeatureTimelineResponse>({
    queryKey: ['feature-timeline', projectId, featureId, environmentKey],
    queryFn: async () => {
      const now = new Date();
      const from = now.toISOString();
      const to = new Date(now.getTime() + 60 * 60 * 1000).toISOString(); // +1 hour
      const location = Intl.DateTimeFormat().resolvedOptions().timeZone; // User's timezone
      
      const response = await apiClient.getFeatureTimeline(
        featureId,
        environmentKey,
        from,
        to,
        location
      );
      return response.data;
    },
    enabled: !!featureId && !!projectId,
  });

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <CircularProgress size={16} />
        <Typography variant="caption" color="text.secondary">
          Loading timeline...
        </Typography>
      </Box>
    );
  }

  if (error) {
    // Check if it's a 404 error (no schedule)
    const is404 = (error as { response?: { status?: number } })?.response?.status === 404;
    
    if (is404) {
      return (
        <Box>
          <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
            Timeline for the next hour
          </Typography>
          <Typography variant="caption" color="text.secondary">
            No schedule
          </Typography>
        </Box>
      );
    }
    
    return (
      <Alert severity="warning" sx={{ py: 0.5 }}>
        <Typography variant="caption">
          Timeline unavailable
        </Typography>
      </Alert>
    );
  }

  if (!timeline || !timeline.events || timeline.events.length === 0) {
    return (
      <Box>
        <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
          Timeline for the next hour
        </Typography>
        <Typography variant="caption" color="text.secondary">
          No scheduled changes
        </Typography>
      </Box>
    );
  }

  // Process timeline events into segments like TimelineChart does
  const now = new Date();
  const fromTime = now.getTime();
  const toTime = now.getTime() + 60 * 60 * 1000; // +1 hour
  
  // Sort events by time
  const sortedEvents = [...timeline.events].sort((a, b) => 
    new Date(a.time).getTime() - new Date(b.time).getTime()
  );

  // Process events into segments
  const segments: Array<{ start: number; end: number; enabled: boolean }> = [];
  let currentState = featureEnabled; // Start with current enabled state from props
  let segmentStart = fromTime;

  for (let i = 0; i < sortedEvents.length; i++) {
    const event = sortedEvents[i];
    const eventTime = new Date(event.time).getTime();

    // If this is the first event
    if (i === 0) {
      // If first event is not at the very beginning, add initial disabled segment
      if (eventTime > fromTime) {
        segments.push({
          start: fromTime,
          end: eventTime,
          enabled: false
        });
      }
      // Set the initial state based on the first event
      currentState = event.enabled;
      segmentStart = eventTime;
    } else {
      // For subsequent events, always close previous segment and start new one
      segments.push({
        start: segmentStart,
        end: eventTime,
        enabled: currentState
      });
      currentState = event.enabled;
      segmentStart = eventTime;
    }
  }

  // Always add final segment from the last event to the end
  if (sortedEvents.length > 0) {
    const lastEventTime = new Date(sortedEvents[sortedEvents.length - 1].time).getTime();
    segments.push({
      start: lastEventTime,
      end: toTime,
      enabled: currentState
    });
  } else {
    // If no events, create a single disabled segment for the entire range
    segments.push({
      start: fromTime,
      end: toTime,
      enabled: false
    });
  }

  // Find current segment and next change
  const currentTime = now.getTime();
  const currentSegment = segments.find(segment => 
    currentTime >= segment.start && currentTime < segment.end
  );
  
  const nextChange = segments.find(segment => segment.start > currentTime);
  const nextChangeTime = nextChange ? new Date(nextChange.start) : null;
  const timeUntilNext = nextChangeTime ? nextChangeTime.getTime() - currentTime : null;

  return (
    <Box>
      <Typography variant="subtitle2" sx={{ mb: 1, color: 'text.secondary' }}>
        Timeline for the next hour
      </Typography>
      
      {/* Timeline visualization */}
      <Box
        sx={{
          height: 8,
          bgcolor: 'divider',
          position: 'relative',
          overflow: 'hidden',
          mb: 1,
        }}
      >
        {segments.map((segment, index) => {
          const leftPercent = ((segment.start - fromTime) / (toTime - fromTime)) * 100;
          const widthPercent = ((segment.end - segment.start) / (toTime - fromTime)) * 100;
          
          return (
            <Box
              key={index}
              sx={{
                position: 'absolute',
                left: `${leftPercent}%`,
                width: `${widthPercent}%`,
                height: '100%',
                bgcolor: segment.enabled ? 'success.main' : 'grey.400',
              }}
            />
          );
        })}
        
        {/* Current time indicator */}
        <Box
          sx={{
            position: 'absolute',
            left: `${((currentTime - fromTime) / (toTime - fromTime)) * 100}%`,
            top: 0,
            bottom: 0,
            width: 2,
            bgcolor: 'primary.main',
            zIndex: 10,
          }}
        />
      </Box>
    </Box>
  );
};

export default SimpleTimelinePreview;
