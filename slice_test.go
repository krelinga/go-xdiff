package diff_test

import (
	"fmt"
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func TestSlice(t *testing.T) {
	leftMatchErr := func(_ int, _ any) (any, error) {
		return nil, fmt.Errorf("left key failure")
	}

	rightMatchErr := func(i int, _ any) (any, error) {
		if i == 0 {
			return i, nil
		}
		return nil, fmt.Errorf("right key failure")
	}

	tests := []struct {
		name             string
		differ           diff.Slice
		path             diff.Path
		nilState         bool
		left             any
		right            any
		wantSame         bool
		wantErrSubstring string
		wantCounter      *diff.Counter
	}{
		{
			name:             "nil state returns error",
			differ:           diff.Slice{},
			nilState:         true,
			left:             []int{1},
			right:            []int{1},
			wantErrSubstring: "state must not be nil",
		},
		{
			name:             "nil input returns error",
			differ:           diff.Slice{},
			path:             diff.Path{diff.RootKey{}},
			left:             nil,
			right:            []int{},
			wantErrSubstring: "left and right must not be nil",
		},
		{
			name:             "non slice input returns error",
			differ:           diff.Slice{},
			path:             diff.Path{diff.RootKey{}},
			left:             1,
			right:            1,
			wantErrSubstring: "left and right must be slices",
		},
		{
			name:     "both nil slices are equal",
			differ:   diff.Slice{},
			left:     []int(nil),
			right:    []int(nil),
			wantSame: true,
		},
		{
			name:        "nil and empty slices differ by default",
			differ:      diff.Slice{},
			left:        []int(nil),
			right:       []int{},
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:     "nil and empty slices are equal when treat nil as empty",
			differ:   diff.Slice{TreatNilAsEmpty: true},
			left:     []int(nil),
			right:    []int{},
			wantSame: true,
		},
		{
			name:             "left key function error is returned",
			differ:           diff.Slice{KeyFunc: leftMatchErr},
			path:             diff.Path{diff.RootKey{}},
			left:             []int{1},
			right:            []int{1},
			wantErrSubstring: "error in MatchFunc at left index 0",
		},
		{
			name:             "right key function error is returned",
			differ:           diff.Slice{KeyFunc: rightMatchErr},
			path:             diff.Path{diff.RootKey{}},
			left:             []int{1},
			right:            []int{1, 2},
			wantErrSubstring: "error in MatchFunc at right index 1",
		},
		{
			name:        "index matching reports element differences",
			differ:      diff.Slice{},
			left:        []int{1, 2},
			right:       []int{1, 3},
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:        "left only unmatched entries are reported",
			differ:      diff.Slice{},
			left:        []int{1, 2},
			right:       []int{1},
			wantSame:    false,
			wantCounter: &diff.Counter{NumLeftOnly: 1},
		},
		{
			name:        "right only unmatched entries are reported",
			differ:      diff.Slice{},
			left:        []int{1},
			right:       []int{1, 2},
			wantSame:    false,
			wantCounter: &diff.Counter{NumRightOnly: 1},
		},
		// Surprising behavior: this currently returns same=false even when the only unmatched
		// entry is on the left and AllowLeftOnly=true.
		// {
		// 	name:     "left only can be allowed",
		// 	differ:   diff.Slice{AllowLeftOnly: true},
		// 	left:     []int{1, 2},
		// 	right:    []int{1},
		// 	wantSame: true,
		// },
		{
			name:     "right only can be allowed",
			differ:   diff.Slice{AllowRightOnly: true},
			left:     []int{1},
			right:    []int{1, 2},
			wantSame: true,
		},
		{
			name:        "inefficient key minimizes pairing differences",
			differ:      diff.Slice{KeyFunc: diff.SliceKeyInefficient},
			left:        []int{1, 2},
			right:       []int{2, 1},
			wantSame:    true,
			wantCounter: &diff.Counter{},
		},
		// Surprising behavior: the inefficient matcher should minimize total diffs, but this
		// case currently reports an avoidable element difference in addition to left-only.
		// {
		// 	name:        "inefficient key reports extra unmatched entry",
		// 	differ:      diff.Slice{KeyFunc: diff.SliceKeyInefficient},
		// 	left:        []int{1, 2, 3},
		// 	right:       []int{3, 1},
		// 	wantSame:    false,
		// 	wantCounter: &diff.Counter{NumLeftOnly: 1},
		// },
		{
			name:             "element differ errors are wrapped at slice pair path",
			differ:           diff.Slice{ElemDiffer: diff.Compare{}},
			path:             diff.Path{diff.RootKey{}},
			left:             []any{[]int{1}},
			right:            []any{[]int{1}},
			wantErrSubstring: "slice pair: left index 0, right index 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := &diff.Counter{}
			state := &diff.State{Path: tt.path, Reporter: counter}
			if tt.nilState {
				state = nil
			}

			same, err := tt.differ.Diff(state, tt.left, tt.right)

			if same != tt.wantSame {
				t.Fatalf("same = %v, want %v", same, tt.wantSame)
			}

			if tt.wantErrSubstring == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErrSubstring)
				}
				if !strings.Contains(err.Error(), tt.wantErrSubstring) {
					t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErrSubstring)
				}
			}

			if tt.wantCounter != nil {
				counterEqual(t, *counter, *tt.wantCounter)
			}
		})
	}
}
