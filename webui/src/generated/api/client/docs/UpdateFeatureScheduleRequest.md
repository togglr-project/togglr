# UpdateFeatureScheduleRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**starts_at** | **string** |  | [optional] [default to undefined]
**ends_at** | **string** |  | [optional] [default to undefined]
**cron_expr** | **string** |  | [optional] [default to undefined]
**timezone** | **string** |  | [default to undefined]
**action** | [**FeatureScheduleAction**](FeatureScheduleAction.md) |  | [default to undefined]

## Example

```typescript
import { UpdateFeatureScheduleRequest } from './api';

const instance: UpdateFeatureScheduleRequest = {
    starts_at,
    ends_at,
    cron_expr,
    timezone,
    action,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
