# Feature


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**project_id** | **string** |  | [default to undefined]
**key** | **string** |  | [default to undefined]
**name** | **string** |  | [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**kind** | [**FeatureKind**](FeatureKind.md) |  | [default to undefined]
**rollout_key** | **string** |  | [optional] [default to undefined]
**enabled** | **boolean** | Whether the feature is enabled in the specified environment | [default to undefined]
**default_value** | **string** | Default value for the feature in the specified environment | [default to undefined]
**created_at** | **string** |  | [default to undefined]
**updated_at** | **string** |  | [default to undefined]

## Example

```typescript
import { Feature } from './api';

const instance: Feature = {
    id,
    project_id,
    key,
    name,
    description,
    kind,
    rollout_key,
    enabled,
    default_value,
    created_at,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
