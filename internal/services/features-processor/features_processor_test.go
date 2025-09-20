package featuresprocessor

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/etoggle/internal/domain"
)

func TestService_Evaluate(t *testing.T) {
	projectID := domain.ProjectID("proj1")
	featureKey := "my_feature"

	variantA := domain.FlagVariant{ID: "v1", Name: "A", RolloutPercent: 100}
	rule := domain.Rule{
		ID:        "r1",
		ProjectID: projectID,
		FeatureID: "f1",
		Conditions: domain.BooleanExpression{
			Condition: &domain.Condition{
				Attribute: "country",
				Operator:  domain.OpEq,
				Value:     "RU",
			},
		},
		Action:        domain.RuleActionAssign,
		FlagVariantID: &variantA.ID,
		CreatedAt:     time.Now(),
	}

	feature := domain.FeatureExtended{
		Feature: domain.Feature{
			ID:             "f1",
			ProjectID:      projectID,
			Key:            featureKey,
			Name:           "Test Feature",
			Kind:           domain.FeatureKindMultivariant,
			DefaultVariant: "default",
			Enabled:        true,
			CreatedAt:      time.Now(),
		},
		FlagVariants: []domain.FlagVariant{variantA},
		Rules:        []domain.Rule{rule},
	}

	holder := Holder{
		projectID: ProjectFeatures{
			featureKey: MakeFeaturePrepared(feature),
		},
	}

	svc := New(nil, nil, nil, 0)
	svc.holder = holder

	reqCtx := map[domain.RuleAttribute]any{"country": "RU"}
	value, enabled, found := svc.Evaluate(projectID, featureKey, reqCtx)

	assert.True(t, found)
	assert.True(t, enabled)
	assert.Equal(t, "A", value)
}

func TestService_Evaluate_TableDriven(t *testing.T) {
	projectID := domain.ProjectID("proj1")
	featureKey := "my_feature"

	variantA := domain.FlagVariant{ID: "v1", Name: "A", RolloutPercent: 100}
	rule := domain.Rule{
		ID:        "r1",
		ProjectID: projectID,
		FeatureID: "f1",
		Conditions: domain.BooleanExpression{
			Condition: &domain.Condition{
				Attribute: "country",
				Operator:  domain.OpEq,
				Value:     "RU",
			},
		},
		Action:        domain.RuleActionAssign,
		FlagVariantID: &variantA.ID,
		CreatedAt:     time.Now(),
	}

	baseFeature := domain.FeatureExtended{
		Feature: domain.Feature{
			ID:             "f1",
			ProjectID:      projectID,
			Key:            featureKey,
			Name:           "Test Feature",
			Kind:           domain.FeatureKindMultivariant,
			DefaultVariant: "default",
			Enabled:        true,
			CreatedAt:      time.Now(),
		},
		FlagVariants: []domain.FlagVariant{variantA},
		Rules:        []domain.Rule{rule},
	}

	tests := []struct {
		name          string
		feature       domain.FeatureExtended
		reqCtx        map[domain.RuleAttribute]any
		expectedValue string
		expectedEn    bool
		expectedFound bool
	}{
		{
			name:          "condition matches → variant A",
			feature:       baseFeature,
			reqCtx:        map[domain.RuleAttribute]any{"country": "RU"},
			expectedValue: "A",
			expectedEn:    true,
			expectedFound: true,
		},
		{
			name:          "condition does not match → default",
			feature:       baseFeature,
			reqCtx:        map[domain.RuleAttribute]any{"country": "US"},
			expectedValue: "default",
			expectedEn:    true,
			expectedFound: true,
		},
		{
			name: "feature disabled",
			feature: func() domain.FeatureExtended {
				f := baseFeature
				f.Enabled = false
				return f
			}(),
			reqCtx:        map[domain.RuleAttribute]any{"country": "RU"},
			expectedValue: "",
			expectedEn:    false,
			expectedFound: true,
		},
		{
			name:          "feature not found",
			feature:       domain.FeatureExtended{},
			reqCtx:        map[domain.RuleAttribute]any{"country": "RU"},
			expectedValue: "",
			expectedEn:    false,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holder := Holder{}
			if tt.expectedFound {
				holder = Holder{
					projectID: ProjectFeatures{
						featureKey: MakeFeaturePrepared(tt.feature),
					},
				}
			}

			svc := New(nil, nil, nil, 0)
			svc.holder = holder

			value, enabled, found := svc.Evaluate(projectID, featureKey, tt.reqCtx)

			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedEn, enabled)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

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
			name: "schedule disable overrides enabled feature",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched2",
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
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched3",
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-1 * time.Hour),
					},
					{
						ID:        "sched4",
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
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched5",
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now,
					},
					{
						ID:        "sched6",
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
			ok := IsFeatureActiveNow(MakeFeaturePrepared(tt.feature), tt.now)
			assert.Equal(t, tt.expected, ok)
		})
	}
}

func TestIsScheduleActive(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Date(2025, 9, 16, 10, 0, 0, 0, loc)

	featureCreatedAt := now.Add(-24 * time.Hour) // Feature created 24 hours ago

	tests := []struct {
		name     string
		schedule domain.FeatureSchedule
		now      time.Time
		expected bool
	}{
		{
			name: "active with no limits",
			schedule: domain.FeatureSchedule{
				Action:   domain.FeatureScheduleActionEnable,
				Timezone: "UTC",
			},
			now:      now,
			expected: true,
		},
		{
			name: "inactive before starts_at",
			schedule: domain.FeatureSchedule{
				Action:   domain.FeatureScheduleActionEnable,
				StartsAt: ptrTime(now.Add(1 * time.Hour)),
				Timezone: "UTC",
			},
			now:      now,
			expected: false,
		},
		{
			name: "inactive after ends_at",
			schedule: domain.FeatureSchedule{
				Action:   domain.FeatureScheduleActionEnable,
				EndsAt:   ptrTime(now.Add(-1 * time.Hour)),
				Timezone: "UTC",
			},
			now:      now,
			expected: false,
		},
		{
			name: "active between starts_at and ends_at",
			schedule: domain.FeatureSchedule{
				Action:   domain.FeatureScheduleActionEnable,
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
				ID:       "cron1",
				Action:   domain.FeatureScheduleActionEnable,
				CronExpr: ptrString("0 10 * * *"),
				Timezone: "UTC",
			},
			now:      now,
			expected: true,
		},
		{
			name: "active cron expr matches hour (fired at 9:00, still active at 10:00)",
			schedule: domain.FeatureSchedule{
				ID:       "cron2",
				Action:   domain.FeatureScheduleActionEnable,
				CronExpr: ptrString("0 9 * * *"),
				Timezone: "UTC",
			},
			now:      now,
			expected: true,
		},
		{
			name: "disable action returns false",
			schedule: domain.FeatureSchedule{
				Action:   domain.FeatureScheduleActionDisable,
				Timezone: "UTC",
			},
			now:      now,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crons := CronsMap{}
			if tt.schedule.CronExpr != nil {
				cronSched, err := ParseSchedule(*tt.schedule.CronExpr)
				require.NoError(t, err)

				crons[tt.schedule.ID] = cronSched
			}
			compatible, action := IsScheduleActive(tt.schedule, crons, tt.now, featureCreatedAt)
			ok := compatible && action == domain.FeatureScheduleActionEnable
			assert.Equal(t, tt.expected, ok)
		})
	}
}

func TestMatchCondition(t *testing.T) {
	tests := []struct {
		name     string
		reqCtx   map[domain.RuleAttribute]any
		cond     domain.Condition
		expected bool
	}{
		{
			name:     "eq operator matches",
			reqCtx:   map[domain.RuleAttribute]any{"country": "US"},
			cond:     domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "US"},
			expected: true,
		},
		{
			name:     "eq operator does not match",
			reqCtx:   map[domain.RuleAttribute]any{"country": "CA"},
			cond:     domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "US"},
			expected: false,
		},
		{
			name:     "neq operator matches",
			reqCtx:   map[domain.RuleAttribute]any{"age": 30},
			cond:     domain.Condition{Attribute: "age", Operator: domain.OpNotEq, Value: 40},
			expected: true,
		},
		{
			name:     "in operator matches",
			reqCtx:   map[domain.RuleAttribute]any{"role": "admin"},
			cond:     domain.Condition{Attribute: "role", Operator: domain.OpIn, Value: []string{"user", "admin"}},
			expected: true,
		},
		{
			name:     "not_in operator matches",
			reqCtx:   map[domain.RuleAttribute]any{"role": "guest"},
			cond:     domain.Condition{Attribute: "role", Operator: domain.OpNotIn, Value: []string{"user", "admin"}},
			expected: true,
		},
		{
			name:     "gt operator matches",
			reqCtx:   map[domain.RuleAttribute]any{"age": 25},
			cond:     domain.Condition{Attribute: "age", Operator: domain.OpGt, Value: 20},
			expected: true,
		},
		{
			name:     "gte operator matches equal",
			reqCtx:   map[domain.RuleAttribute]any{"age": 20},
			cond:     domain.Condition{Attribute: "age", Operator: domain.OpGte, Value: 20},
			expected: true,
		},
		{
			name:     "lt operator does not match",
			reqCtx:   map[domain.RuleAttribute]any{"age": 25},
			cond:     domain.Condition{Attribute: "age", Operator: domain.OpLt, Value: 20},
			expected: false,
		},
		{
			name:     "regex matches",
			reqCtx:   map[domain.RuleAttribute]any{"email": "test@example.com"},
			cond:     domain.Condition{Attribute: "email", Operator: domain.OpRegex, Value: `.+@example\.com`},
			expected: true,
		},
		{
			name:     "regex invalid pattern",
			reqCtx:   map[domain.RuleAttribute]any{"email": "test@example.com"},
			cond:     domain.Condition{Attribute: "email", Operator: domain.OpRegex, Value: `([a-z`}, // некорректный regex
			expected: false,
		},
		{
			name:     "percentage rollout 100% always matches",
			reqCtx:   map[domain.RuleAttribute]any{"user": "alice"},
			cond:     domain.Condition{Attribute: "user", Operator: domain.OpPercentage, Value: 100},
			expected: true,
		},
		{
			name:     "percentage rollout 0% never matches",
			reqCtx:   map[domain.RuleAttribute]any{"user": "bob"},
			cond:     domain.Condition{Attribute: "user", Operator: domain.OpPercentage, Value: 0},
			expected: false,
		},
		{
			name:     "attribute not in reqCtx returns false",
			reqCtx:   map[domain.RuleAttribute]any{},
			cond:     domain.Condition{Attribute: "missing", Operator: domain.OpEq, Value: "x"},
			expected: false,
		},
		{
			name:     "unknown operator returns false",
			reqCtx:   map[domain.RuleAttribute]any{"foo": "bar"},
			cond:     domain.Condition{Attribute: "foo", Operator: "unknown", Value: "bar"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchCondition(tt.reqCtx, tt.cond)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStableHash(t *testing.T) {
	tests := []struct {
		str      string
		expected int
	}{
		{"", 0},
		{"a", StableHash("a")},       // должно быть детерминированно
		{"abc", StableHash("abc")},   // стабильное значение
		{"abc", StableHash("abc")},   // одинаковые входы → одинаковый результат
		{"abcd", StableHash("abcd")}, // другое значение, чем у "abc"
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got := StableHash(tt.str)
			assert.Equal(t, tt.expected, got)
			assert.GreaterOrEqual(t, got, 0, "StableHash must be non-negative")
		})
	}

	t.Run("different strings produce different hashes (usually)", func(t *testing.T) {
		h1 := StableHash("abc")
		h2 := StableHash("xyz")
		assert.NotEqual(t, h1, h2)
	})
}

func TestInList(t *testing.T) {
	tests := []struct {
		name            string
		actual          any
		value           any
		caseInsensitive bool
		expected        bool
	}{
		{
			name:     "match in []string, case sensitive",
			actual:   "foo",
			value:    []string{"bar", "foo"},
			expected: true,
		},
		{
			name:     "no match in []string, case sensitive",
			actual:   "FOO",
			value:    []string{"bar", "foo"},
			expected: false,
		},
		{
			name:            "match in []string, case insensitive",
			actual:          "FOO",
			value:           []string{"bar", "foo"},
			caseInsensitive: true,
			expected:        true,
		},
		{
			name:     "match in []any, case sensitive",
			actual:   "123",
			value:    []any{"456", "123"},
			expected: true,
		},
		{
			name:     "no match in []any, case sensitive",
			actual:   "123",
			value:    []any{"456", "789"},
			expected: false,
		},
		{
			name:            "match in []any, case insensitive",
			actual:          "FOO",
			value:           []any{"bar", "foo"},
			caseInsensitive: true,
			expected:        true,
		},
		{
			name:     "unsupported value type returns false",
			actual:   "foo",
			value:    123, // не []any и не []string
			expected: false,
		},
		{
			name:     "numbers compared as strings",
			actual:   42,
			value:    []any{1, 2, 42},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InList(tt.actual, tt.value, tt.caseInsensitive)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestEvaluateExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     domain.BooleanExpression
		reqCtx   map[domain.RuleAttribute]any
		expected bool
	}{
		{
			name: "single condition true",
			expr: domain.BooleanExpression{
				Condition: &domain.Condition{
					Attribute: "country",
					Operator:  domain.OpEq,
					Value:     "RU",
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"country": "RU"},
			expected: true,
		},
		{
			name: "single condition false",
			expr: domain.BooleanExpression{
				Condition: &domain.Condition{
					Attribute: "country",
					Operator:  domain.OpEq,
					Value:     "RU",
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"country": "US"},
			expected: false,
		},
		{
			name: "AND group all true",
			expr: domain.BooleanExpression{
				Group: &domain.ConditionGroup{
					Operator: domain.LogicalOpAND,
					Children: []domain.BooleanExpression{
						{Condition: &domain.Condition{Attribute: "age", Operator: domain.OpGt, Value: 18}},
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "RU"}},
					},
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"age": 25, "country": "RU"},
			expected: true,
		},
		{
			name: "AND group one false",
			expr: domain.BooleanExpression{
				Group: &domain.ConditionGroup{
					Operator: domain.LogicalOpAND,
					Children: []domain.BooleanExpression{
						{Condition: &domain.Condition{Attribute: "age", Operator: domain.OpGt, Value: 18}},
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "RU"}},
					},
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"age": 25, "country": "US"},
			expected: false,
		},
		{
			name: "OR group one true",
			expr: domain.BooleanExpression{
				Group: &domain.ConditionGroup{
					Operator: domain.LogicalOpOR,
					Children: []domain.BooleanExpression{
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "RU"}},
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "BY"}},
					},
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"country": "BY"},
			expected: true,
		},
		{
			name: "OR group all false",
			expr: domain.BooleanExpression{
				Group: &domain.ConditionGroup{
					Operator: domain.LogicalOpOR,
					Children: []domain.BooleanExpression{
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "RU"}},
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "BY"}},
					},
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"country": "US"},
			expected: false,
		},
		{
			name: "AND NOT group left true right false",
			expr: domain.BooleanExpression{
				Group: &domain.ConditionGroup{
					Operator: domain.LogicalOpANDNot,
					Children: []domain.BooleanExpression{
						{Condition: &domain.Condition{Attribute: "age", Operator: domain.OpGt, Value: 18}},
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "US"}},
					},
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"age": 25, "country": "RU"},
			expected: true, // left true && !right
		},
		{
			name: "AND NOT group left true right true",
			expr: domain.BooleanExpression{
				Group: &domain.ConditionGroup{
					Operator: domain.LogicalOpANDNot,
					Children: []domain.BooleanExpression{
						{Condition: &domain.Condition{Attribute: "age", Operator: domain.OpGt, Value: 18}},
						{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "RU"}},
					},
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"age": 25, "country": "RU"},
			expected: false, // left true && !right → false
		},
		{
			name: "nested OR inside AND",
			expr: domain.BooleanExpression{
				Group: &domain.ConditionGroup{
					Operator: domain.LogicalOpAND,
					Children: []domain.BooleanExpression{
						{Condition: &domain.Condition{Attribute: "age", Operator: domain.OpGt, Value: 18}},
						{
							Group: &domain.ConditionGroup{
								Operator: domain.LogicalOpOR,
								Children: []domain.BooleanExpression{
									{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "RU"}},
									{Condition: &domain.Condition{Attribute: "country", Operator: domain.OpEq, Value: "BY"}},
								},
							},
						},
					},
				},
			},
			reqCtx:   map[domain.RuleAttribute]any{"age": 30, "country": "BY"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := EvaluateExpression(tt.expr, tt.reqCtx)
			assert.Equal(t, tt.expected, ok)
		})
	}
}

func TestService_NextState(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().UTC()

	tests := []struct {
		name            string
		feature         domain.FeatureExtended
		expectedEnabled bool
		expectedTime    time.Time
		hasNextState    bool
	}{
		{
			name: "feature with no schedules returns zero values",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-1 * time.Hour),
				},
			},
			expectedEnabled: false,
			expectedTime:    time.Time{},
			hasNextState:    false,
		},
		{
			name: "schedule without cron, returns end time",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched1",
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(-1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-1 * time.Hour),
					},
				},
			},
			expectedEnabled: false, // After enable schedule ends, feature becomes inactive
			expectedTime:    now.Add(1 * time.Hour),
			hasNextState:    true,
		},
		{
			name: "schedule with cron, returns next cron trigger",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched2",
						Action:    domain.FeatureScheduleActionEnable,
						CronExpr:  ptrString("0 14 * * *"), // 2 PM daily
						Timezone:  "UTC",
						CreatedAt: now.Add(-1 * time.Hour),
					},
				},
			},
			expectedEnabled: true,
			expectedTime:    time.Date(2025, 9, 20, 14, 0, 0, 0, loc), // The next occurrence is in 4 days
			hasNextState:    true,
		},
		{
			name: "schedule with cron and duration, returns action end time",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:           "sched3",
						Action:       domain.FeatureScheduleActionEnable,
						CronExpr:     ptrString("0 14 * * *"), // 2 PM daily
						CronDuration: ptrDuration(2 * time.Hour),
						Timezone:     "UTC",
						CreatedAt:    now.Add(-1 * time.Hour),
					},
				},
			},
			expectedEnabled: false,                                    // opposite action after duration
			expectedTime:    time.Date(2025, 9, 20, 16, 0, 0, 0, loc), // Next occurrence + 2 hours
			hasNextState:    true,
		},
		{
			name: "schedule not yet started, returns start time",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched6",
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(1 * time.Hour)),
						EndsAt:    ptrTime(now.Add(3 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-1 * time.Hour),
					},
				},
			},
			expectedEnabled: true,
			expectedTime:    now.Add(1 * time.Hour),
			hasNextState:    true,
		},
		{
			name: "schedule already ended, returns zero values",
			feature: domain.FeatureExtended{
				Feature: domain.Feature{
					Enabled:   true,
					CreatedAt: now.Add(-2 * time.Hour),
				},
				Schedules: []domain.FeatureSchedule{
					{
						ID:        "sched7",
						Action:    domain.FeatureScheduleActionEnable,
						StartsAt:  ptrTime(now.Add(-3 * time.Hour)),
						EndsAt:    ptrTime(now.Add(-1 * time.Hour)),
						Timezone:  "UTC",
						CreatedAt: now.Add(-2 * time.Hour),
					},
				},
			},
			expectedEnabled: false,
			expectedTime:    time.Time{},
			hasNextState:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(nil, nil, nil, 0)

			enabled, timestamp := svc.NextStateAt(tt.feature, now)

			if timestamp.IsZero() && tt.hasNextState {
				t.Logf("Expected non-zero timestamp but got zero. Feature: %+v", tt.feature)
			}

			assert.Equal(t, tt.expectedEnabled, enabled)
			if tt.hasNextState {
				assert.Equal(t, tt.expectedTime, timestamp)
			} else {
				assert.True(t, timestamp.IsZero())
			}
		})
	}
}

func TestIsScheduleActive_SimpleCronBuilderCases(t *testing.T) {
	loc := time.UTC
	createdAt := time.Date(2025, 1, 1, 0, 0, 0, 0, loc)

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	cases := []struct {
		name         string
		expr         string
		duration     *time.Duration
		checkTime    time.Time
		expectActive bool
		expectAction domain.FeatureScheduleAction
	}{
		{
			name:         "repeat every 15 minutes",
			expr:         "*/15 * * * *",
			checkTime:    time.Date(2025, 1, 15, 9, 30, 0, 0, loc), // 09:30 divisible by 15
			expectActive: true,
			expectAction: domain.FeatureScheduleActionEnable,
		},
		{
			name:         "daily at 09:30",
			expr:         "30 9 * * *",
			checkTime:    time.Date(2025, 1, 15, 9, 30, 0, 0, loc),
			expectActive: true,
			expectAction: domain.FeatureScheduleActionEnable,
		},
		{
			name:         "monthly on 1st at 10:00",
			expr:         "0 10 1 * *",
			checkTime:    time.Date(2025, 1, 1, 10, 0, 0, 0, loc),
			expectActive: true,
			expectAction: domain.FeatureScheduleActionEnable,
		},
		{
			name:         "yearly on Jan 1 at 00:00",
			expr:         "0 0 1 1 *",
			checkTime:    time.Date(2025, 1, 1, 0, 0, 0, 0, loc),
			expectActive: true,
			expectAction: domain.FeatureScheduleActionEnable,
		},
		{
			name:         "with duration still active",
			expr:         "0 * * * *", // hourly
			duration:     ptrDuration(30 * time.Minute),
			checkTime:    time.Date(2025, 1, 15, 9, 15, 0, 0, loc), // within 30m after 09:00
			expectActive: true,
			expectAction: domain.FeatureScheduleActionEnable,
		},
		{
			name:         "with duration expired",
			expr:         "0 * * * *", // hourly
			duration:     ptrDuration(10 * time.Minute),
			checkTime:    time.Date(2025, 1, 15, 9, 15, 0, 0, loc), // 15m > 10m
			expectActive: true,
			expectAction: domain.FeatureScheduleActionDisable,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			sched, err := parser.Parse(tt.expr)
			assert.NoError(t, err)

			crons := CronsMap{
				"test": sched,
			}

			fs := domain.FeatureSchedule{
				ID:           "test",
				ProjectID:    "proj1",
				FeatureID:    "f1",
				CronExpr:     &tt.expr,
				CronDuration: tt.duration,
				Timezone:     "UTC",
				Action:       domain.FeatureScheduleActionEnable,
				CreatedAt:    createdAt,
			}

			active, action := IsScheduleActive(fs, crons, tt.checkTime, createdAt)
			assert.Equal(t, tt.expectActive, active)
			assert.Equal(t, tt.expectAction, action)
		})
	}
}

func ptrTime(t time.Time) *time.Time             { return &t }
func ptrString(s string) *string                 { return &s }
func ptrDuration(d time.Duration) *time.Duration { return &d }
