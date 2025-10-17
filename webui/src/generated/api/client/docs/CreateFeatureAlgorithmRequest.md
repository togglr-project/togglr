# CreateFeatureAlgorithmRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**algorithm_slug** | **string** | Algorithm slug from the /api/v1/algorithms list | [default to undefined]
**environment_id** | **number** | Environment ID for which the algorithm applies | [default to undefined]
**settings** | **{ [key: string]: number; }** | Numeric algorithm settings overriding defaults | [default to undefined]
**enabled** | **boolean** |  | [default to false]

## Example

```typescript
import { CreateFeatureAlgorithmRequest } from './api';

const instance: CreateFeatureAlgorithmRequest = {
    algorithm_slug,
    environment_id,
    settings,
    enabled,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
