package diff

import (
	"fmt"
	"reflect"
)

// Default is a Differ with built-in default comparisons.
type Default struct{}

func (_ Default) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, fmt.Errorf("state must not be nil")
	}

	if left == nil && right == nil {
		return true, nil
	}

	if left == nil || right == nil {
		state.Different()
		return false, nil
	}

	leftType := reflect.TypeOf(left)
	rightType := reflect.TypeOf(right)
	if leftType.Comparable() && rightType.Comparable() {
		return Compare{}.Diff(state, left, right)
	} else if leftType.Kind() == reflect.Pointer && rightType.Kind() == reflect.Pointer {
		return Pointer{}.Diff(state, left, right)
	} else if leftType.Kind() == reflect.Struct && rightType.Kind() == reflect.Struct {
		return Struct{}.Diff(state, left, right)
	} else if leftType.Kind() == reflect.Map && rightType.Kind() == reflect.Map {
		return Map{}.Diff(state, left, right)
	} else if leftType.Kind() == reflect.Slice && rightType.Kind() == reflect.Slice {
		return Slice{}.Diff(state, left, right)
	}

	return false, WrapError(state.Path, fmt.Errorf("default comparisons are not supported for this type: left=%s right=%s", leftType, rightType))
}
