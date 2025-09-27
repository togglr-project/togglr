package guard_engine

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	appcontext "github.com/togglr-project/togglr/internal/context"
	"github.com/togglr-project/togglr/internal/contract"
	"github.com/togglr-project/togglr/internal/domain"
	mockcontract "github.com/togglr-project/togglr/test_mocks/internal_/contract"
)

func TestCheckGuardedOperation_Simple1(t *testing.T) {
	tests := []struct {
		name           string
		request        contract.GuardRequest
		setupMocks     func(*mockcontract.MockGuardService, *mockcontract.MockPendingChangesUseCase)
		expectedResult func(*testing.T, *domain.PendingChange, bool, bool, error)
	}{
		{
			name: "Feature Update - Not Guarded - Should Proceed",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Feature update via API",
				Origin:        "feature-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "Old Feature",
					},
					Enabled: false,
				},
				NewEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "New Feature",
					},
					Enabled: true,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(false, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.Nil(t, pc)
				assert.False(t, conflict)
				assert.True(t, proceed)
			},
		},
		{
			name: "Feature Update - Guarded - Should Create Pending Change",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Feature update via API",
				Origin:        "feature-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "Old Feature",
					},
					Enabled: false,
				},
				NewEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "New Feature",
					},
					Enabled: true,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-123",
					Change: domain.PendingChangePayload{
						Entities: []domain.EntityChange{
							{
								Entity:   string(domain.EntityFeature),
								EntityID: "feature-123",
								Action:   domain.EntityActionUpdate,
								Changes: map[string]domain.ChangeValue{
									"name":    {Old: "Old Feature", New: "New Feature"},
									"enabled": {Old: false, New: true},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Feature update via API",
							Client: "ui",
							Origin: "feature-update",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-123"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)
			},
		},
		{
			name: "FeatureParams Update - Value to Pointer - Should Work",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Feature params update via API",
				Origin:        "feature-params-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: domain.FeatureParams{
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Enabled:       false,
					DefaultValue:  "old_value",
				},
				NewEntity: &domain.FeatureParams{
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Enabled:       true,
					DefaultValue:  "new_value",
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-params-123",
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
						Meta: domain.PendingChangeMeta{
							Reason: "Feature params update via API",
							Client: "ui",
							Origin: "feature-params-update",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-params-123"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)
			},
		},
		{
			name: "No Changes Detected - Should Proceed",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Feature update via API",
				Origin:        "feature-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "Same Feature",
					},
					Enabled: true,
				},
				NewEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "Same Feature",
					},
					Enabled: true,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				// No mocks needed - should return proceed=true before checking guards
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.Nil(t, pc)
				assert.False(t, conflict)
				assert.True(t, proceed)
			},
		},
		{
			name: "Rule Update - Pointer to Value - Should Work",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Rule update via API",
				Origin:        "rule-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.Rule{
					ID:           "rule-123",
					ProjectID:    "project-123",
					FeatureID:    "feature-123",
					IsCustomized: false,
					Action:       domain.RuleActionAssign,
					Priority:     100,
				},
				NewEntity: domain.Rule{
					ID:           "rule-123",
					ProjectID:    "project-123",
					FeatureID:    "feature-123",
					IsCustomized: true,
					Action:       domain.RuleActionInclude,
					Priority:     200,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-rule-123",
					Change: domain.PendingChangePayload{
						Entities: []domain.EntityChange{
							{
								Entity:   string(domain.EntityRule),
								EntityID: "rule-123",
								Action:   domain.EntityActionUpdate,
								Changes: map[string]domain.ChangeValue{
									"is_customized": {Old: false, New: true},
									"action":        {Old: "assign", New: "include"},
									"priority":      {Old: uint8(100), New: uint8(200)},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Rule update via API",
							Client: "ui",
							Origin: "rule-update",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-rule-123"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)
			},
		},
		{
			name: "FlagVariant Insert - Value - Should Work",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Flag variant create via API",
				Origin:        "flag-variant-create",
				Action:        domain.EntityActionInsert,
				OldEntity:     nil,
				NewEntity: domain.FlagVariant{
					ID:             "variant-123",
					ProjectID:      "project-123",
					FeatureID:      "feature-123",
					Name:           "New Variant",
					RolloutPercent: 75,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-variant-123",
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
									"rollout_percent": {New: uint8(75)},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Flag variant create via API",
							Client: "ui",
							Origin: "flag-variant-create",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-variant-123"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)
			},
		},
		{
			name: "FeatureSchedule Update - Value to Pointer - Should Work",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Schedule update via API",
				Origin:        "schedule-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: domain.FeatureSchedule{
					ID:            "schedule-123",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
				},
				NewEntity: &domain.FeatureSchedule{
					ID:            "schedule-123",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "Europe/Moscow",
					Action:        domain.FeatureScheduleActionDisable,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-schedule-123",
					Change: domain.PendingChangePayload{
						Entities: []domain.EntityChange{
							{
								Entity:   string(domain.EntityFeatureSchedule),
								EntityID: "schedule-123",
								Action:   domain.EntityActionUpdate,
								Changes: map[string]domain.ChangeValue{
									"timezone": {Old: "UTC", New: "Europe/Moscow"},
									"action":   {Old: "enable", New: "disable"},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Schedule update via API",
							Client: "ui",
							Origin: "schedule-update",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-schedule-123"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)
			},
		},
		{
			name: "Conflict Detected - Should Return Conflict",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Feature update via API",
				Origin:        "feature-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "Test Feature",
					},
					Enabled: false,
				},
				NewEntity: &domain.Feature{
					BasicFeature: domain.BasicFeature{
						ID:   "feature-123",
						Name: "Updated Feature",
					},
					Enabled: true,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(true, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.Nil(t, pc)
				assert.True(t, conflict)
				assert.False(t, proceed)
			},
		},
		{
			name: "Unknown Entity Type - Should Error",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Unknown entity via API",
				Origin:        "unknown-entity",
				Action:        domain.EntityActionUpdate,
				OldEntity:     "invalid-entity",
				NewEntity:     "invalid-entity",
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				// No mocks needed - should fail before checking guards
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unknown entity type")
				assert.Nil(t, pc)
				assert.False(t, conflict)
				assert.False(t, proceed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockGuardService := mockcontract.NewMockGuardService(t)
			mockPendingUseCase := mockcontract.NewMockPendingChangesUseCase(t)

			// Setup mocks
			tt.setupMocks(mockGuardService, mockPendingUseCase)

			// Create service
			service := New(mockGuardService, mockPendingUseCase)

			// Create context with user info
			ctx := context.Background()
			ctx = appcontext.WithUserID(ctx, domain.UserID(1))
			ctx = appcontext.WithUsername(ctx, "testuser")

			// Execute
			pc, conflict, proceed, err := service.CheckGuardedOperation(ctx, tt.request)

			// Assert
			tt.expectedResult(t, pc, conflict, proceed, err)
		})
	}
}

func TestCheckGuardedOperation_FeatureSchedule_Comprehensive(t *testing.T) {
	tests := []struct {
		name           string
		request        contract.GuardRequest
		setupMocks     func(*mockcontract.MockGuardService, *mockcontract.MockPendingChangesUseCase)
		expectedResult func(*testing.T, *domain.PendingChange, bool, bool, error)
	}{
		{
			name: "FeatureSchedule Update - All Fields Changed",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Complete schedule update via API",
				Origin:        "schedule-complete-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.FeatureSchedule{
					ID:            "schedule-123",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndsAt:        timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
					CronDuration:  durationPtr(24 * time.Hour),
					CronExpr:      stringPtr("0 0 * * *"),
				},
				NewEntity: &domain.FeatureSchedule{
					ID:            "schedule-123",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 2,                                                        // Changed
					Timezone:      "Europe/Moscow",                                          // Changed
					Action:        domain.FeatureScheduleActionDisable,                      // Changed
					StartsAt:      timePtr(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)),    // Changed
					EndsAt:        timePtr(time.Date(2024, 11, 30, 18, 30, 0, 0, time.UTC)), // Changed
					CronDuration:  durationPtr(12 * time.Hour),                              // Changed
					// IsRecurring:    true,                                                     // Changed
					CronExpr: stringPtr("0 12 * * 1-5"), // Changed
					// IsActive:       false,                                                    // Changed
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-schedule-complete-123",
					Change: domain.PendingChangePayload{
						Entities: []domain.EntityChange{
							{
								Entity:   string(domain.EntityFeatureSchedule),
								EntityID: "schedule-123",
								Action:   domain.EntityActionUpdate,
								Changes: map[string]domain.ChangeValue{
									"environment_id": {Old: domain.EnvironmentID(1), New: domain.EnvironmentID(2)},
									"timezone":       {Old: "UTC", New: "Europe/Moscow"},
									"action":         {Old: "enable", New: "disable"},
									"starts_at":      {Old: "2024-01-01T00:00:00Z", New: "2024-06-01T12:00:00Z"},
									"ends_at":        {Old: "2024-12-31T23:59:59Z", New: "2024-11-30T18:30:00Z"},
									"cron_duration":  {Old: "24h0m0s", New: "12h0m0s"},
									"cron_expr":      {Old: "0 0 * * *", New: "0 12 * * 1-5"},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Complete schedule update via API",
							Client: "ui",
							Origin: "schedule-complete-update",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-schedule-complete-123"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)

				// Verify all fields are captured in changes
				changes := pc.Change.Entities[0].Changes
				assert.Contains(t, changes, "environment_id")
				assert.Contains(t, changes, "timezone")
				assert.Contains(t, changes, "action")
				assert.Contains(t, changes, "starts_at")
				assert.Contains(t, changes, "ends_at")
				assert.Contains(t, changes, "cron_duration")
				// assert.Contains(t, changes, "is_recurring")
				assert.Contains(t, changes, "cron_expr")
				// assert.Contains(t, changes, "is_active")
			},
		},
		{
			name: "FeatureSchedule Insert - All Fields Set",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Create new schedule via API",
				Origin:        "schedule-create",
				Action:        domain.EntityActionInsert,
				OldEntity:     nil,
				NewEntity: &domain.FeatureSchedule{
					ID:            "schedule-new-456",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "America/New_York",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      timePtr(time.Date(2024, 3, 15, 9, 30, 0, 0, time.UTC)),
					EndsAt:        timePtr(time.Date(2024, 8, 15, 17, 45, 0, 0, time.UTC)),
					CronDuration:  durationPtr(8 * time.Hour),
					// IsRecurring:    true,
					CronExpr: stringPtr("30 9 * * 1-5"),
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-schedule-insert-456",
					Change: domain.PendingChangePayload{
						Entities: []domain.EntityChange{
							{
								Entity:   string(domain.EntityFeatureSchedule),
								EntityID: "schedule-new-456",
								Action:   domain.EntityActionInsert,
								Changes: map[string]domain.ChangeValue{
									"project_id":     {New: "project-123"},
									"feature_id":     {New: "feature-123"},
									"environment_id": {New: domain.EnvironmentID(1)},
									"timezone":       {New: "America/New_York"},
									"action":         {New: "enable"},
									"starts_at":      {New: "2024-03-15T09:30:00Z"},
									"ends_at":        {New: "2024-08-15T17:45:00Z"},
									"cron_duration":  {New: "8h0m0s"},
									// "is_recurring":   {New: true},
									"cron_expr": {New: "30 9 * * 1-5"},
									// "is_active": {New: true},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Create new schedule via API",
							Client: "ui",
							Origin: "schedule-create",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-schedule-insert-456"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)

				// Verify all fields are captured in changes
				changes := pc.Change.Entities[0].Changes
				assert.Contains(t, changes, "project_id")
				assert.Contains(t, changes, "feature_id")
				assert.Contains(t, changes, "environment_id")
				assert.Contains(t, changes, "timezone")
				assert.Contains(t, changes, "action")
				assert.Contains(t, changes, "starts_at")
				assert.Contains(t, changes, "ends_at")
				assert.Contains(t, changes, "cron_duration")
				// assert.Contains(t, changes, "is_recurring")
				assert.Contains(t, changes, "cron_expr")
				// assert.Contains(t, changes, "is_active")
			},
		},
		{
			name: "FeatureSchedule Update - Only Time Fields Changed",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Update schedule times via API",
				Origin:        "schedule-time-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.FeatureSchedule{
					ID:            "schedule-789",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndsAt:        timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
					CronDuration:  durationPtr(24 * time.Hour),
					CronExpr:      nil,
				},
				NewEntity: &domain.FeatureSchedule{
					ID:            "schedule-789",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      timePtr(time.Date(2024, 2, 1, 8, 0, 0, 0, time.UTC)),    // Changed
					EndsAt:        timePtr(time.Date(2024, 11, 30, 16, 0, 0, 0, time.UTC)), // Changed
					CronDuration:  durationPtr(12 * time.Hour),                             // Changed
					CronExpr:      nil,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-schedule-time-789",
					Change: domain.PendingChangePayload{
						Entities: []domain.EntityChange{
							{
								Entity:   string(domain.EntityFeatureSchedule),
								EntityID: "schedule-789",
								Action:   domain.EntityActionUpdate,
								Changes: map[string]domain.ChangeValue{
									"starts_at":     {Old: "2024-01-01T00:00:00Z", New: "2024-02-01T08:00:00Z"},
									"ends_at":       {Old: "2024-12-31T23:59:59Z", New: "2024-11-30T16:00:00Z"},
									"cron_duration": {Old: "24h0m0s", New: "12h0m0s"},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Update schedule times via API",
							Client: "ui",
							Origin: "schedule-time-update",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-schedule-time-789"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)

				// Verify only time-related fields are captured
				changes := pc.Change.Entities[0].Changes
				assert.Contains(t, changes, "starts_at")
				assert.Contains(t, changes, "ends_at")
				assert.Contains(t, changes, "cron_duration")
				assert.NotContains(t, changes, "timezone")
				assert.NotContains(t, changes, "action")
				// assert.NotContains(t, changes, "is_recurring")
				assert.NotContains(t, changes, "cron_expr")
				// assert.NotContains(t, changes, "is_active")
			},
		},
		{
			name: "FeatureSchedule Update - No Changes - Should Proceed",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "No changes schedule update via API",
				Origin:        "schedule-no-changes",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.FeatureSchedule{
					ID:            "schedule-101",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      nil,
					EndsAt:        nil,
					CronDuration:  nil,
					CronExpr:      nil,
				},
				NewEntity: &domain.FeatureSchedule{
					ID:            "schedule-101",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      nil,
					EndsAt:        nil,
					CronDuration:  nil,
					CronExpr:      nil,
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				// No mocks needed - should return proceed=true before checking guards
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.Nil(t, pc)
				assert.False(t, conflict)
				assert.True(t, proceed)
			},
		},
		{
			name: "FeatureSchedule Update - Only String Fields Changed",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "Update schedule strings via API",
				Origin:        "schedule-strings-update",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.FeatureSchedule{
					ID:            "schedule-202",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      nil,
					EndsAt:        nil,
					CronDuration:  nil,
					CronExpr:      stringPtr("0 0 * * *"),
				},
				NewEntity: &domain.FeatureSchedule{
					ID:            "schedule-202",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "Asia/Tokyo",                        // Changed
					Action:        domain.FeatureScheduleActionDisable, // Changed
					StartsAt:      nil,
					EndsAt:        nil,
					CronDuration:  nil,
					CronExpr:      stringPtr("0 9 * * 1-5"), // Changed
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				guardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
				pendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
				pendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
				pendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
					ID: "pending-schedule-strings-202",
					Change: domain.PendingChangePayload{
						Entities: []domain.EntityChange{
							{
								Entity:   string(domain.EntityFeatureSchedule),
								EntityID: "schedule-202",
								Action:   domain.EntityActionUpdate,
								Changes: map[string]domain.ChangeValue{
									"timezone":  {Old: "UTC", New: "Asia/Tokyo"},
									"action":    {Old: "enable", New: "disable"},
									"cron_expr": {Old: "0 0 * * *", New: "0 9 * * 1-5"},
								},
							},
						},
						Meta: domain.PendingChangeMeta{
							Reason: "Update schedule strings via API",
							Client: "ui",
							Origin: "schedule-strings-update",
						},
					},
				}, nil)
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, pc)
				assert.Equal(t, domain.PendingChangeID("pending-schedule-strings-202"), pc.ID)
				assert.False(t, conflict)
				assert.False(t, proceed)

				// Verify only string fields are captured
				changes := pc.Change.Entities[0].Changes
				assert.Contains(t, changes, "timezone")
				assert.Contains(t, changes, "action")
				assert.Contains(t, changes, "cron_expr")
				// assert.NotContains(t, changes, "is_recurring")
				// assert.NotContains(t, changes, "is_active")
				assert.NotContains(t, changes, "starts_at")
				assert.NotContains(t, changes, "ends_at")
				assert.NotContains(t, changes, "cron_duration")
			},
		},
		{
			name: "FeatureSchedule Update - No Changes - Should Proceed",
			request: contract.GuardRequest{
				ProjectID:     "project-123",
				EnvironmentID: 1,
				FeatureID:     "feature-123",
				Reason:        "No changes schedule update via API",
				Origin:        "schedule-no-changes",
				Action:        domain.EntityActionUpdate,
				OldEntity: &domain.FeatureSchedule{
					ID:            "schedule-303",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndsAt:        timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
					CronDuration:  durationPtr(24 * time.Hour),
					// IsRecurring:    true,
					CronExpr: stringPtr("0 0 * * *"),
				},
				NewEntity: &domain.FeatureSchedule{
					ID:            "schedule-303",
					ProjectID:     "project-123",
					FeatureID:     "feature-123",
					EnvironmentID: 1,
					Timezone:      "UTC",
					Action:        domain.FeatureScheduleActionEnable,
					StartsAt:      timePtr(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndsAt:        timePtr(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
					CronDuration:  durationPtr(24 * time.Hour),
					// IsRecurring:    true,
					CronExpr: stringPtr("0 0 * * *"),
				},
			},
			setupMocks: func(guardService *mockcontract.MockGuardService, pendingUseCase *mockcontract.MockPendingChangesUseCase) {
				// No mocks needed - should return proceed=true before checking guards
			},
			expectedResult: func(t *testing.T, pc *domain.PendingChange, conflict, proceed bool, err error) {
				assert.NoError(t, err)
				assert.Nil(t, pc)
				assert.False(t, conflict)
				assert.True(t, proceed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockGuardService := mockcontract.NewMockGuardService(t)
			mockPendingUseCase := mockcontract.NewMockPendingChangesUseCase(t)

			// Setup mocks
			tt.setupMocks(mockGuardService, mockPendingUseCase)

			// Create service
			service := New(mockGuardService, mockPendingUseCase)

			// Create context with user info
			ctx := context.Background()
			ctx = appcontext.WithUserID(ctx, domain.UserID(1))
			ctx = appcontext.WithUsername(ctx, "testuser")

			// Execute
			pc, conflict, proceed, err := service.CheckGuardedOperation(ctx, tt.request)

			// Assert
			tt.expectedResult(t, pc, conflict, proceed, err)
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
