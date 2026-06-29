package diff

import (
	"fmt"
	"reflect"
)

// Struct is a Differ for structs and pointers to structs.
type Struct struct {
	Fields map[string]Differ
}

func (s Struct) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, fmt.Errorf("state must not be nil")
	}

	for fieldName, differ := range s.Fields {
		if differ == nil {
			return false, NewError(state.Path, fmt.Errorf("field differ for %q must not be nil", fieldName))
		}
	}

	leftIsNil := isNilStructValue(left)
	rightIsNil := isNilStructValue(right)
	if leftIsNil && rightIsNil {
		return true, nil
	}
	if leftIsNil || rightIsNil {
		state.Different()
		return false, nil
	}

	leftValue, leftType, err := structValue(left)
	if err != nil {
		return false, wrapStructError(state, err)
	}

	rightValue, rightType, err := structValue(right)
	if err != nil {
		return false, wrapStructError(state, err)
	}

	if leftType != rightType {
		return false, NewError(state.Path, fmt.Errorf("left and right must have the same type: left=%s right=%s", leftType, rightType))
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

		state.Push(fieldKey, leftField, rightField)
		fieldSame, fieldErr := fieldDiffer.Diff(state, leftField, rightField)
		state.Pop()
		if fieldErr != nil {
			return false, fieldErr
		}
		if !fieldSame {
			allSame = false
		}
	}

	return allSame, nil
}

func isNilStructValue(value any) bool {
	if value == nil {
		return true
	}

	valueReflect := reflect.ValueOf(value)
	for valueReflect.Kind() == reflect.Pointer {
		if valueReflect.IsNil() {
			return true
		}
		valueReflect = valueReflect.Elem()
	}

	return false
}

func structValue(value any) (reflect.Value, reflect.Type, error) {
	valueReflect := reflect.ValueOf(value)
	valueType := valueReflect.Type()
	for valueType.Kind() == reflect.Pointer {
		valueReflect = valueReflect.Elem()
		valueType = valueReflect.Type()
	}

	if valueType.Kind() != reflect.Struct {
		return reflect.Value{}, nil, fmt.Errorf("value must be a struct or pointer to struct: %T", value)
	}

	return valueReflect, valueType, nil
}

func wrapStructError(state *State, err error) error {
	if err == nil {
		return nil
	}

	return NewError(state.Path, err)
}
