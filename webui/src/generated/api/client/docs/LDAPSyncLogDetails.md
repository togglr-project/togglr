# LDAPSyncLogDetails


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
**stack_trace** | **string** |  | [optional] [default to undefined]
**ldap_error_code** | **number** |  | [optional] [default to undefined]
**ldap_error_message** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { LDAPSyncLogDetails } from './api';

const instance: LDAPSyncLogDetails = {
    id,
    timestamp,
    level,
    message,
    username,
    details,
    sync_session_id,
    stack_trace,
    ldap_error_code,
    ldap_error_message,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
