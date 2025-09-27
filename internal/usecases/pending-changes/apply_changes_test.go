package pending_changes

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/togglr-project/togglr/internal/domain"
	mockcontract "github.com/togglr-project/togglr/test_mocks/internal_/contract"
)

func TestApplyChanges(t *testing.T) {
	tests := []struct {
		name          string
		pendingChange domain.PendingChange
		setupMocks    func(*mockcontract.MockFeaturesRepository, *mockcontract.MockFeatureParamsRepository, *mockcontract.MockRulesRepository, *mockcontract.MockFlagVariantsRepository, *mockcontract.MockFeatureSchedulesRepository, *mockcontract.MockFeatureTagsRepository)
		expectedError string
	}{
		{
			name: "Feature Update - Name and Description",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeature),
							EntityID: "feature-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"name":        {Old: "Old Feature", New: "New Feature"},
								"description": {Old: "Old Description", New: "New Description"},
							},
						},
					},
				},
				EnvironmentID: 1,
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				featuresRepo.On("GetByID", mock.Anything, domain.FeatureID("feature-123")).Return(domain.BasicFeature{
					ID:          "feature-123",
					ProjectID:   "project-123",
					Name:        "Old Feature",
					Description: "Old Description",
				}, nil)
				featuresRepo.On("Update", mock.Anything, domain.EnvironmentID(1), mock.MatchedBy(func(f domain.BasicFeature) bool {
					return f.Name == "New Feature" && f.Description == "New Description"
				})).Return(domain.BasicFeature{}, nil)
			},
			expectedError: "",
		},
		{
			name: "Feature Delete",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeature),
							EntityID: "feature-123",
							Action:   domain.EntityActionDelete,
							Changes:  map[string]domain.ChangeValue{},
						},
					},
				},
				EnvironmentID: 1,
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				featuresRepo.On("Delete", mock.Anything, domain.EnvironmentID(1), domain.FeatureID("feature-123")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "FeatureParams Update - Enabled and DefaultValue",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureParams),
							EntityID: "feature-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"enabled":       {Old: false, New: true},
								"default_value": {Old: "old_value", New: "new_value"},
							},
						},
					},
				},
				EnvironmentID: 1,
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				featureParamsRepo.On("GetByFeatureWithEnv", mock.Anything, domain.FeatureID("feature-123"), domain.EnvironmentID(1)).Return(domain.FeatureParams{
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Enabled:       false,
					DefaultValue:  "old_value",
				}, nil)
				featuresRepo.On("GetByID", mock.Anything, domain.FeatureID("feature-123")).Return(domain.BasicFeature{
					ID:        "feature-123",
					ProjectID: "project-123",
				}, nil)
				featureParamsRepo.On("Update", mock.Anything, domain.ProjectID("project-123"), mock.MatchedBy(func(p domain.FeatureParams) bool {
					return p.Enabled == true && p.DefaultValue == "new_value"
				})).Return(domain.FeatureParams{}, nil)
			},
			expectedError: "",
		},
		{
			name: "FeatureParams Create - Not Found",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureParams),
							EntityID: "feature-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"enabled": {Old: nil, New: true},
							},
						},
					},
				},
				EnvironmentID: 1,
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				featureParamsRepo.On("GetByFeatureWithEnv", mock.Anything, domain.FeatureID("feature-123"), domain.EnvironmentID(1)).Return(domain.FeatureParams{}, domain.ErrEntityNotFound)
				featuresRepo.On("GetByID", mock.Anything, domain.FeatureID("feature-123")).Return(domain.BasicFeature{
					ID:        "feature-123",
					ProjectID: "project-123",
				}, nil)
				featureParamsRepo.On("Update", mock.Anything, domain.ProjectID("project-123"), mock.Anything).Return(domain.FeatureParams{}, domain.ErrEntityNotFound)
				featureParamsRepo.On("Create", mock.Anything, domain.ProjectID("project-123"), mock.MatchedBy(func(p domain.FeatureParams) bool {
					return p.Enabled == true && p.FeatureID == "feature-123" && p.EnvironmentID == 1
				})).Return(domain.FeatureParams{}, nil)
			},
			expectedError: "",
		},
		{
			name: "Rule Update - All Fields",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityRule),
							EntityID: "rule-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"is_customized":   {Old: false, New: true},
								"action":          {Old: "assign", New: "include"},
								"priority":        {Old: float64(100), New: float64(200)},
								"flag_variant_id": {Old: "variant-1", New: "variant-2"},
								"segment_id":      {Old: "segment-1", New: "segment-2"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				rulesRepo.On("GetByID", mock.Anything, domain.RuleID("rule-123")).Return(domain.Rule{
					ID:           "rule-123",
					IsCustomized: false,
					Action:       domain.RuleActionAssign,
					Priority:     100,
				}, nil)
				rulesRepo.On("Update", mock.Anything, mock.MatchedBy(func(r domain.Rule) bool {
					return r.IsCustomized == true && r.Action == domain.RuleActionInclude && r.Priority == 200
				})).Return(domain.Rule{}, nil)
			},
			expectedError: "",
		},
		{
			name: "Rule Insert - All Fields",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityRule),
							EntityID: "rule-123",
							Action:   domain.EntityActionInsert,
							Changes: map[string]domain.ChangeValue{
								"project_id":      {New: "project-123"},
								"feature_id":      {New: "feature-123"},
								"is_customized":   {New: true},
								"action":          {New: "exclude"},
								"priority":        {New: float64(150)},
								"flag_variant_id": {New: "variant-123"},
								"segment_id":      {New: "segment-123"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				rulesRepo.On("Create", mock.Anything, mock.MatchedBy(func(r domain.Rule) bool {
					return r.ProjectID == "project-123" && r.FeatureID == "feature-123" && r.IsCustomized == true && r.Action == domain.RuleActionExclude && r.Priority == 150
				})).Return(domain.Rule{}, nil)
			},
			expectedError: "",
		},
		{
			name: "Rule Delete",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityRule),
							EntityID: "rule-123",
							Action:   domain.EntityActionDelete,
							Changes:  map[string]domain.ChangeValue{},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				rulesRepo.On("Delete", mock.Anything, domain.RuleID("rule-123")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "FlagVariant Update - Name and RolloutPercent",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFlagVariant),
							EntityID: "variant-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"name":            {Old: "Old Variant", New: "New Variant"},
								"rollout_percent": {Old: float64(50), New: float64(75)},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				flagVariantsRepo.On("GetByID", mock.Anything, domain.FlagVariantID("variant-123")).Return(domain.FlagVariant{
					ID:             "variant-123",
					Name:           "Old Variant",
					RolloutPercent: 50,
				}, nil)
				flagVariantsRepo.On("Update", mock.Anything, mock.MatchedBy(func(v domain.FlagVariant) bool {
					return v.Name == "New Variant" && v.RolloutPercent == 75
				})).Return(domain.FlagVariant{}, nil)
			},
			expectedError: "",
		},
		{
			name: "FlagVariant Insert - All Fields",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFlagVariant),
							EntityID: "variant-123",
							Action:   domain.EntityActionInsert,
							Changes: map[string]domain.ChangeValue{
								"project_id":      {New: "project-123"},
								"feature_id":      {New: "feature-123"},
								"name":            {New: "New Variant"},
								"rollout_percent": {New: float64(80)},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				flagVariantsRepo.On("Create", mock.Anything, mock.MatchedBy(func(v domain.FlagVariant) bool {
					return v.ProjectID == "project-123" && v.FeatureID == "feature-123" && v.Name == "New Variant" && v.RolloutPercent == 80
				})).Return(domain.FlagVariant{}, nil)
			},
			expectedError: "",
		},
		{
			name: "FlagVariant Delete",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFlagVariant),
							EntityID: "variant-123",
							Action:   domain.EntityActionDelete,
							Changes:  map[string]domain.ChangeValue{},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				flagVariantsRepo.On("Delete", mock.Anything, domain.FlagVariantID("variant-123")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "FeatureSchedule Update - All Fields",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureSchedule),
							EntityID: "schedule-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"starts_at":     {Old: "2024-01-01T00:00:00Z", New: "2024-06-01T12:00:00Z"},
								"ends_at":       {Old: "2024-12-31T23:59:59Z", New: "2024-11-30T18:30:00Z"},
								"cron_expr":     {Old: "0 0 * * *", New: "0 12 * * 1-5"},
								"cron_duration": {Old: "24h0m0s", New: "12h0m0s"},
								"timezone":      {Old: "UTC", New: "Europe/Moscow"},
								"action":        {Old: "enable", New: "disable"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				schedulesRepo.On("GetByID", mock.Anything, domain.FeatureScheduleID("schedule-123")).Return(domain.FeatureSchedule{
					ID:           "schedule-123",
					StartsAt:     timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndsAt:       timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
					CronExpr:     stringPtr("0 0 * * *"),
					CronDuration: durationPtr(24 * time.Hour),
					Timezone:     "UTC",
					Action:       domain.FeatureScheduleActionEnable,
				}, nil)
				schedulesRepo.On("Update", mock.Anything, mock.MatchedBy(func(s domain.FeatureSchedule) bool {
					return s.Timezone == "Europe/Moscow" && s.Action == domain.FeatureScheduleActionDisable
				})).Return(domain.FeatureSchedule{}, nil)
			},
			expectedError: "",
		},
		{
			name: "FeatureSchedule Insert - All Fields",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureSchedule),
							EntityID: "schedule-123",
							Action:   domain.EntityActionInsert,
							Changes: map[string]domain.ChangeValue{
								"project_id":     {New: "project-123"},
								"feature_id":     {New: "feature-123"},
								"environment_id": {New: float64(1)},
								"starts_at":      {New: "2024-03-15T09:30:00Z"},
								"ends_at":        {New: "2024-08-15T17:45:00Z"},
								"cron_expr":      {New: "30 9 * * 1-5"},
								"cron_duration":  {New: "8h0m0s"},
								"timezone":       {New: "America/New_York"},
								"action":         {New: "enable"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				schedulesRepo.On("Create", mock.Anything, mock.MatchedBy(func(s domain.FeatureSchedule) bool {
					return s.ProjectID == "project-123" && s.FeatureID == "feature-123" && s.EnvironmentID == 1 && s.Timezone == "America/New_York" && s.Action == domain.FeatureScheduleActionEnable
				})).Return(domain.FeatureSchedule{}, nil)
			},
			expectedError: "",
		},
		{
			name: "FeatureSchedule Delete",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureSchedule),
							EntityID: "schedule-123",
							Action:   domain.EntityActionDelete,
							Changes:  map[string]domain.ChangeValue{},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				schedulesRepo.On("Delete", mock.Anything, domain.FeatureScheduleID("schedule-123")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "FeatureTag Insert",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureTag),
							EntityID: "tag-123",
							Action:   domain.EntityActionInsert,
							Changes: map[string]domain.ChangeValue{
								"feature_id": {New: "feature-123"},
								"tag_id":     {New: "tag-123"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				featureTagsRepo.On("AddFeatureTag", mock.Anything, domain.FeatureID("feature-123"), domain.TagID("tag-123")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "FeatureTag Delete",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureTag),
							EntityID: "tag-123",
							Action:   domain.EntityActionDelete,
							Changes: map[string]domain.ChangeValue{
								"feature_id": {New: "feature-123"},
								"tag_id":     {New: "tag-123"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				featureTagsRepo.On("RemoveFeatureTag", mock.Anything, domain.FeatureID("feature-123"), domain.TagID("tag-123")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Multiple Entities - Mixed Actions",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeature),
							EntityID: "feature-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"name": {Old: "Old Feature", New: "New Feature"},
							},
						},
						{
							Entity:   string(domain.EntityRule),
							EntityID: "rule-123",
							Action:   domain.EntityActionInsert,
							Changes: map[string]domain.ChangeValue{
								"project_id":    {New: "project-123"},
								"feature_id":    {New: "feature-123"},
								"is_customized": {New: true},
								"action":        {New: "include"},
								"priority":      {New: float64(100)},
							},
						},
						{
							Entity:   string(domain.EntityFlagVariant),
							EntityID: "variant-123",
							Action:   domain.EntityActionDelete,
							Changes:  map[string]domain.ChangeValue{},
						},
					},
				},
				EnvironmentID: 1,
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				// Feature update
				featuresRepo.On("GetByID", mock.Anything, domain.FeatureID("feature-123")).Return(domain.BasicFeature{
					ID:        "feature-123",
					ProjectID: "project-123",
					Name:      "Old Feature",
				}, nil)
				featuresRepo.On("Update", mock.Anything, domain.EnvironmentID(1), mock.MatchedBy(func(f domain.BasicFeature) bool {
					return f.Name == "New Feature"
				})).Return(domain.BasicFeature{}, nil)

				// Rule insert
				rulesRepo.On("Create", mock.Anything, mock.MatchedBy(func(r domain.Rule) bool {
					return r.ProjectID == "project-123" && r.FeatureID == "feature-123" && r.IsCustomized == true
				})).Return(domain.Rule{}, nil)

				// FlagVariant delete
				flagVariantsRepo.On("Delete", mock.Anything, domain.FlagVariantID("variant-123")).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Unsupported Entity Type",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   "unsupported_entity",
							EntityID: "entity-123",
							Action:   domain.EntityActionUpdate,
							Changes: map[string]domain.ChangeValue{
								"field": {Old: "old", New: "new"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				// No mocks needed - should fail before calling any repository
			},
			expectedError: "unsupported entity type: unsupported_entity",
		},
		{
			name: "FeatureTag Missing FeatureID",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureTag),
							EntityID: "tag-123",
							Action:   domain.EntityActionInsert,
							Changes: map[string]domain.ChangeValue{
								"tag_id": {New: "tag-123"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				// No mocks needed - should fail before calling repository
			},
			expectedError: "feature_id is required for feature_tag change",
		},
		{
			name: "FeatureTag Missing TagID",
			pendingChange: domain.PendingChange{
				ID: "pending-123",
				Change: domain.PendingChangePayload{
					Entities: []domain.EntityChange{
						{
							Entity:   string(domain.EntityFeatureTag),
							EntityID: "",
							Action:   domain.EntityActionInsert,
							Changes: map[string]domain.ChangeValue{
								"feature_id": {New: "feature-123"},
							},
						},
					},
				},
			},
			setupMocks: func(featuresRepo *mockcontract.MockFeaturesRepository, featureParamsRepo *mockcontract.MockFeatureParamsRepository, rulesRepo *mockcontract.MockRulesRepository, flagVariantsRepo *mockcontract.MockFlagVariantsRepository, schedulesRepo *mockcontract.MockFeatureSchedulesRepository, featureTagsRepo *mockcontract.MockFeatureTagsRepository) {
				// No mocks needed - should fail before calling repository
			},
			expectedError: "tag_id is required for feature_tag change",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockFeaturesRepo := mockcontract.NewMockFeaturesRepository(t)
			mockFeatureParamsRepo := mockcontract.NewMockFeatureParamsRepository(t)
			mockRulesRepo := mockcontract.NewMockRulesRepository(t)
			mockFlagVariantsRepo := mockcontract.NewMockFlagVariantsRepository(t)
			mockSchedulesRepo := mockcontract.NewMockFeatureSchedulesRepository(t)
			mockFeatureTagsRepo := mockcontract.NewMockFeatureTagsRepository(t)

			// Setup mocks
			tt.setupMocks(mockFeaturesRepo, mockFeatureParamsRepo, mockRulesRepo, mockFlagVariantsRepo, mockSchedulesRepo, mockFeatureTagsRepo)

			// Create service
			service := &Service{
				featuresRepo:      mockFeaturesRepo,
				featureParamsRepo: mockFeatureParamsRepo,
				rulesRepo:         mockRulesRepo,
				flagVariantsRepo:  mockFlagVariantsRepo,
				schedulesRepo:     mockSchedulesRepo,
				featureTagsRepo:   mockFeatureTagsRepo,
			}

			// Execute
			err := service.applyChanges(context.Background(), tt.pendingChange)

			// Assert
			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

// Helper functions for creating pointers
func timePtr(t time.Time) *time.Time {
	return &t
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}

func stringPtr(s string) *string {
	return &s
}
