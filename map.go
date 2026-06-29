package diff

import (
	"fmt"
	"reflect"
)

type Map struct {
	// The differ to use for comparing the values of the map.
	// If nil, the default differ will be used.
	ValueDiffer Differ

	// If true, unmatched keys in the left map will be ignored.
	// If false, unmatched keys in the left map will be reported as differences.
	AllowLeftOnly bool

	// If true, unmatched keys in the right map will be ignored.
	// If false, unmatched keys in the right map will be reported as differences.
	AllowRightOnly bool

	// If true, nil maps will be treated as empty maps.
	// If false, nil maps will be treated as nil and will be reported as differences if the other map is not nil.
	TreatNilAsEmpty bool
}

func (m Map) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, fmt.Errorf("state must not be nil")
	}

	if left == nil || right == nil {
		return false, WrapError(state.Path, fmt.Errorf("left and right must not be nil"))
	}

	leftVal := reflect.ValueOf(left)
	rightVal := reflect.ValueOf(right)

	if leftVal.Kind() != reflect.Map || rightVal.Kind() != reflect.Map {
		return false, WrapError(state.Path, fmt.Errorf("left and right must be maps"))
	}

	if leftVal.IsNil() && rightVal.IsNil() {
		return true, nil
	}

	if !m.TreatNilAsEmpty && (leftVal.IsNil() || rightVal.IsNil()) {
		state.Different()
		return false, nil
	}

	valueDiffer := m.ValueDiffer
	if valueDiffer == nil {
		valueDiffer = Default{}
	}

	leftKeys := leftVal.MapKeys()
	rightKeys := rightVal.MapKeys()

	leftKeySet := make(map[any]struct{}, len(leftKeys))
	for _, key := range leftKeys {
		leftKeySet[key.Interface()] = struct{}{}
	}

	rightKeySet := make(map[any]struct{}, len(rightKeys))
	for _, key := range rightKeys {
		rightKeySet[key.Interface()] = struct{}{}
	}

	allSame := true

	for _, leftKey := range leftKeys {
		leftKeyInterface := leftKey.Interface()
		if _, ok := rightKeySet[leftKeyInterface]; !ok {
			if !m.AllowLeftOnly {
				state.LeftOnly(NewMapKey(leftKeyInterface), leftVal.MapIndex(leftKey).Interface())
				allSame = false
			}
		} else {
			leftValue := leftVal.MapIndex(leftKey).Interface()
			rightValue := rightVal.MapIndex(leftKey).Interface()

			same, err := state.DiffChild(NewMapKey(leftKeyInterface), leftValue, rightValue, valueDiffer)
			if err != nil {
				return false, err
			}
			if !same {
				allSame = false
			}
		}
	}

	for _, rightKey := range rightKeys {
		rightKeyInterface := rightKey.Interface()
		if _, ok := leftKeySet[rightKeyInterface]; !ok {
			if !m.AllowRightOnly {
				state.RightOnly(NewMapKey(rightKeyInterface), rightVal.MapIndex(rightKey).Interface())
				allSame = false
			}
		}
	}

	return allSame, nil
}