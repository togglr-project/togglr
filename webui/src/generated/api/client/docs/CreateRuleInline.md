# CreateRuleInline


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | Client-provided UUID for the rule | [default to undefined]
**conditions** | [**RuleConditionExpression**](RuleConditionExpression.md) |  | [default to undefined]
**segment_id** | **string** |  | [optional] [default to undefined]
**is_customized** | **boolean** |  | [default to undefined]
**action** | [**RuleAction**](RuleAction.md) |  | [default to undefined]
**flag_variant_id** | **string** |  | [optional] [default to undefined]
**priority** | **number** |  | [optional] [default to undefined]
**environment_key** | **string** | Environment key (dev, stage, prod) for this rule | [default to undefined]

## Example

```typescript
import { CreateRuleInline } from './api';

const instance: CreateRuleInline = {
    id,
    conditions,
    segment_id,
    is_customized,
    action,
    flag_variant_id,
    priority,
    environment_key,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
