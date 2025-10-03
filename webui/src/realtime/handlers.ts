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

// Helper function to invalidate queries by pattern
const invalidateByPattern = (qc: QueryClient, pattern: (query: any) => boolean, description: string) => {
  try {
    qc.invalidateQueries({ predicate: pattern });
  } catch (e) {
    console.error(`[Realtime] Error invalidating ${description}:`, e);
  }
};

// Helper function to refetch queries by pattern
const refetchByPattern = (qc: QueryClient, pattern: (query: any) => boolean, description: string) => {
  try {
    qc.refetchQueries({ predicate: pattern });
  } catch (e) {
    console.error(`[Realtime] Error refetching ${description}:`, e);
  }
};

// Helper function to invalidate queries by key pattern
const invalidateByKey = (qc: QueryClient, keyPattern: any[], description: string) => {
  try {
    qc.invalidateQueries({ queryKey: keyPattern, exact: false });
  } catch (e) {
    console.error(`[Realtime] Error invalidating ${description}:`, e);
  }
};

// Helper function to refetch queries by key pattern
const refetchByKey = (qc: QueryClient, keyPattern: any[], description: string) => {
  try {
    qc.refetchQueries({ queryKey: keyPattern, exact: false });
  } catch (e) {
    console.error(`[Realtime] Error refetching ${description}:`, e);
  }
};

// Helper function to invalidate pending changes count for menu badge
const invalidatePendingChangesCount = (qc: QueryClient, projectId: string) => {
  try {
    // Invalidate all pending changes count queries for this project
    qc.invalidateQueries({
      predicate: (query) => {
        const key = query.queryKey;
        return Boolean(key[0] === 'pendingChanges' && 
               key[1] === 'list' && 
               key[2] && 
               typeof key[2] === 'object' && 
               (key[2] as any).projectId === projectId && 
               key[3] === 'count');
      }
    });
    console.log(`[Realtime] Invalidated pending changes count for project ${projectId}`);
  } catch (e) {
    console.error(`[Realtime] Error invalidating pending changes count:`, e);
  }
};

// Helper function to refetch pending changes count for menu badge
const refetchPendingChangesCount = (qc: QueryClient, projectId: string) => {
  try {
    // Refetch all pending changes count queries for this project
    qc.refetchQueries({
      predicate: (query) => {
        const key = query.queryKey;
        return Boolean(key[0] === 'pendingChanges' && 
               key[1] === 'list' && 
               key[2] && 
               typeof key[2] === 'object' && 
               (key[2] as any).projectId === projectId && 
               key[3] === 'count');
      }
    });
    console.log(`[Realtime] Refetched pending changes count for project ${projectId}`);
  } catch (e) {
    console.error(`[Realtime] Error refetching pending changes count:`, e);
  }
};

// Handle feature update events
const handleFeatureUpdate = (qc: QueryClient, projectId: string, envKey: string, featureId: string) => {
  console.log(`[Realtime] Handling feature update: ${featureId} in project ${projectId}`);
  
  // 1. Invalidate all project-features queries (with any parameters)
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-features' && query.queryKey[1] === projectId,
    'project-features queries'
  );
  
  // 2. Invalidate feature-details queries (for preview panel)
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-details' && query.queryKey[1] === featureId,
    'feature-details queries (preview panel)'
  );
  
  // 3. Invalidate feature-changes queries (for preview panel)
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-changes' && query.queryKey[1] === featureId && query.queryKey[2] === projectId,
    'feature-changes queries (preview panel)'
  );
  
  // 4. Invalidate feature-names queries
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-names' && query.queryKey[1] === projectId,
    'feature-names queries'
  );
  
  // 4. Invalidate all queries that contain projectId (broader pattern)
  invalidateByPattern(
    qc,
    (query) => query.queryKey.includes(projectId),
    'all project queries'
  );
  
  // 5. Refetch all project-features queries
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-features' && query.queryKey[1] === projectId,
    'project-features queries'
  );
  
  // 6. Refetch feature-details queries
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-details' && query.queryKey[1] === featureId,
    'feature-details queries'
  );
  
  // 7. Refetch feature-changes queries (for preview panel)
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-changes' && query.queryKey[1] === featureId && query.queryKey[2] === projectId,
    'feature-changes queries (preview panel)'
  );
  
  // 8. Refetch all queries that contain projectId
  refetchByPattern(
    qc,
    (query) => query.queryKey.includes(projectId),
    'all project queries'
  );
};

// Handle feature pending events
const handleFeaturePending = (qc: QueryClient, projectId: string, envKey: string, featureId: string) => {
  console.log(`[Realtime] Handling feature pending: ${featureId} in project ${projectId}`);
  
  // Invalidate pending changes queries
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'pending-changes' && query.queryKey[1] === projectId,
    'pending-changes queries'
  );
  
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-pending-changes' && query.queryKey[1] === projectId,
    'project-pending-changes queries'
  );
  
  // Refetch pending changes queries
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'pending-changes' && query.queryKey[1] === projectId,
    'pending-changes queries'
  );
  
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-pending-changes' && query.queryKey[1] === projectId,
    'project-pending-changes queries'
  );
  
  // Invalidate pending changes count for menu badge
  invalidatePendingChangesCount(qc, projectId);
  
  // Refetch pending changes count for menu badge
  refetchPendingChangesCount(qc, projectId);
};

// Handle feature deletion events
const handleFeatureDeleted = (qc: QueryClient, projectId: string, envKey: string, featureId: string) => {
  console.log(`[Realtime] Handling feature deletion: ${featureId} in project ${projectId}`);
  
  // Remove from cache and invalidate
  qc.removeQueries({ 
    predicate: (query) => query.queryKey[0] === 'project-features' && query.queryKey[1] === projectId 
  });
  
  qc.removeQueries({ 
    predicate: (query) => query.queryKey[0] === 'feature-details' && query.queryKey[1] === featureId 
  });
  
  // Invalidate and refetch project-features
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-features' && query.queryKey[1] === projectId,
    'project-features queries'
  );
  
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-features' && query.queryKey[1] === projectId,
    'project-features queries'
  );
};

// Handle child entity events (feature_params, feature_tag, etc.)
const handleChildEntityEvent = (qc: QueryClient, projectId: string, envKey: string, eventType: string) => {
  console.log(`[Realtime] Handling child entity event: ${eventType} in project ${projectId}`);
  
  // Skip approved/rejected events for pending changes to avoid hiding pending chips
  if (eventType.includes('approved') || eventType.includes('rejected')) {
    console.log(`[Realtime] Skipping ${eventType} event to avoid hiding pending chips`);
    return;
  }
  
  // For child entities, we need to invalidate all project-features queries
  // because the parent feature might have changed
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-features' && query.queryKey[1] === projectId,
    'project-features queries (child entity)'
  );
  
  // Invalidate all feature-details queries for this project
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-details',
    'feature-details queries (child entity)'
  );
  
  // Invalidate all feature-changes queries for this project
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-changes' && query.queryKey[2] === projectId,
    'feature-changes queries (child entity)'
  );
  
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-names' && query.queryKey[1] === projectId,
    'feature-names queries (child entity)'
  );
  
  // If this is a pending event, also invalidate pending changes queries
  if (eventType.includes('pending')) {
    console.log(`[Realtime] Child entity pending event detected: ${eventType}`);
    
    // Invalidate pending changes queries
    invalidateByPattern(
      qc,
      (query) => query.queryKey[0] === 'pending-changes' && query.queryKey[1] === projectId,
      'pending-changes queries (child entity pending)'
    );
    
    invalidateByPattern(
      qc,
      (query) => query.queryKey[0] === 'project-pending-changes' && query.queryKey[1] === projectId,
      'project-pending-changes queries (child entity pending)'
    );
    
    // Invalidate pending changes count for menu badge
    invalidatePendingChangesCount(qc, projectId);
  }
  
  // Refetch all project-features queries
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-features' && query.queryKey[1] === projectId,
    'project-features queries (child entity)'
  );
  
  // Refetch all feature-details queries
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-details',
    'feature-details queries (child entity)'
  );
  
  // Refetch all feature-changes queries for this project
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'feature-changes' && query.queryKey[2] === projectId,
    'feature-changes queries (child entity)'
  );
  
  // If this is a pending event, also refetch pending changes queries
  if (eventType.includes('pending')) {
    // Refetch pending changes queries
    refetchByPattern(
      qc,
      (query) => query.queryKey[0] === 'pending-changes' && query.queryKey[1] === projectId,
      'pending-changes queries (child entity pending)'
    );
    
    refetchByPattern(
      qc,
      (query) => query.queryKey[0] === 'project-pending-changes' && query.queryKey[1] === projectId,
      'project-pending-changes queries (child entity pending)'
    );
    
    // Refetch pending changes count for menu badge
    refetchPendingChangesCount(qc, projectId);
  }
};

// Handle pending change events
const handlePendingChangeEvent = (qc: QueryClient, projectId: string, envKey: string) => {
  console.log(`[Realtime] Handling pending change event in project ${projectId}`);
  
  // Invalidate pending changes queries
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'pending-changes' && query.queryKey[1] === projectId,
    'pending-changes queries'
  );
  
  invalidateByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-pending-changes' && query.queryKey[1] === projectId,
    'project-pending-changes queries'
  );
  
  // Refetch pending changes queries
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'pending-changes' && query.queryKey[1] === projectId,
    'pending-changes queries'
  );
  
  refetchByPattern(
    qc,
    (query) => query.queryKey[0] === 'project-pending-changes' && query.queryKey[1] === projectId,
    'project-pending-changes queries'
  );
  
  // Invalidate pending changes count for menu badge
  invalidatePendingChangesCount(qc, projectId);
  
  // Refetch pending changes count for menu badge
  refetchPendingChangesCount(qc, projectId);
};

// Main event handler
export function handleEvent(qc: QueryClient | undefined, evt: RealtimeEvent) {
  if (!qc) {
    console.warn('[Realtime] QueryClient not available');
    return;
  }

  const { project_id: projectId, environment: envKey, entity, entity_id: entityId, type } = evt;

  // Check if this is a feature-related event
  const isFeatureEvent = type.startsWith('feature_');
  
  if (isFeatureEvent) {
    if (entity === 'feature') {
      // Direct feature events
      switch (type) {
        case 'feature_update':
        case 'feature_updated':
          handleFeatureUpdate(qc, projectId, envKey, entityId);
          break;
        case 'feature_pending':
          handleFeaturePending(qc, projectId, envKey, entityId);
          break;
        case 'feature_deleted':
          handleFeatureDeleted(qc, projectId, envKey, entityId);
          break;
        default:
          // For unknown feature types, treat as update
          handleFeatureUpdate(qc, projectId, envKey, entityId);
      }
    } else {
      // Feature-related events on child entities (feature_params, feature_tag, etc.)
      handleChildEntityEvent(qc, projectId, envKey, type);
    }
  } else {
    // Non-feature events
    switch (entity) {
      case 'pending_change':
        handlePendingChangeEvent(qc, projectId, envKey);
        break;
      default:
        console.log(`[Realtime] Unknown entity type: ${entity}`);
    }
  }
}