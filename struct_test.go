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

func TestStruct(t *testing.T) {
	tests := []struct {
		name          string
		differ        diff.Struct
		path          diff.Path
		prepare       func() *recordingDiffer
		left          any
		right         any
		wantSame      bool
		wantErrSubstr string
		wantCounter   *diff.Counter
		check         func(t *testing.T, recorder *recordingDiffer)
	}{
		{
			name:        "default differ compares exported fields only",
			differ:      diff.Struct{},
			left:        structFixture{Name: "same", Count: 1, hidden: 1},
			right:       structFixture{Name: "same", Count: 2, hidden: 2},
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:     "field override differ is used",
			left:     structFixture{Name: "same", Count: 1},
			right:    structFixture{Name: "same", Count: 2},
			wantSame: true,
			prepare: func() *recordingDiffer {
				return &recordingDiffer{DelegateTo: fakeDiffer(func(_ *diff.State, _, _ any) (bool, error) {
					return true, nil
				})}
			},
			check: func(t *testing.T, recorder *recordingDiffer) {
				if !recorder.Called {
					t.Fatalf("expected custom field differ to be called")
				}
				if len(recorder.Path) != 1 {
					t.Fatalf("path length = %d, want 1", len(recorder.Path))
				}
				fieldKey, ok := recorder.Path[0].(diff.FieldKey)
				if !ok {
					t.Fatalf("path element type = %T, want diff.FieldKey", recorder.Path[0])
				}
				if fieldKey.Name != "Count" {
					t.Fatalf("field key name = %q, want %q", fieldKey.Name, "Count")
				}
			},
		},
		{
			name:          "pointers to structs return error",
			differ:        diff.Struct{},
			left:          &structFixture{Name: "same", Count: 1},
			right:         &structFixture{Name: "same", Count: 2},
			wantSame:      false,
			wantErrSubstr: "value must be a struct",
		},
		{
			name:          "nil field differ returns error",
			differ:        diff.Struct{Fields: map[string]diff.Differ{"Count": nil}},
			path:          diff.Path{diff.RootKey{}},
			left:          structFixture{},
			right:         structFixture{},
			wantSame:      false,
			wantErrSubstr: "field differ for \"Count\" must not be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var recorder *recordingDiffer
			if tt.prepare != nil {
				recorder = tt.prepare()
				tt.differ = diff.Struct{Fields: map[string]diff.Differ{"Count": recorder}}
			}

			gotCounter := &diff.Counter{}
			same, err := tt.differ.Diff(&diff.State{Path: tt.path, Reporter: gotCounter}, tt.left, tt.right)

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

			if tt.wantCounter != nil {
				counterEqual(t, *gotCounter, *tt.wantCounter)
			}

			if tt.check != nil {
				tt.check(t, recorder)
			}
		})
	}
}
