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

func TestFeatureTagAnonymousStruct(t *testing.T) {
	// Create mocks
	mockGuardService := mockcontract.NewMockGuardService(t)
	mockPendingUseCase := mockcontract.NewMockPendingChangesUseCase(t)

	// Setup mocks
	mockGuardService.On("IsFeatureGuarded", mock.Anything, domain.FeatureID("feature-123")).Return(true, nil)
	mockPendingUseCase.On("GetProjectActiveUserCount", mock.Anything, domain.ProjectID("project-123")).Return(2, nil)
	mockPendingUseCase.On("CheckEntityConflict", mock.Anything, mock.AnythingOfType("[]domain.EntityChange")).Return(false, nil)
	mockPendingUseCase.On("Create", mock.Anything, domain.ProjectID("project-123"), domain.EnvironmentID(1), "testuser", mock.AnythingOfType("*int"), mock.AnythingOfType("domain.PendingChangePayload")).Return(domain.PendingChange{
		ID: "pending-feature-tag-delete-456",
		Change: domain.PendingChangePayload{
			Entities: []domain.EntityChange{
				{
					Entity:   string(domain.EntityFeatureTag),
					EntityID: "tag-456",
					Action:   domain.EntityActionDelete,
					Changes: map[string]domain.ChangeValue{
						"feature_id": {New: "feature-123"},
						"tag_id":     {New: "tag-456"},
					},
				},
			},
			Meta: domain.PendingChangeMeta{
				Reason: "Remove tag from feature via API",
				Client: "ui",
				Origin: "feature-tag-remove",
			},
		},
	}, nil)

	// Create service
	service := New(mockGuardService, mockPendingUseCase)

	// Create context with user info
	ctx := context.Background()
	ctx = appcontext.WithUserID(ctx, domain.UserID(1))
	ctx = appcontext.WithUsername(ctx, "testuser")

	// Test FeatureTag delete with anonymous struct
	request := contract.GuardRequest{
		ProjectID:     "project-123",
		EnvironmentID: 1,
		FeatureID:     "feature-123",
		Reason:        "Remove tag from feature via API",
		Origin:        "feature-tag-remove",
		Action:        domain.EntityActionDelete,
		OldEntity: struct {
			FeatureID string
			TagID     string
		}{
			FeatureID: "feature-123",
			TagID:     "tag-456",
		},
		NewEntity: nil,
	}

	// Execute
	pc, conflict, proceed, err := service.CheckGuardedOperation(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, pc)
	assert.Equal(t, domain.PendingChangeID("pending-feature-tag-delete-456"), pc.ID)
	assert.False(t, conflict)
	assert.False(t, proceed)

	// Verify FeatureTag changes are captured
	changes := pc.Change.Entities[0].Changes
	assert.Contains(t, changes, "feature_id")
	assert.Contains(t, changes, "tag_id")
	assert.Equal(t, "feature-123", changes["feature_id"].New)
	assert.Equal(t, "tag-456", changes["tag_id"].New)
}
