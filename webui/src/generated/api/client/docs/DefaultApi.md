# DefaultApi

All URIs are relative to *http://localhost*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addFeatureTag**](#addfeaturetag) | **POST** /api/v1/features/{feature_id}/tags | Add tag to feature|
|[**addProject**](#addproject) | **POST** /api/v1/projects/add | Add new project|
|[**approvePendingChange**](#approvependingchange) | **POST** /api/v1/pending_changes/{pending_change_id}/approve | Approve a pending change|
|[**archiveProject**](#archiveproject) | **DELETE** /api/v1/projects/{project_id} | Archive a project|
|[**cancelLDAPSync**](#cancelldapsync) | **DELETE** /api/v1/ldap/sync/cancel | Cancel ongoing synchronization|
|[**cancelPendingChange**](#cancelpendingchange) | **POST** /api/v1/pending_changes/{pending_change_id}/cancel | Cancel a pending change|
|[**confirm2FA**](#confirm2fa) | **POST** /api/v1/users/me/2fa/confirm | Approve enable 2FA (code from app)|
|[**consumeSAMLAssertion**](#consumesamlassertion) | **POST** /api/v1/saml/acs | Assertion Consumer Service (ACS) endpoint|
|[**createCategory**](#createcategory) | **POST** /api/v1/categories | Create new category|
|[**createEnvironment**](#createenvironment) | **POST** /api/v1/projects/{project_id}/environments | Create environment|
|[**createFeatureAlgorithm**](#createfeaturealgorithm) | **POST** /api/v1/features/{feature_id}/algorithms/{environment_id} | Create or attach algorithm to feature in environment|
|[**createFeatureFlagVariant**](#createfeatureflagvariant) | **POST** /api/v1/features/{feature_id}/variants | Create flag variant for feature|
|[**createFeatureRule**](#createfeaturerule) | **POST** /api/v1/features/{feature_id}/rules | Create rule for feature|
|[**createFeatureSchedule**](#createfeatureschedule) | **POST** /api/v1/features/{feature_id}/schedules | Create schedule for feature|
|[**createNotificationSetting**](#createnotificationsetting) | **POST** /api/v1/projects/{project_id}/env/{environment_key}/notification-settings | Create a new notification setting|
|[**createProjectFeature**](#createprojectfeature) | **POST** /api/v1/projects/{project_id}/features | Create feature for project|
|[**createProjectMembership**](#createprojectmembership) | **POST** /api/v1/projects/{project_id}/memberships | Add membership to project|
|[**createProjectSegment**](#createprojectsegment) | **POST** /api/v1/projects/{project_id}/segments | Create segment for project|
|[**createProjectTag**](#createprojecttag) | **POST** /api/v1/projects/{project_id}/tags | Create new tag for project|
|[**createRuleAttribute**](#createruleattribute) | **POST** /api/v1/rule_attributes | Create rule attribute|
|[**createUser**](#createuser) | **POST** /api/v1/users | Create a new user (superuser only)|
|[**deleteCategory**](#deletecategory) | **DELETE** /api/v1/categories/{category_id} | Delete category|
|[**deleteEnvironment**](#deleteenvironment) | **DELETE** /api/v1/environments/{environment_id} | Delete environment|
|[**deleteFeature**](#deletefeature) | **DELETE** /api/v1/features/{feature_id} | Delete feature|
|[**deleteFeatureAlgorithm**](#deletefeaturealgorithm) | **DELETE** /api/v1/features/{feature_id}/algorithms/{environment_id} | Delete feature algorithm from feature|
|[**deleteFeatureSchedule**](#deletefeatureschedule) | **DELETE** /api/v1/feature-schedules/{schedule_id} | Delete feature schedule by ID|
|[**deleteLDAPConfig**](#deleteldapconfig) | **DELETE** /api/v1/ldap/config | Delete LDAP configuration|
|[**deleteNotificationSetting**](#deletenotificationsetting) | **DELETE** /api/v1/projects/{project_id}/env/{environment_key}/notification-settings/{setting_id} | Delete a notification setting|
|[**deleteProjectMembership**](#deleteprojectmembership) | **DELETE** /api/v1/projects/{project_id}/memberships/{membership_id} | Delete membership|
|[**deleteProjectTag**](#deleteprojecttag) | **DELETE** /api/v1/projects/{project_id}/tags/{tag_id} | Delete tag|
|[**deleteRuleAttribute**](#deleteruleattribute) | **DELETE** /api/v1/rule_attributes/{name} | Delete rule attribute|
|[**deleteSegment**](#deletesegment) | **DELETE** /api/v1/segments/{segment_id} | Delete segment|
|[**deleteUser**](#deleteuser) | **DELETE** /api/v1/users/{user_id} | Delete a user (superuser only, cannot delete superusers)|
|[**disable2FA**](#disable2fa) | **POST** /api/v1/users/me/2fa/disable | Disable 2FA (using email-confirmation)|
|[**forgotPassword**](#forgotpassword) | **POST** /api/v1/auth/forgot-password | Request a password reset|
|[**getAuditLogEntry**](#getauditlogentry) | **GET** /api/v1/audit/{id} | Get audit log entry by ID|
|[**getCategory**](#getcategory) | **GET** /api/v1/categories/{category_id} | Get category details|
|[**getCurrentUser**](#getcurrentuser) | **GET** /api/v1/users/me | Get current user information|
|[**getDashboardOverview**](#getdashboardoverview) | **GET** /api/v1/dashboard/overview | Project Dashboard overview|
|[**getEnvironment**](#getenvironment) | **GET** /api/v1/environments/{environment_id} | Get environment|
|[**getFeature**](#getfeature) | **GET** /api/v1/features/{feature_id} | Get feature with rules and variants|
|[**getFeatureAlgorithm**](#getfeaturealgorithm) | **GET** /api/v1/features/{feature_id}/algorithms/{environment_id} | Get algorithm configuration for a feature in environment|
|[**getFeatureSchedule**](#getfeatureschedule) | **GET** /api/v1/feature-schedules/{schedule_id} | Get feature schedule by ID|
|[**getFeatureTimeline**](#getfeaturetimeline) | **GET** /api/v1/features/{feature_id}/timeline | Get feature timeline within period|
|[**getLDAPConfig**](#getldapconfig) | **GET** /api/v1/ldap/config | Get LDAP configuration|
|[**getLDAPStatistics**](#getldapstatistics) | **GET** /api/v1/ldap/statistics | Get LDAP statistics|
|[**getLDAPSyncLogDetails**](#getldapsynclogdetails) | **GET** /api/v1/ldap/sync/logs/{id} | Get synchronization log details|
|[**getLDAPSyncLogs**](#getldapsynclogs) | **GET** /api/v1/ldap/sync/logs | Get synchronization logs|
|[**getLDAPSyncProgress**](#getldapsyncprogress) | **GET** /api/v1/ldap/sync/progress | Get synchronization progress|
|[**getLDAPSyncStatus**](#getldapsyncstatus) | **GET** /api/v1/ldap/sync/status | Get synchronization status|
|[**getNotificationSetting**](#getnotificationsetting) | **GET** /api/v1/projects/{project_id}/env/{environment_key}/notification-settings/{setting_id} | Get a specific notification setting|
|[**getPendingChange**](#getpendingchange) | **GET** /api/v1/pending_changes/{pending_change_id} | Get pending change by ID|
|[**getProject**](#getproject) | **GET** /api/v1/projects/{project_id} | Get project details|
|[**getProjectMembership**](#getprojectmembership) | **GET** /api/v1/projects/{project_id}/memberships/{membership_id} | Get membership|
|[**getProjectSetting**](#getprojectsetting) | **GET** /api/v1/projects/{project_id}/settings/{setting_name} | Get project setting by name|
|[**getProjectTag**](#getprojecttag) | **GET** /api/v1/projects/{project_id}/tags/{tag_id} | Get tag details|
|[**getRolePermissions**](#getrolepermissions) | **GET** /api/v1/roles/{role_id}/permissions | Get permissions for a role|
|[**getSAMLMetadata**](#getsamlmetadata) | **GET** /api/v1/saml/metadata | Get SAML metadata|
|[**getSSOProviders**](#getssoproviders) | **GET** /api/v1/auth/sso/providers | Get available SSO providers|
|[**getSegment**](#getsegment) | **GET** /api/v1/segments/{segment_id} | Get segment by ID|
|[**getUnreadNotificationsCount**](#getunreadnotificationscount) | **GET** /api/v1/notifications/unread-count | Get unread notifications count|
|[**getUserNotifications**](#getusernotifications) | **GET** /api/v1/notifications | Get user notifications|
|[**initiateTOTPApproval**](#initiatetotpapproval) | **POST** /api/v1/pending_changes/{pending_change_id}/initiate-totp | Initiate TOTP approval session|
|[**listAlgorithms**](#listalgorithms) | **GET** /api/v1/algorithms | List of algorithms|
|[**listAllFeatureSchedules**](#listallfeatureschedules) | **GET** /api/v1/projects/{project_id}/env/{environment_key}/feature-schedules | List all feature schedules for project|
|[**listCategories**](#listcategories) | **GET** /api/v1/categories | Get categories list|
|[**listFeatureAlgorithms**](#listfeaturealgorithms) | **GET** /api/v1/projects/{project_id}/feature-algorithms | List feature algorithms for a feature|
|[**listFeatureFlagVariants**](#listfeatureflagvariants) | **GET** /api/v1/features/{feature_id}/variants | List flag variants for feature|
|[**listFeatureRules**](#listfeaturerules) | **GET** /api/v1/features/{feature_id}/rules | List rules for feature|
|[**listFeatureSchedules**](#listfeatureschedules) | **GET** /api/v1/features/{feature_id}/schedules | List schedules for feature|
|[**listFeatureTags**](#listfeaturetags) | **GET** /api/v1/features/{feature_id}/tags | List feature tags|
|[**listNotificationSettings**](#listnotificationsettings) | **GET** /api/v1/projects/{project_id}/env/{environment_key}/notification-settings | List all notification settings for a project|
|[**listPendingChanges**](#listpendingchanges) | **GET** /api/v1/pending_changes | List pending changes|
|[**listPermissions**](#listpermissions) | **GET** /api/v1/permissions | List all permissions|
|[**listProjectAuditLogs**](#listprojectauditlogs) | **GET** /api/v1/projects/{project_id}/audit | List audit log entries for project|
|[**listProjectChanges**](#listprojectchanges) | **GET** /api/v1/projects/{project_id}/changes | Get project changes history|
|[**listProjectEnvironments**](#listprojectenvironments) | **GET** /api/v1/projects/{project_id}/environments | List project environments|
|[**listProjectFeatures**](#listprojectfeatures) | **GET** /api/v1/projects/{project_id}/features | List features for project|
|[**listProjectMemberships**](#listprojectmemberships) | **GET** /api/v1/projects/{project_id}/memberships | List memberships for project|
|[**listProjectSegments**](#listprojectsegments) | **GET** /api/v1/projects/{project_id}/segments | List segments for project|
|[**listProjectSettings**](#listprojectsettings) | **GET** /api/v1/projects/{project_id}/settings | List project settings|
|[**listProjectTags**](#listprojecttags) | **GET** /api/v1/projects/{project_id}/tags | Get tags list for project|
|[**listProjects**](#listprojects) | **GET** /api/v1/projects | Get projects list|
|[**listRolePermissions**](#listrolepermissions) | **GET** /api/v1/roles/permissions | List permissions for all roles|
|[**listRoles**](#listroles) | **GET** /api/v1/roles | List all roles|
|[**listRuleAttributes**](#listruleattributes) | **GET** /api/v1/rule_attributes | List of rule attributes|
|[**listSegmentDesyncFeatureIDs**](#listsegmentdesyncfeatureids) | **GET** /api/v1/segments/{segment_id}/desync-features | Get desync feature IDs by segment ID|
|[**listUsers**](#listusers) | **GET** /api/v1/users | List all users (superuser only)|
|[**login**](#login) | **POST** /api/v1/auth/login | Authenticate user and get access token|
|[**markAllNotificationsAsRead**](#markallnotificationsasread) | **PUT** /api/v1/notifications/read-all | Mark all notifications as read|
|[**markNotificationAsRead**](#marknotificationasread) | **PUT** /api/v1/notifications/{notification_id}/read | Mark notification as read|
|[**refreshToken**](#refreshtoken) | **POST** /api/v1/auth/refresh | Refresh access token|
|[**rejectPendingChange**](#rejectpendingchange) | **POST** /api/v1/pending_changes/{pending_change_id}/reject | Reject a pending change|
|[**removeFeatureTag**](#removefeaturetag) | **DELETE** /api/v1/features/{feature_id}/tags | Remove tag from feature|
|[**reset2FA**](#reset2fa) | **POST** /api/v1/users/me/2fa/reset | Reset/generate secret 2FA (using email-confirmation)|
|[**resetPassword**](#resetpassword) | **POST** /api/v1/auth/reset-password | Reset password using token|
|[**sSOCallback**](#ssocallback) | **POST** /api/v1/auth/sso/callback | Handle SSO callback from Keycloak|
|[**sSOInitiate**](#ssoinitiate) | **GET** /api/v1/auth/sso/initiate | Initiate SSO login flow|
|[**sendTestNotification**](#sendtestnotification) | **POST** /api/v1/projects/{project_id}/env/{environment_key}/notification-settings/{setting_id}/test | Send test notification|
|[**setSuperuserStatus**](#setsuperuserstatus) | **PUT** /api/v1/users/{user_id}/superuser | Set or unset superuser status (superuser only, cannot modify admin user)|
|[**setUserActiveStatus**](#setuseractivestatus) | **PUT** /api/v1/users/{user_id}/active | Set or unset user active status (superuser only)|
|[**setup2FA**](#setup2fa) | **POST** /api/v1/users/me/2fa/setup | Begin setup 2FA (generate secret and QR-code)|
|[**syncCustomizedFeatureRule**](#synccustomizedfeaturerule) | **PUT** /api/v1/features/{feature_id}/rules/{rule_id}/sync | Synchronize customized feature rule|
|[**syncLDAPUsers**](#syncldapusers) | **POST** /api/v1/ldap/sync/users | Start user synchronization|
|[**testFeatureTimeline**](#testfeaturetimeline) | **POST** /api/v1/features/{feature_id}/timeline/test | Test feature timeline with mock schedules|
|[**testLDAPConnection**](#testldapconnection) | **POST** /api/v1/ldap/test-connection | Test LDAP connection|
|[**toggleFeature**](#togglefeature) | **PUT** /api/v1/features/{feature_id}/toggle | Toggle feature enabled state|
|[**updateCategory**](#updatecategory) | **PUT** /api/v1/categories/{category_id} | Update category|
|[**updateEnvironment**](#updateenvironment) | **PUT** /api/v1/environments/{environment_id} | Update environment|
|[**updateFeature**](#updatefeature) | **PUT** /api/v1/features/{feature_id} | Update feature with rules and variants|
|[**updateFeatureAlgorithm**](#updatefeaturealgorithm) | **PATCH** /api/v1/features/{feature_id}/algorithms/{environment_id} | Update feature algorithm configuration|
|[**updateFeatureSchedule**](#updatefeatureschedule) | **PUT** /api/v1/feature-schedules/{schedule_id} | Update feature schedule by ID|
|[**updateLDAPConfig**](#updateldapconfig) | **POST** /api/v1/ldap/config | Create or update LDAP configuration|
|[**updateLicenseAcceptance**](#updatelicenseacceptance) | **PUT** /api/v1/users/me/license-acceptance | Update license acceptance status|
|[**updateNotificationSetting**](#updatenotificationsetting) | **PUT** /api/v1/projects/{project_id}/env/{environment_key}/notification-settings/{setting_id} | Update a notification setting|
|[**updateProject**](#updateproject) | **PUT** /api/v1/projects/{project_id} | Update project name and description|
|[**updateProjectMembership**](#updateprojectmembership) | **PUT** /api/v1/projects/{project_id}/memberships/{membership_id} | Update membership|
|[**updateProjectSetting**](#updateprojectsetting) | **PUT** /api/v1/projects/{project_id}/settings/{setting_name} | Update project setting|
|[**updateProjectTag**](#updateprojecttag) | **PUT** /api/v1/projects/{project_id}/tags/{tag_id} | Update tag|
|[**updateSegment**](#updatesegment) | **PUT** /api/v1/segments/{segment_id} | Update segment|
|[**userChangeMyPassword**](#userchangemypassword) | **POST** /api/v1/users/me/change-password | Change my password|
|[**verify2FA**](#verify2fa) | **POST** /api/v1/auth/2fa/verify | Verify 2FA-code on login|

# **addFeatureTag**
> addFeatureTag(addFeatureTagRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    AddFeatureTagRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let addFeatureTagRequest: AddFeatureTagRequest; //

const { status, data } = await apiInstance.addFeatureTag(
    featureId,
    addFeatureTagRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **addFeatureTagRequest** | **AddFeatureTagRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|


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
|**201** | Tag added to feature |  -  |
|**202** | Change is pending approval (for guarded features) |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature or tag not found |  -  |
|**409** | Conflict - tag already associated or change cannot be applied due to pending change or lock |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

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
|**403** | Permission denied |  -  |
|**409** | Project already exists |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **approvePendingChange**
> SuccessResponse approvePendingChange(approvePendingChangeRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    ApprovePendingChangeRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let pendingChangeId: string; // (default to undefined)
let approvePendingChangeRequest: ApprovePendingChangeRequest; //

const { status, data } = await apiInstance.approvePendingChange(
    pendingChangeId,
    approvePendingChangeRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **approvePendingChangeRequest** | **ApprovePendingChangeRequest**|  | |
| **pendingChangeId** | [**string**] |  | defaults to undefined|


### Return type

**SuccessResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pending change approved successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized or invalid credentials |  -  |
|**403** | Permission denied |  -  |
|**404** | Pending change not found |  -  |
|**409** | Conflict - pending change is not in pending status |  -  |
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

# **cancelPendingChange**
> SuccessResponse cancelPendingChange(cancelPendingChangeRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CancelPendingChangeRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let pendingChangeId: string; // (default to undefined)
let cancelPendingChangeRequest: CancelPendingChangeRequest; //

const { status, data } = await apiInstance.cancelPendingChange(
    pendingChangeId,
    cancelPendingChangeRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **cancelPendingChangeRequest** | **CancelPendingChangeRequest**|  | |
| **pendingChangeId** | [**string**] |  | defaults to undefined|


### Return type

**SuccessResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pending change cancelled successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Pending change not found |  -  |
|**409** | Conflict - pending change is not in pending status |  -  |
|**500** | Internal server error |  -  |
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

# **createCategory**
> CategoryResponse createCategory(createCategoryRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateCategoryRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let createCategoryRequest: CreateCategoryRequest; //

const { status, data } = await apiInstance.createCategory(
    createCategoryRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createCategoryRequest** | **CreateCategoryRequest**|  | |


### Return type

**CategoryResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Category created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**409** | Category already exists |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createEnvironment**
> EnvironmentResponse createEnvironment(createEnvironmentRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateEnvironmentRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let createEnvironmentRequest: CreateEnvironmentRequest; //

const { status, data } = await apiInstance.createEnvironment(
    projectId,
    createEnvironmentRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createEnvironmentRequest** | **CreateEnvironmentRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**EnvironmentResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Environment created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**409** | Environment already exists |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createFeatureAlgorithm**
> createFeatureAlgorithm(createFeatureAlgorithmRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateFeatureAlgorithmRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let environmentId: number; // (default to undefined)
let createFeatureAlgorithmRequest: CreateFeatureAlgorithmRequest; //

const { status, data } = await apiInstance.createFeatureAlgorithm(
    featureId,
    environmentId,
    createFeatureAlgorithmRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFeatureAlgorithmRequest** | **CreateFeatureAlgorithmRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentId** | [**number**] |  | defaults to undefined|


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
|**201** | Feature algorithm created |  -  |
|**202** | Change is pending approval (for guarded features) |  -  |
|**400** | Invalid request data |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**409** | Algorithm already exists for this feature/environment |  -  |
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
let environmentKey: string; // (default to undefined)
let createFlagVariantRequest: CreateFlagVariantRequest; //

const { status, data } = await apiInstance.createFeatureFlagVariant(
    featureId,
    environmentKey,
    createFlagVariantRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFlagVariantRequest** | **CreateFlagVariantRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
|**202** | Change is pending approval (for guarded features) |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
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
let environmentKey: string; // (default to undefined)
let createRuleRequest: CreateRuleRequest; //

const { status, data } = await apiInstance.createFeatureRule(
    featureId,
    environmentKey,
    createRuleRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createRuleRequest** | **CreateRuleRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
|**202** | Change is pending approval (for guarded features) |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature or related resource not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
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
let environmentKey: string; // (default to undefined)
let createFeatureScheduleRequest: CreateFeatureScheduleRequest; //

const { status, data } = await apiInstance.createFeatureSchedule(
    featureId,
    environmentKey,
    createFeatureScheduleRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFeatureScheduleRequest** | **CreateFeatureScheduleRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
|**202** | Change is pending approval (for guarded features) |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature or related resource not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createNotificationSetting**
> NotificationSetting createNotificationSetting(createNotificationSettingRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateNotificationSettingRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)
let createNotificationSettingRequest: CreateNotificationSettingRequest; //

const { status, data } = await apiInstance.createNotificationSetting(
    projectId,
    environmentKey,
    createNotificationSettingRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createNotificationSettingRequest** | **CreateNotificationSettingRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


### Return type

**NotificationSetting**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Notification setting created successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not authorized to modify this project |  -  |
|**404** | Project not found |  -  |
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

# **createProjectMembership**
> Membership createProjectMembership(createMembershipRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateMembershipRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let createMembershipRequest: CreateMembershipRequest; //

const { status, data } = await apiInstance.createProjectMembership(
    projectId,
    createMembershipRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createMembershipRequest** | **CreateMembershipRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**Membership**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Membership created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createProjectSegment**
> SegmentResponse createProjectSegment(createSegmentRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateSegmentRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let createSegmentRequest: CreateSegmentRequest; //

const { status, data } = await apiInstance.createProjectSegment(
    projectId,
    createSegmentRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createSegmentRequest** | **CreateSegmentRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**SegmentResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Segment created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createProjectTag**
> ProjectTagResponse createProjectTag(createProjectTagRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateProjectTagRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let createProjectTagRequest: CreateProjectTagRequest; //

const { status, data } = await apiInstance.createProjectTag(
    projectId,
    createProjectTagRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createProjectTagRequest** | **CreateProjectTagRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**ProjectTagResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Tag created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**409** | Tag already exists |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createRuleAttribute**
> createRuleAttribute(createRuleAttributeRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    CreateRuleAttributeRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let createRuleAttributeRequest: CreateRuleAttributeRequest; //

const { status, data } = await apiInstance.createRuleAttribute(
    createRuleAttributeRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createRuleAttributeRequest** | **CreateRuleAttributeRequest**|  | |


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
|**204** | Attribute created |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**409** | Attribute already exists |  -  |
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

# **deleteCategory**
> deleteCategory()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let categoryId: string; // (default to undefined)

const { status, data } = await apiInstance.deleteCategory(
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] |  | defaults to undefined|


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
|**204** | Category deleted successfully |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Category not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteEnvironment**
> deleteEnvironment()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let environmentId: number; // (default to undefined)

const { status, data } = await apiInstance.deleteEnvironment(
    environmentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **environmentId** | [**number**] |  | defaults to undefined|


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
|**204** | Environment deleted |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Environment not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteFeature**
> PendingChangeResponse deleteFeature()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.deleteFeature(
    featureId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


### Return type

**PendingChangeResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Feature deleted successfully |  -  |
|**202** | Deletion is pending approval (for guarded features) |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteFeatureAlgorithm**
> PendingChangeResponse deleteFeatureAlgorithm()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let environmentId: number; // (default to undefined)

const { status, data } = await apiInstance.deleteFeatureAlgorithm(
    featureId,
    environmentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentId** | [**number**] |  | defaults to undefined|


### Return type

**PendingChangeResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**202** | Change is pending approval (for guarded features) |  -  |
|**204** | Feature algorithm deleted |  -  |
|**404** | Not found |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**409** | Conflict |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteFeatureSchedule**
> PendingChangeResponse deleteFeatureSchedule()


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

**PendingChangeResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Feature schedule deleted |  -  |
|**202** | Change is pending approval (for guarded features) |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Schedule not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
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

# **deleteNotificationSetting**
> deleteNotificationSetting()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)
let settingId: number; // (default to undefined)

const { status, data } = await apiInstance.deleteNotificationSetting(
    projectId,
    environmentKey,
    settingId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|
| **settingId** | [**number**] |  | defaults to undefined|


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
|**204** | Notification setting deleted successfully |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not authorized to delete this notification setting |  -  |
|**404** | Notification setting not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteProjectMembership**
> deleteProjectMembership()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let membershipId: string; // (default to undefined)

const { status, data } = await apiInstance.deleteProjectMembership(
    projectId,
    membershipId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **membershipId** | [**string**] |  | defaults to undefined|


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
|**204** | Membership deleted |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Membership not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteProjectTag**
> deleteProjectTag()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let tagId: string; // (default to undefined)

const { status, data } = await apiInstance.deleteProjectTag(
    projectId,
    tagId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **tagId** | [**string**] |  | defaults to undefined|


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
|**204** | Tag deleted successfully |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Tag not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRuleAttribute**
> deleteRuleAttribute()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let name: string; // (default to undefined)

const { status, data } = await apiInstance.deleteRuleAttribute(
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **name** | [**string**] |  | defaults to undefined|


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
|**204** | Attribute deleted successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Attribute not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteSegment**
> deleteSegment()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let segmentId: string; // (default to undefined)

const { status, data } = await apiInstance.deleteSegment(
    segmentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **segmentId** | [**string**] |  | defaults to undefined|


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
|**204** | Segment deleted successfully |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Segment not found |  -  |
|**500** | Internal server error |  -  |
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

# **getAuditLogEntry**
> AuditLog getAuditLogEntry()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let id: number; // (default to undefined)

const { status, data } = await apiInstance.getAuditLogEntry(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] |  | defaults to undefined|


### Return type

**AuditLog**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Audit log entry |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Audit log entry not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCategory**
> CategoryResponse getCategory()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let categoryId: string; // (default to undefined)

const { status, data } = await apiInstance.getCategory(
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **categoryId** | [**string**] |  | defaults to undefined|


### Return type

**CategoryResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Category details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Category not found |  -  |
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

# **getDashboardOverview**
> DashboardOverviewResponse getDashboardOverview()

Returns aggregated dashboard data for a project: - project health - category health - feature activity (upcoming & recent) - recent activity (batched by request_id) - risky features - pending summary 

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let environmentKey: string; //Environment key (prod/stage/dev) (default to undefined)
let projectId: string; //Optional project ID to filter results (optional) (default to undefined)
let limit: number; //Limit for recent activity entries (optional) (default to 20)

const { status, data } = await apiInstance.getDashboardOverview(
    environmentKey,
    projectId,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **environmentKey** | [**string**] | Environment key (prod/stage/dev) | defaults to undefined|
| **projectId** | [**string**] | Optional project ID to filter results | (optional) defaults to undefined|
| **limit** | [**number**] | Limit for recent activity entries | (optional) defaults to 20|


### Return type

**DashboardOverviewResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Dashboard data |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getEnvironment**
> EnvironmentResponse getEnvironment()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let environmentId: number; // (default to undefined)

const { status, data } = await apiInstance.getEnvironment(
    environmentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **environmentId** | [**number**] |  | defaults to undefined|


### Return type

**EnvironmentResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Environment details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Environment not found |  -  |
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
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.getFeature(
    featureId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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

# **getFeatureAlgorithm**
> FeatureAlgorithm getFeatureAlgorithm()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let environmentId: number; // (default to undefined)

const { status, data } = await apiInstance.getFeatureAlgorithm(
    featureId,
    environmentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentId** | [**number**] |  | defaults to undefined|


### Return type

**FeatureAlgorithm**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature algorithm configuration |  -  |
|**404** | Algorithm not found for feature/environment |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
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

# **getFeatureTimeline**
> FeatureTimelineResponse getFeatureTimeline()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let environmentKey: string; //Target environment key (e.g., dev, stage, prod) (default to undefined)
let from: string; //Start of the period (inclusive) (default to undefined)
let to: string; //End of the period (exclusive) (default to undefined)
let location: string; //Browser\'s location string (default to undefined)

const { status, data } = await apiInstance.getFeatureTimeline(
    featureId,
    environmentKey,
    from,
    to,
    location
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] | Target environment key (e.g., dev, stage, prod) | defaults to undefined|
| **from** | [**string**] | Start of the period (inclusive) | defaults to undefined|
| **to** | [**string**] | End of the period (exclusive) | defaults to undefined|
| **location** | [**string**] | Browser\&#39;s location string | defaults to undefined|


### Return type

**FeatureTimelineResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature timeline |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
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

# **getNotificationSetting**
> NotificationSetting getNotificationSetting()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)
let settingId: number; // (default to undefined)

const { status, data } = await apiInstance.getNotificationSetting(
    projectId,
    environmentKey,
    settingId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|
| **settingId** | [**number**] |  | defaults to undefined|


### Return type

**NotificationSetting**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Notification setting details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not authorized to access this notification setting |  -  |
|**404** | Notification setting not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPendingChange**
> PendingChangeResponse getPendingChange()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let pendingChangeId: string; // (default to undefined)

const { status, data } = await apiInstance.getPendingChange(
    pendingChangeId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **pendingChangeId** | [**string**] |  | defaults to undefined|


### Return type

**PendingChangeResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pending change details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Pending change not found |  -  |
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

# **getProjectMembership**
> Membership getProjectMembership()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let membershipId: string; // (default to undefined)

const { status, data } = await apiInstance.getProjectMembership(
    projectId,
    membershipId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **membershipId** | [**string**] |  | defaults to undefined|


### Return type

**Membership**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Membership |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Membership not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getProjectSetting**
> ProjectSettingResponse getProjectSetting()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let settingName: string; // (default to undefined)

const { status, data } = await apiInstance.getProjectSetting(
    projectId,
    settingName
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **settingName** | [**string**] |  | defaults to undefined|


### Return type

**ProjectSettingResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Project setting details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Setting not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getProjectTag**
> ProjectTagResponse getProjectTag()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let tagId: string; // (default to undefined)

const { status, data } = await apiInstance.getProjectTag(
    projectId,
    tagId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **tagId** | [**string**] |  | defaults to undefined|


### Return type

**ProjectTagResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Tag details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Tag not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRolePermissions**
> Array<Permission> getRolePermissions()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let roleId: string; // (default to undefined)

const { status, data } = await apiInstance.getRolePermissions(
    roleId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **roleId** | [**string**] |  | defaults to undefined|


### Return type

**Array<Permission>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of permissions for role |  -  |
|**401** | Unauthorized |  -  |
|**404** | Role not found |  -  |
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

# **getSegment**
> SegmentResponse getSegment()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let segmentId: string; // (default to undefined)

const { status, data } = await apiInstance.getSegment(
    segmentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **segmentId** | [**string**] |  | defaults to undefined|


### Return type

**SegmentResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Segment details |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Segment not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getUnreadNotificationsCount**
> UnreadCountResponse getUnreadNotificationsCount()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.getUnreadNotificationsCount();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**UnreadCountResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Unread notifications count |  -  |
|**401** | Unauthorized |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getUserNotifications**
> UserNotificationsResponse getUserNotifications()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let limit: number; // (optional) (default to 50)
let offset: number; // (optional) (default to 0)

const { status, data } = await apiInstance.getUserNotifications(
    limit,
    offset
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **limit** | [**number**] |  | (optional) defaults to 50|
| **offset** | [**number**] |  | (optional) defaults to 0|


### Return type

**UserNotificationsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of user notifications |  -  |
|**401** | Unauthorized |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **initiateTOTPApproval**
> InitiateTOTPApprovalResponse initiateTOTPApproval(initiateTOTPApprovalRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    InitiateTOTPApprovalRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let pendingChangeId: string; // (default to undefined)
let initiateTOTPApprovalRequest: InitiateTOTPApprovalRequest; //

const { status, data } = await apiInstance.initiateTOTPApproval(
    pendingChangeId,
    initiateTOTPApprovalRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **initiateTOTPApprovalRequest** | **InitiateTOTPApprovalRequest**|  | |
| **pendingChangeId** | [**string**] |  | defaults to undefined|


### Return type

**InitiateTOTPApprovalResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | TOTP approval session initiated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Pending change not found |  -  |
|**409** | Conflict - pending change is not in pending status |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listAlgorithms**
> ListAlgorithmsResponse listAlgorithms()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listAlgorithms();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ListAlgorithmsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of tags for project |  -  |
|**401** | Unauthorized |  -  |
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

let projectId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.listAllFeatureSchedules(
    projectId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
|**200** | List of feature schedules for project |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listCategories**
> Array<Category> listCategories()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listCategories();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<Category>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of categories |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listFeatureAlgorithms**
> ListFeatureAlgorithmsResponse listFeatureAlgorithms()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; //Filter by environment key (default to undefined)

const { status, data } = await apiInstance.listFeatureAlgorithms(
    projectId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] | Filter by environment key | defaults to undefined|


### Return type

**ListFeatureAlgorithmsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of feature algorithms |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
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
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.listFeatureFlagVariants(
    featureId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.listFeatureRules(
    featureId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.listFeatureSchedules(
    featureId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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

# **listFeatureTags**
> Array<ProjectTag> listFeatureTags()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)

const { status, data } = await apiInstance.listFeatureTags(
    featureId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|


### Return type

**Array<ProjectTag>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of feature tags |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listNotificationSettings**
> ListNotificationSettingsResponse listNotificationSettings()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.listNotificationSettings(
    projectId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


### Return type

**ListNotificationSettingsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of notification settings |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not authorized to access this project |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listPendingChanges**
> PendingChangesListResponse listPendingChanges()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let environmentId: number; // (optional) (default to undefined)
let environmentKey: string; //Target environment key (e.g., dev, stage, prod). If provided, takes precedence over environment_id. (optional) (default to undefined)
let projectId: string; // (optional) (default to undefined)
let status: 'pending' | 'approved' | 'rejected' | 'cancelled'; // (optional) (default to undefined)
let userId: number; // (optional) (default to undefined)
let page: number; // (optional) (default to 1)
let perPage: number; // (optional) (default to 20)
let sortBy: 'created_at' | 'status' | 'requested_by'; // (optional) (default to 'created_at')
let sortDesc: boolean; // (optional) (default to true)

const { status, data } = await apiInstance.listPendingChanges(
    environmentId,
    environmentKey,
    projectId,
    status,
    userId,
    page,
    perPage,
    sortBy,
    sortDesc
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **environmentId** | [**number**] |  | (optional) defaults to undefined|
| **environmentKey** | [**string**] | Target environment key (e.g., dev, stage, prod). If provided, takes precedence over environment_id. | (optional) defaults to undefined|
| **projectId** | [**string**] |  | (optional) defaults to undefined|
| **status** | [**&#39;pending&#39; | &#39;approved&#39; | &#39;rejected&#39; | &#39;cancelled&#39;**]**Array<&#39;pending&#39; &#124; &#39;approved&#39; &#124; &#39;rejected&#39; &#124; &#39;cancelled&#39;>** |  | (optional) defaults to undefined|
| **userId** | [**number**] |  | (optional) defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to 1|
| **perPage** | [**number**] |  | (optional) defaults to 20|
| **sortBy** | [**&#39;created_at&#39; | &#39;status&#39; | &#39;requested_by&#39;**]**Array<&#39;created_at&#39; &#124; &#39;status&#39; &#124; &#39;requested_by&#39;>** |  | (optional) defaults to 'created_at'|
| **sortDesc** | [**boolean**] |  | (optional) defaults to true|


### Return type

**PendingChangesListResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of pending changes |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listPermissions**
> Array<Permission> listPermissions()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listPermissions();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<Permission>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of all permissions |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectAuditLogs**
> ListProjectAuditLogs200Response listProjectAuditLogs()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; //Filter by environment (optional) (default to undefined)
let entity: string; //Filter by entity type (e.g. \"feature\", \"segment\") (optional) (default to undefined)
let entityId: string; //Filter by specific entity (optional) (default to undefined)
let actor: string; //Filter by actor username (optional) (default to undefined)
let from: string; //Start of time range (optional) (default to undefined)
let to: string; //End of time range (optional) (default to undefined)
let sortBy: 'environment_key' | 'entity' | 'entity_id' | 'actor' | 'action' | 'username' | 'created_at'; //Sort by field (optional) (default to undefined)
let sortOrder: SortOrder; //Sort order (optional) (default to undefined)
let page: number; //Page number (starts from 1) (optional) (default to 1)
let perPage: number; //Items per page (optional) (default to 20)

const { status, data } = await apiInstance.listProjectAuditLogs(
    projectId,
    environmentKey,
    entity,
    entityId,
    actor,
    from,
    to,
    sortBy,
    sortOrder,
    page,
    perPage
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] | Filter by environment | (optional) defaults to undefined|
| **entity** | [**string**] | Filter by entity type (e.g. \&quot;feature\&quot;, \&quot;segment\&quot;) | (optional) defaults to undefined|
| **entityId** | [**string**] | Filter by specific entity | (optional) defaults to undefined|
| **actor** | [**string**] | Filter by actor username | (optional) defaults to undefined|
| **from** | [**string**] | Start of time range | (optional) defaults to undefined|
| **to** | [**string**] | End of time range | (optional) defaults to undefined|
| **sortBy** | [**&#39;environment_key&#39; | &#39;entity&#39; | &#39;entity_id&#39; | &#39;actor&#39; | &#39;action&#39; | &#39;username&#39; | &#39;created_at&#39;**]**Array<&#39;environment_key&#39; &#124; &#39;entity&#39; &#124; &#39;entity_id&#39; &#124; &#39;actor&#39; &#124; &#39;action&#39; &#124; &#39;username&#39; &#124; &#39;created_at&#39;>** | Sort by field | (optional) defaults to undefined|
| **sortOrder** | **SortOrder** | Sort order | (optional) defaults to undefined|
| **page** | [**number**] | Page number (starts from 1) | (optional) defaults to 1|
| **perPage** | [**number**] | Items per page | (optional) defaults to 20|


### Return type

**ListProjectAuditLogs200Response**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of audit log entries |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectChanges**
> ListChangesResponse listProjectChanges()

Get history of changes made to project features, rules, and other entities grouped by request_id

### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; //Project ID (default to undefined)
let page: number; //Page number (starts from 1) (optional) (default to 1)
let perPage: number; //Items per page (optional) (default to 20)
let sortBy: 'created_at' | 'actor' | 'entity'; //Sort by field (optional) (default to 'created_at')
let sortOrder: SortOrder; //Sort order (optional) (default to undefined)
let actor: string; //Filter by actor (system, sdk, user:<user_id>) (optional) (default to undefined)
let entity: EntityType; //Filter by entity type (optional) (default to undefined)
let action: AuditAction; //Filter by action type (optional) (default to undefined)
let featureId: string; //Filter by specific feature ID (optional) (default to undefined)
let from: string; //Filter changes from this date (ISO 8601 format) (optional) (default to undefined)
let to: string; //Filter changes until this date (ISO 8601 format) (optional) (default to undefined)

const { status, data } = await apiInstance.listProjectChanges(
    projectId,
    page,
    perPage,
    sortBy,
    sortOrder,
    actor,
    entity,
    action,
    featureId,
    from,
    to
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] | Project ID | defaults to undefined|
| **page** | [**number**] | Page number (starts from 1) | (optional) defaults to 1|
| **perPage** | [**number**] | Items per page | (optional) defaults to 20|
| **sortBy** | [**&#39;created_at&#39; | &#39;actor&#39; | &#39;entity&#39;**]**Array<&#39;created_at&#39; &#124; &#39;actor&#39; &#124; &#39;entity&#39;>** | Sort by field | (optional) defaults to 'created_at'|
| **sortOrder** | **SortOrder** | Sort order | (optional) defaults to undefined|
| **actor** | [**string**] | Filter by actor (system, sdk, user:&lt;user_id&gt;) | (optional) defaults to undefined|
| **entity** | **EntityType** | Filter by entity type | (optional) defaults to undefined|
| **action** | **AuditAction** | Filter by action type | (optional) defaults to undefined|
| **featureId** | [**string**] | Filter by specific feature ID | (optional) defaults to undefined|
| **from** | [**string**] | Filter changes from this date (ISO 8601 format) | (optional) defaults to undefined|
| **to** | [**string**] | Filter changes until this date (ISO 8601 format) | (optional) defaults to undefined|


### Return type

**ListChangesResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of change groups for the project |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectEnvironments**
> ListEnvironmentsResponse listProjectEnvironments()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)

const { status, data } = await apiInstance.listProjectEnvironments(
    projectId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**ListEnvironmentsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of environments |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectFeatures**
> ListFeaturesResponse listProjectFeatures()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; //Environment key (dev, stage, prod) to filter features (default to undefined)
let kind: 'simple' | 'multivariant'; //Filter by feature kind (optional) (default to undefined)
let enabled: boolean; //Filter by enabled state (optional) (default to undefined)
let textSelector: string; //Case-insensitive text search across key, name, description, rollout_key (optional) (default to undefined)
let tagIds: string; //Filter by tag IDs (comma-separated) (optional) (default to undefined)
let sortBy: 'name' | 'key' | 'enabled' | 'kind' | 'created_at' | 'updated_at'; //Sort by field (optional) (default to undefined)
let sortOrder: SortOrder; //Sort order (optional) (default to undefined)
let page: number; //Page number (starts from 1) (optional) (default to 1)
let perPage: number; //Items per page (optional) (default to 20)

const { status, data } = await apiInstance.listProjectFeatures(
    projectId,
    environmentKey,
    kind,
    enabled,
    textSelector,
    tagIds,
    sortBy,
    sortOrder,
    page,
    perPage
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] | Environment key (dev, stage, prod) to filter features | defaults to undefined|
| **kind** | [**&#39;simple&#39; | &#39;multivariant&#39;**]**Array<&#39;simple&#39; &#124; &#39;multivariant&#39;>** | Filter by feature kind | (optional) defaults to undefined|
| **enabled** | [**boolean**] | Filter by enabled state | (optional) defaults to undefined|
| **textSelector** | [**string**] | Case-insensitive text search across key, name, description, rollout_key | (optional) defaults to undefined|
| **tagIds** | [**string**] | Filter by tag IDs (comma-separated) | (optional) defaults to undefined|
| **sortBy** | [**&#39;name&#39; | &#39;key&#39; | &#39;enabled&#39; | &#39;kind&#39; | &#39;created_at&#39; | &#39;updated_at&#39;**]**Array<&#39;name&#39; &#124; &#39;key&#39; &#124; &#39;enabled&#39; &#124; &#39;kind&#39; &#124; &#39;created_at&#39; &#124; &#39;updated_at&#39;>** | Sort by field | (optional) defaults to undefined|
| **sortOrder** | **SortOrder** | Sort order | (optional) defaults to undefined|
| **page** | [**number**] | Page number (starts from 1) | (optional) defaults to 1|
| **perPage** | [**number**] | Items per page | (optional) defaults to 20|


### Return type

**ListFeaturesResponse**

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

# **listProjectMemberships**
> Array<Membership> listProjectMemberships()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)

const { status, data } = await apiInstance.listProjectMemberships(
    projectId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|


### Return type

**Array<Membership>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of memberships |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectSegments**
> ListSegmentsResponse listProjectSegments()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let textSelector: string; //Case-insensitive text search across name, description (optional) (default to undefined)
let sortBy: 'name' | 'created_at' | 'updated_at'; //Sort by field (optional) (default to undefined)
let sortOrder: SortOrder; //Sort order (optional) (default to undefined)
let page: number; //Page number (starts from 1) (optional) (default to 1)
let perPage: number; //Items per page (optional) (default to 20)

const { status, data } = await apiInstance.listProjectSegments(
    projectId,
    textSelector,
    sortBy,
    sortOrder,
    page,
    perPage
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **textSelector** | [**string**] | Case-insensitive text search across name, description | (optional) defaults to undefined|
| **sortBy** | [**&#39;name&#39; | &#39;created_at&#39; | &#39;updated_at&#39;**]**Array<&#39;name&#39; &#124; &#39;created_at&#39; &#124; &#39;updated_at&#39;>** | Sort by field | (optional) defaults to undefined|
| **sortOrder** | **SortOrder** | Sort order | (optional) defaults to undefined|
| **page** | [**number**] | Page number (starts from 1) | (optional) defaults to 1|
| **perPage** | [**number**] | Items per page | (optional) defaults to 20|


### Return type

**ListSegmentsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of segments for the project |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectSettings**
> ListProjectSettingsResponse listProjectSettings()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let page: number; // (optional) (default to 1)
let perPage: number; // (optional) (default to 20)

const { status, data } = await apiInstance.listProjectSettings(
    projectId,
    page,
    perPage
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to 1|
| **perPage** | [**number**] |  | (optional) defaults to 20|


### Return type

**ListProjectSettingsResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of project settings |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Project not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listProjectTags**
> Array<ProjectTag> listProjectTags()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let categoryId: string; //Filter by category ID (optional) (default to undefined)

const { status, data } = await apiInstance.listProjectTags(
    projectId,
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **categoryId** | [**string**] | Filter by category ID | (optional) defaults to undefined|


### Return type

**Array<ProjectTag>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of tags for project |  -  |
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

# **listRolePermissions**
> Array<ListRolePermissions200ResponseInner> listRolePermissions()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listRolePermissions();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<ListRolePermissions200ResponseInner>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Map of role to permissions |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listRoles**
> Array<Role> listRoles()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listRoles();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<Role>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of roles |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listRuleAttributes**
> Array<RuleAttributeEntity> listRuleAttributes()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.listRuleAttributes();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<RuleAttributeEntity>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of rule attributes |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listSegmentDesyncFeatureIDs**
> Array<string> listSegmentDesyncFeatureIDs()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let segmentId: string; // (default to undefined)

const { status, data } = await apiInstance.listSegmentDesyncFeatureIDs(
    segmentId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **segmentId** | [**string**] |  | defaults to undefined|


### Return type

**Array<string>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature IDs |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Segment not found |  -  |
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

# **markAllNotificationsAsRead**
> markAllNotificationsAsRead()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

const { status, data } = await apiInstance.markAllNotificationsAsRead();
```

### Parameters
This endpoint does not have any parameters.


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
|**204** | All notifications marked as read |  -  |
|**401** | Unauthorized |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **markNotificationAsRead**
> markNotificationAsRead()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let notificationId: number; // (default to undefined)

const { status, data } = await apiInstance.markNotificationAsRead(
    notificationId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **notificationId** | [**number**] |  | defaults to undefined|


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
|**204** | Notification marked as read |  -  |
|**401** | Unauthorized |  -  |
|**404** | Notification not found |  -  |
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

# **rejectPendingChange**
> SuccessResponse rejectPendingChange(rejectPendingChangeRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    RejectPendingChangeRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let pendingChangeId: string; // (default to undefined)
let rejectPendingChangeRequest: RejectPendingChangeRequest; //

const { status, data } = await apiInstance.rejectPendingChange(
    pendingChangeId,
    rejectPendingChangeRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **rejectPendingChangeRequest** | **RejectPendingChangeRequest**|  | |
| **pendingChangeId** | [**string**] |  | defaults to undefined|


### Return type

**SuccessResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Pending change rejected successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Pending change not found |  -  |
|**409** | Conflict - pending change is not in pending status |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **removeFeatureTag**
> PendingChangeResponse removeFeatureTag()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let tagId: string; // (default to undefined)

const { status, data } = await apiInstance.removeFeatureTag(
    featureId,
    tagId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **tagId** | [**string**] |  | defaults to undefined|


### Return type

**PendingChangeResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Tag removed from feature |  -  |
|**202** | Change is pending approval (for guarded features) |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature or tag not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
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

# **sendTestNotification**
> sendTestNotification()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)
let settingId: number; // (default to undefined)

const { status, data } = await apiInstance.sendTestNotification(
    projectId,
    environmentKey,
    settingId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|
| **settingId** | [**number**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Notification successfully sent |  -  |
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

# **syncCustomizedFeatureRule**
> RuleResponse syncCustomizedFeatureRule()


### Example

```typescript
import {
    DefaultApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let ruleId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)

const { status, data } = await apiInstance.syncCustomizedFeatureRule(
    featureId,
    ruleId,
    environmentKey
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **featureId** | [**string**] |  | defaults to undefined|
| **ruleId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


### Return type

**RuleResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature rule synchronized |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature or related resource not found |  -  |
|**500** | Internal server error |  -  |
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

# **testFeatureTimeline**
> FeatureTimelineResponse testFeatureTimeline(testFeatureTimelineRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    TestFeatureTimelineRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let environmentKey: string; //Target environment key (e.g., dev, stage, prod) (default to undefined)
let from: string; //Start of the period (inclusive) (default to undefined)
let to: string; //End of the period (exclusive) (default to undefined)
let location: string; //Browser\'s location string (default to undefined)
let testFeatureTimelineRequest: TestFeatureTimelineRequest; //

const { status, data } = await apiInstance.testFeatureTimeline(
    featureId,
    environmentKey,
    from,
    to,
    location,
    testFeatureTimelineRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **testFeatureTimelineRequest** | **TestFeatureTimelineRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] | Target environment key (e.g., dev, stage, prod) | defaults to undefined|
| **from** | [**string**] | Start of the period (inclusive) | defaults to undefined|
| **to** | [**string**] | End of the period (exclusive) | defaults to undefined|
| **location** | [**string**] | Browser\&#39;s location string | defaults to undefined|


### Return type

**FeatureTimelineResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Feature timeline with test schedules |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**500** | Internal server error |  -  |
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
let environmentKey: string; // (default to undefined)
let toggleFeatureRequest: ToggleFeatureRequest; //

const { status, data } = await apiInstance.toggleFeature(
    featureId,
    environmentKey,
    toggleFeatureRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **toggleFeatureRequest** | **ToggleFeatureRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
|**202** | Feature is guarded and change is pending approval |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**409** | Feature is already locked by another pending change |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateCategory**
> CategoryResponse updateCategory(updateCategoryRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateCategoryRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let categoryId: string; // (default to undefined)
let updateCategoryRequest: UpdateCategoryRequest; //

const { status, data } = await apiInstance.updateCategory(
    categoryId,
    updateCategoryRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateCategoryRequest** | **UpdateCategoryRequest**|  | |
| **categoryId** | [**string**] |  | defaults to undefined|


### Return type

**CategoryResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Category updated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Category not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateEnvironment**
> EnvironmentResponse updateEnvironment(updateEnvironmentRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateEnvironmentRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let environmentId: number; // (default to undefined)
let updateEnvironmentRequest: UpdateEnvironmentRequest; //

const { status, data } = await apiInstance.updateEnvironment(
    environmentId,
    updateEnvironmentRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateEnvironmentRequest** | **UpdateEnvironmentRequest**|  | |
| **environmentId** | [**number**] |  | defaults to undefined|


### Return type

**EnvironmentResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Environment updated |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Environment not found |  -  |
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
let environmentKey: string; // (default to undefined)
let createFeatureRequest: CreateFeatureRequest; //

const { status, data } = await apiInstance.updateFeature(
    featureId,
    environmentKey,
    createFeatureRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createFeatureRequest** | **CreateFeatureRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|


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
|**202** | Change is pending approval (for guarded features) |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Feature not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateFeatureAlgorithm**
> FeatureAlgorithm updateFeatureAlgorithm(updateFeatureAlgorithmRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateFeatureAlgorithmRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let featureId: string; // (default to undefined)
let environmentId: number; // (default to undefined)
let updateFeatureAlgorithmRequest: UpdateFeatureAlgorithmRequest; //

const { status, data } = await apiInstance.updateFeatureAlgorithm(
    featureId,
    environmentId,
    updateFeatureAlgorithmRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateFeatureAlgorithmRequest** | **UpdateFeatureAlgorithmRequest**|  | |
| **featureId** | [**string**] |  | defaults to undefined|
| **environmentId** | [**number**] |  | defaults to undefined|


### Return type

**FeatureAlgorithm**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Updated feature algorithm |  -  |
|**202** | Change is pending approval (for guarded features) |  -  |
|**404** | Algorithm not found |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**409** | Conflict |  -  |
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
|**202** | Change is pending approval (for guarded features) |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Schedule not found |  -  |
|**409** | Conflict - change cannot be applied due to existing pending change or lock |  -  |
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

# **updateNotificationSetting**
> NotificationSetting updateNotificationSetting(updateNotificationSettingRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateNotificationSettingRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let environmentKey: string; // (default to undefined)
let settingId: number; // (default to undefined)
let updateNotificationSettingRequest: UpdateNotificationSettingRequest; //

const { status, data } = await apiInstance.updateNotificationSetting(
    projectId,
    environmentKey,
    settingId,
    updateNotificationSettingRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateNotificationSettingRequest** | **UpdateNotificationSettingRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|
| **environmentKey** | [**string**] |  | defaults to undefined|
| **settingId** | [**number**] |  | defaults to undefined|


### Return type

**NotificationSetting**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Notification setting updated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden - Not authorized to modify this notification setting |  -  |
|**404** | Notification setting not found |  -  |
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

# **updateProjectMembership**
> Membership updateProjectMembership(updateMembershipRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateMembershipRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let membershipId: string; // (default to undefined)
let updateMembershipRequest: UpdateMembershipRequest; //

const { status, data } = await apiInstance.updateProjectMembership(
    projectId,
    membershipId,
    updateMembershipRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateMembershipRequest** | **UpdateMembershipRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|
| **membershipId** | [**string**] |  | defaults to undefined|


### Return type

**Membership**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Membership updated |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Membership not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateProjectSetting**
> ProjectSettingResponse updateProjectSetting(updateProjectSettingRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateProjectSettingRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let settingName: string; // (default to undefined)
let updateProjectSettingRequest: UpdateProjectSettingRequest; //

const { status, data } = await apiInstance.updateProjectSetting(
    projectId,
    settingName,
    updateProjectSettingRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateProjectSettingRequest** | **UpdateProjectSettingRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|
| **settingName** | [**string**] |  | defaults to undefined|


### Return type

**ProjectSettingResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Project setting updated |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Setting not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateProjectTag**
> ProjectTagResponse updateProjectTag(updateProjectTagRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateProjectTagRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let projectId: string; // (default to undefined)
let tagId: string; // (default to undefined)
let updateProjectTagRequest: UpdateProjectTagRequest; //

const { status, data } = await apiInstance.updateProjectTag(
    projectId,
    tagId,
    updateProjectTagRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateProjectTagRequest** | **UpdateProjectTagRequest**|  | |
| **projectId** | [**string**] |  | defaults to undefined|
| **tagId** | [**string**] |  | defaults to undefined|


### Return type

**ProjectTagResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Tag updated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Tag not found |  -  |
|**500** | Internal server error |  -  |
|**0** | Unexpected error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateSegment**
> SegmentResponse updateSegment(updateSegmentRequest)


### Example

```typescript
import {
    DefaultApi,
    Configuration,
    UpdateSegmentRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DefaultApi(configuration);

let segmentId: string; // (default to undefined)
let updateSegmentRequest: UpdateSegmentRequest; //

const { status, data } = await apiInstance.updateSegment(
    segmentId,
    updateSegmentRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **updateSegmentRequest** | **UpdateSegmentRequest**|  | |
| **segmentId** | [**string**] |  | defaults to undefined|


### Return type

**SegmentResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Segment updated successfully |  -  |
|**400** | Bad request |  -  |
|**401** | Unauthorized |  -  |
|**403** | Permission denied |  -  |
|**404** | Segment not found |  -  |
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

