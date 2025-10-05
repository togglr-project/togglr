//nolint:nestif // fix it
package apibackend

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

//nolint:maintidx,gocognit // fix later
func (r *RestAPI) GetDashboardOverview(
	ctx context.Context,
	params generatedapi.GetDashboardOverviewParams,
) (generatedapi.GetDashboardOverviewRes, error) {
	var (
		projectIDPtr *string
		limit        = params.Limit.Or(20)
	)
	if params.ProjectID.IsSet() {
		pid := params.ProjectID.Value.String()
		projectIDPtr = &pid
	}

	overview, err := r.dashboardUseCase.Overview(ctx, params.EnvironmentKey, projectIDPtr, limit)
	if err != nil {
		slog.Error("get dashboard overview failed", "error", err)

		return nil, err
	}

	resp := generatedapi.DashboardOverviewResponse{}

	// Map projects health
	if len(overview.Projects) > 0 {
		items := make([]generatedapi.ProjectHealth, 0, len(overview.Projects))
		for i := range overview.Projects {
			ph := overview.Projects[i]
			var item generatedapi.ProjectHealth
			if id, err := uuid.Parse(ph.ProjectID); err == nil {
				item.ProjectID = generatedapi.NewOptUUID(id)
			}
			item.ProjectName = generatedapi.NewOptString(ph.ProjectName)
			if id, err := uuid.Parse(ph.EnvironmentID); err == nil {
				item.EnvironmentID = generatedapi.NewOptUUID(id)
			}
			item.EnvironmentKey = generatedapi.NewOptString(ph.EnvironmentKey)
			item.TotalFeatures = generatedapi.NewOptUint(ph.TotalFeatures)
			item.EnabledFeatures = generatedapi.NewOptUint(ph.EnabledFeatures)
			item.DisabledFeatures = generatedapi.NewOptUint(ph.DisabledFeatures)
			item.AutoDisableManagedFeatures = generatedapi.NewOptUint(ph.AutoDisableManagedFeatures)
			item.UncategorizedFeatures = generatedapi.NewOptUint(ph.UncategorizedFeatures)
			item.GuardedFeatures = generatedapi.NewOptUint(ph.GuardedFeatures)
			item.PendingFeatures = generatedapi.NewOptUint(ph.PendingFeatures)
			item.PendingGuardedFeatures = generatedapi.NewOptUint(ph.PendingGuardedFeatures)
			// Health status enum
			switch string(ph.HealthStatus) {
			case "green":
				item.HealthStatus = generatedapi.NewOptProjectHealthHealthStatus(generatedapi.ProjectHealthHealthStatusGreen)
			case "yellow":
				item.HealthStatus = generatedapi.NewOptProjectHealthHealthStatus(generatedapi.ProjectHealthHealthStatusYellow)
			case "red":
				item.HealthStatus = generatedapi.NewOptProjectHealthHealthStatus(generatedapi.ProjectHealthHealthStatusRed)
			}
			items = append(items, item)
		}
		resp.Projects = items
	}

	// Map category health
	if len(overview.Categories) > 0 {
		items := make([]generatedapi.CategoryHealth, 0, len(overview.Categories))
		for i := range overview.Categories {
			ch := overview.Categories[i]
			var item generatedapi.CategoryHealth
			if id, err := uuid.Parse(ch.ProjectID); err == nil {
				item.ProjectID = generatedapi.NewOptUUID(id)
			}
			item.ProjectName = generatedapi.NewOptString(ch.ProjectName)
			if id, err := uuid.Parse(ch.EnvironmentID); err == nil {
				item.EnvironmentID = generatedapi.NewOptUUID(id)
			}
			item.EnvironmentKey = generatedapi.NewOptString(ch.EnvironmentKey)
			if id, err := uuid.Parse(ch.CategoryID); err == nil {
				item.CategoryID = generatedapi.NewOptUUID(id)
			}
			item.CategoryName = generatedapi.NewOptString(ch.CategoryName)
			item.CategorySlug = generatedapi.NewOptString(ch.CategorySlug)
			item.TotalFeatures = generatedapi.NewOptUint(ch.TotalFeatures)
			item.EnabledFeatures = generatedapi.NewOptUint(ch.EnabledFeatures)
			item.DisabledFeatures = generatedapi.NewOptUint(ch.DisabledFeatures)
			item.PendingFeatures = generatedapi.NewOptUint(ch.PendingFeatures)
			item.GuardedFeatures = generatedapi.NewOptUint(ch.GuardedFeatures)
			item.AutoDisableManagedFeatures = generatedapi.NewOptUint(ch.AutoDisableManagedFeatures)
			item.PendingGuardedFeatures = generatedapi.NewOptUint(ch.PendingGuardedFeatures)
			switch string(ch.HealthStatus) {
			case "green":
				item.HealthStatus = generatedapi.NewOptCategoryHealthHealthStatus(generatedapi.CategoryHealthHealthStatusGreen)
			case "yellow":
				item.HealthStatus = generatedapi.NewOptCategoryHealthHealthStatus(generatedapi.CategoryHealthHealthStatusYellow)
			case "red":
				item.HealthStatus = generatedapi.NewOptCategoryHealthHealthStatus(generatedapi.CategoryHealthHealthStatusRed)
			}
			items = append(items, item)
		}
		resp.Categories = items
	}

	// Map recent activity
	if len(overview.RecentActivity) > 0 {
		items := make([]generatedapi.RecentActivity, 0, len(overview.RecentActivity))
		for i := range overview.RecentActivity {
			ra := overview.RecentActivity[i]
			var item generatedapi.RecentActivity
			if id, err := uuid.Parse(ra.ProjectID); err == nil {
				item.ProjectID = generatedapi.NewOptUUID(id)
			}
			if id, err := uuid.Parse(ra.EnvironmentID); err == nil {
				item.EnvironmentID = generatedapi.NewOptUUID(id)
			}
			item.EnvironmentKey = generatedapi.NewOptString(ra.EnvironmentKey)
			item.ProjectName = generatedapi.NewOptString(ra.ProjectName)
			if id, err := uuid.Parse(ra.RequestID); err == nil {
				item.RequestID = generatedapi.NewOptUUID(id)
			}
			item.Actor = generatedapi.NewOptString(ra.Actor)
			item.CreatedAt = generatedapi.NewOptDateTime(ra.CreatedAt)
			switch ra.Status {
			case "applied":
				item.Status = generatedapi.NewOptRecentActivityStatus(generatedapi.RecentActivityStatusApplied)
			case "pending":
				item.Status = generatedapi.NewOptRecentActivityStatus(generatedapi.RecentActivityStatusPending)
			case "rejected":
				item.Status = generatedapi.NewOptRecentActivityStatus(generatedapi.RecentActivityStatusRejected)
			}
			if len(ra.Changes) > 0 {
				changes := make([]generatedapi.RecentActivityChangesItem, 0, len(ra.Changes))
				for _, ch := range ra.Changes {
					var c generatedapi.RecentActivityChangesItem
					c.Entity = generatedapi.NewOptString(ch.Entity)
					if id, err := uuid.Parse(ch.EntityID); err == nil {
						c.EntityID = generatedapi.NewOptUUID(id)
					}
					c.Action = generatedapi.NewOptString(ch.Action)
					changes = append(changes, c)
				}
				item.Changes = changes
			}
			items = append(items, item)
		}
		resp.RecentActivity = items
	}

	// Map risky features
	if len(overview.RiskyFeatures) > 0 {
		items := make([]generatedapi.RiskyFeature, 0, len(overview.RiskyFeatures))
		for i := range overview.RiskyFeatures {
			rf := overview.RiskyFeatures[i]
			var item generatedapi.RiskyFeature
			if id, err := uuid.Parse(rf.ProjectID); err == nil {
				item.ProjectID = generatedapi.NewOptUUID(id)
			}
			item.ProjectName = generatedapi.NewOptString(rf.ProjectName)
			if id, err := uuid.Parse(rf.EnvironmentID); err == nil {
				item.EnvironmentID = generatedapi.NewOptUUID(id)
			}
			item.EnvironmentKey = generatedapi.NewOptString(rf.EnvironmentKey)
			if id, err := uuid.Parse(rf.FeatureID); err == nil {
				item.FeatureID = generatedapi.NewOptUUID(id)
			}
			item.FeatureName = generatedapi.NewOptString(rf.FeatureName)
			item.Enabled = generatedapi.NewOptBool(rf.Enabled)
			item.HasPending = generatedapi.NewOptBool(rf.HasPending)
			item.RiskyTags = generatedapi.NewOptString(rf.RiskyTags)
			items = append(items, item)
		}
		resp.RiskyFeatures = items
	}

	// Map pending summary
	if len(overview.PendingSummary) > 0 {
		items := make([]generatedapi.PendingSummary, 0, len(overview.PendingSummary))
		for i := range overview.PendingSummary {
			ps := overview.PendingSummary[i]
			var item generatedapi.PendingSummary
			if id, err := uuid.Parse(ps.ProjectID); err == nil {
				item.ProjectID = generatedapi.NewOptUUID(id)
			}
			item.ProjectName = generatedapi.NewOptString(ps.ProjectName)
			if id, err := uuid.Parse(ps.EnvironmentID); err == nil {
				item.EnvironmentID = generatedapi.NewOptUUID(id)
			}
			item.EnvironmentKey = generatedapi.NewOptString(ps.EnvironmentKey)
			item.TotalPending = generatedapi.NewOptUint(ps.TotalPending)
			item.PendingFeatureChanges = generatedapi.NewOptUint(ps.PendingFeatureChanges)
			item.PendingGuardedChanges = generatedapi.NewOptUint(ps.PendingGuardedChanges)
			if ps.OldestRequestAt != nil {
				item.OldestRequestAt = generatedapi.NewOptDateTime(*ps.OldestRequestAt)
			}
			items = append(items, item)
		}
		resp.PendingSummary = items
	}

	// Feature activity (optional, currently empty)
	if len(overview.Upcoming) > 0 || len(overview.Recent) > 0 {
		fa := generatedapi.DashboardOverviewResponseFeatureActivity{}
		if len(overview.Upcoming) > 0 {
			upItems := make([]generatedapi.FeatureUpcoming, 0, len(overview.Upcoming))
			for _, up := range overview.Upcoming {
				var fu generatedapi.FeatureUpcoming
				if id, err := uuid.Parse(up.FeatureID); err == nil {
					fu.FeatureID = generatedapi.NewOptUUID(id)
				}
				fu.FeatureName = generatedapi.NewOptString(up.FeatureName)
				switch up.NextState {
				case "enabled":
					fu.NextState = generatedapi.NewOptFeatureUpcomingNextState(generatedapi.FeatureUpcomingNextStateEnabled)
				case "disabled":
					fu.NextState = generatedapi.NewOptFeatureUpcomingNextState(generatedapi.FeatureUpcomingNextStateDisabled)
				}
				fu.At = generatedapi.NewOptDateTime(up.At)
				upItems = append(upItems, fu)
			}
			fa.Upcoming = upItems
		}
		if len(overview.Recent) > 0 {
			rcItems := make([]generatedapi.FeatureRecent, 0, len(overview.Recent))
			for _, rc := range overview.Recent {
				var fr generatedapi.FeatureRecent
				if id, err := uuid.Parse(rc.FeatureID); err == nil {
					fr.FeatureID = generatedapi.NewOptUUID(id)
				}
				fr.FeatureName = generatedapi.NewOptString(rc.FeatureName)
				fr.Action = generatedapi.NewOptString(rc.Action)
				fr.Actor = generatedapi.NewOptString(rc.Actor)
				fr.At = generatedapi.NewOptDateTime(rc.At)
				rcItems = append(rcItems, fr)
			}
			fa.Recent = rcItems
		}
		resp.FeatureActivity = generatedapi.NewOptDashboardOverviewResponseFeatureActivity(fa)
	}

	return &resp, nil
}
