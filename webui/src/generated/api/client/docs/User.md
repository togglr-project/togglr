# User


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **number** |  | [default to undefined]
**username** | **string** |  | [default to undefined]
**email** | **string** |  | [default to undefined]
**is_superuser** | **boolean** |  | [default to undefined]
**is_active** | **boolean** |  | [default to undefined]
**is_external** | **boolean** |  | [default to undefined]
**is_tmp_password** | **boolean** |  | [default to undefined]
**two_fa_enabled** | **boolean** |  | [default to undefined]
**license_accepted** | **boolean** | Flag indicating whether the user has accepted the license agreement | [default to undefined]
**created_at** | **string** |  | [default to undefined]
**last_login** | **string** |  | [optional] [default to undefined]
**project_permissions** | **{ [key: string]: Array&lt;string&gt;; }** | Map of project_id to list of permission keys for that project. Contains only projects where user has membership. | [optional] [default to undefined]

## Example

```typescript
import { User } from './api';

const instance: User = {
    id,
    username,
    email,
    is_superuser,
    is_active,
    is_external,
    is_tmp_password,
    two_fa_enabled,
    license_accepted,
    created_at,
    last_login,
    project_permissions,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
