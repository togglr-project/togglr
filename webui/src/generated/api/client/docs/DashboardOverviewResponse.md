# DashboardOverviewResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**projects** | [**Array&lt;ProjectHealth&gt;**](ProjectHealth.md) | Project-level health overview | [optional] [default to undefined]
**categories** | [**Array&lt;CategoryHealth&gt;**](CategoryHealth.md) | Per-category health | [optional] [default to undefined]
**feature_activity** | [**DashboardOverviewResponseFeatureActivity**](DashboardOverviewResponseFeatureActivity.md) |  | [optional] [default to undefined]
**recent_activity** | [**Array&lt;RecentActivity&gt;**](RecentActivity.md) | Recent batched changes | [optional] [default to undefined]
**risky_features** | [**Array&lt;RiskyFeature&gt;**](RiskyFeature.md) | Features with risky tags (critical, guarded, auto-disable) | [optional] [default to undefined]
**pending_summary** | [**Array&lt;PendingSummary&gt;**](PendingSummary.md) | Summary of pending changes | [optional] [default to undefined]

## Example

```typescript
import { DashboardOverviewResponse } from './api';

const instance: DashboardOverviewResponse = {
    projects,
    categories,
    feature_activity,
    recent_activity,
    risky_features,
    pending_summary,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
