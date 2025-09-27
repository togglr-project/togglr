# PendingChangeResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**environment_key** | **string** |  | [default to undefined]
**project_id** | **string** |  | [default to undefined]
**requested_by** | **string** |  | [default to undefined]
**request_user_id** | **number** |  | [optional] [default to undefined]
**change** | [**PendingChangePayload**](PendingChangePayload.md) |  | [default to undefined]
**status** | **string** |  | [default to undefined]
**created_at** | **string** |  | [default to undefined]
**approved_by** | **string** |  | [optional] [default to undefined]
**approved_user_id** | **number** |  | [optional] [default to undefined]
**approved_at** | **string** |  | [optional] [default to undefined]
**rejected_by** | **string** |  | [optional] [default to undefined]
**rejected_at** | **string** |  | [optional] [default to undefined]
**rejection_reason** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { PendingChangeResponse } from './api';

const instance: PendingChangeResponse = {
    id,
    environment_key,
    project_id,
    requested_by,
    request_user_id,
    change,
    status,
    created_at,
    approved_by,
    approved_user_id,
    approved_at,
    rejected_by,
    rejected_at,
    rejection_reason,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
