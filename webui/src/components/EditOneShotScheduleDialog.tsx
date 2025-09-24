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
  Select,
  Tooltip,
  IconButton
} from '@mui/material';
import { Help as HelpIcon } from '@mui/icons-material';
// @ts-expect-error - timezone-support types are not available
import { findTimeZone, getZonedTime } from 'timezone-support';
import type { FeatureSchedule, FeatureScheduleAction } from '../generated/api/client';

interface EditOneShotScheduleData {
  timezone: string;
  starts_at: string;
  ends_at: string;
  action: FeatureScheduleAction;
}

interface EditOneShotScheduleDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: EditOneShotScheduleData) => void;
  initialData?: FeatureSchedule;
  existingSchedules?: Array<{ id: string; starts_at?: string; ends_at?: string }>;
}


// Helper functions for timezone conversion
const pad2 = (n: number) => (n < 10 ? '0' + n : '' + n);

const toDatetimeLocalInZone = (iso?: string, tz?: string): string => {
  if (!iso) return '';
  try {
    const date = new Date(iso);
    if (isNaN(date.getTime())) return '';
    const tzObj = findTimeZone(tz || 'UTC');
    const z = getZonedTime(date, tzObj);
    const yyyy = z.year;
    const MM = pad2(z.month);
    const dd = pad2(z.day);
    const HH = pad2(z.hours);
    const mm = pad2(z.minutes);
    return `${yyyy}-${MM}-${dd}T${HH}:${mm}`;
  } catch {
    return '';
  }
};


const EditOneShotScheduleDialog: React.FC<EditOneShotScheduleDialogProps> = ({
  open,
  onClose,
  onSubmit,
  initialData,
  existingSchedules = []
}) => {
  const [data, setData] = useState<EditOneShotScheduleData>(() => ({
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone, // Use browser timezone
    starts_at: initialData?.starts_at || '',
    ends_at: initialData?.ends_at || '',
    action: initialData?.action || 'enable'
  }));

  const [errors, setErrors] = useState<string[]>([]);

  // Reset form when dialog opens
  useEffect(() => {
    if (open && initialData) {
      const browserTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
      const startsLocal = initialData.starts_at ? toDatetimeLocalInZone(initialData.starts_at, browserTimezone) : '';
      const endsLocal = initialData.ends_at ? toDatetimeLocalInZone(initialData.ends_at, browserTimezone) : '';
      
      setData({
        timezone: browserTimezone,
        starts_at: startsLocal,
        ends_at: endsLocal,
        action: initialData.action
      });
      setErrors([]);
    }
  }, [open, initialData]);

  // Validate form
  useEffect(() => {
    const validationErrors = validateOneShotScheduleData(data, existingSchedules, initialData?.id);
    setErrors(validationErrors);
  }, [data, existingSchedules, initialData?.id]);

  const handleSubmit = () => {
    const validationErrors = validateOneShotScheduleData(data, existingSchedules, initialData?.id);
    if (validationErrors.length === 0) {
      const browserTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
      
      // Convert datetime-local to ISO (same logic as cron-like builder)
      let convertedStarts = '';
      let convertedEnds = '';
      
      if (data.starts_at) {
        const [dateStr, timeStr] = data.starts_at.split('T');
        if (dateStr && timeStr) {
          const [year, month, day] = dateStr.split('-').map(Number);
          const [hours, minutes] = timeStr.split(':').map(Number);
          const date = new Date(year, month - 1, day, hours, minutes, 0, 0);
          if (!isNaN(date.getTime())) {
            convertedStarts = date.toISOString();
          }
        }
      }
      
      if (data.ends_at) {
        const [dateStr, timeStr] = data.ends_at.split('T');
        if (dateStr && timeStr) {
          const [year, month, day] = dateStr.split('-').map(Number);
          const [hours, minutes] = timeStr.split(':').map(Number);
          const date = new Date(year, month - 1, day, hours, minutes, 0, 0);
          if (!isNaN(date.getTime())) {
            convertedEnds = date.toISOString();
          }
        }
      }

      onSubmit({
        timezone: browserTimezone,
        starts_at: convertedStarts,
        ends_at: convertedEnds,
        action: data.action
      });
    }
  };



  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>Edit one-shot schedule</DialogTitle>
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

          <Alert severity="warning" sx={{ mb: 3 }}>
            <Typography variant="body2">
              <strong>Important:</strong> Activate/Deactivate schedules only work when the feature's Master Enable switch is ON. 
              When Master Enable is OFF, the feature is completely disabled and schedules are ignored.
            </Typography>
          </Alert>

          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                  Action
                </Typography>
                <Tooltip title="Master Enable is the global on/off switch: when OFF the feature is always disabled and schedules are ignored. When ON, the feature is either controlled manually (if no schedules are defined) or automatically by schedules. You can create either one repeating schedule (using the friendly 'repeat' builder) or one-or-more non-overlapping one-shot intervals. Repeating schedules define periodic active windows (duration is required); one-shot intervals define exact start/end periods. Baseline (state outside scheduled windows) depends on schedule types: repeating-Activate → baseline OFF (activate only during windows); repeating-Deactivate → baseline ON (deactivate only during windows); for one-shots baseline is ON if any Deactivate interval exists, otherwise OFF. Newer schedules override older ones; if two are created at the same instant, Deactivate wins.">
                  <IconButton size="small">
                    <HelpIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              </Box>
              <FormControl fullWidth>
                <InputLabel>Action</InputLabel>
                <Select
                  value={data.action}
                  onChange={(e) => setData(prev => ({ ...prev, action: e.target.value as FeatureScheduleAction }))}
                  label="Action"
                  size="small"
                >
                  <MenuItem value="enable">Activate feature</MenuItem>
                  <MenuItem value="disable">Deactivate feature</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} md={6}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 1 }}>
                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                  User's timezone by default
                </Typography>
                <Tooltip title="By default, the timezone is detected from your browser">
                  <IconButton size="small">
                    <HelpIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              </Box>
              <TextField
                fullWidth
                size="small"
                label="Timezone"
                value={Intl.DateTimeFormat().resolvedOptions().timeZone}
                disabled
              />
            </Grid>

            {/* Start Date and Time */}
            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600, mt: 2 }}>
                Start Date & Time
              </Typography>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                size="small"
                label="Start Date"
                type="date"
                InputLabelProps={{ shrink: true }}
                value={data.starts_at ? data.starts_at.split('T')[0] : ''}
                onChange={(e) => {
                  const dateStr = e.target.value;
                  const currentTime = data.starts_at ? data.starts_at.split('T')[1] || '00:00' : '00:00';
                  setData(prev => ({ ...prev, starts_at: dateStr ? `${dateStr}T${currentTime}` : '' }));
                }}
                helperText="Schedule start date"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                size="small"
                label="Start Time"
                type="time"
                InputLabelProps={{ shrink: true }}
                value={data.starts_at ? data.starts_at.split('T')[1] || '00:00' : '00:00'}
                onChange={(e) => {
                  const timeStr = e.target.value;
                  const currentDate = data.starts_at ? data.starts_at.split('T')[0] : '';
                  if (currentDate && timeStr) {
                    setData(prev => ({ ...prev, starts_at: `${currentDate}T${timeStr}` }));
                  }
                }}
                inputProps={{ step: 60 }}
                helperText="Schedule start time"
              />
            </Grid>

            {/* End Date and Time */}
            <Grid item xs={12}>
              <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600, mt: 2 }}>
                End Date & Time
              </Typography>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                size="small"
                label="End Date"
                type="date"
                InputLabelProps={{ shrink: true }}
                value={data.ends_at ? data.ends_at.split('T')[0] : ''}
                onChange={(e) => {
                  const dateStr = e.target.value;
                  const currentTime = data.ends_at ? data.ends_at.split('T')[1] || '23:59' : '23:59';
                  setData(prev => ({ ...prev, ends_at: dateStr ? `${dateStr}T${currentTime}` : '' }));
                }}
                helperText="Schedule end date"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                size="small"
                label="End Time"
                type="time"
                InputLabelProps={{ shrink: true }}
                value={data.ends_at ? data.ends_at.split('T')[1] || '23:59' : '23:59'}
                onChange={(e) => {
                  const timeStr = e.target.value;
                  const currentDate = data.ends_at ? data.ends_at.split('T')[0] : '';
                  if (currentDate && timeStr) {
                    setData(prev => ({ ...prev, ends_at: `${currentDate}T${timeStr}` }));
                  }
                }}
                inputProps={{ step: 60 }}
                helperText="Schedule end time"
              />
            </Grid>
          </Grid>

          <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
            <Button
              size="small"
              variant="outlined"
              onClick={() => {
                const now = new Date();
                const dateStr = now.toISOString().slice(0, 10);
                const timeStr = now.toTimeString().slice(0, 5);
                setData(prev => ({ ...prev, starts_at: `${dateStr}T${timeStr}` }));
              }}
            >
              Now
            </Button>
            <Button
              size="small"
              variant="outlined"
              onClick={() => {
                const tomorrow = new Date();
                tomorrow.setDate(tomorrow.getDate() + 1);
                const dateStr = tomorrow.toISOString().slice(0, 10);
                const timeStr = tomorrow.toTimeString().slice(0, 5);
                setData(prev => ({ ...prev, starts_at: `${dateStr}T${timeStr}` }));
              }}
            >
              Tomorrow
            </Button>
            <Button
              size="small"
              variant="outlined"
              onClick={() => {
                const now = new Date();
                const tomorrow = new Date();
                tomorrow.setDate(tomorrow.getDate() + 1);
                const nowDateStr = now.toISOString().slice(0, 10);
                const nowTimeStr = now.toTimeString().slice(0, 5);
                const tomorrowDateStr = tomorrow.toISOString().slice(0, 10);
                const tomorrowTimeStr = tomorrow.toTimeString().slice(0, 5);
                setData(prev => ({ 
                  ...prev, 
                  starts_at: `${nowDateStr}T${nowTimeStr}`, 
                  ends_at: `${tomorrowDateStr}T${tomorrowTimeStr}` 
                }));
              }}
            >
              Now - Tomorrow
            </Button>
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} size="small">Cancel</Button>
        <Button
          variant="contained"
          onClick={handleSubmit}
          disabled={errors.length > 0}
          size="small"
        >
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};

/**
 * Validates one-shot schedule data
 */
function validateOneShotScheduleData(
  data: EditOneShotScheduleData,
  existingSchedules: Array<{ id: string; starts_at?: string; ends_at?: string }>,
  currentScheduleId?: string
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

      // Check for overlap with other schedules (excluding current one)
      const hasOverlap = existingSchedules.some(schedule => {
        if (schedule.id === currentScheduleId) return false; // Skip current schedule
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

export default EditOneShotScheduleDialog;
