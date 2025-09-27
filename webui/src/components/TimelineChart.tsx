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

        <Box sx={{ 
          width: '100%',
          overflow: 'hidden'
        }}>
          <Box sx={{ 
            position: 'relative', 
            px: 1,
            width: '100%'
          }}>
            {/* Time axis */}
            <Box sx={{ 
              height: 20, 
              position: 'relative', 
              mb: 3,
              borderBottom: `1px solid ${theme.palette.divider}`,
              ml: { xs: '200px', sm: '250px' }, // Responsive margin
              width: { xs: 'calc(100% - 200px)', sm: 'calc(100% - 250px)' },
              minWidth: 0
            }}>
              {/* Current time indicator */}
              {timeRange.now >= timeRange.start && timeRange.now <= timeRange.end && (
                <Box sx={{
                  position: 'absolute',
                  left: `${((timeRange.now - timeRange.start) / totalDuration) * 100}%`,
                  top: 0,
                  bottom: 0,
                  width: 2,
                  backgroundColor: theme.palette.primary.main,
                  zIndex: 10
                }} />
              )}
              
              {/* Time labels */}
              <Box sx={{ 
                position: 'absolute', 
                top: '100%', 
                left: 0, 
                right: 0, 
                display: 'flex', 
                justifyContent: 'space-between',
                fontSize: '0.75rem',
                color: 'primary.light',
                mt: 1,
                px: 1
              }}>
                <span>{formatTime(timeRange.start)}</span>
                {timeRange.now >= timeRange.start && timeRange.now <= timeRange.end ? (
                  <span></span>
                ) : (
                  <span>{formatTime(timeRange.start + totalDuration / 2)}</span>
                )}
                <span>{formatTime(timeRange.end)}</span>
              </Box>
              
              {/* Time grid lines */}
              {Array.from({ length: 5 }, (_, i) => {
                const time = timeRange.start + (totalDuration / 4) * i;
                const leftPercent = ((time - timeRange.start) / totalDuration) * 100;
                return (
                  <Box
                    key={i}
                    sx={{
                      position: 'absolute',
                      left: `${leftPercent}%`,
                      top: 0,
                      bottom: 0,
                      width: 1,
                      backgroundColor: theme.palette.divider,
                      opacity: 0.5
                    }}
                  />
                );
              })}
              
              {/* Extended grid lines to cover feature timelines */}
              {Array.from({ length: 5 }, (_, i) => {
                const time = timeRange.start + (totalDuration / 4) * i;
                const leftPercent = ((time - timeRange.start) / totalDuration) * 100;
                return (
                  <Box
                    key={`extended-${i}`}
                    sx={{
                      position: 'absolute',
                      left: `${leftPercent}%`,
                      top: '100%',
                      height: '100%', // Use relative height instead of fixed
                      width: 1,
                      backgroundColor: theme.palette.divider,
                      opacity: 0.3
                    }}
                  />
                );
              })}
            </Box>

            {/* Feature timelines */}
            <Box sx={{ 
              display: 'flex', 
              flexDirection: 'column', 
              gap: 2, 
              mt: 2,
              width: '100%'
            }}>
              {processedFeatures.map(({ feature, segments }) => (
                <Box key={feature.id} sx={{ 
                  display: 'flex', 
                  alignItems: 'center', 
                  minHeight: 50,
                  width: '100%',
                  overflow: 'hidden'
                }}>
                  {/* Feature name */}
                  <Box sx={{ 
                    minWidth: { xs: 200, sm: 250 }, 
                    pr: 3, 
                    display: 'flex', 
                    flexDirection: 'column',
                    justifyContent: 'center',
                    flexShrink: 0
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
                    ml: 1,
                    minWidth: 0,
                    width: '100%'
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
