package diff

import "fmt"

type All []Differ

func (a All) Diff(state *State, left, right any) (same bool, err error) {
	if state == nil {
		return false, fmt.Errorf("state must not be nil")
	}
	allSame := true
	for index, differ := range a {
		if differ == nil {
			return false, WrapError(state.Path, fmt.Errorf("differ at index %d must not be nil", index))
		}
		same, err := state.DiffChild(state.Path[len(state.Path)-1], left, right, differ)
		if err != nil {
			return false, err
		}
		if !same {
			allSame = false
		}
	}
	return allSame, nil
}
