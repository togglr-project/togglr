# UpdateNotificationSettingRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**type** | **string** | Type of notification channel (email, mattermost, slack, etc.) | [optional] [default to undefined]
**config** | **string** | Configuration for the notification channel (JSONB in database) | [optional] [default to undefined]
**enabled** | **boolean** |  | [optional] [default to undefined]

## Example

```typescript
import { UpdateNotificationSettingRequest } from './api';

const instance: UpdateNotificationSettingRequest = {
    type,
    config,
    enabled,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
