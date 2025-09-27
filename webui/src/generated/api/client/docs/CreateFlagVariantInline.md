# CreateFlagVariantInline


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | Client-provided UUID for the variant | [default to undefined]
**name** | **string** |  | [default to undefined]
**rollout_percent** | **number** |  | [default to undefined]
**environment_key** | **string** | Environment key (dev, stage, prod) for this variant | [default to undefined]

## Example

```typescript
import { CreateFlagVariantInline } from './api';

const instance: CreateFlagVariantInline = {
    id,
    name,
    rollout_percent,
    environment_key,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
