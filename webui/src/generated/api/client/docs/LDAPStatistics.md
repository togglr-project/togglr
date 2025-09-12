# LDAPStatistics


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ldap_users** | **number** |  | [optional] [default to undefined]
**local_users** | **number** |  | [optional] [default to undefined]
**active_users** | **number** |  | [optional] [default to undefined]
**inactive_users** | **number** |  | [optional] [default to undefined]
**sync_history** | [**Array&lt;LDAPStatisticsSyncHistoryInner&gt;**](LDAPStatisticsSyncHistoryInner.md) |  | [optional] [default to undefined]
**sync_success_rate** | **number** |  | [optional] [default to undefined]

## Example

```typescript
import { LDAPStatistics } from './api';

const instance: LDAPStatistics = {
    ldap_users,
    local_users,
    active_users,
    inactive_users,
    sync_history,
    sync_success_rate,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
