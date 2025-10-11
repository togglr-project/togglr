# NotificationSetting


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **number** |  | [default to undefined]
**project_id** | **string** |  | [default to undefined]
**environment_id** | **number** |  | [default to undefined]
**environment_key** | **string** |  | [default to undefined]
**type** | **string** | Type of notification channel (email, mattermost, slack, etc.) | [default to undefined]
**config** | **string** | Configuration for the notification channel (JSONB in database) | [default to undefined]
**enabled** | **boolean** |  | [default to undefined]
**created_at** | **string** |  | [default to undefined]
**updated_at** | **string** |  | [default to undefined]

## Example

```typescript
import { NotificationSetting } from './api';

const instance: NotificationSetting = {
    id,
    project_id,
    environment_id,
    environment_key,
    type,
    config,
    enabled,
    created_at,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
