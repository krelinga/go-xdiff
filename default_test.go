package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func counterEqual(t *testing.T, got diff.Counter, want diff.Counter) {
	t.Helper()
	if got.Total() != want.Total() {
		t.Errorf("Total() = %d, want %d", got.Total(), want.Total())
	}
	if got.NumDifferent != want.NumDifferent {
		t.Errorf("NumDifferent = %d, want %d", got.NumDifferent, want.NumDifferent)
	}
	if got.NumLeftOnly != want.NumLeftOnly {
		t.Errorf("NumLeftOnly = %d, want %d", got.NumLeftOnly, want.NumLeftOnly)
	}
	if got.NumRightOnly != want.NumRightOnly {
		t.Errorf("NumRightOnly = %d, want %d", got.NumRightOnly, want.NumRightOnly)
	}
}

func TestDefault(t *testing.T) {
	stateWithPath := &diff.State{Path: diff.Path{diff.RootKey{}}}

	tests := []struct {
		name             string
		state            *diff.State
		left             any
		right            any
		wantSame         bool
		wantErrSubstring string
		wantDifferent    int
	}{
		{
			name:             "nil state returns error",
			state:            nil,
			left:             1,
			right:            1,
			wantSame:         false,
			wantErrSubstring: "state must not be nil",
		},
		{
			name:     "both nil are equal",
			state:    &diff.State{},
			left:     nil,
			right:    nil,
			wantSame: true,
		},
		{
			name:          "one nil is treated as a difference",
			state:         &diff.State{Reporter: &diff.Counter{}, Path: stateWithPath.Path},
			left:          nil,
			right:         1,
			wantSame:      false,
			wantDifferent: 1,
		},
		{
			name:     "comparable values equal",
			state:    &diff.State{},
			left:     42,
			right:    42,
			wantSame: true,
		},
		{
			name:          "comparable values different delegates and reports",
			state:         &diff.State{Reporter: &diff.Counter{}},
			left:          42,
			right:         43,
			wantSame:      false,
			wantDifferent: 1,
		},
		{
			name:             "non comparable values return unsupported error",
			state:            stateWithPath,
			left:             func() {},
			right:            func() {},
			wantSame:         false,
			wantErrSubstring: "default comparisons are not supported for this type",
		},
		{
			name:             "comparable mismatched types delegate to compare error",
			state:            stateWithPath,
			left:             1,
			right:            "1",
			wantSame:         false,
			wantErrSubstring: "left and right must have the same type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := diff.Default{}
			same, err := d.Diff(tt.state, tt.left, tt.right)

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

			if tt.wantDifferent > 0 {
				reporter, ok := tt.state.Reporter.(*diff.Counter)
				if !ok {
					t.Fatalf("expected diff.Counter in state")
				}
				counterEqual(t, *reporter, diff.Counter{NumDifferent: tt.wantDifferent})
			}
		})
	}
}
