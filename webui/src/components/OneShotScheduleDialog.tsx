import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  MenuItem,
  Grid,
  Typography,
  Alert,
  Box,
  FormControl,
  InputLabel,
  Select
} from '@mui/material';
// @ts-ignore
import { listTimeZones } from 'timezone-support';
import type { FeatureScheduleAction } from '../generated/api/client';

interface OneShotScheduleData {
  timezone: string;
  starts_at: string;
  ends_at: string;
  action: FeatureScheduleAction;
}

interface OneShotScheduleDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: OneShotScheduleData) => void;
  initialData?: Partial<OneShotScheduleData>;
  existingSchedules?: Array<{ starts_at?: string; ends_at?: string }>;
}

const allTimezones = listTimeZones();

const OneShotScheduleDialog: React.FC<OneShotScheduleDialogProps> = ({
  open,
  onClose,
  onSubmit,
  initialData,
  existingSchedules = []
}) => {
  const [data, setData] = useState<OneShotScheduleData>(() => ({
    timezone: initialData?.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
    starts_at: initialData?.starts_at || '',
    ends_at: initialData?.ends_at || '',
    action: initialData?.action || 'enable'
  }));

  const [errors, setErrors] = useState<string[]>([]);

  // Сброс формы при открытии
  useEffect(() => {
    if (open) {
      setData({
        timezone: initialData?.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
        starts_at: initialData?.starts_at || '',
        ends_at: initialData?.ends_at || '',
        action: initialData?.action || 'enable'
      });
      setErrors([]);
    }
  }, [open, initialData]);

  // Валидация при изменении данных
  useEffect(() => {
    const validationErrors = validateOneShotSchedule(data, existingSchedules);
    setErrors(validationErrors);
  }, [data, existingSchedules]);

  const handleSubmit = () => {
    const validationErrors = validateOneShotSchedule(data, existingSchedules);
    if (validationErrors.length === 0) {
      onSubmit(data);
    }
  };

  const formatLocalDateTime = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  };

  const getCurrentDateTime = () => {
    return formatLocalDateTime(new Date());
  };

  const getTomorrowDateTime = () => {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    return formatLocalDateTime(tomorrow);
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>Create one-shot schedule</DialogTitle>
      <DialogContent dividers>
        <Box sx={{ mt: 1 }}>
          {errors.length > 0 && (
            <Alert severity="error" sx={{ mb: 2 }}>
              <Typography variant="subtitle2">Validation errors:</Typography>
              <ul>
                {errors.map((error, index) => (
                  <li key={index}>{error}</li>
                ))}
              </ul>
            </Alert>
          )}

          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth>
                <InputLabel>Action</InputLabel>
                <Select
                  value={data.action}
                  onChange={(e) => setData(prev => ({ ...prev, action: e.target.value as FeatureScheduleAction }))}
                  label="Action"
                >
                  <MenuItem value="enable">Enable feature</MenuItem>
                  <MenuItem value="disable">Disable feature</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                select
                fullWidth
                label="Timezone"
                value={data.timezone}
                onChange={(e) => setData(prev => ({ ...prev, timezone: e.target.value }))}
              >
                {allTimezones.map((tz: string) => (
                  <MenuItem key={tz} value={tz}>{tz}</MenuItem>
                ))}
              </TextField>
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Start time"
                type="datetime-local"
                InputLabelProps={{ shrink: true }}
                value={data.starts_at}
                onChange={(e) => setData(prev => ({ ...prev, starts_at: e.target.value }))}
                helperText="Schedule start time"
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="End time"
                type="datetime-local"
                InputLabelProps={{ shrink: true }}
                value={data.ends_at}
                onChange={(e) => setData(prev => ({ ...prev, ends_at: e.target.value }))}
                helperText="Schedule end time"
              />
            </Grid>
          </Grid>

          <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
            <Button
              size="small"
              variant="outlined"
              onClick={() => {
                const now = getCurrentDateTime();
                setData(prev => ({ ...prev, starts_at: now }));
              }}
            >
              Now
            </Button>
            <Button
              size="small"
              variant="outlined"
              onClick={() => {
                const tomorrow = getTomorrowDateTime();
                setData(prev => ({ ...prev, starts_at: tomorrow }));
              }}
            >
              Tomorrow
            </Button>
            <Button
              size="small"
              variant="outlined"
              onClick={() => {
                const now = getCurrentDateTime();
                const tomorrow = getTomorrowDateTime();
                setData(prev => ({ ...prev, starts_at: now, ends_at: tomorrow }));
              }}
            >
              Now - Tomorrow
            </Button>
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          variant="contained"
          onClick={handleSubmit}
          disabled={errors.length > 0}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};

/**
 * Валидирует данные разового расписания
 */
function validateOneShotSchedule(
  data: OneShotScheduleData,
  existingSchedules: Array<{ starts_at?: string; ends_at?: string }>
): string[] {
  const errors: string[] = [];

  if (!data.timezone) {
    errors.push('Timezone is required');
  }

  if (!data.starts_at) {
    errors.push('Start time is required');
  }

  if (!data.ends_at) {
    errors.push('End time is required');
  }

  if (data.starts_at && data.ends_at) {
    const startDate = new Date(data.starts_at);
    const endDate = new Date(data.ends_at);

    if (isNaN(startDate.getTime())) {
      errors.push('Invalid start time');
    }

    if (isNaN(endDate.getTime())) {
      errors.push('Invalid end time');
    }

    if (!isNaN(startDate.getTime()) && !isNaN(endDate.getTime())) {
      if (startDate >= endDate) {
        errors.push('Start time must be before end time');
      }

      // Check for overlap with existing schedules
      const hasOverlap = existingSchedules.some(schedule => {
        if (!schedule.starts_at || !schedule.ends_at) return false;
        
        const existingStart = new Date(schedule.starts_at);
        const existingEnd = new Date(schedule.ends_at);
        
        if (isNaN(existingStart.getTime()) || isNaN(existingEnd.getTime())) return false;

        // Check for interval overlap
        return !(endDate <= existingStart || startDate >= existingEnd);
      });

      if (hasOverlap) {
        errors.push('Schedule overlaps with existing schedule');
      }
    }
  }

  return errors;
}

export default OneShotScheduleDialog;
