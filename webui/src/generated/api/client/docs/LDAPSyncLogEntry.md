# LDAPSyncLogEntry


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **number** |  | [default to undefined]
**timestamp** | **string** |  | [default to undefined]
**level** | **string** |  | [default to undefined]
**message** | **string** |  | [default to undefined]
**username** | **string** |  | [optional] [default to undefined]
**details** | **string** |  | [optional] [default to undefined]
**sync_session_id** | **string** |  | [default to undefined]

## Example

```typescript
import { LDAPSyncLogEntry } from './api';

const instance: LDAPSyncLogEntry = {
    id,
    timestamp,
    level,
    message,
    username,
    details,
    sync_session_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
