package apibackend

import (
	"time"

	"github.com/togglr-project/togglr/internal/domain"
	generatedapi "github.com/togglr-project/togglr/internal/generated/server"
)

func optString2FlagVariantIDRef(optString generatedapi.OptUUID) *domain.FlagVariantID {
	if !optString.IsSet() {
		return nil
	}

	id := domain.FlagVariantID(optString.Value.String())

	return &id
}

func flagVariantRef2OptString(flagVariantID *domain.FlagVariantID) generatedapi.OptString {
	if flagVariantID == nil {
		return generatedapi.OptString{}
	}

	return generatedapi.NewOptString(flagVariantID.String())
}

func ruleAttribute2OptString(attr domain.RuleAttribute) generatedapi.OptString {
	if attr == "" {
		return generatedapi.OptString{}
	}

	return generatedapi.NewOptString(attr.String())
}

func optNilDurationToPtr(v generatedapi.OptNilDuration) *time.Duration {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	t, _ := v.Get()

	return &t
}

func optNilDateTimeToPtr(v generatedapi.OptNilDateTime) *time.Time {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	t, _ := v.Get()

	return &t
}

func ptrToOptNilDateTime(p *time.Time) generatedapi.OptNilDateTime {
	if p == nil {
		var o generatedapi.OptNilDateTime

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilDateTime(*p)
}

func ptrToOptNilDuration(p *time.Duration) generatedapi.OptNilDuration {
	if p == nil {
		var o generatedapi.OptNilDuration

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilDuration(*p)
}

func optNilStringToPtr(v generatedapi.OptNilString) *string {
	if !v.IsSet() || v.IsNull() {
		return nil
	}

	s, _ := v.Get()

	return &s
}

func ptrToOptNilString(p *string) generatedapi.OptNilString {
	if p == nil {
		var o generatedapi.OptNilString

		o.SetToNull()

		return o
	}

	return generatedapi.NewOptNilString(*p)
}

func timePtrString(p *time.Time) interface{} {
	if p == nil {
		return nil
	}

	return p.Format(time.RFC3339)
}

// durationPtrString converts *time.Duration to string (e.g., "30m0s"), or nil if pointer is nil.
func durationPtrString(p *time.Duration) interface{} {
	if p == nil {
		return nil
	}

	return p.String()
}

// stringPtrValue dereferences *string into its value, or nil if pointer is nil.
func stringPtrValue(p *string) interface{} {
	if p == nil {
		return nil
	}

	return *p
}

// buildFeatureScheduleChangeDiff compares two FeatureSchedule values and returns a map of changed fields
// formatted for pending-change payloads. Values are converted to JSON-friendly representations.
func buildFeatureScheduleChangeDiff(
	oldSch domain.FeatureSchedule,
	newSch domain.FeatureSchedule,
) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	if timePtrString(oldSch.StartsAt) != timePtrString(newSch.StartsAt) {
		changes["starts_at"] = domain.ChangeValue{
			Old: timePtrString(oldSch.StartsAt),
			New: timePtrString(newSch.StartsAt),
		}
	}

	if timePtrString(oldSch.EndsAt) != timePtrString(newSch.EndsAt) {
		changes["ends_at"] = domain.ChangeValue{
			Old: timePtrString(oldSch.EndsAt),
			New: timePtrString(newSch.EndsAt),
		}
	}

	if stringPtrValue(oldSch.CronExpr) != stringPtrValue(newSch.CronExpr) {
		changes["cron_expr"] = domain.ChangeValue{
			Old: stringPtrValue(oldSch.CronExpr),
			New: stringPtrValue(newSch.CronExpr),
		}
	}

	if durationPtrString(oldSch.CronDuration) != durationPtrString(newSch.CronDuration) {
		changes["cron_duration"] = domain.ChangeValue{
			Old: durationPtrString(oldSch.CronDuration),
			New: durationPtrString(newSch.CronDuration),
		}
	}

	if oldSch.Timezone != newSch.Timezone {
		changes["timezone"] = domain.ChangeValue{Old: oldSch.Timezone, New: newSch.Timezone}
	}

	if oldSch.Action != newSch.Action {
		changes["action"] = domain.ChangeValue{Old: oldSch.Action.String(), New: newSch.Action.String()}
	}

	return changes
}
