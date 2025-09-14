# CreateFeatureRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**key** | **string** |  | [default to undefined]
**name** | **string** |  | [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**kind** | [**FeatureKind**](FeatureKind.md) |  | [default to undefined]
**default_variant** | **string** |  | [default to undefined]
**enabled** | **boolean** |  | [optional] [default to undefined]

## Example

```typescript
import { CreateFeatureRequest } from './api';

const instance: CreateFeatureRequest = {
    key,
    name,
    description,
    kind,
    default_variant,
    enabled,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
