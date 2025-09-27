package guard_engine

import (
	"context"
	"testing"

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
