# PendingSummary


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**project_id** | **string** |  | [optional] [default to undefined]
**environment_id** | **string** |  | [optional] [default to undefined]
**environment_key** | **string** |  | [optional] [default to undefined]
**total_pending** | **number** |  | [optional] [default to undefined]
**pending_feature_changes** | **number** |  | [optional] [default to undefined]
**pending_guarded_changes** | **number** |  | [optional] [default to undefined]
**oldest_request_at** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { PendingSummary } from './api';

const instance: PendingSummary = {
    project_id,
    environment_id,
    environment_key,
    total_pending,
    pending_feature_changes,
    pending_guarded_changes,
    oldest_request_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
