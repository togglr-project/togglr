package pending_changes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/togglr-project/togglr/internal/domain"
)

func TestApplyChangesToEntity(t *testing.T) {
	// Test applying string changes to BasicFeature
	entity := &domain.BasicFeature{
		ID:          "feature-1",
		Name:        "Old Name",
		Description: "Old Description",
		RolloutKey:  "user_id",
	}

	changes := map[string]domain.ChangeValue{
		"name":        {New: "New Name"},
		"description": {New: "New Description"},
	}

	err := ApplyChangesToEntity(entity, changes)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", entity.Name)
	assert.Equal(t, "New Description", entity.Description)
}
