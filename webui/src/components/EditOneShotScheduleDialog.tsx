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
// @ts-ignore
import { listTimeZones, findTimeZone, getZonedTime, getUTCOffset } from 'timezone-support';
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

const allTimezones = listTimeZones();

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
  } catch (_) {
    return '';
  }
};

const fromDatetimeLocalInZoneToISO = (val?: string, tz?: string): string | undefined => {
  if (!val) return undefined;
  try {
    const [datePart, timePart] = val.split('T');
    if (!datePart || !timePart) return undefined;
    const [y, m, d] = datePart.split('-').map((s) => parseInt(s, 10));
    const [hh, mm] = timePart.split(':').map((s) => parseInt(s, 10));
    if (!y || !m || !d || isNaN(hh) || isNaN(mm)) return undefined;
    const tzObj = findTimeZone(tz || 'UTC');
    // Build a Date from the provided wall-time components as if they were in UTC,
    // then subtract the timezone offset to get the real UTC instant.
    const wallAsUTC = new Date(Date.UTC(y, m - 1, d, hh, mm, 0, 0));
    const offsetMinutes = getUTCOffset(tzObj, wallAsUTC);
    const utcDate = new Date(wallAsUTC.getTime() - offsetMinutes * 60 * 1000);
    return utcDate.toISOString();
  } catch (_) {
    return undefined;
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
    timezone: initialData?.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
    starts_at: initialData?.starts_at || '',
    ends_at: initialData?.ends_at || '',
    action: initialData?.action || 'enable'
  }));

  const [errors, setErrors] = useState<string[]>([]);

  // Reset form when dialog opens
  useEffect(() => {
    if (open && initialData) {
      const tz = initialData.timezone || 'UTC';
      const startsLocal = initialData.starts_at ? toDatetimeLocalInZone(initialData.starts_at, tz) : '';
      const endsLocal = initialData.ends_at ? toDatetimeLocalInZone(initialData.ends_at, tz) : '';
      
      setData({
        timezone: tz,
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
      // Convert datetime-local back to ISO
      const safeTz = allTimezones.includes(data.timezone) ? data.timezone : 'UTC';
      let convertedStarts = fromDatetimeLocalInZoneToISO(data.starts_at, safeTz);
      let convertedEnds = fromDatetimeLocalInZoneToISO(data.ends_at, safeTz);
      
      if (!convertedStarts && data.starts_at) {
        try { convertedStarts = new Date(data.starts_at).toISOString(); } catch { /* ignore */ }
      }
      if (!convertedEnds && data.ends_at) {
        try { convertedEnds = new Date(data.ends_at).toISOString(); } catch { /* ignore */ }
      }

      onSubmit({
        timezone: data.timezone,
        starts_at: convertedStarts || '',
        ends_at: convertedEnds || '',
        action: data.action
      });
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
                <Tooltip title="By default, the timezone is detected from your browser (Intl.DateTimeFormat). You can change it if you need a different zone for this schedule.">
                  <IconButton size="small">
                    <HelpIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              </Box>
              <TextField
                select
                fullWidth
                label="Timezone"
                value={data.timezone}
                onChange={(e) => {
                  const newTz = e.target.value;
                  setData(prev => ({
                    ...prev,
                    timezone: newTz,
                    // Reconvert dates to new timezone
                    starts_at: prev.starts_at ? toDatetimeLocalInZone(prev.starts_at, newTz) : prev.starts_at,
                    ends_at: prev.ends_at ? toDatetimeLocalInZone(prev.ends_at, newTz) : prev.ends_at,
                  }));
                }}
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
