# CreateRuleInline


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | Client-provided UUID for the rule | [default to undefined]
**conditions** | [**Array&lt;RuleCondition&gt;**](RuleCondition.md) |  | [default to undefined]
**action** | [**RuleAction**](RuleAction.md) |  | [default to undefined]
**flag_variant_id** | **string** |  | [optional] [default to undefined]
**priority** | **number** |  | [optional] [default to undefined]

## Example

```typescript
import { CreateRuleInline } from './api';

const instance: CreateRuleInline = {
    id,
    conditions,
    action,
    flag_variant_id,
    priority,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
