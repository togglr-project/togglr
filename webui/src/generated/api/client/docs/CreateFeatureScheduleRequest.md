# CreateFeatureScheduleRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**starts_at** | **string** |  | [optional] [default to undefined]
**ends_at** | **string** |  | [optional] [default to undefined]
**cron_expr** | **string** |  | [optional] [default to undefined]
**cron_duration** | **string** | Duration for cron-based schedules. When cron triggers, feature will be enabled/disabled for this duration. Format: \&#39;1h30m\&#39;, \&#39;45m\&#39;, \&#39;2h\&#39;, etc. | [optional] [default to undefined]
**timezone** | **string** |  | [default to 'UTC']
**action** | [**FeatureScheduleAction**](FeatureScheduleAction.md) |  | [default to undefined]

## Example

```typescript
import { CreateFeatureScheduleRequest } from './api';

const instance: CreateFeatureScheduleRequest = {
    starts_at,
    ends_at,
    cron_expr,
    cron_duration,
    timezone,
    action,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
