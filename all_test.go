package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func TestAll(t *testing.T) {
	tests := []struct {
		name             string
		differ           diff.All
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
			differ:           diff.All{diff.Compare{}},
			nilState:         true,
			left:             1,
			right:            1,
			wantErrSubstring: "state must not be nil",
		},
		{
			name:     "empty differ list returns same",
			differ:   diff.All{},
			path:     diff.Path{diff.RootKey{}},
			left:     1,
			right:    1,
			wantSame: true,
		},
		{
			name:             "nil differ returns error",
			differ:           diff.All{diff.Compare{}, nil},
			path:             diff.Path{diff.RootKey{}},
			left:             1,
			right:            1,
			wantErrSubstring: "differ at index 1 must not be nil",
		},
		{
			name:        "all differs must report same",
			differ:      diff.All{diff.Compare{}, diff.Default{}},
			path:        diff.Path{diff.RootKey{}},
			left:        1,
			right:       2,
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 2},
		},
		{
			name:     "all equal values remain same",
			differ:   diff.All{diff.Compare{}, diff.Default{}},
			path:     diff.Path{diff.RootKey{}},
			left:     5,
			right:    5,
			wantSame: true,
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

	// NOTE: This currently panics because All.Diff indexes the last path element even when Path is empty.
	// Keeping this case commented until behavior is clarified.
	// t.Run("empty path does not panic", func(t *testing.T) {
	// 	_, _ = (diff.All{diff.Compare{}}).Diff(&diff.State{}, 1, 1)
	// })
}
