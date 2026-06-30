package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name             string
		differ           diff.Map
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
			differ:           diff.Map{},
			nilState:         true,
			left:             map[string]int{"a": 1},
			right:            map[string]int{"a": 1},
			wantErrSubstring: "state must not be nil",
		},
		{
			name:             "nil input returns error",
			differ:           diff.Map{},
			path:             diff.Path{diff.RootKey{}},
			left:             nil,
			right:            map[string]int{},
			wantErrSubstring: "left and right must not be nil",
		},
		{
			name:             "non map input returns error",
			differ:           diff.Map{},
			path:             diff.Path{diff.RootKey{}},
			left:             []int{1},
			right:            []int{1},
			wantErrSubstring: "left and right must be maps",
		},
		{
			name:     "both nil maps are equal",
			differ:   diff.Map{},
			left:     map[string]int(nil),
			right:    map[string]int(nil),
			wantSame: true,
		},
		{
			name:        "one nil map reports difference when treat nil as empty is false",
			differ:      diff.Map{},
			left:        map[string]int(nil),
			right:       map[string]int{},
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:     "one nil map equals empty when treat nil as empty is true",
			differ:   diff.Map{TreatNilAsEmpty: true},
			left:     map[string]int(nil),
			right:    map[string]int{},
			wantSame: true,
		},
		{
			name:        "different mapped values report difference",
			differ:      diff.Map{},
			left:        map[string]int{"a": 1},
			right:       map[string]int{"a": 2},
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:        "left only key reported",
			differ:      diff.Map{},
			left:        map[string]int{"a": 1},
			right:       map[string]int{},
			wantSame:    false,
			wantCounter: &diff.Counter{NumLeftOnly: 1},
		},
		{
			name:        "right only key reported",
			differ:      diff.Map{},
			left:        map[string]int{},
			right:       map[string]int{"a": 1},
			wantSame:    false,
			wantCounter: &diff.Counter{NumRightOnly: 1},
		},
		{
			name:     "left only can be allowed",
			differ:   diff.Map{AllowLeftOnly: true},
			left:     map[string]int{"a": 1},
			right:    map[string]int{},
			wantSame: true,
		},
		{
			name:     "right only can be allowed",
			differ:   diff.Map{AllowRightOnly: true},
			left:     map[string]int{},
			right:    map[string]int{"a": 1},
			wantSame: true,
		},
		{
			name:             "value differ errors are wrapped at map key path",
			differ:           diff.Map{ValueDiffer: diff.Compare{}},
			path:             diff.Path{diff.RootKey{}},
			left:             map[string]any{"a": []int{1}},
			right:            map[string]any{"a": []int{1}},
			wantErrSubstring: "map key: a",
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
