# Rule


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**feature_id** | **string** |  | [default to undefined]
**conditions** | [**RuleConditionExpression**](RuleConditionExpression.md) |  | [default to undefined]
**segment_id** | **string** |  | [optional] [default to undefined]
**is_customized** | **boolean** |  | [default to undefined]
**action** | [**RuleAction**](RuleAction.md) |  | [default to undefined]
**flag_variant_id** | **string** |  | [optional] [default to undefined]
**priority** | **number** |  | [default to undefined]
**created_at** | **string** |  | [default to undefined]

## Example

```typescript
import { Rule } from './api';

const instance: Rule = {
    id,
    feature_id,
    conditions,
    segment_id,
    is_customized,
    action,
    flag_variant_id,
    priority,
    created_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
