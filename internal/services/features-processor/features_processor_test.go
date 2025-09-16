package featuresprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/etoggle/internal/domain"
)

func TestIsScheduleActive(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2025, 9, 16, 10, 0, 0, 0, loc)

	tests := []struct {
		name     string
		schedule domain.FeatureSchedule
		now      time.Time
		expected bool
	}{
		{
			name: "active with no limits",
			schedule: domain.FeatureSchedule{
				Timezone: "UTC",
			},
			now:      now,
			expected: true,
		},
		{
			name: "inactive before starts_at",
			schedule: domain.FeatureSchedule{
				StartsAt: ptrTime(now.Add(1 * time.Hour)),
				Timezone: "UTC",
			},
			now:      now,
			expected: false,
		},
		{
			name: "inactive after ends_at",
			schedule: domain.FeatureSchedule{
				EndsAt:   ptrTime(now.Add(-1 * time.Hour)),
				Timezone: "UTC",
			},
			now:      now,
			expected: false,
		},
		{
			name: "active between starts_at and ends_at",
			schedule: domain.FeatureSchedule{
				StartsAt: ptrTime(now.Add(-1 * time.Hour)),
				EndsAt:   ptrTime(now.Add(1 * time.Hour)),
				Timezone: "UTC",
			},
			now:      now,
			expected: true,
		},
		{
			name: "active cron expr matches hour",
			schedule: domain.FeatureSchedule{
				CronExpr: ptrString("0 10 * * *"),
				Timezone: "UTC",
			},
			now:      now,
			expected: true,
		},
		{
			name: "inactive cron expr not matching hour",
			schedule: domain.FeatureSchedule{
				CronExpr: ptrString("0 9 * * *"),
				Timezone: "UTC",
			},
			now:      now,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := IsScheduleActive(tt.schedule, tt.now)
			assert.Equal(t, tt.expected, ok)
		})
	}
}

func ptrTime(t time.Time) *time.Time { return &t }
func ptrString(s string) *string     { return &s }
