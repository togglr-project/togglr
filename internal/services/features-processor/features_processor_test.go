package featuresprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/etoggle/internal/domain"
)

func TestIsFeatureActiveNow(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2025, 9, 16, 12, 0, 0, 0, loc)

	tests := []struct {
		name     string
		feature  domain.FeatureExtended
		now      time.Time
		expected bool
	}{
		{
			name: "enabled feature with no schedules stays enabled",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{Enabled: true},
			},
			now:      now,
			expected: true,
		},
		{
			name: "disabled feature with no schedules stays disabled",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{Enabled: false},
			},
			now:      now,
			expected: false,
		},
		{
			name: "schedule enable overrides disabled feature",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{Enabled: false},
				Schedules: []domain.FeatureSchedule{
					{
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-30 * time.Minute),
					},
				},
			},
			now:      now,
			expected: true,
		},
		{
			name: "schedule disable overrides enabled feature",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{Enabled: true},
				Schedules: []domain.FeatureSchedule{
					{
						Action:    domain.FeatureScheduleActionDisable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-30 * time.Minute),
					},
				},
			},
			now:      now,
			expected: false,
		},
		{
			name: "newer schedule overrides older one",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{Enabled: true},
				Schedules: []domain.FeatureSchedule{
					{
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-1 * time.Hour),
					},
					{
						Action:    domain.FeatureScheduleActionDisable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-10 * time.Minute),
					},
				},
			},
			now:      now,
			expected: false, // disable более свежее → перекрывает
		},
		{
			name: "same created_at, disable wins over enable",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{Enabled: true},
				Schedules: []domain.FeatureSchedule{
					{
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now,
					},
					{
						Action:    domain.FeatureScheduleActionDisable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now,
					},
				},
			},
			now:      now,
			expected: false, // при равенстве disable выше
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := IsFeatureActiveNow(tt.feature, tt.now)
			assert.Equal(t, tt.expected, ok)
		})
	}
}

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
