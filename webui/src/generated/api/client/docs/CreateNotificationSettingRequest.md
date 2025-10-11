# CreateNotificationSettingRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**type** | [**NotificationChannelType**](NotificationChannelType.md) |  | [default to undefined]
**config** | **string** | Configuration for the notification channel (JSONB in database) | [default to undefined]
**enabled** | **boolean** |  | [optional] [default to true]

## Example

```typescript
import { CreateNotificationSettingRequest } from './api';

const instance: CreateNotificationSettingRequest = {
    type,
    config,
    enabled,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
