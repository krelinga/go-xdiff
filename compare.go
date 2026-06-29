package diff

import (
	"fmt"
	"reflect"
)

// Compare is a Differ that compares comparable values of the same type using ==.
type Compare struct{}

func (_ Compare) Diff(state *State, left, right any) (same bool, err error) {
	if left == nil && right == nil {
		return true, nil
	}
	if left == nil || right == nil {
		return false, NewError(state.path, fmt.Errorf("one value is nil while the other is non-nil"))
	}

	leftType := reflect.TypeOf(left)
	rightType := reflect.TypeOf(right)
	if leftType != rightType {
		return false, NewError(state.path, fmt.Errorf("left and right must have the same type: left=%s right=%s", leftType, rightType))
	}

	if !leftType.Comparable() {
		return false, NewError(state.path, fmt.Errorf("type %s is not comparable", leftType))
	}

	if left == right {
		return true, nil
	}

	if state != nil {
		state.Different()
	}
	return false, nil
}
