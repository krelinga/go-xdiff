package diff

import (
	"fmt"
	"reflect"
)

// Struct is a Differ for structs.
type Struct struct {
	Fields map[string]Differ
}

func (s Struct) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, fmt.Errorf("state must not be nil")
	}

	for fieldName, differ := range s.Fields {
		if differ == nil {
			return false, WrapError(state.Path, fmt.Errorf("field differ for %q must not be nil", fieldName))
		}
	}

	if left == nil || right == nil {
		return false, WrapError(state.Path, fmt.Errorf("left and right must not be nil"))
	}

	leftValue, leftType, err := structValue(left)
	if err != nil {
		return false, WrapError(state.Path, err)
	}

	rightValue, rightType, err := structValue(right)
	if err != nil {
		return false, WrapError(state.Path, err)
	}

	if leftType != rightType {
		return false, WrapError(state.Path, fmt.Errorf("left and right must have the same type: left=%s right=%s", leftType, rightType))
	}

	var defaultDiffer Differ = Default{}
	allSame := true
	for index := range leftType.NumField() {
		fieldType := leftType.Field(index)
		if !fieldType.IsExported() {
			continue
		}

		fieldKey := NewFieldKey(fieldType.Name)
		leftField := leftValue.Field(index).Interface()
		rightField := rightValue.Field(index).Interface()

		fieldDiffer := defaultDiffer
		if differ, ok := s.Fields[fieldType.Name]; ok {
			fieldDiffer = differ
		}

		fieldSame, fieldErr := state.DiffChild(fieldKey, leftField, rightField, fieldDiffer)
		if fieldErr != nil {
			return false, fieldErr
		}
		if !fieldSame {
			allSame = false
		}
	}

	return allSame, nil
}

func structValue(value any) (reflect.Value, reflect.Type, error) {
	valueReflect := reflect.ValueOf(value)
	valueType := valueReflect.Type()

	if valueType.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("value must be a struct: %T", value)
	}

	return valueReflect, valueType, nil
}
