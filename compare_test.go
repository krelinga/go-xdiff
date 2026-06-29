package diff_test

import (
	"strings"
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
		wantErrSubstr string
		wantDifferent int
	}{
		{
			name:          "nil state returns error",
			state:         nil,
			left:          1,
			right:         1,
			wantSame:      false,
			wantErr:       true,
			wantErrSubstr: "state must not be nil",
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
		{
			name:          "mismatched comparable types return error",
			state:         &diff.State{Path: diff.Path{diff.RootKey{}}},
			left:          1,
			right:         "1",
			wantSame:      false,
			wantErr:       true,
			wantErrSubstr: "left and right must have the same type",
		},
		{
			name:          "non comparable type returns error",
			state:         &diff.State{Path: diff.Path{diff.RootKey{}}},
			left:          []int{1},
			right:         []int{1},
			wantSame:      false,
			wantErr:       true,
			wantErrSubstr: "type []int is not comparable",
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

			if tt.wantErrSubstr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErrSubstr)) {
				t.Fatalf("err = %v, want substring %q", err, tt.wantErrSubstr)
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
