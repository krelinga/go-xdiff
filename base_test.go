package diff_test

import (
	"slices"
	"testing"

	diff "github.com/krelinga/go-xdiff"
)

// A helper function to compare two diff.Counter values for equality in tests.
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

// A recordingDiffer is a wrapper around a Differ that records whether it was called and the path at which it was called.
type recordingDiffer struct {
	// Should be set at construction time.
	DelegateTo diff.Differ

	// Should be read after Diff is called to determine if the differ was called.
	Called bool
	Path   diff.Path
}

func (r *recordingDiffer) Diff(state *diff.State, left, right any) (bool, error) {
	r.Called = true
	r.Path = slices.Clone(state.Path)
	return r.DelegateTo.Diff(state, left, right)
}

// A fakeDiffer is a Differ that calls a provided function to perform the diffing logic, or (more-commonly in tests) to force a specific behavior (e.g., returning an error or a specific same/different result).
type fakeDiffer func(state *diff.State, left, right any) (bool, error)

func (f fakeDiffer) Diff(state *diff.State, left, right any) (bool, error) {
	return f(state, left, right)
}
