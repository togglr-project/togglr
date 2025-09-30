import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/apiClient';

// Hook to get feature names by IDs
export const useFeatureNames = (featureIds: string[], projectId?: string, environmentKey: string = 'prod') => {
  return useQuery({
    queryKey: ['feature-names', featureIds, projectId, environmentKey],
    queryFn: async () => {
      if (!projectId || featureIds.length === 0) {
        return {};
      }

      // Get all features for the project
      const response = await apiClient.listProjectFeatures(projectId, environmentKey, undefined, undefined, undefined, undefined, undefined);
      
      // Create a map of feature ID to name
      const featureNames: Record<string, string> = {};
      response.data.items?.forEach(feature => {
        if (feature.id && feature.name) {
          featureNames[feature.id] = feature.name;
        }
      });

      return featureNames;
    },
    enabled: !!projectId && featureIds.length > 0,
    staleTime: 0, // No caching - always fetch fresh data
    refetchOnWindowFocus: true, // Refetch when window gains focus
  });
};

// Helper function to get entity display name
export const getEntityDisplayName = (
  entity: { entity: string; entity_id: string },
  featureNames?: Record<string, string>
): string => {
  if (entity.entity === 'feature' && featureNames?.[entity.entity_id]) {
    return featureNames[entity.entity_id];
  }
  
  // Fallback to entity type and truncated ID
  return `${entity.entity} (${entity.entity_id.slice(0, 8)}...)`;
};
