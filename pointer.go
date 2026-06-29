package diff

import (
	"fmt"
	"reflect"
)

type Pointer struct {
	Elem Differ
	ByAddress bool
}

func (p Pointer) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, NewError(state.Path, fmt.Errorf("state must not be nil"))
	}

	if left == nil || right == nil {
		return false, NewError(state.Path, fmt.Errorf("left and right must not be nil"))
	}
	leftVal := reflect.ValueOf(left)
	rightVal := reflect.ValueOf(right)

	if leftVal.Kind() != reflect.Ptr || rightVal.Kind() != reflect.Ptr {
		return false, NewError(state.Path, fmt.Errorf("left and right must be pointers"))
	}

	if leftVal.IsNil() && rightVal.IsNil() {
		return true, nil
	}

	if leftVal.IsNil() || rightVal.IsNil() {
		state.Different()
		return false, nil
	}

	if p.ByAddress && leftVal.Pointer() == rightVal.Pointer() {
		return true, nil
	}

	elemDiffer := p.Elem
	if elemDiffer == nil {
		elemDiffer = Default{}
	}

	state.Push(PointerKey{}, leftVal.Elem().Interface(), rightVal.Elem().Interface())
	same, err = elemDiffer.Diff(state, leftVal.Elem().Interface(), rightVal.Elem().Interface())
	state.Pop()
	return same, err
}