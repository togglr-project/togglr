import { QueryClient } from '@tanstack/react-query';

export type RealtimeEvent = {
  source: 'pending' | 'audit' | string;
  type: string; // e.g., feature_updated
  timestamp: string;
  project_id: string;
  environment: string; // environment key
  entity: string; // e.g., 'feature', 'pending_change'
  entity_id: string;
  action: string; // e.g., 'updated'
};

const invalidate = (qc: QueryClient, keys: any[]) => {
  try { qc.invalidateQueries({ queryKey: keys, exact: false }); } catch {}
};

export function handleEvent(qc: QueryClient | undefined, evt: RealtimeEvent) {
  if (!qc) return;

  const envKey = evt.environment;
  const projectId = evt.project_id;

  switch (evt.entity) {
    case 'feature': {
      // Invalidate lists and detail views that likely include this feature
      invalidate(qc, ['project-features', projectId, envKey]);
      invalidate(qc, ['feature-details', evt.entity_id, envKey]);
      invalidate(qc, ['feature-timelines', projectId, envKey]);
      // some pages might use a generic 'dashboard' key
      invalidate(qc, ['dashboard', projectId]);
      break;
    }
    case 'pending_change': {
      invalidate(qc, ['pending-changes', projectId, envKey]);
      break;
    }
    default: {
      // Fallback: invalidate project-level queries
      invalidate(qc, [projectId]);
      break;
    }
  }
}
