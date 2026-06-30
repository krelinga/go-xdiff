package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name             string
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
			nilState:         true,
			left:             1,
			right:            1,
			wantSame:         false,
			wantErrSubstring: "state must not be nil",
		},
		{
			name:     "both nil are equal",
			left:     nil,
			right:    nil,
			wantSame: true,
		},
		{
			name:        "one nil is treated as a difference",
			path:        diff.Path{diff.RootKey{}},
			left:        nil,
			right:       1,
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:     "comparable values equal",
			left:     42,
			right:    42,
			wantSame: true,
		},
		{
			name:        "comparable values different delegates and reports",
			left:        42,
			right:       43,
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:             "non comparable values return unsupported error",
			path:             diff.Path{diff.RootKey{}},
			left:             func() {},
			right:            func() {},
			wantSame:         false,
			wantErrSubstring: "default comparisons are not supported for this type",
		},
		{
			name:             "comparable mismatched types delegate to compare error",
			path:             diff.Path{diff.RootKey{}},
			left:             1,
			right:            "1",
			wantSame:         false,
			wantErrSubstring: "left and right must have the same type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := diff.Default{}
			counter := &diff.Counter{}
			state := &diff.State{Path: tt.path, Reporter: counter}
			if tt.nilState {
				state = nil
			}

			same, err := d.Diff(state, tt.left, tt.right)

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
