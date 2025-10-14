# FeatureAlgorithm


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**feature_id** | **string** |  | [default to undefined]
**environment_id** | **number** |  | [default to undefined]
**algorithm_slug** | **string** |  | [default to undefined]
**enabled** | **boolean** |  | [default to undefined]
**settings** | **{ [key: string]: number; }** | Numeric settings for the feature algorithm | [default to undefined]

## Example

```typescript
import { FeatureAlgorithm } from './api';

const instance: FeatureAlgorithm = {
    feature_id,
    environment_id,
    algorithm_slug,
    enabled,
    settings,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
