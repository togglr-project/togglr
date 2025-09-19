package featuresprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/etoggle/internal/domain"
)

func TestBuildFeatureTimeline(t *testing.T) {
	t.SkipNow()
	now := time.Date(2025, 9, 16, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		feature   domain.FeatureExtended
		from      time.Time
		to        time.Time
		wantTimes []time.Time
		wantEn    []bool
	}{
		{
			name: "enabled feature without schedules",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f1",
					Key:       "always_on",
					Enabled:   true,
					CreatedAt: now.Add(-time.Hour),
				},
				Schedules: nil,
			},
			from: now.Add(-30 * time.Minute),
			to:   now.Add(30 * time.Minute),
			wantTimes: []time.Time{
				now.Add(-30 * time.Minute), // from
				now.Add(30 * time.Minute),  // to
				now.Add(-time.Hour),        // created_at
			},
			wantEn: []bool{true, true, true},
		},
		{
			name: "feature with start and end schedule",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f2",
					Key:       "with_range",
					Enabled:   false,
					CreatedAt: now.Add(-time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched1",
						StartsAt:  ptrTime(now.Add(5 * time.Minute)),
						EndsAt:    ptrTime(now.Add(15 * time.Minute)),
						Action:    domain.FeatureScheduleActionEnable,
						Timezone:  "UTC",
						CreatedAt: now,
					},
				},
			},
			from: now,
			to:   now.Add(30 * time.Minute),
			wantTimes: []time.Time{
				now,                       // from
				now.Add(5 * time.Minute),  // starts_at
				now.Add(15 * time.Minute), // ends_at
				now.Add(30 * time.Minute), // to
			},
			wantEn: []bool{false, true, false, false},
		},
		{
			name: "feature with cron schedule",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f3",
					Key:       "cron_based",
					Enabled:   false,
					CreatedAt: now,
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "cron1",
						CronExpr:  ptrStr("*/10 * * * *"), // каждые 10 минут
						Action:    domain.FeatureScheduleActionEnable,
						Timezone:  "UTC",
						CreatedAt: now,
					},
				},
			},
			from: now,
			to:   now.Add(30 * time.Minute),
			wantTimes: []time.Time{
				now, // from
				time.Date(2025, 9, 16, 12, 10, 0, 0, time.UTC), // cron 1
				time.Date(2025, 9, 16, 12, 20, 0, 0, time.UTC), // cron 2
				time.Date(2025, 9, 16, 12, 30, 0, 0, time.UTC), // cron 3
				now.Add(30 * time.Minute),                      // to
			},
			wantEn: []bool{true, true, true, true, true},
		},
		{
			name: "disabled feature without schedules",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f4",
					Key:       "always_off",
					Enabled:   false,
					CreatedAt: now.Add(-time.Hour),
				},
			},
			from: now,
			to:   now.Add(30 * time.Minute),
			wantTimes: []time.Time{
				now,                       // from
				now.Add(30 * time.Minute), // to
			},
			wantEn: []bool{false, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(nil, nil, nil, 0)
			got, err := svc.BuildFeatureTimeline(tt.feature, tt.from, tt.to)
			assert.NoError(t, err)
			assert.Equal(t, len(tt.wantTimes), len(got), "unexpected number of events")

			for i := range tt.wantTimes {
				assert.Equal(t, tt.wantTimes[i], got[i].Time, "event %d time mismatch", i)
				assert.Equal(t, tt.wantEn[i], got[i].Enabled, "event %d enabled mismatch", i)
			}
		})
	}
}

func ptrStr(s string) *string {
	return &s
}
