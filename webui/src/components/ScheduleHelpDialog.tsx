import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Divider
} from '@mui/material';
import { Close as CloseIcon } from '@mui/icons-material';

interface ScheduleHelpDialogProps {
  open: boolean;
  onClose: () => void;
}

const ScheduleHelpDialog: React.FC<ScheduleHelpDialogProps> = ({ open, onClose }) => {
  return (
    <Dialog 
      open={open} 
      onClose={onClose}
      maxWidth="md"
      fullWidth
    >
      <DialogTitle>
        Understanding Feature Enablement and Schedules
        <Button
          onClick={onClose}
          sx={{ position: 'absolute', right: 8, top: 8 }}
          startIcon={<CloseIcon />}
        >
          Close
        </Button>
      </DialogTitle>
      <DialogContent>
        <Box sx={{ '& > * + *': { mt: 2 } }}>
          <Typography variant="h6" gutterBottom>
            How the master switch and schedules control a feature
          </Typography>
          
          <Typography variant="body1">
            Each feature has a <strong>Master Enable</strong> switch and (optionally) one or more <strong>schedules</strong> that automatically change the feature state.
          </Typography>

          <Box>
            <Typography variant="h6" gutterBottom>
              1. Master Enable (global switch)
            </Typography>
            <Box component="ul" sx={{ pl: 2, m: 0 }}>
              <li><strong>If Master Enable = OFF</strong> &rarr; the feature is completely turned off. Schedules and rules are ignored.</li>
              <li><strong>If Master Enable = ON</strong> &rarr; the feature is "live" and its state is determined either manually (if no schedules exist) or by schedules (if schedules are configured).</li>
            </Box>
          </Box>

          <Box>
            <Typography variant="h6" gutterBottom>
              2. Two schedule types
            </Typography>
            <Box component="ul" sx={{ pl: 2, m: 0 }}>
              <li><strong>Repeating schedule</strong> (the "repeat" builder you see in the UI). Internally it is stored as a cron-like rule and a required duration. Exactly <strong>one repeating schedule</strong> is allowed per feature.</li>
              <li><strong>One-shot schedules</strong> (fixed start/end intervals). You may create <strong>one or several non-overlapping</strong> one-shot schedules for a feature. The UI prevents overlaps.</li>
            </Box>
          </Box>

          <Box>
            <Typography variant="h6" gutterBottom>
              3. How the current state is determined (algorithm)
            </Typography>
            <Typography variant="body2" sx={{ mb: 1 }}>
              When Master Enable = ON:
            </Typography>
            <Box component="ul" sx={{ pl: 2, m: 0 }}>
              <li>If <strong>no schedules</strong> are configured &rarr; the feature remains in its manual state (what you set in the UI). Schedules do not change anything because none exist.</li>
              <li>If <strong>a repeating schedule</strong> exists (the cron-like mode) &rarr; the feature's baseline and active windows depend on the repeating schedule action:
                <Box component="ul" sx={{ pl: 2, mt: 1 }}>
                  <li>If the schedule action is <strong>Activate</strong> &rarr; the baseline (outside scheduled windows) is <strong>OFF</strong>; the feature becomes <strong>ON</strong> only during scheduled windows (each window lasts the configured duration).</li>
                  <li>If the schedule action is <strong>Deactivate</strong> &rarr; the baseline (outside scheduled windows) is <strong>ON</strong>; the feature becomes <strong>OFF</strong> only during scheduled windows.</li>
                </Box>
              </li>
              <li>If <strong>one-shot schedules</strong> exist (one or more intervals) &rarr; the baseline is derived from the collection of intervals:
                <Box component="ul" sx={{ pl: 2, mt: 1 }}>
                  <li>If <strong>all</strong> one-shot intervals are <code>Activate</code> &rarr; baseline = <strong>OFF</strong> (feature is OFF except during the activate intervals).</li>
                  <li>If <strong>any</strong> one-shot interval is <code>Deactivate</code> (or there is a mix of Activate/Deactivate) &rarr; baseline = <strong>ON</strong> (feature is ON except during the deactivate intervals).</li>
                  <li>At any moment if a one-shot interval is active, that interval's action (Activate &rarr; ON, Deactivate &rarr; OFF) determines visibility for that moment.</li>
                </Box>
              </li>
            </Box>
          </Box>

          <Box>
            <Typography variant="h6" gutterBottom>
              4. Conflicts and precedence
            </Typography>
            <Box component="ul" sx={{ pl: 2, m: 0 }}>
              <li>You cannot mix repeating (cron) and one-shot schedules for the same feature (the UI and DB prevent that).</li>
              <li>If multiple schedules somehow apply at the same time, the system picks the schedule with the <strong>latest creation time</strong> (<code>created_at</code>). If two schedules have the same creation time, <strong>Deactivate</strong> wins over <strong>Activate</strong>.</li>
            </Box>
          </Box>

          <Box>
            <Typography variant="h6" gutterBottom>
              5. Required fields and notes
            </Typography>
            <Box component="ul" sx={{ pl: 2, m: 0 }}>
              <li>Every repeating schedule must include a <strong>duration</strong> (how long the action lasts after each trigger).</li>
              <li>Timezone is required and determines the local time at which scheduled actions run.</li>
              <li>The UI shows a friendly preview (human description + timeline) so you can confirm what you created before saving.</li>
            </Box>
          </Box>

          <Divider sx={{ my: 2 }} />

          <Typography variant="h6" gutterBottom>
            Examples
          </Typography>
          <Box component="ul" sx={{ pl: 2, m: 0 }}>
            <li>Feature <code>X</code> with Master Enable = OFF &rarr; <code>X</code> is always off, regardless of schedules.</li>
            <li>Feature <code>Y</code> with Master Enable = ON and <strong>no schedules</strong> &rarr; <code>Y</code> stays in whatever manual state you set in the UI.</li>
            <li>Feature <code>Z</code> with Master Enable = ON and <strong>one repeating schedule</strong> action = Activate, duration = 30m (daily at 09:30) &rarr; <code>Z</code> is OFF by default and turns ON for 30 minutes starting at 09:30 local time.</li>
            <li>Feature <code>A</code> with Master Enable = ON and two one-shot intervals: (<code>Deactivate</code> on Sep 20 18:00–18:30) and (<code>Activate</code> on Sep 25 09:00–10:00) &rarr; baseline = ON (because a Deactivate interval exists). At Sep 20 18:05 &rarr; <code>A</code> is OFF. At Sep 25 09:15 &rarr; <code>A</code> is ON. At other times &rarr; <code>A</code> is ON.</li>
          </Box>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Got it
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ScheduleHelpDialog;
