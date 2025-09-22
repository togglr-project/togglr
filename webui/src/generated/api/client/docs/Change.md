# Change


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **number** | Audit log entry ID | [default to undefined]
**entity** | [**EntityType**](EntityType.md) |  | [default to undefined]
**entity_id** | **string** | ID of the changed entity | [default to undefined]
**action** | [**AuditAction**](AuditAction.md) |  | [default to undefined]
**old_value** | **object** | Previous value of the entity (null for create actions) | [default to undefined]
**new_value** | **object** | New value of the entity (null for delete actions) | [default to undefined]

## Example

```typescript
import { Change } from './api';

const instance: Change = {
    id,
    entity,
    entity_id,
    action,
    old_value,
    new_value,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
