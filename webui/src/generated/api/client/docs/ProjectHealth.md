# ProjectHealth


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**project_id** | **string** |  | [optional] [default to undefined]
**project_name** | **string** |  | [optional] [default to undefined]
**environment_id** | **string** |  | [optional] [default to undefined]
**environment_key** | **string** |  | [optional] [default to undefined]
**total_features** | **number** |  | [optional] [default to undefined]
**enabled_features** | **number** |  | [optional] [default to undefined]
**disabled_features** | **number** |  | [optional] [default to undefined]
**auto_disable_managed_features** | **number** |  | [optional] [default to undefined]
**uncategorized_features** | **number** |  | [optional] [default to undefined]
**guarded_features** | **number** |  | [optional] [default to undefined]
**pending_features** | **number** |  | [optional] [default to undefined]
**pending_guarded_features** | **number** |  | [optional] [default to undefined]
**health_status** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { ProjectHealth } from './api';

const instance: ProjectHealth = {
    project_id,
    project_name,
    environment_id,
    environment_key,
    total_features,
    enabled_features,
    disabled_features,
    auto_disable_managed_features,
    uncategorized_features,
    guarded_features,
    pending_features,
    pending_guarded_features,
    health_status,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
