# CreateRuleRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**conditions** | [**RuleConditionExpression**](RuleConditionExpression.md) |  | [default to undefined]
**segment_id** | **string** |  | [optional] [default to undefined]
**is_customized** | **boolean** |  | [default to undefined]
**action** | [**RuleAction**](RuleAction.md) |  | [default to undefined]
**flag_variant_id** | **string** |  | [optional] [default to undefined]
**priority** | **number** |  | [optional] [default to undefined]

## Example

```typescript
import { CreateRuleRequest } from './api';

const instance: CreateRuleRequest = {
    conditions,
    segment_id,
    is_customized,
    action,
    flag_variant_id,
    priority,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
