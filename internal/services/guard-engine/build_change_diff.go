package guard_engine

import (
	"reflect"

	"github.com/togglr-project/togglr/internal/domain"
)

func BuildChangeDiff(old, new any) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	vOld := reflect.ValueOf(old)
	vNew := reflect.ValueOf(new)

	// Обрабатываем указатели
	if vOld.Kind() == reflect.Ptr {
		vOld = vOld.Elem()
	}
	if vNew.Kind() == reflect.Ptr {
		vNew = vNew.Elem()
	}

	typ := vOld.Type()

	// Обрабатываем поля структуры
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		// Если поле имеет тег editable, обрабатываем его
		if f.Tag.Get("editable") == "true" {
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

		// Если поле является вложенной структурой (например, BasicFeature в Feature),
		// рекурсивно обрабатываем его поля
		if f.Anonymous || (f.Type.Kind() == reflect.Struct && f.Tag.Get("editable") == "") {
			oldField := vOld.Field(i)
			newField := vNew.Field(i)

			// Обрабатываем вложенную структуру
			nestedChanges := buildNestedChangeDiff(oldField, newField)
			for k, v := range nestedChanges {
				changes[k] = v
			}
		}
	}

	return changes
}

// buildNestedChangeDiff обрабатывает вложенные структуры
func buildNestedChangeDiff(oldField, newField reflect.Value) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	// Обрабатываем указатели
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

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Tag.Get("editable") != "true" {
			continue
		}

		fieldName := f.Tag.Get("db")
		if fieldName == "" {
			continue // Пропускаем поля без db тега
		}

		oldVal := oldField.Field(i).Interface()
		newVal := newField.Field(i).Interface()

		if !reflect.DeepEqual(oldVal, newVal) {
			changes[fieldName] = domain.ChangeValue{Old: oldVal, New: newVal}
		}
	}

	return changes
}

// BuildInsertChanges создает изменения для insert action, включая все поля с тегом editable
// и обязательные поля для создания записи (например, ID, ProjectID, FeatureID)
func BuildInsertChanges(entity any) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()

	// Обрабатываем поля структуры
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fieldName := f.Tag.Get("db")

		// Если поле имеет тег editable, включаем его
		if f.Tag.Get("editable") == "true" && fieldName != "" {
			val := v.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}
			continue
		}

		// Если поле является вложенной структурой (например, BasicFeature в Feature),
		// рекурсивно обрабатываем его поля
		if f.Anonymous || (f.Type.Kind() == reflect.Struct && f.Tag.Get("editable") == "") {
			field := v.Field(i)

			// Обрабатываем вложенную структуру
			nestedChanges := buildNestedInsertChanges(field)
			for k, v := range nestedChanges {
				changes[k] = v
			}
		}

		// Включаем обязательные поля для создания записи
		// (поля с тегом pk или поля, которые обычно нужны для создания)
		if fieldName != "" && (f.Tag.Get("pk") == "true" || isRequiredForInsert(fieldName)) {
			val := v.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}
		}
	}

	return changes
}

// buildNestedInsertChanges обрабатывает вложенные структуры для insert
func buildNestedInsertChanges(field reflect.Value) map[string]domain.ChangeValue {
	changes := make(map[string]domain.ChangeValue)

	// Обрабатываем указатели
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}

	if !field.IsValid() {
		return changes
	}

	typ := field.Type()

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fieldName := f.Tag.Get("db")
		if fieldName == "" {
			continue
		}

		// Включаем все поля с тегом editable
		if f.Tag.Get("editable") == "true" {
			val := field.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}
			continue
		}

		// Включаем обязательные поля для создания записи
		if f.Tag.Get("pk") == "true" || isRequiredForInsert(fieldName) {
			val := field.Field(i).Interface()
			changes[fieldName] = domain.ChangeValue{New: val}
		}
	}

	return changes
}

// isRequiredForInsert проверяет, является ли поле обязательным для insert операции
func isRequiredForInsert(fieldName string) bool {
	requiredFields := map[string]bool{
		"project_id":     true,
		"feature_id":     true,
		"environment_id": true,
	}

	return requiredFields[fieldName]
}
