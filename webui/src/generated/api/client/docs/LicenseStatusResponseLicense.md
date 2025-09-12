# LicenseStatusResponseLicense


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** | License ID | [optional] [default to undefined]
**type** | [**LicenseType**](LicenseType.md) |  | [optional] [default to undefined]
**issued_at** | **string** | When the license was issued | [optional] [default to undefined]
**expires_at** | **string** | When the license expires | [optional] [default to undefined]
**is_valid** | **boolean** | Whether the license is currently valid | [optional] [default to undefined]
**is_expired** | **boolean** | Whether the license has expired | [optional] [default to undefined]
**days_until_expiry** | **number** | Number of days until license expires (negative if expired) | [optional] [default to undefined]
**license_text** | **string** | The full license text | [optional] [default to undefined]
**features** | [**Array&lt;LicenseFeature&gt;**](LicenseFeature.md) | List of features available in this license | [optional] [default to undefined]

## Example

```typescript
import { LicenseStatusResponseLicense } from './api';

const instance: LicenseStatusResponseLicense = {
    id,
    type,
    issued_at,
    expires_at,
    is_valid,
    is_expired,
    days_until_expiry,
    license_text,
    features,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
