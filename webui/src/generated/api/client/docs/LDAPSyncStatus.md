# LDAPSyncStatus


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**status** | **string** |  | [default to undefined]
**is_running** | **boolean** |  | [default to undefined]
**last_sync_time** | **string** |  | [optional] [default to undefined]
**total_users** | **number** |  | [default to undefined]
**synced_users** | **number** |  | [default to undefined]
**errors** | **number** |  | [default to undefined]
**warnings** | **number** |  | [default to undefined]
**last_sync_duration** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { LDAPSyncStatus } from './api';

const instance: LDAPSyncStatus = {
    status,
    is_running,
    last_sync_time,
    total_users,
    synced_users,
    errors,
    warnings,
    last_sync_duration,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
