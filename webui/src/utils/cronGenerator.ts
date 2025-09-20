/**
 * Utilities for generating cron expressions from user interface
 */

export type ScheduleType = 
  | 'repeat_every' 
  | 'daily' 
  | 'monthly' 
  | 'yearly';

export interface RepeatEveryParams {
  interval: number;
  unit: 'minutes' | 'hours';
}

export interface DailyParams {
  time: string; // HH:MM format
}

export interface MonthlyParams {
  dayOfMonth: number; // 1-31
  time: string; // HH:MM format
}

export interface YearlyParams {
  month: number; // 1-12
  day: number; // 1-31
  time: string; // HH:MM format
}

export interface ScheduleBuilderData {
  timezone: string;
  startsAt?: string; // ISO string for start date
  endsAt?: string;   // ISO string for end date (optional)
  scheduleType: ScheduleType;
  repeatEvery?: RepeatEveryParams;
  daily?: DailyParams;
  monthly?: MonthlyParams;
  yearly?: YearlyParams;
  duration: {
    value: number;
    unit: 'minutes' | 'hours' | 'days';
  };
  action: 'enable' | 'disable';
}

/**
 * Generates cron expression based on schedule parameters
 */
export function generateCronExpression(data: ScheduleBuilderData): string {
  const { scheduleType, repeatEvery, daily, monthly, yearly } = data;

  switch (scheduleType) {
    case 'repeat_every':
      if (!repeatEvery) throw new Error('Repeat every parameters required');
      return generateRepeatEveryCron(repeatEvery);

    case 'daily':
      if (!daily) throw new Error('Daily parameters required');
      return generateDailyCron(daily);

    case 'monthly':
      if (!monthly) throw new Error('Monthly parameters required');
      return generateMonthlyCron(monthly);

    case 'yearly':
      if (!yearly) throw new Error('Yearly parameters required');
      return generateYearlyCron(yearly);

    default:
      throw new Error('Unknown schedule type');
  }
}

function generateRepeatEveryCron(params: RepeatEveryParams): string {
  const { interval, unit } = params;
  
  if (unit === 'minutes') {
    // Every N minutes: */N * * * *
    return `*/${interval} * * * *`;
  } else {
    // Every N hours: 0 */N * * *
    return `0 */${interval} * * *`;
  }
}

function generateDailyCron(params: DailyParams): string {
  const { time } = params;
  const [hours, minutes] = time.split(':').map(Number);
  
  // Daily at specific time: MM HH * * *
  return `${minutes} ${hours} * * *`;
}

function generateMonthlyCron(params: MonthlyParams): string {
  const { dayOfMonth, time } = params;
  const [hours, minutes] = time.split(':').map(Number);
  
  // Monthly on specific day: MM HH D * *
  return `${minutes} ${hours} ${dayOfMonth} * *`;
}

function generateYearlyCron(params: YearlyParams): string {
  const { month, day, time } = params;
  const [hours, minutes] = time.split(':').map(Number);
  
  // Yearly on specific date: MM HH D M *
  return `${minutes} ${hours} ${day} ${month} *`;
}

/**
 * Generates human-readable schedule description
 */
export function generateScheduleDescription(data: ScheduleBuilderData): string {
  const { scheduleType, repeatEvery, daily, monthly, yearly, duration, action } = data;
  
  const actionText = action === 'enable' ? 'activate' : 'deactivate';
  const durationText = formatDuration(duration);

  switch (scheduleType) {
    case 'repeat_every':
      if (!repeatEvery) return '';
      const { interval, unit } = repeatEvery;
      const unitText = unit === 'minutes' ? 'minutes' : 'hours';
      return `Repeat every ${interval} ${unitText}, ${actionText} for ${durationText}`;

    case 'daily':
      if (!daily) return '';
      return `Daily at ${daily.time}, ${actionText} for ${durationText}`;

    case 'monthly':
      if (!monthly) return '';
      return `Monthly on ${monthly.dayOfMonth} at ${monthly.time}, ${actionText} for ${durationText}`;

    case 'yearly':
      if (!yearly) return '';
      const monthNames = [
        'January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'
      ];
      return `Yearly on ${monthNames[yearly.month - 1]} ${yearly.day} at ${yearly.time}, ${actionText} for ${durationText}`;

    default:
      return '';
  }
}

function formatDuration(duration: { value: number; unit: string }): string {
  const { value, unit } = duration;
  
  switch (unit) {
    case 'minutes':
      return value === 1 ? '1 minute' : `${value} minutes`;
    case 'hours':
      return value === 1 ? '1 hour' : `${value} hours`;
    case 'days':
      return value === 1 ? '1 day' : `${value} days`;
    default:
      return `${value} ${unit}`;
  }
}

/**
 * Validates schedule parameters
 */
export function validateScheduleData(data: ScheduleBuilderData): string[] {
  const errors: string[] = [];

  if (!data.timezone) {
    errors.push('Timezone is required');
  }

  if (!data.scheduleType) {
    errors.push('Schedule type is required');
  }

  if (data.duration.value <= 0) {
    errors.push('Duration must be greater than 0');
  }

  switch (data.scheduleType) {
    case 'repeat_every':
      if (!data.repeatEvery) {
        errors.push('Repeat parameters are required');
      } else if (data.repeatEvery.interval <= 0) {
        errors.push('Interval must be greater than 0');
      } else {
        // Duration must be less than repeat interval and unit must be <= interval unit
        const toMinutes = (value: number, unit: 'minutes' | 'hours' | 'days') =>
          unit === 'minutes' ? value : unit === 'hours' ? value * 60 : value * 24 * 60;

        const repeatMinutes = toMinutes(data.repeatEvery.interval, data.repeatEvery.unit);
        const durationMinutes = toMinutes(data.duration.value, data.duration.unit);

        if (data.repeatEvery.unit === 'minutes' && data.duration.unit !== 'minutes') {
          errors.push('For minute-based intervals, duration unit must be minutes');
        }
        if (data.repeatEvery.unit === 'hours' && data.duration.unit === 'days') {
          errors.push('For hour-based intervals, duration unit must be hours or minutes');
        }
        if (durationMinutes >= repeatMinutes) {
          errors.push('Duration must be less than the repeat interval');
        }
      }
      break;

    case 'daily':
      if (!data.daily?.time) {
        errors.push('Time for daily schedule is required');
      }
      break;

    case 'monthly':
      if (!data.monthly) {
        errors.push('Monthly schedule parameters are required');
      } else {
        if (data.monthly.dayOfMonth < 1 || data.monthly.dayOfMonth > 31) {
          errors.push('Day of month must be between 1 and 31');
        }
        if (!data.monthly.time) {
          errors.push('Time for monthly schedule is required');
        }
      }
      break;

    case 'yearly':
      if (!data.yearly) {
        errors.push('Yearly schedule parameters are required');
      } else {
        if (data.yearly.month < 1 || data.yearly.month > 12) {
          errors.push('Month must be between 1 and 12');
        }
        if (data.yearly.day < 1 || data.yearly.day > 31) {
          errors.push('Day must be between 1 and 31');
        }
        if (!data.yearly.time) {
          errors.push('Time for yearly schedule is required');
        }
      }
      break;
  }

  return errors;
}
