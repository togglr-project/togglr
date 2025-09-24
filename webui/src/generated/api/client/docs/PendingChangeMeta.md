# PendingChangeMeta


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**reason** | **string** |  | [default to undefined]
**client** | **string** |  | [default to undefined]
**origin** | **string** |  | [default to undefined]
**single_user_project** | **boolean** | True if the project has only 1 active user (enables auto-approve) | [optional] [default to undefined]

## Example

```typescript
import { PendingChangeMeta } from './api';

const instance: PendingChangeMeta = {
    reason,
    client,
    origin,
    single_user_project,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
