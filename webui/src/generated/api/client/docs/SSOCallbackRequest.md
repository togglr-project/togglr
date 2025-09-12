# SSOCallbackRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**provider** | **string** | Name of the SSO provider | [default to undefined]
**response** | **string** | Response from SSO provider (code for OIDC, SAML response for SAML) | [default to undefined]
**state** | **string** | State parameter for CSRF protection | [default to undefined]

## Example

```typescript
import { SSOCallbackRequest } from './api';

const instance: SSOCallbackRequest = {
    provider,
    response,
    state,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
