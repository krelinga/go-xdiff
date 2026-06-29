package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

type structFixture struct {
	Name   string
	Count  int
	hidden int
}

type stubDiffer struct {
	called bool
	path   diff.Path
	same   bool
	err    error
}

func (s *stubDiffer) Diff(state *diff.State, _, _ any) (bool, error) {
	s.called = true
	s.path = append(diff.Path(nil), state.Path...)
	return s.same, s.err
}

func TestStruct(t *testing.T) {
	tests := []struct {
		name          string
		differ        diff.Struct
		prepare       func() *stubDiffer
		state         *diff.State
		left          any
		right         any
		wantSame      bool
		wantErrSubstr string
		wantDifferent int
		check         func(t *testing.T, stub *stubDiffer)
	}{
		{
			name:          "default differ compares exported fields only",
			differ:        diff.Struct{},
			state:         &diff.State{Reporter: &testReporter{}},
			left:          structFixture{Name: "same", Count: 1, hidden: 1},
			right:         structFixture{Name: "same", Count: 2, hidden: 2},
			wantSame:      false,
			wantDifferent: 1,
		},
		{
			name:     "field override differ is used",
			state:    &diff.State{},
			left:     structFixture{Name: "same", Count: 1},
			right:    structFixture{Name: "same", Count: 2},
			wantSame: true,
			prepare: func() *stubDiffer {
				return &stubDiffer{same: true}
			},
			check: func(t *testing.T, stub *stubDiffer) {
				if !stub.called {
					t.Fatalf("expected custom field differ to be called")
				}
				if len(stub.path) != 1 {
					t.Fatalf("path length = %d, want 1", len(stub.path))
				}
				fieldKey, ok := stub.path[0].(diff.FieldKey)
				if !ok {
					t.Fatalf("path element type = %T, want diff.FieldKey", stub.path[0])
				}
				if fieldKey.Name != "Count" {
					t.Fatalf("field key name = %q, want %q", fieldKey.Name, "Count")
				}
			},
		},
		{
			name:          "pointers to structs are supported",
			differ:        diff.Struct{},
			state:         &diff.State{Reporter: &testReporter{}},
			left:          &structFixture{Name: "same", Count: 1},
			right:         &structFixture{Name: "same", Count: 2},
			wantSame:      false,
			wantDifferent: 1,
		},
		{
			name:          "nil field differ returns error",
			differ:        diff.Struct{Fields: map[string]diff.Differ{"Count": nil}},
			state:         &diff.State{Path: diff.Path{diff.RootKey{}}},
			left:          structFixture{},
			right:         structFixture{},
			wantSame:      false,
			wantErrSubstr: "field differ for \"Count\" must not be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stub *stubDiffer
			if tt.prepare != nil {
				stub = tt.prepare()
				tt.differ = diff.Struct{Fields: map[string]diff.Differ{"Count": stub}}
			}

			same, err := tt.differ.Diff(tt.state, tt.left, tt.right)

			if same != tt.wantSame {
				t.Fatalf("same = %v, want %v", same, tt.wantSame)
			}

			if tt.wantErrSubstr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErrSubstr)
				}
				if !strings.Contains(err.Error(), tt.wantErrSubstr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErrSubstr)
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

			if tt.check != nil {
				tt.check(t, stub)
			}
		})
	}
}
