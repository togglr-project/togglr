# FeatureExtended


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**project_id** | **string** |  | [default to undefined]
**key** | **string** |  | [default to undefined]
**name** | **string** |  | [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**kind** | [**FeatureKind**](FeatureKind.md) |  | [default to undefined]
**default_variant** | **string** |  | [default to undefined]
**enabled** | **boolean** |  | [default to undefined]
**rollout_key** | **string** |  | [optional] [default to undefined]
**created_at** | **string** |  | [default to undefined]
**updated_at** | **string** |  | [default to undefined]
**is_active** | **boolean** | Indicates if the feature is currently active (taking schedules) | [default to undefined]
**next_state** | **boolean** | Indicates the next state the feature will transition to based on schedule (null if no schedule) | [optional] [default to undefined]
**next_state_time** | **string** | Timestamp when the feature will transition to the next state (null if no schedule) | [optional] [default to undefined]

## Example

```typescript
import { FeatureExtended } from './api';

const instance: FeatureExtended = {
    id,
    project_id,
    key,
    name,
    description,
    kind,
    default_variant,
    enabled,
    rollout_key,
    created_at,
    updated_at,
    is_active,
    next_state,
    next_state_time,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
