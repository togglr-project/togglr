-- Add cron_duration field to feature_schedules table
ALTER TABLE feature_schedules ADD COLUMN cron_duration VARCHAR(30);

-- Add comment to explain the field
COMMENT ON COLUMN feature_schedules.cron_duration IS 'Duration for cron-based schedules. When cron triggers, feature will be enabled/disabled for this duration.';
