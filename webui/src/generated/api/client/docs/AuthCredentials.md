# AuthCredentials


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**method** | **string** |  | [default to undefined]
**credential** | **string** |  | [default to undefined]
**session_id** | **string** | Session ID for TOTP approval (required when method is \&#39;totp\&#39;) | [optional] [default to undefined]

## Example

```typescript
import { AuthCredentials } from './api';

const instance: AuthCredentials = {
    method,
    credential,
    session_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
