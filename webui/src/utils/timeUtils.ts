/**
 * Converts seconds to a human-readable time format
 * @param seconds - number of seconds (can be undefined)
 * @returns formatted time string or 'N/A' if value is undefined
 */
export const formatDuration = (seconds: number | undefined): string => {
  if (seconds === undefined || seconds === null || seconds <= 0) {
    return 'N/A';
  }
  
  if (seconds < 60) {
    return `${Math.round(seconds)}s`;
  }
  
  if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = Math.round(seconds % 60);
    if (remainingSeconds === 0) {
      return `${minutes}m`;
    }
    return `${minutes}m ${remainingSeconds}s`;
  }
  
  if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600);
    const remainingMinutes = Math.floor((seconds % 3600) / 60);
    if (remainingMinutes === 0) {
      return `${hours}h`;
    }
    return `${hours}h ${remainingMinutes}m`;
  }
  
  const days = Math.floor(seconds / 86400);
  const remainingHours = Math.floor((seconds % 86400) / 3600);
  if (remainingHours === 0) {
    return `${days}d`;
  }
  return `${days}d ${remainingHours}h`;
};

/**
 * Converts seconds to a short format for UI display
 * @param seconds - number of seconds (can be undefined)
 * @returns short formatted time string or 'N/A' if value is undefined
 */
export const formatDurationShort = (seconds: number | undefined): string => {
  if (seconds === undefined || seconds === null || seconds <= 0) {
    return 'N/A';
  }
  
  if (seconds < 60) {
    return `${Math.round(seconds)}s`;
  }
  
  if (seconds < 3600) {
    const minutes = Math.floor(seconds / 60);
    return `${minutes}m`;
  }
  
  if (seconds < 86400) {
    const hours = Math.floor(seconds / 3600);
    return `${hours}h`;
  }
  
  const days = Math.floor(seconds / 86400);
  return `${days}d`;
};

/**
 * Formats next state time for display
 * @param nextStateTime - ISO string timestamp or undefined
 * @returns formatted time string or null if no time provided
 */
export const formatNextStateTime = (nextStateTime: string | undefined): string | null => {
  if (!nextStateTime) {
    return null;
  }
  
  try {
    const date = new Date(nextStateTime);
    const now = new Date();
    const diffMs = date.getTime() - now.getTime();
    
    // If the time is in the past, show "Overdue"
    if (diffMs < 0) {
      return 'Overdue';
    }
    
    // If it's within the next hour, show relative time
    if (diffMs < 3600000) { // 1 hour
      const minutes = Math.floor(diffMs / 60000);
      if (minutes < 1) {
        return 'Now';
      }
      return `in ${minutes}m`;
    }
    
    // If it's within the next day, show hours
    if (diffMs < 86400000) { // 24 hours
      const hours = Math.floor(diffMs / 3600000);
      const minutes = Math.floor((diffMs % 3600000) / 60000);
      if (minutes === 0) {
        return `in ${hours}h`;
      }
      return `in ${hours}h ${minutes}m`;
    }
    
    // For longer periods, show date and time
    return date.toLocaleString();
  } catch (error) {
    console.error('Error formatting next state time:', error);
    return 'Invalid date';
  }
};

/**
 * Gets a human-readable description of the next state
 * @param nextState - boolean indicating if next state is enabled, or undefined
 * @param nextStateTime - ISO string timestamp or undefined
 * @returns description string or null if no next state
 */
export const getNextStateDescription = (nextState: boolean | undefined, nextStateTime: string | undefined): string | null => {
  if (nextState === undefined || !nextStateTime) {
    return null;
  }
  
  const timeStr = formatNextStateTime(nextStateTime);
  if (!timeStr) {
    return null;
  }
  
  const action = nextState ? 'enable' : 'disable';
  return `Will ${action} ${timeStr}`;
}; 