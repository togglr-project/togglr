package guard_engine

import (
	"testing"
	"time"

	"github.com/togglr-project/togglr/internal/domain"
)

func TestBuildChangeDiff_FeatureSchedule(t *testing.T) {
	// Test for FeatureSchedule - check all editable fields
	oldSchedule := &domain.FeatureSchedule{
		ID:            "test-id",
		ProjectID:     "project-1",
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		Timezone:      "UTC",
		Action:        domain.FeatureScheduleActionEnable,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	newSchedule := &domain.FeatureSchedule{
		ID:            "test-id",
		ProjectID:     "project-1",
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		Timezone:      "Europe/Moscow",                     // Changed timezone
		Action:        domain.FeatureScheduleActionDisable, // Changed action
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	changes := BuildChangeDiff(oldSchedule, newSchedule)

	// Check that changes contain only editable fields
	expectedChanges := map[string]bool{
		"timezone": true,
		"action":   true,
	}

	if len(changes) != len(expectedChanges) {
		t.Errorf("Expected %d changes, got %d", len(expectedChanges), len(changes))
	}

	for field := range expectedChanges {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected change for field %s", field)
		}
	}

	// Check that non-editable fields are not included
	if _, exists := changes["id"]; exists {
		t.Error("ID field should not be included in changes")
	}
	if _, exists := changes["created_at"]; exists {
		t.Error("CreatedAt field should not be included in changes")
	}
}

func TestBuildInsertChanges_FeatureSchedule(t *testing.T) {
	// Test for FeatureSchedule insert
	schedule := &domain.FeatureSchedule{
		ID:            "test-id",
		ProjectID:     "project-1",
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		Timezone:      "UTC",
		Action:        domain.FeatureScheduleActionEnable,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	changes := BuildInsertChanges(schedule)

	// Check that all necessary fields are included
	expectedFields := map[string]bool{
		"id":             true, // pk field
		"project_id":     true, // required field
		"feature_id":     true, // required field
		"environment_id": true, // required field
		"timezone":       true, // editable field
		"action":         true, // editable field
	}

	for field := range expectedFields {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected field %s in insert changes", field)
		}
	}

	// Check that non-editable fields are not included
	if _, exists := changes["created_at"]; exists {
		t.Error("CreatedAt field should not be included in insert changes")
	}
	if _, exists := changes["updated_at"]; exists {
		t.Error("UpdatedAt field should not be included in insert changes")
	}
}

func TestBuildChangeDiff_Feature(t *testing.T) {
	// Test for Feature - full structure with nested BasicFeature
	oldFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Test Feature",
			Description: "Test Description",
			RolloutKey:  "user_id",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      true,      // This field doesn't have editable tag
		DefaultValue: "default", // This field doesn't have editable tag
	}

	newFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Updated Test Feature",     // Changed name
			Description: "Updated Test Description", // Changed description
			RolloutKey:  "session_id",               // Changed rollout_key
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      false,         // Changed enabled (but this field is not editable)
		DefaultValue: "new-default", // Changed default_value (but this field is not editable)
	}

	changes := BuildChangeDiff(oldFeature, newFeature)

	// Check that changes contain only editable fields from BasicFeature
	expectedChanges := map[string]bool{
		"name":        true, // from BasicFeature
		"description": true, // from BasicFeature
		"rollout_key": true, // from BasicFeature
	}

	if len(changes) != len(expectedChanges) {
		t.Errorf("Expected %d changes, got %d", len(expectedChanges), len(changes))
	}

	for field := range expectedChanges {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected change for field %s", field)
		}
	}

	// Check that non-editable fields are not included
	if _, exists := changes["id"]; exists {
		t.Error("ID field should not be included in changes")
	}
	if _, exists := changes["key"]; exists {
		t.Error("Key field should not be included in changes")
	}
	if _, exists := changes["created_at"]; exists {
		t.Error("CreatedAt field should not be included in changes")
	}
	// enabled and default_value should not be included, since they don't have the editable tag
	if _, exists := changes["enabled"]; exists {
		t.Error("Enabled field should not be included in changes (no editable tag)")
	}
	if _, exists := changes["default_value"]; exists {
		t.Error("DefaultValue field should not be included in changes (no editable tag)")
	}
}

func TestBuildChangeDiff_FeatureParams(t *testing.T) {
	// Test for FeatureParams - here enabled and default_value have the editable tag
	oldParams := &domain.FeatureParams{
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		Enabled:       true,
		DefaultValue:  "default",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	newParams := &domain.FeatureParams{
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		Enabled:       false,         // Changed enabled
		DefaultValue:  "new-default", // Changed default_value
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	changes := BuildChangeDiff(oldParams, newParams)

	// Check that changes contain only editable fields
	expectedChanges := map[string]bool{
		"enabled":       true,
		"default_value": true,
	}

	if len(changes) != len(expectedChanges) {
		t.Errorf("Expected %d changes, got %d", len(expectedChanges), len(changes))
	}

	for field := range expectedChanges {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected change for field %s", field)
		}
	}

	// Check that non-editable fields are not included
	if _, exists := changes["feature_id"]; exists {
		t.Error("FeatureID field should not be included in changes (no editable tag)")
	}
	if _, exists := changes["environment_id"]; exists {
		t.Error("EnvironmentID field should not be included in changes (no editable tag)")
	}
	if _, exists := changes["created_at"]; exists {
		t.Error("CreatedAt field should not be included in changes")
	}
}

func TestBuildInsertChanges_FeatureParams(t *testing.T) {
	// Test for FeatureParams insert
	params := &domain.FeatureParams{
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		Enabled:       true,
		DefaultValue:  "default",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	changes := BuildInsertChanges(params)

	// Check that all necessary fields are included
	expectedFields := map[string]bool{
		"feature_id":     true, // pk field
		"environment_id": true, // required field
		"enabled":        true, // editable field
		"default_value":  true, // editable field
	}

	for field := range expectedFields {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected field %s in insert changes", field)
		}
	}

	// Check that non-editable fields are not included
	if _, exists := changes["created_at"]; exists {
		t.Error("CreatedAt field should not be included in insert changes")
	}
	if _, exists := changes["updated_at"]; exists {
		t.Error("UpdatedAt field should not be included in insert changes")
	}
}

func TestBuildChangeDiff_Rule(t *testing.T) {
	// Test for Rule
	oldRule := &domain.Rule{
		ID:            "rule-1",
		ProjectID:     "project-1",
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		IsCustomized:  false,
		Action:        domain.RuleActionAssign,
		Priority:      1,
		CreatedAt:     time.Now(),
	}

	newRule := &domain.Rule{
		ID:            "rule-1",
		ProjectID:     "project-1",
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		IsCustomized:  true,                     // Changed
		Action:        domain.RuleActionInclude, // Changed
		Priority:      2,                        // Changed
		CreatedAt:     time.Now(),
	}

	changes := BuildChangeDiff(oldRule, newRule)

	// Check that changes contain only editable fields
	expectedChanges := map[string]bool{
		"is_customized": true,
		"action":        true,
		"priority":      true,
	}

	if len(changes) != len(expectedChanges) {
		t.Errorf("Expected %d changes, got %d", len(expectedChanges), len(changes))
	}

	for field := range expectedChanges {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected change for field %s", field)
		}
	}

	// Check that non-editable fields are not included
	if _, exists := changes["id"]; exists {
		t.Error("ID field should not be included in changes")
	}
	if _, exists := changes["project_id"]; exists {
		t.Error("ProjectID field should not be included in changes")
	}
}

func TestBuildInsertChanges_Rule(t *testing.T) {
	// Test for Rule insert
	rule := &domain.Rule{
		ID:            "rule-1",
		ProjectID:     "project-1",
		FeatureID:     "feature-1",
		EnvironmentID: 1,
		IsCustomized:  false,
		Action:        domain.RuleActionAssign,
		Priority:      1,
		CreatedAt:     time.Now(),
	}

	changes := BuildInsertChanges(rule)

	// Check that all necessary fields are included
	expectedFields := map[string]bool{
		"id":             true, // pk field
		"project_id":     true, // required field
		"feature_id":     true, // required field
		"environment_id": true, // required field
		"is_customized":  true, // editable field
		"action":         true, // editable field
		"priority":       true, // editable field
	}

	for field := range expectedFields {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected field %s in insert changes", field)
		}
	}
}

func TestBuildChangeDiff_FlagVariant(t *testing.T) {
	// Test for FlagVariant
	oldVariant := &domain.FlagVariant{
		ID:             "variant-1",
		ProjectID:      "project-1",
		FeatureID:      "feature-1",
		EnvironmentID:  1,
		Name:           "Variant A",
		RolloutPercent: 50,
	}

	newVariant := &domain.FlagVariant{
		ID:             "variant-1",
		ProjectID:      "project-1",
		FeatureID:      "feature-1",
		EnvironmentID:  1,
		Name:           "Variant B", // Changed
		RolloutPercent: 75,          // Changed
	}

	changes := BuildChangeDiff(oldVariant, newVariant)

	// Check that changes contain only editable fields
	expectedChanges := map[string]bool{
		"name":            true,
		"rollout_percent": true,
	}

	if len(changes) != len(expectedChanges) {
		t.Errorf("Expected %d changes, got %d", len(expectedChanges), len(changes))
	}

	for field := range expectedChanges {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected change for field %s", field)
		}
	}

	// Check that non-editable fields are not included
	if _, exists := changes["id"]; exists {
		t.Error("ID field should not be included in changes")
	}
	if _, exists := changes["project_id"]; exists {
		t.Error("ProjectID field should not be included in changes")
	}
}

func TestBuildInsertChanges_FlagVariant(t *testing.T) {
	// Test for FlagVariant insert
	variant := &domain.FlagVariant{
		ID:             "variant-1",
		ProjectID:      "project-1",
		FeatureID:      "feature-1",
		EnvironmentID:  1,
		Name:           "Variant A",
		RolloutPercent: 50,
	}

	changes := BuildInsertChanges(variant)

	// Check that all necessary fields are included
	expectedFields := map[string]bool{
		"id":              true, // pk field
		"project_id":      true, // required field
		"feature_id":      true, // required field
		"environment_id":  true, // required field
		"name":            true, // editable field
		"rollout_percent": true, // editable field
	}

	for field := range expectedFields {
		if _, exists := changes[field]; !exists {
			t.Errorf("Expected field %s in insert changes", field)
		}
	}
}

func TestComputeFeatureChanges_Integration(t *testing.T) {
	// Integration test for computeFeatureChanges
	oldFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Test Feature",
			Description: "Test Description",
			RolloutKey:  "user_id",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      true,
		DefaultValue: "default",
	}

	newFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Updated Test Feature",     // Changed name
			Description: "Updated Test Description", // Changed description
			RolloutKey:  "session_id",               // Changed rollout_key
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      false,         // Changed enabled
		DefaultValue: "new-default", // Changed default_value
	}

	service := &Service{}
	changes, err := service.computeFeatureChanges(oldFeature, newFeature, domain.EntityActionUpdate)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be 2 changes: one for features, one for feature_params
	if len(changes) != 2 {
		t.Errorf("Expected 2 changes, got %d", len(changes))
	}

	// Check changes for features table
	var featuresChanges *EntityChanges
	var paramsChanges *EntityChanges

	for i := range changes {
		if changes[i].EntityType == string(domain.EntityFeature) {
			featuresChanges = &changes[i]
		} else if changes[i].EntityType == string(domain.EntityFeatureParams) {
			paramsChanges = &changes[i]
		}
	}

	if featuresChanges == nil {
		t.Error("Expected changes for features table")
	} else {
		// Check that only editable fields from BasicFeature are included
		expectedFields := map[string]bool{
			"name":        true,
			"description": true,
			"rollout_key": true,
		}
		for field := range expectedFields {
			if _, exists := featuresChanges.Changes[field]; !exists {
				t.Errorf("Expected change for field %s in features table", field)
			}
		}
	}

	if paramsChanges == nil {
		t.Error("Expected changes for feature_params table")
	} else {
		// Check that only editable fields from FeatureParams are included
		expectedFields := map[string]bool{
			"enabled":       true,
			"default_value": true,
		}
		for field := range expectedFields {
			if _, exists := paramsChanges.Changes[field]; !exists {
				t.Errorf("Expected change for field %s in feature_params table", field)
			}
		}
	}
}

func TestComputeFeatureChanges_OnlyBasicFeature(t *testing.T) {
	// Test for the case when only BasicFeature fields are changed
	oldFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Test Feature",
			Description: "Test Description",
			RolloutKey:  "user_id",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      true,
		DefaultValue: "default",
	}

	newFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Updated Test Feature", // Changed only name
			Description: "Test Description",
			RolloutKey:  "user_id",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      true,      // Not changed
		DefaultValue: "default", // Not changed
	}

	service := &Service{}
	changes, err := service.computeFeatureChanges(oldFeature, newFeature, domain.EntityActionUpdate)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be only 1 change for features table
	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].EntityType != string(domain.EntityFeature) {
		t.Errorf("Expected changes for features table, got %s", changes[0].EntityType)
	}

	// Check that only name is included
	if len(changes[0].Changes) != 1 {
		t.Errorf("Expected 1 field change, got %d", len(changes[0].Changes))
	}

	if _, exists := changes[0].Changes["name"]; !exists {
		t.Error("Expected change for name field")
	}
}

func TestComputeFeatureChanges_OnlyParams(t *testing.T) {
	oldFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Test Feature",
			Description: "Test Description",
			RolloutKey:  "user_id",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      true,
		DefaultValue: "default",
	}

	newFeature := &domain.Feature{
		BasicFeature: domain.BasicFeature{
			ID:          "feature-1",
			ProjectID:   "project-1",
			Key:         "test-feature",
			Kind:        domain.FeatureKindSimple,
			Name:        "Test Feature",     // Not changed
			Description: "Test Description", // Not changed
			RolloutKey:  "user_id",          // Not changed
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Enabled:      false,         // Changed only enabled
		DefaultValue: "new-default", // Changed only default_value
	}

	service := &Service{}
	changes, err := service.computeFeatureChanges(oldFeature, newFeature, domain.EntityActionUpdate)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be only 1 change for feature_params table
	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}

	if changes[0].EntityType != string(domain.EntityFeatureParams) {
		t.Errorf("Expected changes for feature_params table, got %s", changes[0].EntityType)
	}

	// Check that enabled and default_value are included
	expectedFields := map[string]bool{
		"enabled":       true,
		"default_value": true,
	}
	for field := range expectedFields {
		if _, exists := changes[0].Changes[field]; !exists {
			t.Errorf("Expected change for field %s", field)
		}
	}
}
