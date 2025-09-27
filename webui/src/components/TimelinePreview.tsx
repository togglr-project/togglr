import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  CircularProgress,
  Alert,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem
} from '@mui/material';
import apiClient from '../api/apiClient';
import TimelineChart from './TimelineChart';

interface TimelinePreviewProps {
  featureId: string;
  environmentKey: string;
  schedules: Array<{
    startsAt?: string;
    endsAt?: string;
    cronExpr?: string;
    timezone: string;
    action: string;
    cronDuration?: string | { value: number; unit: string };
  }>;
}

const TimelinePreview: React.FC<TimelinePreviewProps> = ({ featureId, environmentKey, schedules }) => {
  const [timelineData, setTimelineData] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [timeRange, setTimeRange] = useState<1 | 3 | 7>(1);

  const generateTimeline = async () => {
    if (!featureId || schedules.length === 0) return;

    setIsLoading(true);
    setError(null);

    try {
      const from = new Date();
      const to = new Date();
      to.setDate(to.getDate() + timeRange);


      const response = await apiClient.testFeatureTimeline(
        featureId,
        environmentKey,
        from.toISOString(),
        to.toISOString(),
        Intl.DateTimeFormat().resolvedOptions().timeZone,
        {
          schedules: schedules.map(schedule => ({
            starts_at: schedule.startsAt ? new Date(schedule.startsAt).toISOString() : undefined,
            ends_at: schedule.endsAt ? new Date(schedule.endsAt).toISOString() : undefined,
            cron_expr: schedule.cronExpr,
            timezone: schedule.timezone,
            action: schedule.action as 'enable' | 'disable',
            cron_duration: typeof schedule.cronDuration === 'string' 
              ? schedule.cronDuration 
              : schedule.cronDuration 
                ? `${schedule.cronDuration.value}${schedule.cronDuration.unit.charAt(0)}`
                : undefined
          }))
        }
      );

      console.log('Timeline response:', response);
      setTimelineData(response.data);
    } catch (err: any) {
      console.error('Failed to generate timeline:', err);
      console.error('Error details:', {
        message: err.message,
        code: err.code,
        status: err.response?.status,
        statusText: err.response?.statusText,
        data: err.response?.data
      });
      setError(`Failed to generate timeline preview: ${err.message}`);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    generateTimeline();
  }, [featureId, environmentKey, schedules, timeRange]);

  if (isLoading) {
    return (
      <Paper sx={{ p: 2, mt: 2 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <CircularProgress size={20} />
          <Typography variant="body2">Generating timeline preview...</Typography>
        </Box>
      </Paper>
    );
  }

  if (error) {
    return (
      <Paper sx={{ p: 2, mt: 2 }}>
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
        <Button variant="outlined" size="small" onClick={generateTimeline}>
          Retry
        </Button>
      </Paper>
    );
  }

  if (!timelineData || !timelineData.events || timelineData.events.length === 0) {
    return (
      <Paper sx={{ p: 2, mt: 2 }}>
        <Typography variant="body2" color="text.secondary">
          No timeline events to display
        </Typography>
      </Paper>
    );
  }

  // Create a mock feature for TimelineChart
  const mockFeature = {
    id: featureId,
    name: 'Preview Feature',
    key: 'preview-feature',
    description: 'Timeline preview',
    enabled: true,
    kind: 'simple' as const,
    project_id: '',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    default_variant: '',
    is_active: true
  };

  // Convert timeline events to the format expected by TimelineChart
  const timelineEvents = timelineData.events.map((event: any) => ({
    time: event.time,
    enabled: event.enabled
  }));

  return (
    <Paper sx={{ p: 2, mt: 2 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="subtitle1">
          Timeline Preview
        </Typography>
        
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel>Time Range</InputLabel>
          <Select
            value={timeRange}
            label="Time Range"
            onChange={(e) => setTimeRange(e.target.value as 1 | 3 | 7)}
          >
            <MenuItem value={1}>1 Day</MenuItem>
            <MenuItem value={3}>3 Days</MenuItem>
            <MenuItem value={7}>7 Days</MenuItem>
          </Select>
        </FormControl>
      </Box>
      
      <TimelineChart
        features={[mockFeature]}
        timelines={{ [featureId]: timelineEvents }}
        isLoading={false}
        error={null}
        from={new Date().toISOString()}
        to={new Date(Date.now() + timeRange * 24 * 60 * 60 * 1000).toISOString()}
      />
      
      <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
        Showing {timelineEvents.length} events over the next {timeRange} day{timeRange > 1 ? 's' : ''}
      </Typography>
    </Paper>
  );
};

export default TimelinePreview;
