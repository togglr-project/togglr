import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/apiClient';

// Hook to get all pending changes for a project
export const useProjectPendingChanges = (projectId?: string) => {
  return useQuery({
    queryKey: ['project-pending-changes', projectId],
    queryFn: async () => {
      if (!projectId) {
        console.log('[useProjectPendingChanges] No projectId, returning empty array');
        return [];
      }

      console.log(`[useProjectPendingChanges] Fetching pending changes for project ${projectId}`);
      
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
        
        const pendingChanges = response.data.data || [];
        console.log(`[useProjectPendingChanges] Found ${pendingChanges.length} pending changes:`, pendingChanges);
        
        return pendingChanges;
      } catch (error) {
        console.error('[useProjectPendingChanges] Error fetching pending changes:', error);
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
      (entity.entity === 'feature' || entity.entity === 'feature_params') && entity.entity_id === featureId
    )
  ) || false;
  
  console.log(`[useFeatureHasPendingChanges] Feature ${featureId} in project ${projectId}: hasPending=${hasPending}, pendingChanges count=${pendingChanges?.length || 0}`);
  
  // Debug: log the structure of pending changes
  if (pendingChanges && pendingChanges.length > 0) {
    console.log('[useFeatureHasPendingChanges] Pending changes structure:', pendingChanges);
    pendingChanges.forEach((change, index) => {
      console.log(`[useFeatureHasPendingChanges] Change ${index}:`, {
        id: change.id,
        entities: change.change.entities,
        status: change.status
      });
      
      // Debug each entity in the change
      if (change.change.entities) {
        change.change.entities.forEach((entity, entityIndex) => {
          console.log(`[useFeatureHasPendingChanges] Entity ${entityIndex}:`, {
            entity: entity.entity,
            entity_id: entity.entity_id,
            action: entity.action,
            changes: entity.changes
          });
        });
      }
    });
  }
  
  return hasPending;
};
