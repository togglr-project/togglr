package featuresprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/togglr-project/togglr/internal/domain"
)

func TestBuildFeatureTimeline(t *testing.T) {
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
			},
			wantEn: []bool{true, true},
		},
		{
			name: "feature with start and end schedule",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f2",
					Key:       "with_range",
					Enabled:   true, // Master Enable ON
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
				now,                       // from - baseline OFF
				now.Add(5 * time.Minute),  // starts_at - enable
				now.Add(15 * time.Minute), // ends_at - back to baseline OFF
				now.Add(30 * time.Minute), // to - baseline OFF
			},
			wantEn: []bool{false, true, false, false},
		},
		{
			name: "feature with cron schedule",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f3",
					Key:       "cron_based",
					Enabled:   true, // Master Enable ON
					CreatedAt: now,
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:           "cron1",
						CronExpr:     ptrStr("*/10 * * * *"), // каждые 10 минут
						CronDuration: ptrDuration(time.Minute * 5),
						Action:       domain.FeatureScheduleActionEnable,
						Timezone:     "UTC",
						CreatedAt:    now,
					},
				},
			},
			from: now.Add(-1 * time.Minute), // 11:59, чтобы начать с baseline
			to:   now.Add(25 * time.Minute), // 12:25, чтобы не попасть на 12:30
			wantTimes: []time.Time{
				now.Add(-1 * time.Minute),                      // 11:59 - baseline OFF
				time.Date(2025, 9, 16, 12, 0, 0, 0, time.UTC),  // 12:00 - enable
				time.Date(2025, 9, 16, 12, 5, 0, 0, time.UTC),  // 12:05 - back to baseline OFF
				time.Date(2025, 9, 16, 12, 10, 0, 0, time.UTC), // 12:10 - enable
				time.Date(2025, 9, 16, 12, 15, 0, 0, time.UTC), // 12:15 - back to baseline OFF
				time.Date(2025, 9, 16, 12, 20, 0, 0, time.UTC), // 12:20 - enable
				time.Date(2025, 9, 16, 12, 25, 0, 0, time.UTC), // 12:25 - back to baseline OFF
			},
			wantEn: []bool{false, true, false, true, false, true, false},
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

// TestBuildFeatureTimeline_BaselineLogic тестирует правильность baseline логики в timeline
func TestBuildFeatureTimeline_BaselineLogic(t *testing.T) {
	now := time.Date(2025, 9, 16, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		feature   domain.FeatureExtended
		from      time.Time
		to        time.Time
		wantTimes []time.Time
		wantEn    []bool
		desc      string
	}{
		{
			name: "repeating schedule with enable action - baseline should be OFF",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f1",
					Key:       "repeating_enable",
					Enabled:   true, // Master Enable ON
					CreatedAt: now.Add(-time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:           "cron1",
						CronExpr:     ptrStr("0 12 * * *"), // daily at 12:00
						CronDuration: ptrDuration(30 * time.Minute),
						Action:       domain.FeatureScheduleActionEnable,
						Timezone:     "UTC",
						CreatedAt:    now.Add(-30 * time.Minute),
					},
				},
			},
			from: now.Add(-1 * time.Minute), // 11:59
			to:   now.Add(31 * time.Minute), // 12:31
			wantTimes: []time.Time{
				now.Add(-1 * time.Minute),                      // 11:59 - baseline OFF
				time.Date(2025, 9, 16, 12, 0, 0, 0, time.UTC),  // 12:00 - enable
				time.Date(2025, 9, 16, 12, 30, 0, 0, time.UTC), // 12:30 - back to baseline OFF
				now.Add(31 * time.Minute),                      // 12:31 - baseline OFF
			},
			wantEn: []bool{false, true, false, false},
			desc:   "Repeating enable schedule: baseline OFF, active only during 9:00-9:30 window",
		},
		{
			name: "repeating schedule with disable action - baseline should be ON",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f2",
					Key:       "repeating_disable",
					Enabled:   true, // Master Enable ON
					CreatedAt: now.Add(-time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:           "cron2",
						CronExpr:     ptrStr("0 12 * * *"), // daily at 12:00
						CronDuration: ptrDuration(30 * time.Minute),
						Action:       domain.FeatureScheduleActionDisable,
						Timezone:     "UTC",
						CreatedAt:    now.Add(-30 * time.Minute),
					},
				},
			},
			from: now.Add(-1 * time.Minute), // 11:59
			to:   now.Add(31 * time.Minute), // 12:31
			wantTimes: []time.Time{
				now.Add(-1 * time.Minute),                      // 11:59 - baseline ON
				time.Date(2025, 9, 16, 12, 0, 0, 0, time.UTC),  // 12:00 - disable
				time.Date(2025, 9, 16, 12, 30, 0, 0, time.UTC), // 12:30 - back to baseline ON
				now.Add(31 * time.Minute),                      // 12:31 - baseline ON
			},
			wantEn: []bool{true, false, true, true},
			desc:   "Repeating disable schedule: baseline ON, disabled only during 9:00-9:30 window",
		},
		{
			name: "master enable OFF - feature completely disabled",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					ID:        "f3",
					Key:       "master_off",
					Enabled:   false, // Master Enable OFF
					CreatedAt: now.Add(-time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:           "cron3",
						CronExpr:     ptrStr("0 9 * * *"), // daily at 9:00
						CronDuration: ptrDuration(30 * time.Minute),
						Action:       domain.FeatureScheduleActionEnable,
						Timezone:     "UTC",
						CreatedAt:    now.Add(-30 * time.Minute),
					},
				},
			},
			from: now.Add(-30 * time.Minute), // 11:30
			to:   now.Add(30 * time.Minute),  // 12:30
			wantTimes: []time.Time{
				now.Add(-30 * time.Minute), // 11:30 - disabled
				now.Add(30 * time.Minute),  // 12:30 - disabled
			},
			wantEn: []bool{false, false},
			desc:   "Master Enable OFF: feature completely disabled regardless of schedules",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(nil, nil, nil, 0)
			got, err := svc.BuildFeatureTimeline(tt.feature, tt.from, tt.to)
			assert.NoError(t, err, tt.desc)
			assert.Equal(t, len(tt.wantTimes), len(got), "unexpected number of events: %s", tt.desc)

			for i := range tt.wantTimes {
				assert.Equal(t, tt.wantTimes[i], got[i].Time, "event %d time mismatch: %s", i, tt.desc)
				assert.Equal(t, tt.wantEn[i], got[i].Enabled, "event %d enabled mismatch: %s", i, tt.desc)
			}
		})
	}
}
