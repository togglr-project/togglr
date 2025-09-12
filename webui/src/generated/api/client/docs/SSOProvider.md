# SSOProvider


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **string** | Internal name of the provider | [default to undefined]
**display_name** | **string** | Display name for UI | [default to undefined]
**type** | **string** | Type of SSO provider | [default to undefined]
**icon_url** | **string** | URL to provider icon | [optional] [default to undefined]
**enabled** | **boolean** | Whether the provider is enabled | [default to undefined]

## Example

```typescript
import { SSOProvider } from './api';

const instance: SSOProvider = {
    name,
    display_name,
    type,
    icon_url,
    enabled,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
