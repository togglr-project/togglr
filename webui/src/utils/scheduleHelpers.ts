import type { FeatureSchedule } from '../generated/api/client';

/**
 * Checks if feature has cron schedule
 */
export function hasCronSchedule(schedules: FeatureSchedule[]): boolean {
  return schedules.some(schedule => !!schedule.cron_expr);
}

/**
 * Checks if feature has one-shot schedule
 */
export function hasOneShotSchedule(schedules: FeatureSchedule[]): boolean {
  return schedules.some(schedule => !schedule.cron_expr && (schedule.starts_at || schedule.ends_at));
}

/**
 * Checks if new schedule can be added for feature
 * "Add schedule" button is shown only if:
 * - No schedules exist, OR
 * - Only one-shot schedules exist (no cron)
 * 
 * Logic: EITHER one recurring (cron) schedule OR one or more one-shot schedules
 */
export function canAddSchedule(schedules: FeatureSchedule[]): boolean {
  if (schedules.length === 0) return true;
  return !hasCronSchedule(schedules);
}

/**
 * Checks if recurring schedule can be added
 * Only allowed if no schedules exist at all
 */
export function canAddRecurringSchedule(schedules: FeatureSchedule[]): boolean {
  return schedules.length === 0;
}

/**
 * Checks if one-shot schedule can be added
 * Allowed if no schedules exist OR only one-shot schedules exist
 */
export function canAddOneShotSchedule(schedules: FeatureSchedule[]): boolean {
  if (schedules.length === 0) return true;
  return !hasCronSchedule(schedules);
}

/**
 * Gets schedule type for display
 */
export function getScheduleType(schedule: FeatureSchedule): 'cron' | 'one-shot' {
  return schedule.cron_expr ? 'cron' : 'one-shot';
}
