package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

type testReporter struct {
	different int
}

func (_ *testReporter) Push(_ diff.Key, _, _ any) {}
func (_ *testReporter) Pop()                      {}
func (_ *testReporter) LeftOnly(_ diff.Key, _ any) {
}
func (_ *testReporter) RightOnly(_ diff.Key, _ any) {
}
func (r *testReporter) Different() {
	r.different++
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
			name:             "one nil returns wrapped error",
			state:            stateWithPath,
			left:             nil,
			right:            1,
			wantSame:         false,
			wantErrSubstring: "one value is nil while the other is non-nil",
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
			state:         &diff.State{Reporter: &testReporter{}},
			left:          42,
			right:         43,
			wantSame:      false,
			wantDifferent: 1,
		},
		{
			name:             "non comparable values return unsupported error",
			state:            stateWithPath,
			left:             []int{1},
			right:            []int{1},
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
				reporter, ok := tt.state.Reporter.(*testReporter)
				if !ok {
					t.Fatalf("expected testReporter in state")
				}
				if reporter.different != tt.wantDifferent {
					t.Fatalf("different = %d, want %d", reporter.different, tt.wantDifferent)
				}
			}
		})
	}
}
