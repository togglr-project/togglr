import React, { useState, useEffect } from 'react';
import {
  Box,
  Stepper,
  Step,
  StepLabel,
  Button,
  Typography,
  TextField,
  MenuItem,
  FormControl,
  FormLabel,
  RadioGroup,
  FormControlLabel,
  Radio,
  Grid,
  Paper,
  Alert,
  Divider,
  Tooltip,
  IconButton
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  ArrowForward as ArrowForwardIcon,
  Check as CheckIcon,
  Help as HelpIcon
} from '@mui/icons-material';
// @ts-ignore
import { listTimeZones, findTimeZone, getUTCOffset } from 'timezone-support';
import { isValidCron } from 'cron-validator';
import cronstrue from 'cronstrue';
import {
  type ScheduleBuilderData,
  type ScheduleType,
  generateCronExpression,
  generateScheduleDescription,
  validateScheduleData
} from '../utils/cronGenerator';
import type { FeatureSchedule, FeatureScheduleAction } from '../generated/api/client';

interface EditRecurringScheduleBuilderProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (data: ScheduleBuilderData & { cronExpression: string }) => void;
  initialData?: FeatureSchedule;
}

const allTimezones = listTimeZones();

const steps = [
  'Timezone',
  'Date Range',
  'Schedule Type',
  'Parameters',
  'Duration',
  'Action',
  'Preview'
];

// Helper function to parse cron expression back to schedule data
const parseCronToScheduleData = (cronExpr: string, timezone: string, action: FeatureScheduleAction, cronDuration: string): Partial<ScheduleBuilderData> => {
  const duration = parseDuration(cronDuration);
  // This is a simplified parser - in a real app you'd want more sophisticated parsing
  // For now, we'll return a basic structure and let the user adjust
  const parts = cronExpr.split(' ');
  if (parts.length !== 5) return {};

  const [minutes, hours, dayOfMonth, month, dayOfWeek] = parts;

  // Try to detect pattern
  if (minutes.startsWith('*/') && hours === '*' && dayOfMonth === '*' && month === '*' && dayOfWeek === '*') {
    // Every N minutes
    const interval = parseInt(minutes.substring(2));
    return {
      scheduleType: 'repeat_every',
      repeatEvery: { interval, unit: 'minutes' },
      duration
    };
  } else if (minutes === '0' && hours.startsWith('*/') && dayOfMonth === '*' && month === '*' && dayOfWeek === '*') {
    // Every N hours
    const interval = parseInt(hours.substring(2));
    return {
      scheduleType: 'repeat_every',
      repeatEvery: { interval, unit: 'hours' },
      duration
    };
  } else if (minutes !== '*' && hours !== '*' && dayOfMonth === '*' && month === '*' && dayOfWeek === '*') {
    // Daily at specific time
    return {
      scheduleType: 'daily',
      daily: { time: `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}` },
      duration
    };
  } else if (minutes !== '*' && hours !== '*' && dayOfMonth !== '*' && month === '*' && dayOfWeek === '*') {
    // Monthly on specific day
    return {
      scheduleType: 'monthly',
      monthly: { 
        dayOfMonth: parseInt(dayOfMonth), 
        time: `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}` 
      },
      duration
    };
  } else if (minutes !== '*' && hours !== '*' && dayOfMonth !== '*' && month !== '*' && dayOfWeek === '*') {
    // Yearly on specific date
    return {
      scheduleType: 'yearly',
      yearly: { 
        month: parseInt(month), 
        day: parseInt(dayOfMonth), 
        time: `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}` 
      },
      duration
    };
  }

  return { duration };
};

// Helper function to parse duration
const parseDuration = (duration: string): { value: number; unit: 'minutes' | 'hours' | 'days' } => {
  if (!duration) return { value: 1, unit: 'hours' };
  
  const match = duration.match(/^(\d+)([mhdD])$/i);
  if (match) {
    const value = parseInt(match[1]);
    const unitChar = match[2].toLowerCase();
    const unit = unitChar === 'm' ? 'minutes' : unitChar === 'h' ? 'hours' : 'days';
    return { value, unit };
  }
  
  return { value: 1, unit: 'hours' };
};

const EditRecurringScheduleBuilder: React.FC<EditRecurringScheduleBuilderProps> = ({
  open,
  onClose,
  onSubmit,
  initialData
}) => {
  const [activeStep, setActiveStep] = useState(0);
  const [data, setData] = useState<ScheduleBuilderData>(() => {
    if (initialData?.cron_expr) {
      const parsed = parseCronToScheduleData(initialData.cron_expr, initialData.timezone, initialData.action, initialData.cron_duration || '');
      return {
        timezone: initialData.timezone,
        startsAt: initialData.starts_at,
        endsAt: initialData.ends_at,
        scheduleType: (parsed.scheduleType as ScheduleType) || 'repeat_every',
        repeatEvery: parsed.repeatEvery || { interval: 1, unit: 'minutes' },
        daily: parsed.daily || { time: '09:00' },
        monthly: parsed.monthly || { dayOfMonth: 1, time: '09:00' },
        yearly: parsed.yearly || { month: 1, day: 1, time: '09:00' },
        duration: parsed.duration || parseDuration(initialData.cron_duration || ''),
        action: initialData.action,
        ...parsed
      };
    }
    
    return {
      timezone: initialData?.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
      startsAt: initialData?.starts_at,
      endsAt: initialData?.ends_at,
      scheduleType: 'repeat_every',
      repeatEvery: { interval: 1, unit: 'minutes' },
      daily: { time: '09:00' },
      monthly: { dayOfMonth: 1, time: '09:00' },
      yearly: { month: 1, day: 1, time: '09:00' },
      duration: parseDuration(initialData?.cron_duration || '') || { value: 1, unit: 'hours' },
      action: initialData?.action || 'enable'
    };
  });

  const [errors, setErrors] = useState<string[]>([]);
  const [cronExpression, setCronExpression] = useState<string>('');
  const [cronDescription, setCronDescription] = useState<string>('');

  // Validation and cron generation when data changes
  useEffect(() => {
    const validationErrors = validateScheduleData(data);
    setErrors(validationErrors);

    if (validationErrors.length === 0) {
      try {
        const cron = generateCronExpression(data);
        setCronExpression(cron);
        
        // Check cron validity and generate description
        if (isValidCron(cron, { seconds: false, allowBlankDay: true, alias: true })) {
          setCronDescription(cronstrue.toString(cron));
        } else {
          setCronDescription('');
        }
      } catch (error) {
        setCronExpression('');
        setCronDescription('');
      }
    } else {
      setCronExpression('');
      setCronDescription('');
    }
  }, [data]);

  const handleNext = () => {
    if (activeStep < steps.length - 1) {
      setActiveStep(activeStep + 1);
    }
  };

  const handleBack = () => {
    if (activeStep > 0) {
      setActiveStep(activeStep - 1);
    }
  };

  const handleSubmit = () => {
    if (errors.length === 0 && cronExpression) {
      // Convert dates to the selected timezone
      let finalStartsAt = data.startsAt;
      let finalEndsAt = data.endsAt;

      if (data.startsAt) {
        try {
          const tzObj = findTimeZone(data.timezone);
          if (tzObj) {
            const [year, month, day] = data.startsAt.split('T')[0].split('-').map(Number);
            const localDate = new Date(year, month - 1, day, 0, 0, 0, 0);
            const offset = getUTCOffset(tzObj, localDate);
            finalStartsAt = new Date(localDate.getTime() + offset * 60 * 1000).toISOString();
          }
        } catch (error) {
          console.error('Error converting start date:', error);
        }
      }

      if (data.endsAt) {
        try {
          const tzObj = findTimeZone(data.timezone);
          if (tzObj) {
            const [year, month, day] = data.endsAt.split('T')[0].split('-').map(Number);
            const localDate = new Date(year, month - 1, day, 23, 59, 59, 999); // End of day
            const offset = getUTCOffset(tzObj, localDate);
            finalEndsAt = new Date(localDate.getTime() + offset * 60 * 1000).toISOString();
          }
        } catch (error) {
          console.error('Error converting end date:', error);
        }
      }

      onSubmit({
        ...data,
        startsAt: finalStartsAt,
        endsAt: finalEndsAt,
        cronExpression
      });
    }
  };

  const canProceed = () => {
    switch (activeStep) {
      case 0: // Timezone
        return !!data.timezone;
      case 1: // Date Range
        return !!data.startsAt; // startsAt is required, endsAt is optional
      case 2: // Schedule Type
        return !!data.scheduleType;
      case 3: // Parameters
        switch (data.scheduleType) {
          case 'repeat_every':
            return data.repeatEvery && data.repeatEvery.interval > 0;
          case 'daily':
            return data.daily && !!data.daily.time;
          case 'monthly':
            return data.monthly && data.monthly.dayOfMonth > 0 && !!data.monthly.time;
          case 'yearly':
            return data.yearly && data.yearly.month > 0 && data.yearly.day > 0 && !!data.yearly.time;
          default:
            return false;
        }
      case 4: // Duration
        if (data.duration.value <= 0) return false;
        
        // Additional validation for repeat_every schedule type
        if (data.scheduleType === 'repeat_every' && data.repeatEvery) {
          const convertToMinutes = (value: number, unit: 'minutes' | 'hours' | 'days'): number => {
            switch (unit) {
              case 'minutes': return value;
              case 'hours': return value * 60;
              case 'days': return value * 24 * 60;
              default: return value;
            }
          };
          
          const repeatIntervalMinutes = convertToMinutes(data.repeatEvery.interval, data.repeatEvery.unit);
          const durationMinutes = convertToMinutes(data.duration.value, data.duration.unit);
          
          return durationMinutes < repeatIntervalMinutes;
        }
        
        return true;
      case 5: // Action
        return !!data.action;
      case 6: // Preview
        return errors.length === 0 && !!cronExpression;
      default:
        return false;
    }
  };

  const renderStepContent = () => {
    switch (activeStep) {
      case 0:
        return renderTimezoneStep();
      case 1:
        return renderDateRangeStep();
      case 2:
        return renderScheduleTypeStep();
      case 3:
        return renderParametersStep();
      case 4:
        return renderDurationStep();
      case 5:
        return renderActionStep();
      case 6:
        return renderPreviewStep();
      default:
        return null;
    }
  };

  const renderTimezoneStep = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        Select timezone
      </Typography>
      <TextField
        select
        fullWidth
        label="Timezone"
        value={data.timezone}
        onChange={(e) => setData(prev => ({ ...prev, timezone: e.target.value }))}
        sx={{ mt: 2 }}
      >
        {allTimezones.map((tz: string) => (
          <MenuItem key={tz} value={tz}>{tz}</MenuItem>
        ))}
      </TextField>
    </Box>
  );

  const renderDateRangeStep = () => {
    // Simple date conversion helpers
    const toDateString = (isoString: string): string => {
      if (!isoString) return '';
      try {
        const date = new Date(isoString);
        if (isNaN(date.getTime())) return '';
        return date.toISOString().slice(0, 10);
      } catch (error) {
        return '';
      }
    };

    const fromDateStringToISO = (dateString: string): string => {
      if (!dateString) return '';
      try {
        // Create date at midnight UTC
        const date = new Date(dateString + 'T00:00:00.000Z');
        if (isNaN(date.getTime())) return '';
        return date.toISOString();
      } catch (error) {
        return '';
      }
    };

    return (
      <Box>
        <Typography variant="h6" gutterBottom>
          Set schedule date range
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          Define when the schedule should start and optionally when it should end.
        </Typography>
        
        <Grid container spacing={2}>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              label="Start date"
              type="date"
              value={data.startsAt ? toDateString(data.startsAt) : ''}
              onChange={(e) => {
                const isoString = fromDateStringToISO(e.target.value);
                setData(prev => ({ ...prev, startsAt: isoString }));
              }}
              InputLabelProps={{
                shrink: true,
              }}
              required
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              label="End date (optional)"
              type="date"
              value={data.endsAt ? toDateString(data.endsAt) : ''}
              onChange={(e) => {
                const isoString = e.target.value ? fromDateStringToISO(e.target.value) : '';
                setData(prev => ({ ...prev, endsAt: isoString || undefined }));
              }}
              InputLabelProps={{
                shrink: true,
              }}
            />
          </Grid>
        </Grid>
        
        <Typography variant="caption" color="text.secondary" sx={{ mt: 1, display: 'block' }}>
          Dates will be converted to the selected timezone when saving the schedule.
        </Typography>
      </Box>
    );
  };

  const renderScheduleTypeStep = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        Select schedule type
      </Typography>
      <FormControl component="fieldset" sx={{ mt: 2 }}>
        <RadioGroup
          value={data.scheduleType}
          onChange={(e) => setData(prev => ({ ...prev, scheduleType: e.target.value as ScheduleType }))}
        >
          <FormControlLabel
            value="repeat_every"
            control={<Radio />}
            label="Repeat every N minutes/hours"
          />
          <FormControlLabel
            value="daily"
            control={<Radio />}
            label="At fixed time daily"
          />
          <FormControlLabel
            value="monthly"
            control={<Radio />}
            label="At fixed day monthly"
          />
          <FormControlLabel
            value="yearly"
            control={<Radio />}
            label="Once a year"
          />
        </RadioGroup>
      </FormControl>
    </Box>
  );

  const renderParametersStep = () => {
    switch (data.scheduleType) {
      case 'repeat_every':
        return (
          <Box>
            <Typography variant="h6" gutterBottom>
              Repeat parameters
            </Typography>
            <Grid container spacing={2} sx={{ mt: 1 }}>
              <Grid item xs={6}>
                <TextField
                  fullWidth
                  label="Interval"
                  type="number"
                  value={data.repeatEvery?.interval || ''}
                  onChange={(e) => setData(prev => ({
                    ...prev,
                    repeatEvery: {
                      ...prev.repeatEvery!,
                      interval: parseInt(e.target.value) || 0
                    }
                  }))}
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  select
                  fullWidth
                  label="Unit"
                  value={data.repeatEvery?.unit || 'minutes'}
                  onChange={(e) => setData(prev => ({
                    ...prev,
                    repeatEvery: {
                      ...prev.repeatEvery!,
                      unit: e.target.value as 'minutes' | 'hours'
                    }
                  }))}
                >
                  <MenuItem value="minutes">Minutes</MenuItem>
                  <MenuItem value="hours">Hours</MenuItem>
                </TextField>
              </Grid>
            </Grid>
          </Box>
        );

      case 'daily':
        return (
          <Box>
            <Typography variant="h6" gutterBottom>
              Execution time
            </Typography>
            <TextField
              fullWidth
              label="Time"
              type="time"
              value={data.daily?.time || ''}
              onChange={(e) => setData(prev => ({
                ...prev,
                daily: { time: e.target.value }
              }))}
              InputLabelProps={{ shrink: true }}
              sx={{ mt: 2 }}
            />
          </Box>
        );

      case 'monthly':
        return (
          <Box>
            <Typography variant="h6" gutterBottom>
              Monthly schedule parameters
            </Typography>
            <Grid container spacing={2} sx={{ mt: 1 }}>
              <Grid item xs={6}>
                <TextField
                  fullWidth
                  label="Day of month"
                  type="number"
                  inputProps={{ min: 1, max: 31 }}
                  value={data.monthly?.dayOfMonth || ''}
                  onChange={(e) => setData(prev => ({
                    ...prev,
                    monthly: {
                      ...prev.monthly!,
                      dayOfMonth: parseInt(e.target.value) || 1
                    }
                  }))}
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  fullWidth
                  label="Time"
                  type="time"
                  value={data.monthly?.time || ''}
                  onChange={(e) => setData(prev => ({
                    ...prev,
                    monthly: {
                      ...prev.monthly!,
                      time: e.target.value
                    }
                  }))}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
            </Grid>
          </Box>
        );

      case 'yearly':
        return (
          <Box>
            <Typography variant="h6" gutterBottom>
              Yearly schedule parameters
            </Typography>
            <Grid container spacing={2} sx={{ mt: 1 }}>
              <Grid item xs={4}>
                <TextField
                  select
                  fullWidth
                  label="Month"
                  value={data.yearly?.month || ''}
                  onChange={(e) => setData(prev => ({
                    ...prev,
                    yearly: {
                      ...prev.yearly!,
                      month: parseInt(e.target.value) || 1
                    }
                  }))}
                >
                  {Array.from({ length: 12 }, (_, i) => (
                    <MenuItem key={i + 1} value={i + 1}>
                      {new Date(0, i).toLocaleString('en', { month: 'long' })}
                    </MenuItem>
                  ))}
                </TextField>
              </Grid>
              <Grid item xs={4}>
                <TextField
                  fullWidth
                  label="Day"
                  type="number"
                  inputProps={{ min: 1, max: 31 }}
                  value={data.yearly?.day || ''}
                  onChange={(e) => setData(prev => ({
                    ...prev,
                    yearly: {
                      ...prev.yearly!,
                      day: parseInt(e.target.value) || 1
                    }
                  }))}
                />
              </Grid>
              <Grid item xs={4}>
                <TextField
                  fullWidth
                  label="Time"
                  type="time"
                  value={data.yearly?.time || ''}
                  onChange={(e) => setData(prev => ({
                    ...prev,
                    yearly: {
                      ...prev.yearly!,
                      time: e.target.value
                    }
                  }))}
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
            </Grid>
          </Box>
        );

      default:
        return null;
    }
  };

  const renderDurationStep = () => {
    // Helper functions for unit conversion
    const convertToMinutes = (value: number, unit: 'minutes' | 'hours' | 'days'): number => {
      switch (unit) {
        case 'minutes': return value;
        case 'hours': return value * 60;
        case 'days': return value * 24 * 60;
        default: return value;
      }
    };

    const getMaxDurationForRepeatEvery = () => {
      if (data.scheduleType !== 'repeat_every' || !data.repeatEvery) {
        return { maxValue: Infinity, allowedUnits: ['minutes', 'hours', 'days'] as const };
      }

      const repeatIntervalMinutes = convertToMinutes(data.repeatEvery.interval, data.repeatEvery.unit);
      
      // Duration must be strictly less than repeat interval
      const maxDurationMinutesForDisplay = Math.max(repeatIntervalMinutes - 1, 0);
      
      // Allowed units rule: duration unit must be <= interval unit
      let allowedUnits: ('minutes' | 'hours' | 'days')[] = [];
      if (data.repeatEvery.unit === 'minutes') {
        allowedUnits = ['minutes'];
      } else if (data.repeatEvery.unit === 'hours') {
        allowedUnits = ['minutes', 'hours'];
      }

      return { maxValue: maxDurationMinutesForDisplay, allowedUnits, repeatIntervalMinutes } as any;
    };

    const { maxValue, allowedUnits, repeatIntervalMinutes } = getMaxDurationForRepeatEvery() as unknown as { maxValue: number; allowedUnits: ('minutes'|'hours'|'days')[]; repeatIntervalMinutes?: number };
    const currentDurationMinutes = convertToMinutes(data.duration.value, data.duration.unit);
    const limit = repeatIntervalMinutes ?? Infinity;
    const isValidDuration = currentDurationMinutes < limit && allowedUnits.includes(data.duration.unit);

    return (
      <Box>
        <Typography variant="h6" gutterBottom>
          Execution duration
        </Typography>
        
        {data.scheduleType === 'repeat_every' && data.repeatEvery && (
          <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
            Duration must be less than the repeat interval ({data.repeatEvery.interval} {data.repeatEvery.unit}).
            Maximum allowed: {maxValue} minutes.
          </Typography>
        )}

        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={6}>
            <TextField
              fullWidth
              label="Value"
              type="number"
              value={data.duration.value}
              onChange={(e) => setData(prev => ({
                ...prev,
                duration: {
                  ...prev.duration,
                  value: parseInt(e.target.value) || 0
                }
              }))}
              error={!isValidDuration}
              helperText={!isValidDuration ? 'Duration must be less than repeat interval' : ''}
            />
          </Grid>
          <Grid item xs={6}>
            <TextField
              select
              fullWidth
              label="Unit"
              value={data.duration.unit}
              onChange={(e) => setData(prev => ({
                ...prev,
                duration: {
                  ...prev.duration,
                  unit: e.target.value as 'minutes' | 'hours' | 'days'
                }
              }))}
              error={!isValidDuration}
            >
              {allowedUnits.map(unit => (
                <MenuItem key={unit} value={unit}>
                  {unit.charAt(0).toUpperCase() + unit.slice(1)}
                </MenuItem>
              ))}
            </TextField>
          </Grid>
        </Grid>
      </Box>
    );
  };

  const renderActionStep = () => (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
        <Typography variant="h6" gutterBottom>
          Action
        </Typography>
        <Tooltip title="Master Enable is the global on/off switch: when OFF the feature is always disabled and schedules are ignored. When ON, the feature is either controlled manually (if no schedules are defined) or automatically by schedules. You can create either one repeating schedule (using the friendly 'repeat' builder) or one-or-more non-overlapping one-shot intervals. Repeating schedules define periodic active windows (duration is required); one-shot intervals define exact start/end periods. Baseline (state outside scheduled windows) depends on schedule types: repeating-Activate → baseline OFF (activate only during windows); repeating-Deactivate → baseline ON (deactivate only during windows); for one-shots baseline is ON if any Deactivate interval exists, otherwise OFF. Newer schedules override older ones; if two are created at the same instant, Deactivate wins.">
          <IconButton size="small">
            <HelpIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>
      <FormControl component="fieldset" sx={{ mt: 2 }}>
        <RadioGroup
          value={data.action}
          onChange={(e) => setData(prev => ({ ...prev, action: e.target.value as 'enable' | 'disable' }))}
        >
          <FormControlLabel
            value="enable"
            control={<Radio />}
            label="Activate feature"
          />
          <FormControlLabel
            value="disable"
            control={<Radio />}
            label="Deactivate feature"
          />
        </RadioGroup>
        <Alert severity="warning" sx={{ mt: 2 }}>
          <Typography variant="body2">
            <strong>Important:</strong> Activate/Deactivate schedules only work when the feature's Master Enable switch is ON. 
            When Master Enable is OFF, the feature is completely disabled and schedules are ignored.
          </Typography>
        </Alert>
      </FormControl>
    </Box>
  );

  const renderPreviewStep = () => (
    <Box>
      <Typography variant="h6" gutterBottom>
        Schedule preview
      </Typography>
      
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

      {cronExpression && (
        <Paper sx={{ p: 2, mb: 2 }}>
          <Typography variant="subtitle1" gutterBottom>
            Schedule description:
          </Typography>
          <Typography variant="body1" color="text.secondary">
            {generateScheduleDescription(data)}
          </Typography>
          
          <Divider sx={{ my: 2 }} />
          
          <Typography variant="subtitle1" gutterBottom>
            Cron expression:
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <Typography 
              variant="body2" 
              fontFamily="monospace" 
              sx={{ 
                bgcolor: 'grey.800', 
                color: 'white',
                p: 1, 
                borderRadius: 1,
                flex: 1,
                maxWidth: '300px'
              }}
            >
              {cronExpression}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              (Read-only)
            </Typography>
          </Box>

          <Divider sx={{ my: 2 }} />
          
          <Typography variant="subtitle1" gutterBottom>
            Schedule dates:
          </Typography>
          <Typography variant="body2" color="text.secondary">
            <strong>Start:</strong> {data.startsAt ? new Date(data.startsAt).toLocaleDateString() : 'Not set'}
          </Typography>
          {data.endsAt && (
            <Typography variant="body2" color="text.secondary">
              <strong>End:</strong> {new Date(data.endsAt).toLocaleDateString()}
            </Typography>
          )}
          <Typography variant="body2" color="text.secondary">
            <strong>Timezone:</strong> {data.timezone}
          </Typography>
        </Paper>
      )}
    </Box>
  );

  if (!open) return null;

  return (
    <Box sx={{ width: '100%', maxWidth: 800, mx: 'auto', p: 3 }}>
      <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
        {steps.map((label) => (
          <Step key={label}>
            <StepLabel>{label}</StepLabel>
          </Step>
        ))}
      </Stepper>

      <Box sx={{ mb: 4 }}>
        {renderStepContent()}
      </Box>

      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
        <Button
          startIcon={<ArrowBackIcon />}
          onClick={handleBack}
          disabled={activeStep === 0}
        >
          Back
        </Button>

        {activeStep === steps.length - 1 ? (
          <Button
            variant="contained"
            startIcon={<CheckIcon />}
            onClick={handleSubmit}
            disabled={!canProceed()}
          >
            Save Changes
          </Button>
        ) : (
          <Button
            variant="contained"
            endIcon={<ArrowForwardIcon />}
            onClick={handleNext}
            disabled={!canProceed()}
          >
            Next
          </Button>
        )}
      </Box>
    </Box>
  );
};

export default EditRecurringScheduleBuilder;
