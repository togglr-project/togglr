import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/apiClient';

// Hook to get feature names by IDs
export const useFeatureNames = (featureIds: string[], projectId?: string) => {
  return useQuery({
    queryKey: ['feature-names', featureIds, projectId],
    queryFn: async () => {
      if (!projectId || featureIds.length === 0) {
        return {};
      }

      // Get all features for the project
      const response = await apiClient.listProjectFeatures(projectId, undefined, undefined, undefined, undefined, undefined, undefined);
      
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
    staleTime: 5 * 60 * 1000, // 5 minutes
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
