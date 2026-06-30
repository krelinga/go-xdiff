package diff

import (
	"errors"
	"fmt"
)

func Transform[T, U any](name string, f func(T) U, differ Differ) Differ {
	return transformDiffer[T, U]{
		name:   name,
		f:      f,
		differ: differ,
	}
}

type transformDiffer[T, U any] struct {
	name   string
	f      func(T) U
	differ Differ
}

func (t transformDiffer[T, U]) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, fmt.Errorf("state must not be nil")
	}

	if t.f == nil {
		return false, fmt.Errorf("transform function must not be nil")
	}

	leftVal, ok := left.(T)
	if !ok {
		return false, WrapError(state.Path, fmt.Errorf("left value is not of type %T", leftVal))
	}

	rightVal, ok := right.(T)
	if !ok {
		return false, WrapError(state.Path, fmt.Errorf("right value is not of type %T", rightVal))
	}

	leftTransformed := t.f(leftVal)
	rightTransformed := t.f(rightVal)

	differ := t.differ
	if differ == nil {
		differ = Default{}
	}

	key := NewTransformKey(t.name)
	same, err = state.DiffChild(key, leftTransformed, rightTransformed, differ)
	if _, ok := errors.AsType[Error](err); !ok {
		err = WrapError(append(state.Path, key), err)
	}
	return same, err
}
