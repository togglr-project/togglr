# FeatureSchedule


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**project_id** | **string** |  | [default to undefined]
**feature_id** | **string** |  | [default to undefined]
**starts_at** | **string** |  | [optional] [default to undefined]
**ends_at** | **string** |  | [optional] [default to undefined]
**cron_expr** | **string** |  | [optional] [default to undefined]
**timezone** | **string** |  | [default to undefined]
**action** | [**FeatureScheduleAction**](FeatureScheduleAction.md) |  | [default to undefined]
**created_at** | **string** |  | [default to undefined]

## Example

```typescript
import { FeatureSchedule } from './api';

const instance: FeatureSchedule = {
    id,
    project_id,
    feature_id,
    starts_at,
    ends_at,
    cron_expr,
    timezone,
    action,
    created_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
