# ChangeGroup


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**request_id** | **string** | Request ID that groups related changes | [default to undefined]
**actor** | **string** | Who made the changes (system, sdk, user:&lt;user_id&gt;) | [default to undefined]
**username** | **string** | Username of the user who made the changes | [default to undefined]
**created_at** | **string** | When the changes were made | [default to undefined]
**changes** | [**Array&lt;Change&gt;**](Change.md) | List of changes made in this request | [default to undefined]

## Example

```typescript
import { ChangeGroup } from './api';

const instance: ChangeGroup = {
    request_id,
    actor,
    username,
    created_at,
    changes,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
