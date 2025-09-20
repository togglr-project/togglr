# Understanding Feature Enablement and Schedules

Master Enable is the global on/off switch: when **OFF** the feature is always disabled and schedules are ignored.
When **ON**, the feature is either controlled manually (if no schedules are defined) or automatically by schedules.
You can create either one repeating schedule (using the friendly “repeat” builder) **or** one-or-more non-overlapping one-shot intervals.
Repeating schedules define periodic active windows (duration is required); one-shot intervals define exact start/end periods.
Baseline (state outside scheduled windows) depends on schedule types: repeating-`Activate` → baseline OFF (activate only during windows); repeating-`Deactivate` → baseline ON (deactivate only during windows); for one-shots baseline is ON if any Deactivate interval exists, otherwise OFF.
Newer schedules override older ones; if two are created at the same instant, Deactivate wins.
