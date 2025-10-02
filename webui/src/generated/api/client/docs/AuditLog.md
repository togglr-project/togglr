# AuditLog


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **number** |  | [default to undefined]
**project_id** | **string** |  | [default to undefined]
**environment_id** | **number** |  | [default to undefined]
**environment_key** | **string** |  | [optional] [default to undefined]
**entity** | **string** |  | [default to undefined]
**entity_id** | **string** |  | [default to undefined]
**feature_id** | **string** |  | [optional] [default to undefined]
**action** | **string** |  | [default to undefined]
**actor** | **string** |  | [default to undefined]
**username** | **string** |  | [optional] [default to undefined]
**request_id** | **string** |  | [optional] [default to undefined]
**old_value** | **{ [key: string]: any; }** |  | [optional] [default to undefined]
**new_value** | **{ [key: string]: any; }** |  | [optional] [default to undefined]
**created_at** | **string** |  | [default to undefined]

## Example

```typescript
import { AuditLog } from './api';

const instance: AuditLog = {
    id,
    project_id,
    environment_id,
    environment_key,
    entity,
    entity_id,
    feature_id,
    action,
    actor,
    username,
    request_id,
    old_value,
    new_value,
    created_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
