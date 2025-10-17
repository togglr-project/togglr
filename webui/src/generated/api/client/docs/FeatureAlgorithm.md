# FeatureAlgorithm


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**feature_id** | **string** |  | [default to undefined]
**project_id** | **string** |  | [default to undefined]
**environment_id** | **number** |  | [default to undefined]
**algorithm_slug** | **string** |  | [default to undefined]
**enabled** | **boolean** |  | [default to undefined]
**settings** | **{ [key: string]: number; }** | Numeric settings for the feature algorithm | [default to undefined]
**feature** | [**Feature**](.md) |  | [default to undefined]

## Example

```typescript
import { FeatureAlgorithm } from './api';

const instance: FeatureAlgorithm = {
    id,
    feature_id,
    project_id,
    environment_id,
    algorithm_slug,
    enabled,
    settings,
    feature,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
