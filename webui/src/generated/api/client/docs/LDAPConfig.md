# LDAPConfig


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **boolean** | Whether LDAP integration is enabled | [default to undefined]
**url** | **string** | LDAP server URL | [default to undefined]
**bind_dn** | **string** | DN for binding to LDAP server | [default to undefined]
**bind_password** | **string** | Password for binding to LDAP server | [default to undefined]
**user_base_dn** | **string** | Base DN for user search | [default to undefined]
**user_filter** | **string** | Filter for user search | [default to undefined]
**user_name_attr** | **string** | Attribute for username | [default to undefined]
**user_email_attr** | **string** | Attribute for user email | [default to undefined]
**start_tls** | **boolean** | Whether to use StartTLS | [default to undefined]
**insecure_tls** | **boolean** | Whether to skip TLS certificate verification | [default to undefined]
**timeout** | **string** | Connection timeout | [default to undefined]
**sync_interval** | **number** | Background synchronization interval | [default to undefined]

## Example

```typescript
import { LDAPConfig } from './api';

const instance: LDAPConfig = {
    enabled,
    url,
    bind_dn,
    bind_password,
    user_base_dn,
    user_filter,
    user_name_attr,
    user_email_attr,
    start_tls,
    insecure_tls,
    timeout,
    sync_interval,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
