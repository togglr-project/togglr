# LDAPConnectionTest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**url** | **string** | LDAP server URL | [default to undefined]
**bind_dn** | **string** | DN for binding to LDAP server | [default to undefined]
**bind_password** | **string** | Password for binding to LDAP server | [default to undefined]
**user_base_dn** | **string** | Base DN for user search | [optional] [default to undefined]
**user_filter** | **string** | Filter for user search | [optional] [default to undefined]
**user_name_attr** | **string** | Attribute for username | [optional] [default to undefined]
**start_tls** | **boolean** | Whether to use StartTLS | [optional] [default to undefined]
**insecure_tls** | **boolean** | Whether to skip TLS certificate verification | [optional] [default to undefined]
**timeout** | **string** | Connection timeout | [optional] [default to undefined]

## Example

```typescript
import { LDAPConnectionTest } from './api';

const instance: LDAPConnectionTest = {
    url,
    bind_dn,
    bind_password,
    user_base_dn,
    user_filter,
    user_name_attr,
    start_tls,
    insecure_tls,
    timeout,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
