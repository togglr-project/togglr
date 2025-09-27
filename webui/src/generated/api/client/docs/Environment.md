# Environment


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **number** | Environment ID | [default to undefined]
**project_id** | **string** | Project ID | [default to undefined]
**key** | **string** | Environment key (dev, stage, prod) | [default to undefined]
**name** | **string** | Human-readable environment name | [default to undefined]
**api_key** | **string** | API key for this environment | [default to undefined]
**created_at** | **string** | Creation timestamp | [default to undefined]

## Example

```typescript
import { Environment } from './api';

const instance: Environment = {
    id,
    project_id,
    key,
    name,
    api_key,
    created_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
