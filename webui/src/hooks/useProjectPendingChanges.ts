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
          undefined, // environmentId: all environments
          undefined, // environmentKey
          projectId, // projectId
          'pending' as any, // status
          undefined, // userId
          undefined, // page
          undefined, // perPage
          undefined, // sortBy
          undefined  // sortDesc
        );
        
        return response.data.data || [];
      } catch (error) {
        return [];
      }
    },
    enabled: !!projectId,
    staleTime: 0, // No caching - always fetch fresh data
    refetchOnWindowFocus: true, // Refetch when window gains focus
    // Removed refetchInterval - updates now come via WebSocket
  });
};

// Helper to check if a feature has pending changes
export const useFeatureHasPendingChanges = (featureId: string, projectId?: string) => {
  const { data: pendingChanges } = useProjectPendingChanges(projectId);
  
  const hasPending = pendingChanges?.some(change => 
    change.change.entities?.some(entity => 
      entity.entity === 'feature' && entity.entity_id === featureId
    )
  ) || false;
  
  return hasPending;
};
