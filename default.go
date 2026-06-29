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
	}

	return false, NewError(state.Path, fmt.Errorf("default comparisons are not supported for this type: left=%s right=%s", leftType, rightType))
}
