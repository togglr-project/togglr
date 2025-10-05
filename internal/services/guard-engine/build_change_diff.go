package guard_engine

import (
	"reflect"

	"github.com/togglr-project/togglr/internal/domain"
)

func BuildChangeDiff(oldValue, newValue any) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	vOld := reflect.ValueOf(oldValue)
	vNew := reflect.ValueOf(newValue)

	// Handle pointers
	if vOld.Kind() == reflect.Ptr {
		vOld = vOld.Elem()
	}
	if vNew.Kind() == reflect.Ptr {
		vNew = vNew.Elem()
	}

	typ := vOld.Type()

	// Process struct fields
	for i := range typ.NumField() {
		f := typ.Field(i)

		// If field has editable tag, process it
		const editableTrue = "true"
		if f.Tag.Get("editable") == editableTrue {
			fieldName := f.Tag.Get("db")
			if fieldName == "" {
				panic("db tag is required for editable fields")
			}

			oldVal := vOld.Field(i).Interface()
			newVal := vNew.Field(i).Interface()

			if !reflect.DeepEqual(oldVal, newVal) {
				changes[fieldName] = domain.ChangeValue{Old: oldVal, New: newVal}
			}
		}

		// If field is a nested struct (e.g., BasicFeature in Feature),
		// recursively process its fields
		if f.Anonymous || (f.Type.Kind() == reflect.Struct && f.Tag.Get("editable") == "") {
			oldField := vOld.Field(i)
			newField := vNew.Field(i)

			// Process nested struct
			nestedChanges := buildNestedChangeDiff(oldField, newField)
			for k, v := range nestedChanges {
				changes[k] = v
			}
		}
	}

	return changes
}

// buildNestedChangeDiff processes nested structs.
func buildNestedChangeDiff(oldField, newField reflect.Value) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	// Handle pointers
	if oldField.Kind() == reflect.Ptr {
		oldField = oldField.Elem()
	}
	if newField.Kind() == reflect.Ptr {
		newField = newField.Elem()
	}

	if !oldField.IsValid() || !newField.IsValid() {
		return changes
	}

	typ := oldField.Type()

	for i := range typ.NumField() {
		f := typ.Field(i)
		if f.Tag.Get("editable") != "true" { //nolint:goconst // string literal is fine here
			continue
		}

		fieldName := f.Tag.Get("db")
		if fieldName == "" {
			continue // Skip fields without db tag
		}

		oldVal := oldField.Field(i).Interface()
		newVal := newField.Field(i).Interface()

		if !reflect.DeepEqual(oldVal, newVal) {
			changes[fieldName] = domain.ChangeValue{Old: oldVal, New: newVal}
		}
	}

	return changes
}

// BuildInsertChanges creates changes for insert action, including all fields with editable tag
// and required fields for record creation (e.g., ID, ProjectID, FeatureID).
func BuildInsertChanges(entity any) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()

	// Process struct fields
	for i := range typ.NumField() {
		f := typ.Field(i)
		fieldName := f.Tag.Get("db")

		// If field has editable tag, include it
		if f.Tag.Get("editable") == "true" && fieldName != "" {
			val := v.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}

			continue
		}

		// If field is a nested struct (e.g., BasicFeature in Feature),
		// recursively process its fields
		if f.Anonymous || (f.Type.Kind() == reflect.Struct && f.Tag.Get("editable") == "") {
			field := v.Field(i)

			// Process nested struct
			nestedChanges := buildNestedInsertChanges(field)
			for k, v := range nestedChanges {
				changes[k] = v
			}
		}

		// Include required fields for record creation
		// (fields with pk tag or fields that are usually needed for creation)
		if fieldName != "" && (f.Tag.Get("pk") == "true" || isRequiredForInsert(fieldName)) {
			val := v.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}
		}
	}

	return changes
}

// buildNestedInsertChanges processes nested structs for insert.
func buildNestedInsertChanges(field reflect.Value) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	// Handle pointers
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}

	if !field.IsValid() {
		return changes
	}

	typ := field.Type()

	for i := range typ.NumField() {
		f := typ.Field(i)
		fieldName := f.Tag.Get("db")
		if fieldName == "" {
			continue
		}

		// Include all fields with editable tag
		if f.Tag.Get("editable") == "true" {
			val := field.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}

			continue
		}

		// Include required fields for record creation
		if f.Tag.Get("pk") == "true" || isRequiredForInsert(fieldName) {
			val := field.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}
		}
	}

	return changes
}

// isRequiredForInsert checks if a field is required for insert operation.
func isRequiredForInsert(fieldName string) bool {
	requiredFields := map[string]bool{
		"project_id":     true,
		"feature_id":     true,
		"environment_id": true,
	}

	return requiredFields[fieldName]
}
