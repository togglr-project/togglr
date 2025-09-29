package pending_changes

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/togglr-project/togglr/internal/domain"
)

// ApplyChangesToEntity applies changes to an entity using reflection and db tags.
// It takes the current entity, changes map, and applies the changes to the entity.
func ApplyChangesToEntity(currentEntity any, changes map[string]domain.ChangeValue) error {
	v := reflect.ValueOf(currentEntity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("entity must be a struct, got %s", v.Kind())
	}

	typ := v.Type()

	// Process struct fields
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fieldName := f.Tag.Get("db")
		if fieldName == "" {
			continue // Skip fields without db tag
		}

		// Check if this field has changes
		change, hasChange := changes[fieldName]
		if !hasChange {
			continue
		}

		field := v.Field(i)
		if !field.CanSet() {
			continue // Skip unexported fields
		}

		// Apply the change
		if err := applyChangeToField(field, change.New); err != nil {
			return fmt.Errorf("apply change to field %s: %w", fieldName, err)
		}
	}

	return nil
}

// CreateEntityFromChanges creates a new entity from changes using reflection and db tags.
// It takes the entity type, changes map, and creates a new entity with the changes applied.
func CreateEntityFromChanges(entityType reflect.Type, changes map[string]domain.ChangeValue) (any, error) {
	// Create a new instance of the entity type
	entity := reflect.New(entityType).Interface()

	// Apply changes to the new entity
	if err := ApplyChangesToEntity(entity, changes); err != nil {
		return nil, fmt.Errorf("apply changes to new entity: %w", err)
	}

	return entity, nil
}

// applyChangeToField applies a change value to a specific field using reflection.
func applyChangeToField(field reflect.Value, newValue any) error {
	if newValue == nil {
		// Handle nil values for pointer fields
		if field.Kind() == reflect.Ptr {
			field.Set(reflect.Zero(field.Type()))
		}

		return nil
	}

	fieldType := field.Type()

	// Handle different field types
	switch fieldType.Kind() {
	case reflect.String:
		if str, ok := newValue.(string); ok {
			field.SetString(str)
		} else {
			return fmt.Errorf("expected string, got %T", newValue)
		}

	case reflect.Bool:
		if b, ok := newValue.(bool); ok {
			field.SetBool(b)
		} else {
			return fmt.Errorf("expected bool, got %T", newValue)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, ok := convertToInt64(newValue); ok {
			field.SetInt(num)
		} else {
			return fmt.Errorf("expected number, got %T", newValue)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if num, ok := convertToUint64(newValue); ok {
			field.SetUint(num)
		} else {
			return fmt.Errorf("expected number, got %T", newValue)
		}

	case reflect.Float32, reflect.Float64:
		if num, ok := convertToFloat64(newValue); ok {
			field.SetFloat(num)
		} else {
			return fmt.Errorf("expected number, got %T", newValue)
		}

	case reflect.Ptr:
		// Handle pointer fields
		if field.IsNil() {
			field.Set(reflect.New(fieldType.Elem()))
		}

		// Recursively apply to the pointed-to value
		return applyChangeToField(field.Elem(), newValue)

	case reflect.Struct:
		// Handle custom types (like domain.FeatureID, domain.ProjectID, etc.)
		if str, ok := newValue.(string); ok {
			// Try to set the string value directly
			field.SetString(str)
		} else if fieldType == reflect.TypeOf(domain.BooleanExpression{}) {
			// Handle BooleanExpression specially
			if err := convertToBooleanExpression(newValue, field); err != nil {
				return fmt.Errorf("convert to BooleanExpression: %w", err)
			}
		} else {
			return fmt.Errorf("expected string for custom type, got %T", newValue)
		}

	default:
		return fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}

	return nil
}

// convertToInt64 converts various numeric types to int64.
func convertToInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	case string:
		if num, err := strconv.ParseInt(v, 10, 64); err == nil {
			return num, true
		}
	}

	return 0, false
}

// convertToUint64 converts various numeric types to uint64.
func convertToUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case int:
		if v >= 0 {
			return uint64(v), true
		}
	case int8:
		if v >= 0 {
			return uint64(v), true
		}
	case int16:
		if v >= 0 {
			return uint64(v), true
		}
	case int32:
		if v >= 0 {
			return uint64(v), true
		}
	case int64:
		if v >= 0 {
			return uint64(v), true
		}
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case float32:
		if v >= 0 {
			return uint64(v), true
		}
	case float64:
		if v >= 0 {
			return uint64(v), true
		}
	case string:
		if num, err := strconv.ParseUint(v, 10, 64); err == nil {
			return num, true
		}
	}

	return 0, false
}

// convertToFloat64 converts various numeric types to float64.
func convertToFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, true
		}
	}

	return 0, false
}

// convertToBooleanExpression converts a map[string]interface{} to domain.BooleanExpression.
func convertToBooleanExpression(value any, field reflect.Value) error {
	// Convert to JSON and back to BooleanExpression
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal to JSON: %w", err)
	}

	var expr domain.BooleanExpression
	if err := json.Unmarshal(jsonData, &expr); err != nil {
		return fmt.Errorf("unmarshal to BooleanExpression: %w", err)
	}

	field.Set(reflect.ValueOf(expr))

	return nil
}
