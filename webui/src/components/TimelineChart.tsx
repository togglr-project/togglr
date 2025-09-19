import React, { useMemo } from 'react';
import {
  Box,
  Paper,
  Typography,
  useTheme,
  Tooltip,
  LinearProgress,
  Alert
} from '@mui/material';
import type { FeatureTimelineEvent, FeatureExtended } from '../generated/api/client';

interface TimelineChartProps {
  features: FeatureExtended[];
  timelines: Record<string, FeatureTimelineEvent[]>;
  isLoading?: boolean;
  error?: string;
  from: string;
  to: string;
}

interface TimelineSegment {
  start: number;
  end: number;
  enabled: boolean;
}

interface ProcessedFeature {
  feature: FeatureExtended;
  segments: TimelineSegment[];
  totalDuration: number;
}

const TimelineChart: React.FC<TimelineChartProps> = ({
  features,
  timelines,
  isLoading = false,
  error,
  from,
  to
}) => {
  const theme = useTheme();

  const processedFeatures = useMemo(() => {
    if (!features || features.length === 0) return [];

    const fromTime = new Date(from).getTime();
    const toTime = new Date(to).getTime();

    return features.map(feature => {
      const events = timelines[feature.id] || [];
      
      if (events.length === 0) {
        return {
          feature,
          segments: [],
          totalDuration: 0
        };
      }

      // Sort events by time
      const sortedEvents = [...events].sort((a, b) => 
        new Date(a.time).getTime() - new Date(b.time).getTime()
      );

      // Process events into segments
      const segments: TimelineSegment[] = [];
      let currentState = false; // Assume feature starts as disabled
      let segmentStart = fromTime;

      for (let i = 0; i < sortedEvents.length; i++) {
        const event = sortedEvents[i];
        const eventTime = new Date(event.time).getTime();

        // If this is the first event and it's not at the very beginning
        if (i === 0 && eventTime > fromTime) {
          // Add initial segment (disabled state)
          segments.push({
            start: fromTime,
            end: eventTime,
            enabled: false
          });
        }

        // If state changed, close previous segment and start new one
        if (event.enabled !== currentState) {
          if (i > 0) {
            segments.push({
              start: segmentStart,
              end: eventTime,
              enabled: currentState
            });
          }
          currentState = event.enabled;
          segmentStart = eventTime;
        }
      }

      // Add final segment
      const lastEventTime = sortedEvents.length > 0 
        ? new Date(sortedEvents[sortedEvents.length - 1].time).getTime()
        : fromTime;
      
      segments.push({
        start: lastEventTime,
        end: toTime,
        enabled: currentState
      });

      // Calculate total enabled duration
      const totalDuration = segments.reduce((total, segment) => {
        if (segment.enabled) {
          return total + (segment.end - segment.start);
        }
        return total;
      }, 0);

      return {
        feature,
        segments,
        totalDuration
      };
    }).sort((a, b) => b.totalDuration - a.totalDuration); // Sort by total enabled duration
  }, [features, timelines, from, to]);

  const timeRange = useMemo(() => {
    const fromTime = new Date(from).getTime();
    const toTime = new Date(to).getTime();
    const now = Date.now();
    return { start: fromTime, end: toTime, now };
  }, [from, to]);

  const formatTime = (timestamp: number) => {
    const date = new Date(timestamp);
    return date.toLocaleString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZoneName: 'short'
    });
  };

  const getSegmentColor = (enabled: boolean) => {
    return enabled ? theme.palette.success.main : theme.palette.grey[400];
  };

  if (error) {
    return (
      <Alert severity="error" sx={{ mt: 2 }}>
        Error loading timelines: {error}
      </Alert>
    );
  }

  if (isLoading) {
    return (
      <Box sx={{ mt: 2 }}>
        <LinearProgress />
        <Typography variant="body2" color="text.secondary" sx={{ mt: 1, textAlign: 'center' }}>
          Loading feature timelines...
        </Typography>
      </Box>
    );
  }

  if (processedFeatures.length === 0) {
    return (
      <Box sx={{ mt: 2, p: 3, textAlign: 'center' }}>
        <Typography variant="body2" color="text.secondary">
          No features selected or no timeline data available.
        </Typography>
      </Box>
    );
  }

  const totalDuration = timeRange.end - timeRange.start;

  return (
    <Box sx={{ mt: 2 }}>
      <Paper sx={{ p: 2 }}>
        <Typography variant="h6" gutterBottom>
          Feature Timelines
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          Showing feature states from {formatTime(timeRange.start)} to {formatTime(timeRange.end)}
          <br />
          <strong>Timezone:</strong> {Intl.DateTimeFormat().resolvedOptions().timeZone}
        </Typography>

        <Box sx={{ overflowX: 'auto' }}>
          <Box sx={{ minWidth: 1000, position: 'relative', px: 1 }}>
            {/* Time axis */}
            <Box sx={{ 
              height: 20, 
              position: 'relative', 
              mb: 3,
              borderBottom: `1px solid ${theme.palette.divider}`
            }}>
              {/* Current time indicator */}
              <Box sx={{
                position: 'absolute',
                left: `${((timeRange.now - timeRange.start) / totalDuration) * 100}%`,
                top: 0,
                bottom: 0,
                width: 2,
                backgroundColor: theme.palette.primary.main,
                zIndex: 10
              }} />
              
              {/* Time labels */}
              <Box sx={{ 
                position: 'absolute', 
                top: '100%', 
                left: 0, 
                right: 0, 
                display: 'flex', 
                justifyContent: 'space-between',
                fontSize: '0.75rem',
                color: 'text.secondary',
                mt: 1,
                px: 1
              }}>
                <span>{formatTime(timeRange.start)}</span>
                <span>Now</span>
                <span>{formatTime(timeRange.end)}</span>
              </Box>
            </Box>

            {/* Feature timelines */}
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
              {processedFeatures.map(({ feature, segments }) => (
                <Box key={feature.id} sx={{ display: 'flex', alignItems: 'center', minHeight: 50 }}>
                  {/* Feature name */}
                  <Box sx={{ 
                    minWidth: 250, 
                    pr: 3, 
                    display: 'flex', 
                    flexDirection: 'column',
                    justifyContent: 'center'
                  }}>
                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                      {feature.name}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {feature.key}
                    </Typography>
                  </Box>

                  {/* Timeline */}
                  <Box sx={{ 
                    flex: 1, 
                    height: 32, 
                    position: 'relative',
                    backgroundColor: theme.palette.grey[100],
                    borderRadius: 1,
                    overflow: 'hidden',
                    ml: 1
                  }}>
                    {segments.map((segment, index) => {
                      const leftPercent = ((segment.start - timeRange.start) / totalDuration) * 100;
                      const widthPercent = ((segment.end - segment.start) / totalDuration) * 100;
                      
                      return (
                        <Tooltip
                          key={index}
                          title={
                            <Box>
                              <Typography variant="body2">
                                {segment.enabled ? 'Enabled' : 'Disabled'}
                              </Typography>
                              <Typography variant="caption">
                                {formatTime(segment.start)} - {formatTime(segment.end)}
                              </Typography>
                            </Box>
                          }
                          arrow
                        >
                          <Box
                            sx={{
                              position: 'absolute',
                              left: `${leftPercent}%`,
                              width: `${widthPercent}%`,
                              height: '100%',
                              backgroundColor: getSegmentColor(segment.enabled),
                              cursor: 'pointer',
                              transition: 'opacity 0.2s',
                              '&:hover': {
                                opacity: 0.8
                              }
                            }}
                          />
                        </Tooltip>
                      );
                    })}
                  </Box>
                </Box>
              ))}
            </Box>

            {/* Legend */}
            <Box sx={{ 
              display: 'flex', 
              gap: 3, 
              mt: 2, 
              pt: 2, 
              borderTop: `1px solid ${theme.palette.divider}` 
            }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Box sx={{ 
                  width: 16, 
                  height: 16, 
                  backgroundColor: theme.palette.success.main,
                  borderRadius: 0.5
                }} />
                <Typography variant="caption">Enabled</Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <Box sx={{ 
                  width: 16, 
                  height: 16, 
                  backgroundColor: theme.palette.grey[400],
                  borderRadius: 0.5
                }} />
                <Typography variant="caption">Disabled</Typography>
              </Box>
            </Box>
          </Box>
        </Box>
      </Paper>
    </Box>
  );
};

export default TimelineChart;
