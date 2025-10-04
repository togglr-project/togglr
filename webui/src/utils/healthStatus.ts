/**
 * Utility functions for health status display
 */

export type HealthStatus = 'healthy' | 'failing' | 'degraded';

/**
 * Get color for health status chip
 */
export const getHealthStatusColor = (status: string): 'success' | 'error' | 'warning' | 'info' => {
  switch (status.toLowerCase()) {
    case 'healthy':
      return 'success';
    case 'failing':
      return 'error';
    case 'degraded':
      return 'warning';
    default:
      return 'info';
  }
};

/**
 * Get variant for health status chip
 */
export const getHealthStatusVariant = (status: string): 'filled' | 'outlined' => {
  switch (status.toLowerCase()) {
    case 'healthy':
      return 'filled';
    case 'failing':
      return 'filled';
    case 'degraded':
      return 'filled';
    default:
      return 'outlined';
  }
};
