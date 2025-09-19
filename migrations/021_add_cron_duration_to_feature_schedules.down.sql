-- Remove cron_duration field from feature_schedules table
ALTER TABLE feature_schedules 
DROP COLUMN cron_duration;
