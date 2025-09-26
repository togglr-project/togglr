import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/apiClient';

// Hook to get all pending changes for a project
export const useProjectPendingChanges = (projectId?: string) => {
  return useQuery({
    queryKey: ['project-pending-changes', projectId],
    queryFn: async () => {
      if (!projectId) {
        return [];
      }

      try {
        const response = await apiClient.listPendingChanges(
          projectId,
          'pending' as any,
          undefined,
          undefined,
          undefined,
          undefined,
          undefined
        );
        
        return response.data.data || [];
      } catch (error) {
        return [];
      }
    },
    enabled: !!projectId,
    staleTime: 30 * 1000, // 30 seconds
    refetchInterval: 30 * 1000, // Refetch every 30 seconds
  });
};

// Helper to check if a feature has pending changes
export const useFeatureHasPendingChanges = (featureId: string, projectId?: string) => {
  const { data: pendingChanges } = useProjectPendingChanges(projectId);
  
  return pendingChanges?.some(change => 
    change.change.entities?.some(entity => 
      entity.entity === 'feature' && entity.entity_id === featureId
    )
  ) || false;
};
