import { QueryClient } from '@tanstack/react-query';

export type RealtimeEvent = {
  source: 'pending' | 'audit' | string;
  type: string; // e.g., feature_updated, feature_pending
  timestamp: string;
  project_id: string;
  environment: string; // environment key
  entity: string; // e.g., 'feature', 'pending_change'
  entity_id: string;
  action: string; // e.g., 'updated', 'pending', 'deleted'
};

const invalidate = (qc: QueryClient, keys: any[]) => {
  try { qc.invalidateQueries({ queryKey: keys, exact: false }); } catch {}
};

const refetchQueries = (qc: QueryClient, keys: any[]) => {
  try { 
    qc.refetchQueries({ queryKey: keys, exact: false }); 
    console.log('[Realtime] Refetching queries:', keys);
  } catch (e) {
    console.error('[Realtime] Error refetching queries:', e);
  }
};

const updateFeatureInCache = (qc: QueryClient, projectId: string, envKey: string, featureId: string) => {
  // Since we have no caching, just invalidate all related queries
  invalidate(qc, ['project-features', projectId, envKey]);
  invalidate(qc, ['feature-details', featureId, envKey]);
  invalidate(qc, ['feature-names', projectId, envKey]);
  
  console.log('[Realtime] Invalidated feature queries:', featureId);
};

const removeFeatureFromCache = (qc: QueryClient, projectId: string, envKey: string, featureId: string) => {
  // Since we have no caching, just invalidate all related queries
  invalidate(qc, ['project-features', projectId, envKey]);
  invalidate(qc, ['feature-details', featureId, envKey]);
  invalidate(qc, ['feature-names', projectId, envKey]);
  
  console.log('[Realtime] Invalidated feature queries for deletion:', featureId);
};

const markFeatureAsPending = (qc: QueryClient, projectId: string, envKey: string, featureId: string) => {
  // Update feature - pending changes will be invalidated by the main handler
  updateFeatureInCache(qc, projectId, envKey, featureId);
  
  console.log('[Realtime] Marked feature as pending:', featureId);
};

export function handleEvent(qc: QueryClient | undefined, evt: RealtimeEvent) {
  if (!qc) return;

  const envKey = evt.environment;
  const projectId = evt.project_id;
  const entityId = evt.entity_id;

  console.log('[Realtime] Handling event:', evt);

  // Check if this is a feature-related event by type
  const isFeatureEvent = evt.type.startsWith('feature_');
  
  if (isFeatureEvent) {
    // For feature events, we need to determine if we can update specific feature
    // or if we need to invalidate all features
    
    if (evt.entity === 'feature') {
      // Direct feature events - can update specific feature
      const featureId = entityId;
      
      switch (evt.type) {
        case 'feature_update':
        case 'feature_updated': {
          // Update specific feature
          updateFeatureInCache(qc, projectId, envKey, featureId);
          // Force refetch specific feature queries
          refetchQueries(qc, ['project-features', projectId, envKey]);
          refetchQueries(qc, ['feature-details', featureId, envKey]);
          break;
        }
        case 'feature_pending': {
          // Show pending status
          markFeatureAsPending(qc, projectId, envKey, featureId);
          break;
        }
        case 'feature_deleted': {
          // Remove feature from list
          removeFeatureFromCache(qc, projectId, envKey, featureId);
          break;
        }
        default: {
          // For unknown types, update specific feature
          updateFeatureInCache(qc, projectId, envKey, featureId);
        }
      }
    } else {
      // Feature-related events on child entities - invalidate all features
      console.log('[Realtime] Feature-related child entity changed, invalidating all features:', evt.entity, evt.type);
      invalidate(qc, ['project-features', projectId, envKey]);
      invalidate(qc, ['feature-details']);
      invalidate(qc, ['feature-names', projectId, envKey]);
      
      // Force refetch all feature-related queries
      console.log('[Realtime] Force refetching all feature queries');
      refetchQueries(qc, ['project-features', projectId]);
      refetchQueries(qc, ['feature-details']);
      refetchQueries(qc, ['feature-names', projectId, envKey]);
      
      // Also try to refetch with broader patterns
      refetchQueries(qc, ['project-features']);
      refetchQueries(qc, ['feature-details']);
      
      // As a last resort, try to force update by setting stale data
      console.log('[Realtime] Attempting to force update by setting stale data');
      qc.setQueriesData(
        { queryKey: ['project-features', projectId], exact: false },
        (oldData: any) => {
          if (oldData) {
            console.log('[Realtime] Marking project-features as stale');
            return { ...oldData, _isStale: true };
          }
          return oldData;
        }
      );
    }
    
    // Always invalidate related queries for feature events
    invalidate(qc, ['feature-timelines', projectId, envKey]);
    invalidate(qc, ['dashboard', projectId]);
    invalidate(qc, ['pending-changes', projectId, envKey]);
    invalidate(qc, ['project-pending-changes', projectId]);
    
    // Force refetch all feature-related queries
    refetchQueries(qc, ['feature-timelines', projectId, envKey]);
    refetchQueries(qc, ['dashboard', projectId]);
    refetchQueries(qc, ['pending-changes', projectId, envKey]);
    refetchQueries(qc, ['project-pending-changes', projectId]);
    
    // Also invalidate and refetch the pending changes count for menu badge
    invalidate(qc, ['pending-changes', projectId, 'pending', undefined, 'count']);
    refetchQueries(qc, ['pending-changes', projectId, 'pending', undefined, 'count']);
  } else {
    // Non-feature events
    switch (evt.entity) {
    case 'pending_change': {
      // Invalidate pending changes
      invalidate(qc, ['pending-changes', projectId, envKey]);
      invalidate(qc, ['project-pending-changes', projectId]);
      // Also invalidate the count for menu badge
      invalidate(qc, ['pending-changes', projectId, 'pending', undefined, 'count']);
      refetchQueries(qc, ['pending-changes', projectId, 'pending', undefined, 'count']);
      break;
    }
    default: {
      // Fallback: invalidate project-level queries
      console.log('[Realtime] Unknown entity, invalidating project queries:', evt.entity);
      invalidate(qc, [projectId]);
      invalidate(qc, ['project-features', projectId, envKey]);
      invalidate(qc, ['pending-changes', projectId, envKey]);
      invalidate(qc, ['project-pending-changes', projectId]);
      // Also invalidate the count for menu badge
      invalidate(qc, ['pending-changes', projectId, 'pending', undefined, 'count']);
      refetchQueries(qc, ['pending-changes', projectId, 'pending', undefined, 'count']);
      break;
    }
    }
  }
}