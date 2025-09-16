# DefaultApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addProject**](#addproject) | **POST** /api/v1/projects/add | Add new project|
|[**archiveProject**](#archiveproject) | **DELETE** /api/v1/projects/{project_id} | Archive a project|
|[**cancelLDAPSync**](#cancelldapsync) | **DELETE** /api/v1/ldap/sync/cancel | Cancel ongoing synchronization|
|[**confirm2FA**](#confirm2fa) | **POST** /api/v1/users/me/2fa/confirm | Approve enable 2FA (code from app)|
|[**consumeSAMLAssertion**](#consumesamlassertion) | **POST** /api/v1/saml/acs | Assertion Consumer Service (ACS) endpoint|
|[**createFeatureFlagVariant**](#createfeatureflagvariant) | **POST** /api/v1/features/{feature_id}/variants | Create flag variant for feature|
|[**createFeatureRule**](#createfeaturerule) | **POST** /api/v1/features/{feature_id}/rules | Create rule for feature|
|[**createFeatureSchedule**](#createfeatureschedule) | **POST** /api/v1/features/{feature_id}/schedules | Create schedule for feature|
|[**createProjectFeature**](#createprojectfeature) | **POST** /api/v1/projects/{project_id}/features | Create feature for project|
|[**createUser**](#createuser) | **POST** /api/v1/users | Create a new user (superuser only)|
|[**deleteFeature**](#deletefeature) | **DELETE** /api/v1/features/{feature_id} | Delete feature|
|[**deleteFeatureSchedule**](#deletefeatureschedule) | **DELETE** /api/v1/feature-schedules/{schedule_id} | Delete feature schedule by ID|
|[**deleteLDAPConfig**](#deleteldapconfig) | **DELETE** /api/v1/ldap/config | Delete LDAP configuration|
|[**deleteUser**](#deleteuser) | **DELETE** /api/v1/users/{user_id} | Delete a user (superuser only, cannot delete superusers)|
|[**disable2FA**](#disable2fa) | **POST** /api/v1/users/me/2fa/disable | Disable 2FA (using email-confirmation)|
|[**forgotPassword**](#forgotpassword) | **POST** /api/v1/auth/forgot-password | Request a password reset|
|[**getCurrentUser**](#getcurrentuser) | **GET** /api/v1/users/me | Get current user information|
|[**getFeature**](#getfeature) | **GET** /api/v1/features/{feature_id} | Get feature with rules and variants|
|[**getFeatureSchedule**](#getfeatureschedule) | **GET** /api/v1/feature-schedules/{schedule_id} | Get feature schedule by ID|
|[**getLDAPConfig**](#getldapconfig) | **GET** /api/v1/ldap/config | Get LDAP configuration|
|[**getLDAPStatistics**](#getldapstatistics) | **GET** /api/v1/ldap/statistics | Get LDAP statistics|
|[**getLDAPSyncLogDetails**](#getldapsynclogdetails) | **GET** /api/v1/ldap/sync/logs/{id} | Get synchronization log details|
|[**getLDAPSyncLogs**](#getldapsynclogs) | **GET** /api/v1/ldap/sync/logs | Get synchronization logs|
|[**getLDAPSyncProgress**](#getldapsyncprogress) | **GET** /api/v1/ldap/sync/progress | Get synchronization progress|
|[**getLDAPSyncStatus**](#getldapsyncstatus) | **GET** /api/v1/ldap/sync/status | Get synchronization status|
|[**getLicenseStatus**](#getlicensestatus) | **GET** /api/v1/license/status | Get license status|
|[**getProductInfo**](#getproductinfo) | **GET** /api/v1/product/info | Get product information including client ID|
|[**getProject**](#getproject) | **GET** /api/v1/projects/{project_id} | Get project details|
|[**getSAMLMetadata**](#getsamlmetadata) | **GET** /api/v1/saml/metadata | Get SAML metadata|
|[**getSSOProviders**](#getssoproviders) | **GET** /api/v1/auth/sso/providers | Get available SSO providers|
|[**listAllFeatureSchedules**](#listallfeatureschedules) | **GET** /api/v1/feature-schedules | List all feature schedules|
|[**listFeatureFlagVariants**](#listfeatureflagvariants) | **GET** /api/v1/features/{feature_id}/variants | List flag variants for feature|
|[**listFeatureRules**](#listfeaturerules) | **GET** /api/v1/features/{feature_id}/rules | List rules for feature|
|[**listFeatureSchedules**](#listfeatureschedules) | **GET** /api/v1/features/{feature_id}/schedules | List schedules for feature|
|[**listProjectFeatures**](#listprojectfeatures) | **GET** /api/v1/projects/{project_id}/features | List features for project|
|[**listProjects**](#listprojects) | **GET** /api/v1/projects | Get projects list|
|[**listUsers**](#listusers) | **GET** /api/v1/users | List all users (superuser only)|
|[**login**](#login) | **POST** /api/v1/auth/login | Authenticate user and get access token|
|[**refreshToken**](#refreshtoken) | **POST** /api/v1/auth/refresh | Refresh access token|
|[**reset2FA**](#reset2fa) | **POST** /api/v1/users/me/2fa/reset | Reset/generate secret 2FA (using email-confirmation)|
|[**resetPassword**](#resetpassword) | **POST** /api/v1/auth/reset-password | Reset password using token|
|[**sSOCallback**](#ssocallback) | **POST** /api/v1/auth/sso/callback | Handle SSO callback from Keycloak|
|[**sSOInitiate**](#ssoinitiate) | **GET** /api/v1/auth/sso/initiate | Initiate SSO login flow|
|[**setSuperuserStatus**](#setsuperuserstatus) | **PUT** /api/v1/users/{user_id}/superuser | Set or unset superuser status (superuser only, cannot modify admin user)|
|[**setUserActiveStatus**](#setuseractivestatus) | **PUT** /api/v1/users/{user_id}/active | Set or unset user active status (superuser only)|
|[**setup2FA**](#setup2fa) | **POST** /api/v1/users/me/2fa/setup | Begin setup 2FA (generate secret and QR-code)|
|[**syncLDAPUsers**](#syncldapusers) | **POST** /api/v1/ldap/sync/users | Start user synchronization|
|[**testLDAPConnection**](#testldapconnection) | **POST** /api/v1/ldap/test-connection | Test LDAP connection|
|[**toggleFeature**](#togglefeature) | **PUT** /api/v1/features/{feature_id}/toggle | Toggle feature enabled state|
|[**updateFeature**](#updatefeature) | **PUT** /api/v1/features/{feature_id} | Update feature with rules and variants|
|[**updateFeatureSchedule**](#updatefeatureschedule) | **PUT** /api/v1/feature-schedules/{schedule_id} | Update feature schedule by ID|
|[**updateLDAPConfig**](#updateldapconfig) | **POST** /api/v1/ldap/config | Create or update LDAP configuration|
|[**updateLicense**](#updatelicense) | **PUT** /api/v1/license | Update license|
|[**updateLicenseAcceptance**](#updatelicenseacceptance) | **PUT** /api/v1/users/me/license-acceptance | Update license acceptance status|
|[**updateProject**](#updateproject) | **PUT** /api/v1/projects/{project_id} | Update project name and description|
|[**userChangeMyPassword**](#userchangemypassword) | **POST** /api/v1/users/me/change-password | Change my password|
|[**verify2FA**](#verify2fa) | **POST** /api/v1/auth/2fa/verify | Verify 2FA-code on login|

# **addProject**
> addProject(addProjectRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    AddProjectRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let addProjectRequest: AddProjectRequest; //

const { status, data } = await apiInstance.addProject(
    addProjectRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **addProjectRequest** | **AddProjectRequest**|  | |


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Project created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **archiveProject**
> archiveProject()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)

const { status, data } = await apiInstance.archiveProject(
    projectId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Project archived successfully |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **cancelLDAPSync**
> SuccessResponse cancelLDAPSync()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.cancelLDAPSync();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**SuccessResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Synchronization cancelled |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**404** | No active synchronization |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **confirm2FA**
> confirm2FA(twoFAConfirmRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TwoFAConfirmRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let twoFAConfirmRequest: TwoFAConfirmRequest; //

const { status, data } = await apiInstance.confirm2FA(
    twoFAConfirmRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **twoFAConfirmRequest** | **TwoFAConfirmRequest**|  | |


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | 2FA enabled |  -  |
|**400** | Invalid code |  -  |
|**401** | Unauthorized |  -  |
|**429** | Too many requests |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **consumeSAMLAssertion**
> Error consumeSAMLAssertion()

Finishes the SAML authentication flow.   The Identity Provider sends an HTTP-POST request that contains **SAMLResponse** (mandatory, Base64-encoded `<samlp:Response>` XML) and the optional **RelayState** parameter.   On success the service creates a user session (cookie or JWT) and redirects the browser to the application UI. 

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let sAMLResponse: string; //Base64-encoded IdP `<samlp:Response>` document (default to undefined)
let relayState: string; //Value round-tripped from the initial authentication request (default to undefined)

const { status, data } = await apiInstance.consumeSAMLAssertion(
    sAMLResponse,
    relayState
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **sAMLResponse** | [**string**] | Base64-encoded IdP &#x60;&lt;samlp:Response&gt;&#x60; document | defaults to undefined|
| **relayState** | [**string**] | Value round-tripped from the initial authentication request | defaults to undefined|


### Return type

**Error**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/x-www-form-urlencoded
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**302** | Successful authentication — browser will be redirected |  * Location - Target URL for the redirect (e.g. &#x60;/login/success?token&#x3D;…&#x60;) <br>  * Set-Cookie - Session cookie or JWT issued to the client <br>  |
|**400** | Malformed or expired SAML response |  -  |
|**401** | Authentication failed (invalid issuer or signature) |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createFeatureFlagVariant**
> FlagVariantResponse createFeatureFlagVariant(createFlagVariantRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateFlagVariantRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let createFlagVariantRequest: CreateFlagVariantRequest; //

const { status, data } = await apiInstance.createFeatureFlagVariant(
    featureId,
    createFlagVariantRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFlagVariantRequest** | **CreateFlagVariantRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**FlagVariantResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Flag variant created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createFeatureRule**
> RuleResponse createFeatureRule(createRuleRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateRuleRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let createRuleRequest: CreateRuleRequest; //

const { status, data } = await apiInstance.createFeatureRule(
    featureId,
    createRuleRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createRuleRequest** | **CreateRuleRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**RuleResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Rule created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature or related resource not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createFeatureSchedule**
> FeatureScheduleResponse createFeatureSchedule(createFeatureScheduleRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateFeatureScheduleRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let createFeatureScheduleRequest: CreateFeatureScheduleRequest; //

const { status, data } = await apiInstance.createFeatureSchedule(
    featureId,
    createFeatureScheduleRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFeatureScheduleRequest** | **CreateFeatureScheduleRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**FeatureScheduleResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Schedule created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature or related resource not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createProjectFeature**
> FeatureResponse createProjectFeature(createFeatureRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateFeatureRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let createFeatureRequest: CreateFeatureRequest; //

const { status, data } = await apiInstance.createProjectFeature(
    projectId,
    createFeatureRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFeatureRequest** | **CreateFeatureRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**FeatureResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Feature created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createUser**
> CreateUserResponse createUser(createUserRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateUserRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let createUserRequest: CreateUserRequest; //

const { status, data } = await apiInstance.createUser(
    createUserRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createUserRequest** | **CreateUserRequest**|  | |


### Return type

**CreateUserResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | User created successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not a superuser |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteFeature**
> deleteFeature()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)

const { status, data } = await apiInstance.deleteFeature(
    featureId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Feature deleted successfully |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteFeatureSchedule**
> deleteFeatureSchedule()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let scheduleId: string; // (default to undefined)

const { status, data } = await apiInstance.deleteFeatureSchedule(
    scheduleId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **scheduleId** | [**string**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Feature schedule deleted |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Schedule not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteLDAPConfig**
> SuccessResponse deleteLDAPConfig()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.deleteLDAPConfig();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**SuccessResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | LDAP configuration deleted |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteUser**
> deleteUser()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let userId: number; // (default to undefined)

const { status, data } = await apiInstance.deleteUser(
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**number**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | User deleted successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not a superuser or trying to delete a superuser |  -  |
|**404** | User not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **disable2FA**
> disable2FA(twoFADisableRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TwoFADisableRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let twoFADisableRequest: TwoFADisableRequest; //

const { status, data } = await apiInstance.disable2FA(
    twoFADisableRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **twoFADisableRequest** | **TwoFADisableRequest**|  | |


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | 2FA disabled |  -  |
|**400** | Invalid code |  -  |
|**401** | Unauthorized |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **forgotPassword**
> forgotPassword(forgotPasswordRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    ForgotPasswordRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let forgotPasswordRequest: ForgotPasswordRequest; //

const { status, data } = await apiInstance.forgotPassword(
    forgotPasswordRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **forgotPasswordRequest** | **ForgotPasswordRequest**|  | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Password reset email sent successfully |  -  |
|**400** | Bad request |  -  |
|**403** | External user can\&#39;t change password |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCurrentUser**
> User getCurrentUser()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getCurrentUser();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**User**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | User information |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getFeature**
> FeatureDetailsResponse getFeature()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)

const { status, data } = await apiInstance.getFeature(
    featureId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**FeatureDetailsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature details with rules and variants |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getFeatureSchedule**
> FeatureScheduleResponse getFeatureSchedule()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let scheduleId: string; // (default to undefined)

const { status, data } = await apiInstance.getFeatureSchedule(
    scheduleId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **scheduleId** | [**string**] |  | defaults to undefined|


### Return type

**FeatureScheduleResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature schedule details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Schedule not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getLDAPConfig**
> LDAPConfig getLDAPConfig()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getLDAPConfig();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**LDAPConfig**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | LDAP configuration |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**404** | LDAP configuration not found |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getLDAPStatistics**
> LDAPStatistics getLDAPStatistics()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getLDAPStatistics();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**LDAPStatistics**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | LDAP statistics |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getLDAPSyncLogDetails**
> LDAPSyncLogDetails getLDAPSyncLogDetails()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: number; // (default to undefined)

const { status, data } = await apiInstance.getLDAPSyncLogDetails(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] |  | defaults to undefined|


### Return type

**LDAPSyncLogDetails**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Synchronization log details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**404** | Log not found |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getLDAPSyncLogs**
> LDAPSyncLogs getLDAPSyncLogs()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let limit: number; // (optional) (default to 50)
let level: 'info' | 'warning' | 'error'; // (optional) (default to undefined)
let syncId: string; // (optional) (default to undefined)
let username: string; // (optional) (default to undefined)
let from: string; // (optional) (default to undefined)
let to: string; // (optional) (default to undefined)

const { status, data } = await apiInstance.getLDAPSyncLogs(
    limit,
    level,
    syncId,
    username,
    from,
    to
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **limit** | [**number**] |  | (optional) defaults to 50|
| **level** | [**&#39;info&#39; | &#39;warning&#39; | &#39;error&#39;**]**Array<&#39;info&#39; &#124; &#39;warning&#39; &#124; &#39;error&#39;>** |  | (optional) defaults to undefined|
| **syncId** | [**string**] |  | (optional) defaults to undefined|
| **username** | [**string**] |  | (optional) defaults to undefined|
| **from** | [**string**] |  | (optional) defaults to undefined|
| **to** | [**string**] |  | (optional) defaults to undefined|


### Return type

**LDAPSyncLogs**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Synchronization logs |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getLDAPSyncProgress**
> LDAPSyncProgress getLDAPSyncProgress()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getLDAPSyncProgress();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**LDAPSyncProgress**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Synchronization progress |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getLDAPSyncStatus**
> LDAPSyncStatus getLDAPSyncStatus()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getLDAPSyncStatus();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**LDAPSyncStatus**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Synchronization status |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getLicenseStatus**
> LicenseStatusResponse getLicenseStatus()

Returns the current license status including validity, expiration date, and type

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getLicenseStatus();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**LicenseStatusResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | License status retrieved successfully |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getProductInfo**
> ProductInfoResponse getProductInfo()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getProductInfo();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ProductInfoResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product information retrieved successfully |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getProject**
> ProjectResponse getProject()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)

const { status, data } = await apiInstance.getProject(
    projectId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**ProjectResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Project details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getSAMLMetadata**
> string getSAMLMetadata()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getSAMLMetadata();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/samlmetadata+xml, application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | SAML metadata XML |  -  |
|**404** | SAML metadata not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getSSOProviders**
> SSOProvidersResponse getSSOProviders()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getSSOProviders();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**SSOProvidersResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of available SSO providers |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listAllFeatureSchedules**
> Array<FeatureSchedule> listAllFeatureSchedules()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listAllFeatureSchedules();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<FeatureSchedule>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of feature schedules |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listFeatureFlagVariants**
> Array<FlagVariant> listFeatureFlagVariants()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)

const { status, data } = await apiInstance.listFeatureFlagVariants(
    featureId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**Array<FlagVariant>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of flag variants for the feature |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listFeatureRules**
> Array<Rule> listFeatureRules()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)

const { status, data } = await apiInstance.listFeatureRules(
    featureId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**Array<Rule>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of rules for the feature |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listFeatureSchedules**
> Array<FeatureSchedule> listFeatureSchedules()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)

const { status, data } = await apiInstance.listFeatureSchedules(
    featureId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**Array<FeatureSchedule>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of schedules for the feature |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectFeatures**
> Array<Feature> listProjectFeatures()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)

const { status, data } = await apiInstance.listProjectFeatures(
    projectId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**Array<Feature>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of features for the project |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjects**
> Array<Project> listProjects()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listProjects();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<Project>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of projects |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listUsers**
> Array<User> listUsers()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listUsers();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<User>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of users |  -  |
|**400** | User not found |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **login**
> LoginResponse login(loginRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    LoginRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let loginRequest: LoginRequest; //

const { status, data } = await apiInstance.login(
    loginRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **loginRequest** | **LoginRequest**|  | |


### Return type

**LoginResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successful login |  -  |
|**401** | Invalid credentials |  -  |
|**403** | 2FA required |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **refreshToken**
> RefreshTokenResponse refreshToken(refreshTokenRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    RefreshTokenRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let refreshTokenRequest: RefreshTokenRequest; //

const { status, data } = await apiInstance.refreshToken(
    refreshTokenRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **refreshTokenRequest** | **RefreshTokenRequest**|  | |


### Return type

**RefreshTokenResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Token refreshed successfully |  -  |
|**401** | Invalid refresh token |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reset2FA**
> TwoFASetupResponse reset2FA(twoFAResetRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TwoFAResetRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let twoFAResetRequest: TwoFAResetRequest; //

const { status, data } = await apiInstance.reset2FA(
    twoFAResetRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **twoFAResetRequest** | **TwoFAResetRequest**|  | |


### Return type

**TwoFASetupResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Secret + QR |  -  |
|**400** | Invalid code |  -  |
|**401** | Unauthorized |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **resetPassword**
> resetPassword(resetPasswordRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    ResetPasswordRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let resetPasswordRequest: ResetPasswordRequest; //

const { status, data } = await apiInstance.resetPassword(
    resetPasswordRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **resetPasswordRequest** | **ResetPasswordRequest**|  | |


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Password reset successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Invalid or expired token |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **sSOCallback**
> LoginResponse sSOCallback(sSOCallbackRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    SSOCallbackRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let sSOCallbackRequest: SSOCallbackRequest; //

const { status, data } = await apiInstance.sSOCallback(
    sSOCallbackRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **sSOCallbackRequest** | **SSOCallbackRequest**|  | |


### Return type

**LoginResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | SSO authentication successful |  -  |
|**400** | Invalid SSO token |  -  |
|**401** | SSO authentication failed |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **sSOInitiate**
> SSOInitiateResponse sSOInitiate()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let provider: string; //Name of the SSO provider to use (default to undefined)

const { status, data } = await apiInstance.sSOInitiate(
    provider
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **provider** | [**string**] | Name of the SSO provider to use | defaults to undefined|


### Return type

**SSOInitiateResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | SSO login URL |  -  |
|**400** | Invalid provider |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **setSuperuserStatus**
> User setSuperuserStatus(setSuperuserStatusRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    SetSuperuserStatusRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let userId: number; // (default to undefined)
let setSuperuserStatusRequest: SetSuperuserStatusRequest; //

const { status, data } = await apiInstance.setSuperuserStatus(
    userId,
    setSuperuserStatusRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **setSuperuserStatusRequest** | **SetSuperuserStatusRequest**|  | |
| **userId** | [**number**] |  | defaults to undefined|


### Return type

**User**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Superuser status updated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not a superuser |  -  |
|**404** | User not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **setUserActiveStatus**
> User setUserActiveStatus(setUserActiveStatusRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    SetUserActiveStatusRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let userId: number; // (default to undefined)
let setUserActiveStatusRequest: SetUserActiveStatusRequest; //

const { status, data } = await apiInstance.setUserActiveStatus(
    userId,
    setUserActiveStatusRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **setUserActiveStatusRequest** | **SetUserActiveStatusRequest**|  | |
| **userId** | [**number**] |  | defaults to undefined|


### Return type

**User**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | User active status updated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not a superuser |  -  |
|**404** | User not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **setup2FA**
> TwoFASetupResponse setup2FA()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.setup2FA();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**TwoFASetupResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Secret + QR-code |  -  |
|**401** | Unauthorized |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **syncLDAPUsers**
> LDAPSyncStartResponse syncLDAPUsers()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.syncLDAPUsers();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**LDAPSyncStartResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**202** | Synchronization started |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**409** | Sync already in progress |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **testLDAPConnection**
> LDAPConnectionTestResponse testLDAPConnection(lDAPConnectionTest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    LDAPConnectionTest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let lDAPConnectionTest: LDAPConnectionTest; //

const { status, data } = await apiInstance.testLDAPConnection(
    lDAPConnectionTest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **lDAPConnectionTest** | **LDAPConnectionTest**|  | |


### Return type

**LDAPConnectionTestResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Connection test result |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **toggleFeature**
> FeatureResponse toggleFeature(toggleFeatureRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    ToggleFeatureRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let toggleFeatureRequest: ToggleFeatureRequest; //

const { status, data } = await apiInstance.toggleFeature(
    featureId,
    toggleFeatureRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **toggleFeatureRequest** | **ToggleFeatureRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**FeatureResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature toggled successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateFeature**
> FeatureDetailsResponse updateFeature(createFeatureRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateFeatureRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let createFeatureRequest: CreateFeatureRequest; //

const { status, data } = await apiInstance.updateFeature(
    featureId,
    createFeatureRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFeatureRequest** | **CreateFeatureRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**FeatureDetailsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Updated feature details with rules and variants |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateFeatureSchedule**
> FeatureScheduleResponse updateFeatureSchedule(updateFeatureScheduleRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateFeatureScheduleRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let scheduleId: string; // (default to undefined)
let updateFeatureScheduleRequest: UpdateFeatureScheduleRequest; //

const { status, data } = await apiInstance.updateFeatureSchedule(
    scheduleId,
    updateFeatureScheduleRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateFeatureScheduleRequest** | **UpdateFeatureScheduleRequest**|  | |
| **scheduleId** | [**string**] |  | defaults to undefined|


### Return type

**FeatureScheduleResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature schedule updated |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Schedule not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateLDAPConfig**
> LDAPConfigResponse updateLDAPConfig(lDAPConfig)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    LDAPConfig
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let lDAPConfig: LDAPConfig; //

const { status, data } = await apiInstance.updateLDAPConfig(
    lDAPConfig
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **lDAPConfig** | **LDAPConfig**|  | |


### Return type

**LDAPConfigResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | LDAP configuration updated |  -  |
|**400** | Invalid configuration |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Superuser access required |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateLicense**
> LicenseStatusResponse updateLicense(updateLicenseRequest)

Updates the system license with a new license key

### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateLicenseRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let updateLicenseRequest: UpdateLicenseRequest; //

const { status, data } = await apiInstance.updateLicense(
    updateLicenseRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateLicenseRequest** | **UpdateLicenseRequest**|  | |


### Return type

**LicenseStatusResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | License updated successfully |  -  |
|**400** | Invalid license key |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateLicenseAcceptance**
> updateLicenseAcceptance(updateLicenseAcceptanceRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateLicenseAcceptanceRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let updateLicenseAcceptanceRequest: UpdateLicenseAcceptanceRequest; //

const { status, data } = await apiInstance.updateLicenseAcceptance(
    updateLicenseAcceptanceRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateLicenseAcceptanceRequest** | **UpdateLicenseAcceptanceRequest**|  | |


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | License acceptance status updated successfully |  -  |
|**401** | Unauthorized |  -  |
|**400** | Bad request |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateProject**
> ProjectResponse updateProject(updateProjectRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateProjectRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let updateProjectRequest: UpdateProjectRequest; //

const { status, data } = await apiInstance.updateProject(
    projectId,
    updateProjectRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateProjectRequest** | **UpdateProjectRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**ProjectResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Project updated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **userChangeMyPassword**
> userChangeMyPassword(changeUserPasswordRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    ChangeUserPasswordRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let changeUserPasswordRequest: ChangeUserPasswordRequest; //

const { status, data } = await apiInstance.userChangeMyPassword(
    changeUserPasswordRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **changeUserPasswordRequest** | **ChangeUserPasswordRequest**|  | |


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Password changed successfully |  -  |
|**401** | Unauthorized |  -  |
|**400** | Bad request |  -  |
|**403** | External user can\&#39;t change password |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **verify2FA**
> TwoFAVerifyResponse verify2FA(twoFAVerifyRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TwoFAVerifyRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let twoFAVerifyRequest: TwoFAVerifyRequest; //

const { status, data } = await apiInstance.verify2FA(
    twoFAVerifyRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **twoFAVerifyRequest** | **TwoFAVerifyRequest**|  | |


### Return type

**TwoFAVerifyResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success, returns access/refresh tokens |  -  |
|**400** | Invalid code |  -  |
|**401** | Unauthorized |  -  |
|**429** | Too many requests |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

