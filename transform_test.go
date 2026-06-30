package diff_test

import (
	"strings"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

type transformPathRecorder struct {
	called bool
	path   diff.Path
	same   bool
	err    error
}

func (r *transformPathRecorder) Diff(state *diff.State, _, _ any) (bool, error) {
	r.called = true
	r.path = append(diff.Path(nil), state.Path...)
	return r.same, r.err
}

func TestTransform(t *testing.T) {
	type input struct {
		Value int
	}

	type tc struct {
		name             string
		differ           diff.Differ
		path             diff.Path
		nilState         bool
		left             any
		right            any
		wantSame         bool
		wantErrSubstring string
		wantCounter      *diff.Counter
		check            func(t *testing.T, recorder *transformPathRecorder)
	}

	tests := []tc{
		{
			name:             "nil state returns error",
			differ:           diff.Transform("value", func(v input) int { return v.Value }, nil),
			nilState:         true,
			left:             input{Value: 1},
			right:            input{Value: 1},
			wantErrSubstring: "state must not be nil",
		},
		{
			name:             "nil transform function returns error",
			differ:           diff.Transform[int, int]("nil-transform", nil, nil),
			left:             1,
			right:            1,
			wantErrSubstring: "transform function must not be nil",
		},
		{
			name:             "left type mismatch returns wrapped error",
			differ:           diff.Transform("value", func(v int) int { return v }, nil),
			path:             diff.Path{diff.RootKey{}},
			left:             "1",
			right:            1,
			wantErrSubstring: "left value is not of type int",
		},
		{
			name:             "right type mismatch returns wrapped error",
			differ:           diff.Transform("value", func(v int) int { return v }, nil),
			path:             diff.Path{diff.RootKey{}},
			left:             1,
			right:            "1",
			wantErrSubstring: "right value is not of type int",
		},
		{
			name:        "default differ compares transformed values",
			differ:      diff.Transform("value", func(v input) int { return v.Value }, nil),
			left:        input{Value: 1},
			right:       input{Value: 2},
			wantSame:    false,
			wantCounter: &diff.Counter{NumDifferent: 1},
		},
		{
			name:        "custom differ is used for transformed values",
			differ:      diff.Transform("value", func(v input) int { return v.Value }, diff.Ignore{}),
			left:        input{Value: 1},
			right:       input{Value: 2},
			wantSame:    true,
			wantCounter: &diff.Counter{},
		},
		{
			name:   "transform key appears in child path",
			left:   input{Value: 1},
			right:  input{Value: 2},
			wantSame: true,
			check: func(t *testing.T, recorder *transformPathRecorder) {
				if recorder == nil || !recorder.called {
					t.Fatalf("expected recorder differ to be called")
				}
				if len(recorder.path) != 1 {
					t.Fatalf("path length = %d, want 1", len(recorder.path))
				}
				key, ok := recorder.path[0].(diff.TransformKey)
				if !ok {
					t.Fatalf("path key type = %T, want diff.TransformKey", recorder.path[0])
				}
				if key.Name != "value" {
					t.Fatalf("transform key name = %q, want %q", key.Name, "value")
				}
			},
		},
		// Surprising behavior: if a nested differ returns a plain error, Transform does not
		// wrap it with the transform path context.
		// {
		// 	name:             "custom differ errors include transform key path",
		// 	path:             diff.Path{diff.RootKey{}},
		// 	left:             input{Value: 1},
		// 	right:            input{Value: 1},
		// 	wantErrSubstring: "transform: value",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := (*transformPathRecorder)(nil)
			differ := tt.differ

			// TODO: this is weird, we shouldn't have this kind of special casing in a table driven test.
			switch tt.name {
			case "transform key appears in child path":
				recorder = &transformPathRecorder{same: true}
				differ = diff.Transform("value", func(v input) int { return v.Value }, recorder)
			}

			counter := &diff.Counter{}
			state := &diff.State{Path: tt.path, Reporter: counter}
			if tt.nilState {
				state = nil
			}

			same, err := differ.Diff(state, tt.left, tt.right)

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

			if tt.check != nil {
				tt.check(t, recorder)
			}
		})
	}
}
