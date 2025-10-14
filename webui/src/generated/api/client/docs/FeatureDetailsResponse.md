# FeatureDetailsResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**feature** | [**FeatureExtended**](FeatureExtended.md) |  | [default to undefined]
**variants** | [**Array&lt;FlagVariant&gt;**](FlagVariant.md) |  | [default to undefined]
**rules** | [**Array&lt;Rule&gt;**](Rule.md) |  | [default to undefined]
**tags** | [**Array&lt;ProjectTag&gt;**](ProjectTag.md) |  | [default to undefined]
**algorithms** | [**Array&lt;FeatureAlgorithm&gt;**](FeatureAlgorithm.md) |  | [optional] [default to undefined]

## Example

```typescript
import { FeatureDetailsResponse } from './api';

const instance: FeatureDetailsResponse = {
    feature,
    variants,
    rules,
    tags,
    algorithms,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
