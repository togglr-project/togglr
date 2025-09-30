# WebSocket Realtime Updates

This module handles real-time updates for features and pending changes via WebSocket connections.

## How it works

1. **WebSocket Connection**: Automatically connects to `/api/ws` when on a project page
2. **Event Handling**: Processes incoming events and updates React Query cache
3. **Smart Updates**: Instead of invalidating all queries, updates specific features in cache

## Event Types

### Feature Events

- `feature_update` / `feature_updated`: Updates specific feature in cache
- `feature_pending`: Marks feature as having pending changes
- `feature_deleted`: Removes feature from all lists

### Pending Change Events

- `pending_change`: Invalidates pending changes queries

## Example Events

```json
{
  "source": "audit",
  "type": "feature_update",
  "timestamp": "2025-09-30T20:24:37.320813+03:00",
  "project_id": "76f7c810-8ee7-41d1-9276-9663c2f89ab5",
  "environment": "prod",
  "entity": "feature",
  "entity_id": "969f1e7b-38f9-4934-ac74-29cc943ff184",
  "action": "update"
}
```

```json
{
  "source": "pending",
  "type": "feature_pending",
  "timestamp": "2025-09-30T20:27:30.660517+03:00",
  "project_id": "76f7c810-8ee7-41d1-9276-9663c2f89ab5",
  "environment": "prod",
  "entity": "feature",
  "entity_id": "969f1e7b-38f9-4934-ac74-29cc943ff184",
  "action": "pending"
}
```

## Benefits

- **No more polling**: Removed `refetchInterval` from pending changes queries
- **Faster updates**: Features update immediately when changed
- **Better UX**: Pending status shows instantly without page refresh
- **Reduced server load**: No unnecessary API calls every 30 seconds
