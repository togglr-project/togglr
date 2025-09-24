import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  Grid,
  Box,
  Button,
} from '@mui/material';
import { isValidCron } from 'cron-validator';
import cronstrue from 'cronstrue';
import { listTimeZones, findTimeZone, getZonedTime, getUTCOffset } from 'timezone-support';
import type { FeatureScheduleAction } from '../../generated/api/client';

interface ScheduleFormValues {
  action: FeatureScheduleAction;
  timezone: string;
  starts_at: string;
  ends_at: string;
  cron_expr: string;
  cron_duration: string;
}

interface ScheduleDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: ScheduleFormValues) => void;
  initial?: Partial<ScheduleFormValues>;
  title: string;
}

const allTimezones = listTimeZones();

const emptyForm = (): ScheduleFormValues => ({
  action: 'enable',
  timezone: Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC',
  starts_at: '',
  ends_at: '',
  cron_expr: '',
  cron_duration: '',
});

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

const fromDatetimeLocalInZoneToISO = (val?: string, tz?: string): string | undefined => {
  if (!val) return undefined;
  try {
    const [datePart, timePart] = val.split('T');
    if (!datePart || !timePart) return undefined;
    const [y, m, d] = datePart.split('-').map((s) => parseInt(s, 10));
    const [hh, mm] = timePart.split(':').map((s) => parseInt(s, 10));
    if (!y || !m || !d || isNaN(hh) || isNaN(mm)) return undefined;
    const tzObj = findTimeZone(tz || 'UTC');
    const wallAsUTC = new Date(Date.UTC(y, m - 1, d, hh, mm, 0, 0));
    const offsetMinutes = getUTCOffset(tzObj, wallAsUTC);
    const utcDate = new Date(wallAsUTC.getTime() - offsetMinutes * 60 * 1000);
    return utcDate.toISOString();
  } catch {
    return undefined;
  }
};

const ScheduleDialog: React.FC<ScheduleDialogProps> = ({ 
  open, 
  onClose, 
  onSubmit, 
  initial, 
  title 
}) => {
  const [values, setValues] = useState<ScheduleFormValues>(() => ({ ...emptyForm(), ...initial } as ScheduleFormValues));
  const [cronError, setCronError] = useState<string>('');
  const [cronDesc, setCronDesc] = useState<string>('');
  const [tzError, setTzError] = useState<string>('');

  useEffect(() => {
    // Initialize form values and re-validate cron; convert ISO timestamps to datetime-local per timezone
    const base = { ...emptyForm(), ...initial } as ScheduleFormValues;
    const tz = base.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
    const startsLocal = initial?.starts_at ? toDatetimeLocalInZone(initial.starts_at, tz) : '';
    const endsLocal = initial?.ends_at ? toDatetimeLocalInZone(initial.ends_at, tz) : '';
    setValues({ ...base, starts_at: startsLocal, ends_at: endsLocal, cron_duration: initial?.cron_duration || '' });

    const expr = (initial?.cron_expr || '').trim();
    if (expr) {
      const ok = isValidCron(expr, { seconds: false, allowBlankDay: true, alias: true });
      setCronError(ok ? '' : 'Invalid cron expression');
      setCronDesc(ok ? cronstrue.toString(expr) : '');
    } else {
      setCronError('');
      setCronDesc('');
    }
  }, [initial, open]);

  const handleSubmit = () => {
    const tz = values.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
    const startsISO = fromDatetimeLocalInZoneToISO(values.starts_at, tz);
    const endsISO = fromDatetimeLocalInZoneToISO(values.ends_at, tz);
    
    const submitValues: ScheduleFormValues = {
      ...values,
      starts_at: startsISO || '',
      ends_at: endsISO || '',
    };
    
    onSubmit(submitValues);
    onClose();
  };

  const handleCronChange = (expr: string) => {
    setValues(v => ({ ...v, cron_expr: expr }));
    const trimmed = expr.trim();
    if (trimmed) {
      const ok = isValidCron(trimmed, { seconds: false, allowBlankDay: true, alias: true });
      setCronError(ok ? '' : 'Invalid cron expression');
      setCronDesc(ok ? cronstrue.toString(trimmed) : '');
    } else {
      setCronError('');
      setCronDesc('');
    }
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle>{title}</DialogTitle>
      <DialogContent dividers>
        <Box sx={{ mt: 1 }}>
          <Grid container spacing={2}>
            <Grid item xs={12} md={6}>
              <TextField
                select
                fullWidth
                size="small"
                label="Action"
                value={values.action}
                onChange={(e) => setValues(v => ({ ...v, action: e.target.value as FeatureScheduleAction }))}
              >
                <MenuItem value="enable">Activate</MenuItem>
                <MenuItem value="disable">Deactivate</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                select
                fullWidth
                size="small"
                label="Timezone"
                value={values.timezone}
                onChange={(e) => {
                  const val = e.target.value;
                  setValues(v => ({
                    ...v,
                    timezone: val,
                    starts_at: initial?.starts_at ? toDatetimeLocalInZone(initial.starts_at, val) : v.starts_at,
                    ends_at: initial?.ends_at ? toDatetimeLocalInZone(initial.ends_at, val) : v.ends_at,
                  }));
                  setTzError(allTimezones.includes(val) ? '' : 'Invalid timezone');
                }}
                error={Boolean(tzError)}
                helperText={tzError || 'Choose IANA timezone'}
              >
                {allTimezones.map((tz: string) => (
                  <MenuItem key={tz} value={tz}>{tz}</MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                size="small"
                label="Start Date & Time"
                type="datetime-local"
                value={values.starts_at}
                onChange={(e) => setValues(v => ({ ...v, starts_at: e.target.value }))}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                size="small"
                label="End Date & Time"
                type="datetime-local"
                value={values.ends_at}
                onChange={(e) => setValues(v => ({ ...v, ends_at: e.target.value }))}
                InputLabelProps={{ shrink: true }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                size="small"
                label="Cron Expression"
                value={values.cron_expr}
                onChange={(e) => handleCronChange(e.target.value)}
                error={Boolean(cronError)}
                helperText={cronError || cronDesc || 'Enter cron expression (e.g., "0 9 * * 1-5" for weekdays at 9 AM)'}
                placeholder="0 9 * * 1-5"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                size="small"
                label="Duration (optional)"
                value={values.cron_duration}
                onChange={(e) => setValues(v => ({ ...v, cron_duration: e.target.value }))}
                helperText="Duration for cron-based schedules (e.g., '2h', '30m', '1d')"
                placeholder="2h"
              />
            </Grid>
          </Grid>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button 
          onClick={handleSubmit} 
          variant="contained"
          disabled={!values.starts_at || !values.ends_at || Boolean(cronError) || Boolean(tzError)}
        >
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ScheduleDialog;
