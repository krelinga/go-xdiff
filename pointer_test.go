package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func TestPointer(t *testing.T) {
	one := 1
	two := 2
	tests := []struct {
		name             string
		differ           diff.Pointer
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
			differ:           diff.Pointer{},
			nilState:         true,
			left:             &one,
			right:            &one,
			wantErrSubstring: "state must not be nil",
		},
		{
			name:             "nil input returns error",
			differ:           diff.Pointer{},
			path:             diff.Path{diff.RootKey{}},
			left:             nil,
			right:            &one,
			wantErrSubstring: "left and right must not be nil",
		},
		{
			name:             "non pointer input returns error",
			differ:           diff.Pointer{},
			path:             diff.Path{diff.RootKey{}},
			left:             1,
			right:            1,
			wantErrSubstring: "left and right must be pointers",
		},
		{
			name:     "both nil pointers are equal",
			differ:   diff.Pointer{},
			left:     (*int)(nil),
			right:    (*int)(nil),
			wantSame: true,
		},
		{
			name:        "one nil pointer reports difference",
			differ:      diff.Pointer{},
			left:        (*int)(nil),
			right:       &one,
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:     "same address with by address short circuits",
			differ:   diff.Pointer{ByAddress: true},
			left:     &one,
			right:    &one,
			wantSame: true,
		},
		{
			name:        "different pointed values report difference",
			differ:      diff.Pointer{},
			left:        &one,
			right:       &two,
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:        "custom elem differ is used",
			differ:      diff.Pointer{Elem: diff.Ignore{}},
			left:        &one,
			right:       &two,
			wantSame:    true,
			wantCounter: &diff.Counter{},
		},
		// TODO: add a test for the same value being stored at two different addresses.
		{
			name:             "elem differ error is wrapped at pointer path",
			differ:           diff.Pointer{Elem: diff.Compare{}},
			path:             diff.Path{diff.RootKey{}},
			left:             &[]int{1},
			right:            &[]int{1},
			wantErrSubstring: "pointer dereference",
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
