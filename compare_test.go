package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name          string
		path          diff.Path
		nilState      bool
		left          any
		right         any
		wantSame      bool
		wantErr       bool
		wantErrSubstr string
		wantCounter   *diff.Counter
	}{
		{
			name:          "nil state returns error",
			nilState:      true,
			left:          1,
			right:         1,
			wantSame:      false,
			wantErr:       true,
			wantErrSubstr: "state must not be nil",
		},
		{
			name:     "both nil are equal",
			left:     nil,
			right:    nil,
			wantSame: true,
		},
		{
			name:        "one nil is treated as a difference",
			left:        nil,
			right:       1,
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:     "equal comparable values are equal",
			left:     1,
			right:    1,
			wantSame: true,
		},
		{
			name:        "different comparable values report difference",
			left:        1,
			right:       2,
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:          "mismatched comparable types return error",
			path:          diff.Path{diff.RootKey{}},
			left:          1,
			right:         "1",
			wantSame:      false,
			wantErr:       true,
			wantErrSubstr: "left and right must have the same type",
		},
		{
			name:          "non comparable type returns error",
			path:          diff.Path{diff.RootKey{}},
			left:          []int{1},
			right:         []int{1},
			wantSame:      false,
			wantErr:       true,
			wantErrSubstr: "type []int is not comparable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := &diff.Counter{}
			state := &diff.State{Path: tt.path, Reporter: counter}
			if tt.nilState {
				state = nil
			}

			same, err := (diff.Compare{}).Diff(state, tt.left, tt.right)

			if same != tt.wantSame {
				t.Fatalf("same = %v, want %v", same, tt.wantSame)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("err presence = %v, want %v", err != nil, tt.wantErr)
			}

			if tt.wantErrSubstr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErrSubstr)) {
				t.Fatalf("err = %v, want substring %q", err, tt.wantErrSubstr)
			}

			if tt.wantCounter != nil {
				counterEqual(t, *counter, *tt.wantCounter)
			}
		})
	}
}
