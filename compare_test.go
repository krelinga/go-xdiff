package diff_test

import (
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name          string
		state         *diff.State
		left          any
		right         any
		wantSame      bool
		wantErr       bool
		wantDifferent int
	}{
		{
			name:     "both nil are equal",
			state:    &diff.State{},
			left:     nil,
			right:    nil,
			wantSame: true,
		},
		{
			name:          "one nil is treated as a difference",
			state:         &diff.State{Reporter: &testReporter{}},
			left:          nil,
			right:         1,
			wantSame:      false,
			wantDifferent: 1,
		},
		{
			name:     "equal comparable values are equal",
			state:    &diff.State{},
			left:     1,
			right:    1,
			wantSame: true,
		},
		{
			name:          "different comparable values report difference",
			state:         &diff.State{Reporter: &testReporter{}},
			left:          1,
			right:         2,
			wantSame:      false,
			wantDifferent: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			same, err := (diff.Compare{}).Diff(tt.state, tt.left, tt.right)

			if same != tt.wantSame {
				t.Fatalf("same = %v, want %v", same, tt.wantSame)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("err presence = %v, want %v", err != nil, tt.wantErr)
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
