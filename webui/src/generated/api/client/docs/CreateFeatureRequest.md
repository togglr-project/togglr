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
**rollout_key** | **string** |  | [optional] [default to undefined]
**variants** | [**Array&lt;CreateFlagVariantInline&gt;**](CreateFlagVariantInline.md) | Optional list of flag variants to create along with the feature | [optional] [default to undefined]
**rules** | [**Array&lt;CreateRuleInline&gt;**](CreateRuleInline.md) | Optional list of rules to create along with the feature | [optional] [default to undefined]

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
    rollout_key,
    variants,
    rules,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
